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
	"fmt"
	"net"
	"strconv"
	"sync"

	"go.themix.io/transport"
	"go.themix.io/transport/info"
	"go.themix.io/transport/message"
	"go.uber.org/zap"
)

// Proposer is responsible for proposing requests
type Proposer struct {
	lg          *zap.Logger
	reqc        chan []byte
	tp          transport.Transport
	seq         uint64
	id          info.IDType
	lock        sync.Mutex
	coordinator net.Conn
}

func initProposer(lg *zap.Logger, tp transport.Transport, id info.IDType, reqc chan []byte, coordinator net.Conn) *Proposer {
	proposer := &Proposer{lg: lg, tp: tp, id: id, reqc: reqc, lock: sync.Mutex{}, coordinator: coordinator}
	go proposer.run()
	return proposer
}

func (proposer *Proposer) proceed(seq uint64) {
	proposer.lock.Lock()
	defer proposer.lock.Unlock()

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

// Propose broadcast a propose message with the given request and the current sequence number
func (proposer *Proposer) propose(request []byte) {
	proposer.lock.Lock()

	msg := &message.ConsMessage{
		Type:     message.VAL,
		Proposer: proposer.id,
		From:     proposer.id,
		Sequence: proposer.seq,
		Content:  request}

	if len(request) > 0 {
		proposer.lg.Info("propose",
			zap.Int("proposer", int(msg.Proposer)),
			zap.Int("seq", int(msg.Sequence)),
			zap.Int("content", int(msg.Content[0])))
	}

	// send propose to coordinator
	data := "start " + strconv.FormatUint(msg.Sequence, 10)
	_, err := proposer.coordinator.Write([]byte(data))
	if err != nil {
		fmt.Println("send propose to coordinator failed: ", err.Error())
		return
	}

	proposer.seq++

	proposer.lock.Unlock()

	proposer.tp.Broadcast(msg)
}
