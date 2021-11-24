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
	"crypto/ecdsa"
	"fmt"
	"math"
	"sync"
	"time"

	"go.themix.io/crypto/bls"
	myecdsa "go.themix.io/crypto/ecdsa"
	"go.themix.io/crypto/sha256"
	"go.themix.io/transport"
	"go.themix.io/transport/proto/consmsgpb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// the maximum expected round that terminates consensus, P = 1 - pow(0.5, maxround)
var maxround = 30

type instance struct {
	tp            transport.Transport
	blsSig        *bls.BlsSig
	fastRBC       bool
	hasEcho       bool
	hasVotedZero  bool
	hasVotedOne   bool
	hasSentAux    bool
	hasSentCoin   bool
	zeroEndorsed  bool
	oneEndorsed   bool
	fastAuxZero   bool
	fastAuxOne    bool
	isDecided     bool
	isFinished    bool
	sequence      uint64
	n             uint64
	thld          uint64
	f             uint64
	fastgroup     uint64
	round         uint32
	numEcho       uint64
	numReady      uint64
	numOneSkip    uint64
	numZeroSkip   uint64
	binVals       uint8
	lastCoin      uint8
	hash          []byte
	sin           []bool
	canSkipCoin   []bool
	numBvalZero   []uint64
	numBvalOne    []uint64
	numAuxZero    []uint64
	numAuxOne     []uint64
	numCon        []uint64
	numCoin       []uint64
	echoSigns     *consmsgpb.Collections
	readySigns    *consmsgpb.Collections
	proposal      *consmsgpb.WholeMessage
	valMsgs       []*consmsgpb.WholeMessage
	bvalZeroSigns []*consmsgpb.Collections
	bvalOneSigns  []*consmsgpb.Collections
	auxZeroSigns  []*consmsgpb.Collections
	auxOneSigns   []*consmsgpb.Collections
	coinMsgs      [][]*consmsgpb.WholeMessage
	startR        bool
	startB        bool
	startS        bool
	tmrR          time.Timer
	tmrB          time.Timer
	tmrS          time.Timer
	priv          *ecdsa.PrivateKey
	lg            *zap.Logger
	lock          sync.Mutex
}

func initInstance(lg *zap.Logger, tp transport.Transport, blsSig *bls.BlsSig, pkPath string, sequence uint64, n uint64, thld uint64) *instance {
	inst := &instance{
		lg:            lg,
		tp:            tp,
		blsSig:        blsSig,
		sequence:      sequence,
		n:             n,
		thld:          thld,
		f:             n / 2,
		echoSigns:     &consmsgpb.Collections{Collections: make([][]byte, n)},
		sin:           make([]bool, maxround),
		canSkipCoin:   make([]bool, maxround),
		valMsgs:       make([]*consmsgpb.WholeMessage, n),
		bvalZeroSigns: make([]*consmsgpb.Collections, maxround),
		bvalOneSigns:  make([]*consmsgpb.Collections, maxround),
		auxZeroSigns:  make([]*consmsgpb.Collections, maxround),
		auxOneSigns:   make([]*consmsgpb.Collections, maxround),
		coinMsgs:      make([][]*consmsgpb.WholeMessage, maxround),
		numBvalZero:   make([]uint64, maxround),
		numBvalOne:    make([]uint64, maxround),
		numAuxZero:    make([]uint64, maxround),
		numAuxOne:     make([]uint64, maxround),
		numCon:        make([]uint64, maxround),
		numCoin:       make([]uint64, maxround),
		lock:          sync.Mutex{}}
	inst.fastgroup = uint64(math.Ceil(3*float64(inst.f)/2)) + 1
	inst.priv, _ = myecdsa.LoadKey(pkPath)
	for i := 0; i < maxround; i++ {
		inst.bvalZeroSigns[i] = &consmsgpb.Collections{Collections: make([][]byte, n)}
		inst.bvalOneSigns[i] = &consmsgpb.Collections{Collections: make([][]byte, n)}
		inst.auxZeroSigns[i] = &consmsgpb.Collections{Collections: make([][]byte, n)}
		inst.auxOneSigns[i] = &consmsgpb.Collections{Collections: make([][]byte, n)}
		inst.coinMsgs[i] = make([]*consmsgpb.WholeMessage, n)
		inst.canSkipCoin[i] = true
	}
	return inst
}

// return true if the instance is decided or finished at the first time
func (inst *instance) insertMsg(msg *consmsgpb.WholeMessage) (bool, bool) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Just for test
	if msg.ConsMsg.Round > 0 {
		return false, false
	}

	// if len(msg.ConsMsg.Content) > 0 {
	// 	inst.lg.Info("receive msg",
	// 		zap.String("type", consmsgpb.MessageType_name[int32(msg.ConsMsg.Type)]),
	// 		zap.Int("proposer", int(msg.ConsMsg.Proposer)),
	// 		zap.Int("seq", int(msg.ConsMsg.Sequence)),
	// 		zap.Int("round", int(msg.ConsMsg.Round)),
	// 		zap.Int("from", int(msg.From)))
	// } else {
	// 	inst.lg.Info("receive msg",
	// 		zap.String("type", consmsgpb.MessageType_name[int32(msg.ConsMsg.Type)]),
	// 		zap.Int("proposer", int(msg.ConsMsg.Proposer)),
	// 		zap.Int("seq", int(msg.ConsMsg.Sequence)),
	// 		zap.Int("round", int(msg.ConsMsg.Round)),
	// 		zap.Int("from", int(msg.From)))
	// }

	if inst.isFinished {
		return false, false
	}

	switch msg.ConsMsg.Type {

	/*
	 * upon receiving VAL(v)src
	 * if have not sent any ECHO(*) then
	 * broadcast VAL(v)src, ECHO(v)i
	 * start timer tmrR <- 2*delta
	 */
	case consmsgpb.MessageType_VAL:
		if inst.proposal != nil {
			verify := VerifySign(msg, inst.priv)
			if !verify {
				return false, false
			}
		}
		inst.proposal = msg
		hash, _ := sha256.ComputeHash(msg.ConsMsg.Content)
		inst.hash = hash
		if !inst.hasEcho {
			// broadcast VAL(v)src, ECHO(v)i
			inst.hasEcho = true
			m := &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_ECHO,
					Proposer: msg.ConsMsg.Proposer,
					Sequence: msg.ConsMsg.Sequence,
					Content:  hash,
				}}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
			m = &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_VAL_SIGN,
					Proposer: msg.ConsMsg.Proposer,
					Sequence: msg.ConsMsg.Sequence,
				},
				Signature: msg.Signature,
			}
			inst.tp.Broadcast(m)
		}
		return inst.isFastDecided()

	case consmsgpb.MessageType_VAL_SIGN:
		inst.valMsgs[msg.From] = msg

	/*
	 * upon receiving f+1 ECHO(v), and tmrR expires
	 * if have not received any VAL(v')src (v' != v) then
	 * broadcast READY(v)i
	 */
	case consmsgpb.MessageType_ECHO:
		verify := VerifySign(msg, inst.priv)
		if !verify {
			return false, false
		}
		if inst.echoSigns.Collections[msg.From] != nil {
			return false, false
		}
		inst.numEcho++
		inst.echoSigns.Collections[msg.From] = msg.Signature
		/*
		 * upon receiving ECHO(v) from fast group
		 * broadcast ECHO(v) sent by fast group
		 * deliver(v)
		 */
		if inst.numEcho >= inst.fastgroup && !inst.fastRBC && inst.round == 0 {
			inst.fastRBC = true
			collection := serialCollection(inst.echoSigns)
			m := &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_ECHO_COLLECTION,
					Proposer: msg.ConsMsg.Proposer,
					Round:    msg.ConsMsg.Round,
					Sequence: msg.ConsMsg.Sequence,
				},
				Collection: collection,
			}
			inst.tp.Broadcast(m)
			if !inst.hasVotedZero && !inst.hasVotedOne {
				inst.hasVotedOne = true
				m := &consmsgpb.WholeMessage{
					ConsMsg: &consmsgpb.ConsMessage{
						Type:     consmsgpb.MessageType_BVAL,
						Proposer: msg.ConsMsg.Proposer,
						Sequence: msg.ConsMsg.Sequence,
						Content:  []byte{1}, // vote 1
					},
				}
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
		}
		return inst.isFastDecided()
	case consmsgpb.MessageType_BVAL:
		verify := VerifySign(msg, inst.priv)
		if !verify {
			return false, false
		}
		var b bool
		switch msg.ConsMsg.Content[0] {
		case 0:
			inst.numBvalZero[msg.ConsMsg.Round]++
			inst.bvalZeroSigns[msg.ConsMsg.Round].Collections[msg.From] = msg.Signature
		case 1:
			inst.numBvalOne[msg.ConsMsg.Round]++
			inst.bvalOneSigns[msg.ConsMsg.Round].Collections[msg.From] = msg.Signature
		}
		if inst.round == msg.ConsMsg.Round && !inst.hasVotedZero && inst.numBvalZero[inst.round] > inst.f {
			inst.hasVotedZero = true
			collection := serialCollection(inst.bvalZeroSigns[msg.ConsMsg.Round])
			m := &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_BVAL_ZERO_COLLECTION,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
				},
				Collection: collection,
			}
			inst.tp.Broadcast(m)
		}
		if inst.round == msg.ConsMsg.Round && !inst.zeroEndorsed && inst.numBvalZero[inst.round] >= inst.thld {
			inst.zeroEndorsed = true
			if !inst.hasSentAux {
				inst.hasSentAux = true
				m := &consmsgpb.WholeMessage{
					ConsMsg: &consmsgpb.ConsMessage{
						Type:     consmsgpb.MessageType_AUX,
						Proposer: msg.ConsMsg.Proposer,
						Round:    inst.round,
						Sequence: msg.ConsMsg.Sequence,
						Content:  []byte{0}, // aux 0
					},
				}
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
			b = true
		}
		if inst.round == msg.ConsMsg.Round && !inst.hasVotedOne && inst.numBvalOne[inst.round] > inst.f {
			inst.hasVotedOne = true
			collection := serialCollection(inst.bvalOneSigns[msg.ConsMsg.Round])
			m := &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_BVAL_ONE_COLLECTION,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
				},
				Collection: collection,
			}
			inst.tp.Broadcast(m)
		}
		if inst.round == msg.ConsMsg.Round && !inst.oneEndorsed && inst.numBvalOne[inst.round] >= inst.thld {
			inst.oneEndorsed = true
			if !inst.hasSentAux {
				inst.hasSentAux = true
				m := &consmsgpb.WholeMessage{
					ConsMsg: &consmsgpb.ConsMessage{
						Type:     consmsgpb.MessageType_AUX,
						Proposer: msg.ConsMsg.Proposer,
						Round:    inst.round,
						Sequence: msg.ConsMsg.Sequence,
						Content:  []byte{1}, // aux 1
					},
				}
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
		}
		if b {
			return inst.isFastDecided()
		}
	case consmsgpb.MessageType_AUX:
		verify := VerifySign(msg, inst.priv)
		if !verify {
			return false, false
		}
		switch msg.ConsMsg.Content[0] {
		case 0:
			inst.numAuxZero[msg.ConsMsg.Round]++
			inst.auxZeroSigns[msg.ConsMsg.Round].Collections[msg.From] = msg.Signature
		case 1:
			inst.numAuxOne[msg.ConsMsg.Round]++
			inst.auxOneSigns[msg.ConsMsg.Round].Collections[msg.From] = msg.Signature
		}
		if inst.round == msg.ConsMsg.Round && msg.ConsMsg.Content[0] == 0 && !inst.fastAuxZero && inst.numAuxZero[msg.ConsMsg.Round] >= inst.fastgroup {
			inst.fastAuxZero = true
			collection := serialCollection(inst.auxZeroSigns[msg.ConsMsg.Round])
			inst.tp.Broadcast(&consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_AUX_ZERO_COLLECTION,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
				},
				Collection: collection,
			})
			inst.zeroEndorsed = true
			if inst.canSkipCoin[inst.round] {
				inst.tp.Broadcast(&consmsgpb.WholeMessage{
					ConsMsg: &consmsgpb.ConsMessage{
						Type:     consmsgpb.MessageType_SKIP,
						Proposer: msg.ConsMsg.Proposer,
						Round:    inst.round,
						Sequence: msg.ConsMsg.Sequence,
						Content:  msg.ConsMsg.Content,
					},
				})
			}
			return inst.isFastDecided()
		}
		if inst.round == msg.ConsMsg.Round && msg.ConsMsg.Content[0] == 1 && !inst.fastAuxOne && inst.numAuxOne[msg.ConsMsg.Round] >= inst.fastgroup {
			inst.fastAuxOne = true
			collection := serialCollection(inst.auxOneSigns[msg.ConsMsg.Round])
			inst.tp.Broadcast(&consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_AUX_ONE_COLLECTION,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
				},
				Collection: collection,
			})
			inst.oneEndorsed = true
			if inst.canSkipCoin[inst.round] {
				inst.tp.Broadcast(&consmsgpb.WholeMessage{
					ConsMsg: &consmsgpb.ConsMessage{
						Type:     consmsgpb.MessageType_SKIP,
						Proposer: msg.ConsMsg.Proposer,
						Round:    inst.round,
						Sequence: msg.ConsMsg.Sequence,
						Content:  msg.ConsMsg.Content,
					},
				})
			}
			return inst.isFastDecided()
		}
		// return inst.isFastDecided()
	case consmsgpb.MessageType_SKIP:
		switch msg.ConsMsg.Content[0] {
		case 0:
			inst.numZeroSkip++
		case 1:
			inst.numOneSkip++
		}
		if msg.ConsMsg.Content[0] == 0 && inst.proposal != nil && !inst.isDecided && inst.numZeroSkip >= inst.fastgroup {
			return inst.isFastDecided()
		}
		if msg.ConsMsg.Content[0] == 1 && inst.proposal != nil && !inst.isDecided && inst.numOneSkip >= inst.fastgroup {
			return inst.isFastDecided()
		}
		return inst.isFastDecided()
	case consmsgpb.MessageType_ECHO_COLLECTION:
		if inst.fastRBC || inst.hasVotedZero || inst.hasVotedOne {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection.Collections {
			if inst.echoSigns.Collections[i] != nil || sign == nil || !VerifyCollection(sign, inst.hash, inst.priv) {
				continue
			}
			inst.numEcho++
			inst.echoSigns.Collections[i] = sign
		}
		inst.fastRBC = true
		// inst.tp.Broadcast(msg)
		inst.hasVotedOne = true
		m := &consmsgpb.WholeMessage{
			ConsMsg: &consmsgpb.ConsMessage{
				Type:     consmsgpb.MessageType_BVAL,
				Proposer: msg.ConsMsg.Proposer,
				Sequence: msg.ConsMsg.Sequence,
				Content:  []byte{1}, // vote 1
			},
		}
		GetSign(m, inst.priv)
		inst.tp.Broadcast(m)
		inst.isFastDecided()
	case consmsgpb.MessageType_BVAL_ZERO_COLLECTION:
		if inst.zeroEndorsed || inst.oneEndorsed || inst.hasSentAux {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection.Collections {
			if inst.bvalZeroSigns[msg.ConsMsg.Round].Collections[i] != nil || sign == nil || !VerifyCollection(sign, []byte{0}, inst.priv) {
				continue
			}
			inst.numBvalZero[msg.ConsMsg.Round]++
			inst.bvalZeroSigns[msg.ConsMsg.Round].Collections[i] = sign
		}
		inst.zeroEndorsed = true
		// inst.tp.Broadcast(msg)
		if !inst.hasSentAux {
			inst.hasSentAux = true
			m := &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_AUX,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
					Content:  []byte{0}, // aux 0
				},
			}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
		}
		inst.isFastDecided()
	case consmsgpb.MessageType_BVAL_ONE_COLLECTION:
		if inst.oneEndorsed || inst.zeroEndorsed || inst.hasSentAux {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection.Collections {
			if inst.bvalOneSigns[msg.ConsMsg.Round].Collections[i] != nil || sign == nil || !VerifyCollection(sign, []byte{1}, inst.priv) {
				continue
			}
			inst.numBvalOne[msg.ConsMsg.Round]++
			inst.bvalOneSigns[msg.ConsMsg.Round].Collections[i] = sign
		}
		inst.oneEndorsed = true
		// inst.tp.Broadcast(msg)
		if !inst.hasSentAux {
			inst.hasSentAux = true
			m := &consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_AUX,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
					Content:  []byte{1}, // aux 1
				},
			}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
		}
		inst.isFastDecided()
	case consmsgpb.MessageType_AUX_ZERO_COLLECTION:
		if inst.fastAuxZero || inst.fastAuxOne || inst.round != msg.ConsMsg.Round {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection.Collections {
			if inst.auxZeroSigns[msg.ConsMsg.Round].Collections[i] != nil || sign == nil || !VerifyCollection(sign, []byte{0}, inst.priv) {
				continue
			}
			inst.numBvalZero[msg.ConsMsg.Round]++
			inst.bvalZeroSigns[msg.ConsMsg.Round].Collections[i] = sign
		}
		inst.fastAuxZero = true
		inst.zeroEndorsed = true
		if inst.canSkipCoin[inst.round] {
			inst.tp.Broadcast(&consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_SKIP,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
					Content:  []byte{0},
				},
			})
		}
		return inst.isFastDecided()
	case consmsgpb.MessageType_AUX_ONE_COLLECTION:
		if inst.fastAuxZero || inst.fastAuxOne || inst.round != msg.ConsMsg.Round {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection.Collections {
			if inst.auxOneSigns[msg.ConsMsg.Round].Collections[i] != nil || sign == nil || !VerifyCollection(sign, []byte{1}, inst.priv) {
				continue
			}
			inst.numBvalOne[msg.ConsMsg.Round]++
			inst.bvalOneSigns[msg.ConsMsg.Round].Collections[i] = sign
		}
		inst.fastAuxOne = true
		inst.oneEndorsed = true
		if inst.canSkipCoin[inst.round] {
			inst.tp.Broadcast(&consmsgpb.WholeMessage{
				ConsMsg: &consmsgpb.ConsMessage{
					Type:     consmsgpb.MessageType_SKIP,
					Proposer: msg.ConsMsg.Proposer,
					Round:    inst.round,
					Sequence: msg.ConsMsg.Sequence,
					Content:  []byte{1},
				},
			})
		}
	default:
		return false, false
	}
	return false, false
}

func (inst *instance) isFastDecided() (bool, bool) {
	if inst.isDecided {
		inst.isFinished = true
		return false, true
	}
	if inst.proposal != nil &&
		(inst.numZeroSkip >= inst.fastgroup || inst.numOneSkip >= inst.fastgroup) {
		inst.binVals = 1
		inst.isDecided = true
		return true, false
	} else {
		return false, false
	}
}

func (inst *instance) decidedOne() bool {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	return inst.isDecided && inst.binVals == 1
}

func (inst *instance) getProposal() *consmsgpb.WholeMessage {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	return inst.proposal
}

func serialCollection(collection *consmsgpb.Collections) []byte {
	mar_collection, err := proto.Marshal(collection)
	if err != nil {
		panic("Marshal collection failed")
	}
	return mar_collection
}

func deserialCollection(data []byte) consmsgpb.Collections {
	collection := consmsgpb.Collections{}
	err := proto.Unmarshal(data, &collection)
	if err != nil {
		panic("Unmarshal collection failed")
	}
	return collection
}

func GetSign(msg *consmsgpb.WholeMessage, priv *ecdsa.PrivateKey) {
	content, _ := proto.Marshal(msg.ConsMsg)
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

func VerifySign(msg *consmsgpb.WholeMessage, priv *ecdsa.PrivateKey) bool {
	content, _ := proto.Marshal(msg.ConsMsg)
	hash, err := sha256.ComputeHash(content)
	if err != nil {
		panic("sha256 computeHash failed")
	}
	b, err := myecdsa.VerifyECDSA(&priv.PublicKey, msg.Signature, hash)
	if err != nil {
		fmt.Println("Failed to verify a consmsgpb: ", err)
	}
	return b
}

func VerifyCollection(sign, hash []byte, priv *ecdsa.PrivateKey) bool {
	if hash == nil {
		return false
	}
	b, _ := myecdsa.VerifyECDSA(&priv.PublicKey, sign, hash)
	// if err != nil {
	// 	fmt.Println("Failed to verify a consmsgpb: ", err)
	// }
	return b
}
