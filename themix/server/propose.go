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

package server

import (
	"sync"

	"go.themix.io/client/proto/clientpb"
	"go.themix.io/transport"
	"go.themix.io/transport/proto/consmsgpb"
	"go.uber.org/zap"
)

const channelSize = 4096 * 32

// Proposer is responsible for proposing requests
type Proposer struct {
	lg         *zap.Logger
	reqc       chan []byte
	verifyReq  chan *clientpb.ClientMessage
	verifyResp chan int
	tp         transport.Transport
	seq        uint64
	id         uint32
	lock       sync.Mutex
}

func initProposer(lg *zap.Logger, tp transport.Transport, id uint32, reqc chan []byte, pkPath string) *Proposer {
	proposer := &Proposer{lg: lg, tp: tp, id: id, reqc: reqc, lock: sync.Mutex{}}
	proposer.verifyReq = make(chan *clientpb.ClientMessage, channelSize)
	proposer.verifyResp = make(chan int, channelSize)
	go proposer.run()
	return proposer
}

func (proposer *Proposer) proceed(seq uint64) {
	// proposer.lock.Lock()
	// defer proposer.lock.Unlock()

	if proposer.seq <= seq {
		proposer.reqc <- []byte{} // insert an empty reqeust
	}
}

func (proposer *Proposer) run() {
	var req []byte
	for {
		req = <-proposer.reqc
		proposer.propose(req)
	}
}

// Propose broadcast a propose consmsgpb with the given request and the current sequence number
func (proposer *Proposer) propose(request []byte) {
	msg := &consmsgpb.WholeMessage{
		ConsMsg: &consmsgpb.ConsMessage{
			Type:     consmsgpb.MessageType_VAL,
			Proposer: proposer.id,
			Sequence: proposer.seq,
			Content:  request,
		},
		From: proposer.id,
	}
	if len(request) > 0 {
		proposer.lg.Info("propose",
			zap.Int("proposer", int(msg.ConsMsg.Proposer)),
			zap.Int("seq", int(msg.ConsMsg.Sequence)),
			zap.Int("content", int(msg.ConsMsg.Content[0])))
		// zap.Int("signature", int(msg.Signature[0])))
	}

	proposer.seq++
	proposer.tp.Broadcast(msg)
}
