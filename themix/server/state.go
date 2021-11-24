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

	"go.themix.io/crypto/bls"
	"go.themix.io/transport"
	"go.themix.io/transport/info"
	"go.themix.io/transport/message"
	"go.uber.org/zap"
)

type state struct {
	lg        *zap.Logger
	tp        transport.Transport
	blsSig    *bls.BlsSig
	pkPath    string
	proposer  *Proposer
	id        info.IDType
	n         uint64
	collected uint64
	execs     map[uint64]*asyncCommSubset
	lock      sync.RWMutex
	reqc      chan *message.WholeMessage
	repc      chan []byte
}

func initState(lg *zap.Logger,
	tp transport.Transport,
	blsSig *bls.BlsSig,
	pkPath string,
	id info.IDType,
	proposer *Proposer,
	n uint64, repc chan []byte,
	batchsize int) *state {
	st := &state{
		lg:        lg,
		tp:        tp,
		blsSig:    blsSig,
		pkPath:    pkPath,
		id:        id,
		proposer:  proposer,
		n:         n,
		collected: 0,
		execs:     make(map[uint64]*asyncCommSubset),
		lock:      sync.RWMutex{},
		reqc:      make(chan *message.WholeMessage, 2*int(n)*batchsize),
		repc:      repc,
	}
	go st.run()
	return st
}

func (st *state) insertMsg(msg *message.WholeMessage) {
	st.lock.RLock()

	if exec, ok := st.execs[msg.Sequence]; ok {
		st.lock.RUnlock()
		exec.insertMsg(msg)
	} else {
		if st.collected <= msg.Sequence {
			st.lock.RUnlock()

			exec := initACS(st, st.lg, st.tp, st.blsSig, st.pkPath, st.proposer, msg.Sequence, st.n, st.reqc)

			st.lock.Lock()
			if e, ok := st.execs[msg.Sequence]; ok {
				exec = e
			} else {
				st.execs[msg.Sequence] = exec
			}
			st.lock.Unlock()

			exec.insertMsg(msg)
		}
	}
}

func (st *state) garbageCollect(seq uint64) {

	st.lock.Lock()
	defer st.lock.Unlock()

	delete(st.execs, seq)
	for _, b := st.execs[st.collected]; !b; st.collected++ {
	}
}

// execute requests by a single thread
func (st *state) run() {
	for {
		req := <-st.reqc
		if req.Proposer == st.id {
			st.repc <- []byte{}
		}
	}
}
