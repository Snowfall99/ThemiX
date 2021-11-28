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

package transport

import (
	"go.themix.io/transport/http"
	"go.themix.io/transport/proto/consmsgpb"
	"go.uber.org/zap"
)

type Transport interface {
	Broadcast(msg *consmsgpb.WholeMessage)
}

// InitTransport executes transport layer initiliazation, which returns transport, a channel
// for received ConsMessage, a channel for received requests, and a channel for reply
func InitTransport(lg *zap.Logger, id uint32, port int, peers []http.Peer) (Transport,
	chan *consmsgpb.WholeMessage, chan *consmsgpb.WholeMessage, chan []byte, chan []byte) {
	return http.InitTransport(lg, id, port, peers)
}
