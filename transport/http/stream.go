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
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/perlin-network/noise"
	myecdsa "go.themix.io/crypto/ecdsa"
	"go.themix.io/crypto/sha256"
	"go.themix.io/transport/proto/consmsgpb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	clientPrefix  = "/client"
	streamBufSize = 40960 * 16
)

type Peer struct {
	PeerID    uint32
	Addr      string
	PublicKey *ecdsa.PrivateKey
}

// HTTPTransport is responsible for message exchange among nodes
type HTTPTransport struct {
	id         uint32
	node       *noise.Node
	Peers      map[uint32]*Peer
	PrivateKey ecdsa.PrivateKey
	msgc       chan *consmsgpb.WholeMessage
	proposal   [][]byte
}

type NoiseMessage struct {
	Msg *consmsgpb.WholeMessage
}

func (m NoiseMessage) Marshal() []byte {
	data, err := proto.Marshal(m.Msg)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func UnmarshalNoiseMessage(buf []byte) (NoiseMessage, error) {
	m := NoiseMessage{Msg: new(consmsgpb.WholeMessage)}
	err := proto.Unmarshal(buf, m.Msg)
	if err != nil {
		return NoiseMessage{}, err
	}
	return m, nil
}

// Broadcast msg to all peers
func (tp *HTTPTransport) Broadcast(msg *consmsgpb.WholeMessage) {
	msg.From = tp.id
	Sign(msg, &tp.PrivateKey)
	tp.msgc <- msg

	for _, p := range tp.Peers {
		if p != nil {
			go tp.SendMessage(p.PeerID, msg)
		}
	}
}

// InitTransport executes transport layer initiliazation, which returns transport, a channel
// for received ConsMessage, a channel for received requests, and a channel for reply
func InitTransport(lg *zap.Logger, id uint32, port int, peers []Peer) (*HTTPTransport,
	chan *consmsgpb.WholeMessage, chan []byte, chan []byte) {
	msgc := make(chan *consmsgpb.WholeMessage, streamBufSize)
	tp := &HTTPTransport{id: id, Peers: make(map[uint32]*Peer), msgc: msgc}
	tp.proposal = make([][]byte, len(peers))
	for i, p := range peers {
		if index := uint32(i); index != id {
			tp.Peers[index] = &Peer{PeerID: index, Addr: p.Addr[7:], PublicKey: p.PublicKey}
		} else {
			tp.PrivateKey = *p.PublicKey
			ip := strings.Split(p.Addr, ":")
			port, _ := strconv.ParseUint(ip[2], 10, 64)
			node, _ := noise.NewNode(noise.WithNodeBindHost(net.ParseIP("127.0.0.1")),
				noise.WithNodeBindPort(uint16(port)), noise.WithNodeMaxRecvMessageSize(32*1024*1024))
			tp.node = node
		}
	}
	tp.node.RegisterMessage(NoiseMessage{}, UnmarshalNoiseMessage)
	tp.node.Handle(tp.Handler)
	err := tp.node.Listen()
	if err != nil {
		panic(err)
	}
	log.Printf("listening on %v\n", tp.node.Addr())
	reqc := make(chan []byte, streamBufSize)
	repc := make(chan []byte, streamBufSize)
	rprocessor := &ClientMsgProcessor{lg: lg, id: id, reqc: reqc, repc: repc}
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	mux.Handle(clientPrefix, rprocessor)
	mux.Handle(clientPrefix+"/", rprocessor)
	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux}
	server.SetKeepAlivesEnabled(true)

	go server.ListenAndServe()

	return tp, msgc, reqc, repc
}

func (tp *HTTPTransport) SendMessage(id uint32, msg *consmsgpb.WholeMessage) {
	m := NoiseMessage{Msg: msg}
	err := tp.node.SendMessage(context.TODO(), tp.Peers[id].Addr, m)
	for {
		if err == nil {
			return
		} else {
			fmt.Println("err", err.Error())
		}
		err = tp.node.SendMessage(context.TODO(), tp.Peers[id].Addr, m)
	}
}

func (tp *HTTPTransport) Handler(ctx noise.HandlerContext) error {
	obj, err := ctx.DecodeMessage()
	if err != nil {
		log.Fatal(err)
	}
	msg, ok := obj.(NoiseMessage)
	if !ok {
		log.Fatal(err)
	}
	go tp.OnReceiveMessage(msg.Msg)
	return nil
}

func (tp *HTTPTransport) OnReceiveMessage(msg *consmsgpb.WholeMessage) {
	if msg.From == tp.id {
		tp.msgc <- msg
		return
	}
	if msg.ConsMsg.Type == consmsgpb.MessageType_VAL || msg.ConsMsg.Type == consmsgpb.MessageType_ECHO ||
		msg.ConsMsg.Type == consmsgpb.MessageType_BVAL || msg.ConsMsg.Type == consmsgpb.MessageType_AUX {
		if Verify(msg, tp.Peers[msg.From].PublicKey) {
			tp.msgc <- msg
		}
		return
	}
	tp.msgc <- msg
}

func Verify(msg *consmsgpb.WholeMessage, pub *ecdsa.PrivateKey) bool {
	content, err := proto.Marshal(msg.ConsMsg)
	if err != nil {
		panic(err)
	}
	hash, err := sha256.ComputeHash(content)
	if err != nil {
		panic("sha256 computeHash failed")
	}
	b, err := myecdsa.VerifyECDSA(&pub.PublicKey, msg.Signature, hash)
	if err != nil {
		fmt.Println("Failed to verify a consmsgpb: ", err)
	}
	return b
}

func Sign(msg *consmsgpb.WholeMessage, priv *ecdsa.PrivateKey) {
	content, err := proto.Marshal(msg.ConsMsg)
	if err != nil {
		panic(err)
	}
	hash, err := sha256.ComputeHash(content)
	if err != nil {
		panic("sha256 computeHash failed")
	}
	sig, err := myecdsa.SignECDSA(priv, hash)
	if err != nil {
		panic("myecdsa signECDSA failed")
	}
	msg.Signature = sig
}

// ClientMsgProcessor is responsible for listening and processing requests from clients
type ClientMsgProcessor struct {
	num  int
	lg   *zap.Logger
	id   uint32
	reqc chan []byte
	repc chan []byte
}

func (cmsgProcessor *ClientMsgProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
