// Copyright 2021 The themix authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"

	"go.themix.io/transport/info"
	"go.themix.io/transport/message"
	"go.uber.org/zap"
)

var (
	consensusPrefix = "/cons"
	clientPrefix    = "/client"
	streamBufSize   = 40960
)

type streamWriter struct {
	peerID  info.IDType
	msgc    chan *message.ConsMessage
	encoder *gob.Encoder
	flusher http.Flusher
	// isReady bool
	// mu      sync.Mutex
}

type streamReader struct {
	msgc    chan *message.ConsMessage
	decoder *gob.Decoder
}

type peer struct {
	peerID info.IDType
	addr   string
	reader *streamReader
	writer *streamWriter
}

// HTTPTransport is responsible for message exchange among nodes
type HTTPTransport struct {
	id    info.IDType
	peers map[info.IDType]*peer
	msgc  chan *message.ConsMessage
	mu    sync.Mutex
}

func init() {
	gob.Register(&message.ConsMessage{})
}

// Broadcast msg to all peers
func (tp *HTTPTransport) Broadcast(msg *message.ConsMessage) {

	tp.mu.Lock()
	defer tp.mu.Unlock()

	msg.From = tp.id
	tp.msgc <- msg
	for _, p := range tp.peers {
		if p.writer != nil {
			p.writer.msgc <- msg
		}
	}
}

// Send sends the given message to msg.To
// func (tp *Transport) Send(msg *ConsMessage, to IDType) {
// 	p := tp.peers[to]
// 	p.writer.msgc <- msg
// 	// p.writer.encoder.Encode(*msg)
// 	// p.writer.flusher.Flush()
// }

// InitTransport executes transport layer initiliazation, which returns transport, a channel
// for received ConsMessage, a channel for received requests, and a channel for reply
func InitTransport(lg *zap.Logger, id info.IDType, port int, peers []string) (*HTTPTransport,
	chan *message.ConsMessage, chan []byte, chan []byte) {
	msgc := make(chan *message.ConsMessage, streamBufSize)
	tp := &HTTPTransport{id: id, peers: make(map[info.IDType]*peer), msgc: msgc, mu: sync.Mutex{}}
	for i, p := range peers {
		if index := info.IDType(i); index != id {
			tp.peers[index] = &peer{peerID: index, addr: p}
		}
	}

	tp.connect()

	reqc := make(chan []byte, streamBufSize)
	repc := make(chan []byte, streamBufSize)

	rprocessor := &ClientMsgProcessor{lg: lg, id: id, reqc: reqc, repc: repc}

	mux := http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	mux.Handle(consensusPrefix, tp)
	mux.Handle(consensusPrefix+"/", tp)
	mux.Handle(clientPrefix, rprocessor)
	mux.Handle(clientPrefix+"/", rprocessor)

	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux}
	server.SetKeepAlivesEnabled(true)

	go server.ListenAndServe()

	return tp, msgc, reqc, repc
}

func (tp *HTTPTransport) connect() {
	for _, p := range tp.peers {
		go dial(p, tp.id, tp.msgc)
	}
}

func dial(p *peer, id info.IDType, msgc chan *message.ConsMessage) {
	var r *streamReader
	for {
		req, err := http.NewRequest("GET",
			p.addr+consensusPrefix+"/"+strconv.FormatUint(uint64(id), 10), nil)

		if err != nil {
			log.Fatal(err)
		}

		t := &http.Transport{
			Dial: (&net.Dialer{
				KeepAlive: 120 * time.Second,
			}).Dial,
		}

		resp, err := t.RoundTrip(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		r = &streamReader{msgc: msgc,
			decoder: gob.NewDecoder(resp.Body)}
		break
	}
	p.reader = r
	go r.run()
}

func (sr *streamReader) run() {
	for {
		var m message.ConsMessage
		if err := sr.decoder.Decode(&m); err != nil {
			log.Fatal("decode error:", err)
		}
		sr.msgc <- &m
	}
}

func (sw *streamWriter) run() {
	for {
		m := <-sw.msgc
		err := sw.encoder.Encode(m)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		sw.flusher.Flush()
	}
}

func (tp *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	fromStr := path.Base(r.URL.Path)
	fromID, _ := strconv.ParseUint(fromStr, 10, 64)
	p := tp.peers[info.IDType(fromID)]

	enc := gob.NewEncoder(w)

	p.writer = &streamWriter{msgc: make(chan *message.ConsMessage, streamBufSize),
		encoder: enc, flusher: w.(http.Flusher), peerID: info.IDType(fromID)}

	p.writer.run()
}

// ClientMsgProcessor is responsible for listening and processing requests from clients
type ClientMsgProcessor struct {
	num  int
	lg   *zap.Logger
	id   info.IDType
	reqc chan []byte
	repc chan []byte
}

func (cmsgProcessor *ClientMsgProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// v, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	rprocessor.lg.Error("read client request error:", zap.Error(err))
	// }

	// key := r.RequestURI
	defer r.Body.Close()
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed on PUT", http.StatusBadRequest)
		fmt.Println("Failed on PUT", http.StatusBadRequest)
		return
	}

	if len(v) == 0 {
		fmt.Println("error request size sent to : ", len(v), cmsgProcessor.id)
		v = make([]byte, 4)
	}
	// bs := make([]byte, 4)
	// binary.LittleEndian.PutUint32(bs, uint32(cmsgProcessor.num))

	if len(v) > 0 {
		cmsgProcessor.lg.Info("receive request",
			zap.Int("proposer", int(cmsgProcessor.id)),
			zap.Int("seq", cmsgProcessor.num),
			zap.Int("content", int(v[0])))
	}

	cmsgProcessor.reqc <- v

	rep := <-cmsgProcessor.repc

	cmsgProcessor.lg.Info("reply request",
		zap.Int("proposer", int(cmsgProcessor.id)),
		zap.Int("seq", cmsgProcessor.num),
		zap.Int("content", int(v[0])))

	cmsgProcessor.num++
	w.Write(rep)
}
