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
	"net"
	"runtime"

	"go.themix.io/crypto/bls"
	"go.themix.io/transport"
	"go.themix.io/transport/info"
	"go.uber.org/zap"
)

// Node is a local process
type Node struct {
	reply     chan []byte
	proposer  *Proposer
	transport *transport.Transport
}

// InitNode initiate a node for processing messages
func InitNode(lg *zap.Logger, blsSig *bls.BlsSig, pkPath string, id info.IDType, n uint64, port int, addresses []string, batchsize int, coordinator net.Conn) {

	tp, msgc, reqc, repc := transport.InitTransport(lg, id, port, addresses)

	proposer := initProposer(lg, tp, id, reqc, coordinator)

	state := initState(lg, tp, blsSig, pkPath, id, proposer, n, repc, batchsize, coordinator)

	for i := 0; i < runtime.NumCPU()-1; i++ {
		go func() {
			for {
				msg := <-msgc
				state.insertMsg(msg)
			}
		}()
	}
	for {
		msg := <-msgc
		state.insertMsg(msg)
	}
}
