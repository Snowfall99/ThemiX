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

type asyncCommSubset struct {
	st            *state
	lg            *zap.Logger
	n             uint64
	thld          uint64
	sequence      uint64
	numDecided    uint64
	numFinished   uint64
	numDecidedOne uint64
	instances     []*instance
	proposer      *Proposer
	reqc          chan *message.ConsMessage
	lock          sync.Mutex
}

func initACS(st *state,
	lg *zap.Logger,
	tp transport.Transport,
	blsSig *bls.BlsSig,
	pkPath string,
	proposer *Proposer,
	seq uint64, n uint64,
	reqc chan *message.ConsMessage) *asyncCommSubset {
	re := &asyncCommSubset{
		st:        st,
		lg:        lg,
		proposer:  proposer,
		n:         n,
		sequence:  seq,
		instances: make([]*instance, n),
		reqc:      reqc,
		lock:      sync.Mutex{},
		//coordinator: coordinator
	}
	re.thld = n/2 + 1
	for i := info.IDType(0); i < info.IDType(n); i++ {
		re.instances[i] = initInstance(lg, tp, blsSig, pkPath, seq, n, re.thld)
	}
	return re
}

func (acs *asyncCommSubset) insertMsg(msg *message.ConsMessage) {
	isDecided, isFinished := acs.instances[msg.Proposer].insertMsg(msg)
	if isDecided {
		acs.lock.Lock()
		defer acs.lock.Unlock()

		proposal := acs.instances[msg.Proposer].getProposal()
		acs.reqc <- proposal

		// if !acs.instances[msg.Proposer].decidedOne() && msg.Proposer == acs.proposer.id {
		// 	fmt.Printf("ID %d decided zero at %d\n", msg.Proposer, msg.Sequence)
		// }

		// acs.numDecided++
		// if acs.numDecided == 1 {
		// 	acs.proposer.proceed(acs.sequence)
		// }

		// if acs.instances[msg.Proposer].decidedOne() {
		// 	acs.numDecidedOne++
		// }

		// Just for test
		// if acs.numDecidedOne == acs.thld {
		// 	for i, inst := range acs.instances {
		// 		inst.canVoteZero(info.IDType(i), acs.sequence)
		// 	}
		// }

		// if acs.numDecided == acs.n {
		// 	for _, inst := range acs.instances {
		// 		proposal := inst.getProposal()
		// 		if inst.decidedOne() && len(proposal.Content) != 0 {
		// 			inst.lg.Info("executed",
		// 				zap.Int("proposer", int(proposal.Proposer)),
		// 				zap.Int("seq", int(msg.Sequence)),
		// 				zap.Int("content", int(proposal.Content[0])))
		// 			// zap.Int("content", int(binary.LittleEndian.Uint32(proposal.Content))))
		// 			acs.reqc <- proposal

		// 			// // send execute message to coordinator
		// 			// data := "end " + strconv.FormatUint(msg.Sequence, 10)
		// 			// _, err := acs.coordinator.Write([]byte(data))
		// 			// if err != nil {
		// 			// 	fmt.Println("send execute message to coordinator failed: ", err.Error())
		// 			// 	return
		// 			// }
		// 		} else if proposal.Proposer == acs.proposer.id && len(proposal.Content) != 0 {
		// 			inst.lg.Info("repropose",
		// 				zap.Int("proposer", int(proposal.Proposer)),
		// 				zap.Int("seq", int(proposal.Sequence)),
		// 				zap.Int("content", int(proposal.Content[0])))
		// 			// zap.Int("content", int(binary.LittleEndian.Uint32(proposal.Content))))
		// 			acs.proposer.propose(proposal.Content)
		// 		} else if inst.decidedOne() {
		// 			inst.lg.Info("empty",
		// 				zap.Int("proposer", int(proposal.Proposer)),
		// 				zap.Int("seq", int(proposal.Sequence)))
		// 		}
		// 	}
		// }
	} else if isFinished {
		// acs.lock.Lock()
		// defer acs.lock.Unlock()

		// acs.numFinished++
		// if acs.numFinished == acs.n {
		// 	acs.st.garbageCollect(acs.sequence)
		// }
	}
}
