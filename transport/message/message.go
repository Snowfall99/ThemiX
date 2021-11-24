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

package message

import "go.themix.io/transport/info"

// MessageType is type of consensus message
type MessageType uint8

// VAL is proposal message,
// ECHO is sent upon receiving VAL,
// READY is sent upon receiving 2f+1 ECHO
// BVAL is sent upon voting for a binary value (0 or 1)
// AUX is sent upon receiving f+1 matching BVAL message and timeout,
//	or in optimal case where f+1 BVAL matching BVAL messages are received from SG
// COIN is coin message
const (
	VAL                  MessageType = 0
	ECHO                 MessageType = 1
	READY                MessageType = 2
	BVAL                 MessageType = 3
	AUX                  MessageType = 4
	COIN                 MessageType = 5
	CON                  MessageType = 6
	SKIP                 MessageType = 7
	ECHO_COLLECTION      MessageType = 8
	READY_COLLECTION     MessageType = 9
	BVAL_ZERO_COLLECTION MessageType = 10
	BVAL_ONE_COLLECTION  MessageType = 11
	AUX_ZERO_COLLECTION  MessageType = 12
	AUX_ONE_COLLECTION   MessageType = 13
	VAL_SIGN             MessageType = 14
)

// ConsMessage is the message type for achieving consensus
type ConsMessage struct {
	Type     MessageType
	Proposer info.IDType
	Round    uint8
	Sequence uint64
	Content  []byte
}

// WholeMessage is the message type to hold consmessage and signature
type WholeMessage struct {
	ConsMsg    *ConsMessage
	From       info.IDType
	Signature  []byte
	Collection []byte
}

// Request is the message sent from client to servers
type Request struct {
	From      info.IDType
	Sequence  uint64
	Signature []byte
	Content   []byte
}

// GetName return message type name
func (t MessageType) GetName() string {
	switch t {
	case 0:
		return "VAL"
	case 1:
		return "ECHO"
	case 2:
		return "READY"
	case 3:
		return "BVAL"
	case 4:
		return "AUX"
	case 5:
		return "COIN"
	case 6:
		return "CON"
	case 7:
		return "SKIP"
	case 8:
		return "ECHO_COLLECTION"
	case 9:
		return "READY_COLLECTION"
	case 10:
		return "BVAL_ZERO_COLLECTION"
	case 11:
		return "BVAL_ONE_COLLECTION"
	case 12:
		return "AUX_ZERO_COLLECTION"
	case 13:
		return "AUX_ONE_COLLECTION"
	case 14:
		return "VAL_SIGN"
	}
	return "UNKNOWN"
}
