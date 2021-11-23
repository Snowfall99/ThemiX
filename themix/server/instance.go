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
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"go.themix.io/crypto/bls"
	myecdsa "go.themix.io/crypto/ecdsa"
	"go.themix.io/crypto/sha256"
	"go.themix.io/transport"
	"go.themix.io/transport/message"
	"go.uber.org/zap"
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
	round         uint8
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
	echoSigns     [][]byte
	readySigns    [][]byte
	proposal      *message.ConsMessage
	valMsgs       []*message.ConsMessage
	bvalZeroSigns [][][]byte
	bvalOneSigns  [][][]byte
	auxZeroSigns  [][][]byte
	auxOneSigns   [][][]byte
	coinMsgs      [][]*message.ConsMessage
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
		sin:           make([]bool, maxround),
		canSkipCoin:   make([]bool, maxround),
		valMsgs:       make([]*message.ConsMessage, n),
		echoSigns:     make([][]byte, n),
		readySigns:    make([][]byte, n),
		bvalZeroSigns: make([][][]byte, maxround),
		bvalOneSigns:  make([][][]byte, maxround),
		auxZeroSigns:  make([][][]byte, maxround),
		auxOneSigns:   make([][][]byte, maxround),
		coinMsgs:      make([][]*message.ConsMessage, maxround),
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
		inst.bvalZeroSigns[i] = make([][]byte, n)
		inst.bvalOneSigns[i] = make([][]byte, n)
		inst.auxZeroSigns[i] = make([][]byte, n)
		inst.auxOneSigns[i] = make([][]byte, n)
		inst.coinMsgs[i] = make([]*message.ConsMessage, n)
		inst.canSkipCoin[i] = true
	}
	return inst
}

// return true if the instance is decided or finished at the first time
func (inst *instance) insertMsg(msg *message.ConsMessage) (bool, bool) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Just for test
	if msg.Round > 0 {
		return false, false
	}

	// if len(msg.Content) > 0 {
	// 	inst.lg.Info("receive msg",
	// 		zap.String("type", msg.Type.GetName()),
	// 		zap.Int("proposer", int(msg.Proposer)),
	// 		zap.Int("seq", int(msg.Sequence)),
	// 		zap.Int("round", int(msg.Round)),
	// 		zap.Int("from", int(msg.From)),
	// 		zap.Int("content", int(msg.Content[0])))
	// } else {
	// 	inst.lg.Info("receive msg",
	// 		zap.String("type", msg.Type.GetName()),
	// 		zap.Int("proposer", int(msg.Proposer)),
	// 		zap.Int("seq", int(msg.Sequence)),
	// 		zap.Int("round", int(msg.Round)),
	// 		zap.Int("from", int(msg.From)))
	// }

	if inst.isFinished {
		return false, false
	}

	switch msg.Type {

	/*
	 * upon receiving VAL(v)src
	 * if have not sent any ECHO(*) then
	 * broadcast VAL(v)src, ECHO(v)i
	 * start timer tmrR <- 2*delta
	 */
	case message.VAL:
		if inst.proposal != nil {
			verify := VerifySign(msg.Signature, inst.hash, inst.priv)
			if !verify {
				return false, false
			}
		}
		inst.proposal = msg
		hash, _ := sha256.ComputeHash(msg.Content)
		inst.hash = hash
		if !inst.hasEcho {
			// broadcast VAL(v)src, ECHO(v)i
			inst.hasEcho = true
			m := &message.ConsMessage{
				Type:     message.ECHO,
				Proposer: msg.Proposer,
				Sequence: msg.Sequence,
				Content:  hash}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
			m = &message.ConsMessage{
				Type:      message.VAL_SIGN,
				Proposer:  msg.Proposer,
				Sequence:  msg.Sequence,
				Signature: msg.Signature,
			}
			inst.tp.Broadcast(m)
		}
		return inst.isFastDecided()

	case message.VAL_SIGN:
		inst.valMsgs[msg.From] = msg

	/*
	 * upon receiving f+1 ECHO(v), and tmrR expires
	 * if have not received any VAL(v')src (v' != v) then
	 * broadcast READY(v)i
	 */
	case message.ECHO:
		verify := VerifySign(msg.Signature, inst.hash, inst.priv)
		if !verify {
			return false, false
		}
		if inst.echoSigns[msg.From] != nil {
			return false, false
		}
		inst.numEcho++
		inst.echoSigns[msg.From] = msg.Signature
		/*
		 * upon receiving ECHO(v) from fast group
		 * broadcast ECHO(v) sent by fast group
		 * deliver(v)
		 */
		if inst.numEcho >= inst.fastgroup && !inst.fastRBC && inst.round == 0 {
			inst.fastRBC = true
			collection := serialCollection(inst.echoSigns)
			m := &message.ConsMessage{
				Type:       message.ECHO_COLLECTION,
				Proposer:   msg.Proposer,
				Round:      msg.Round,
				Sequence:   msg.Sequence,
				Collection: collection,
			}
			inst.tp.Broadcast(m)
			if !inst.hasVotedZero && !inst.hasVotedOne {
				inst.hasVotedOne = true
				m := &message.ConsMessage{
					Type:     message.BVAL,
					Proposer: msg.Proposer,
					Sequence: msg.Sequence,
					Content:  []byte{1}, // vote 1
				}
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
		}
		return inst.isFastDecided()
	case message.BVAL:
		verify := VerifySign(msg.Signature, msg.Content, inst.priv)
		if !verify {
			return false, false
		}
		var b bool
		switch msg.Content[0] {
		case 0:
			inst.numBvalZero[msg.Round]++
			inst.bvalZeroSigns[msg.Round][msg.From] = msg.Signature
		case 1:
			inst.numBvalOne[msg.Round]++
			inst.bvalOneSigns[msg.Round][msg.From] = msg.Signature
		}
		if inst.round == msg.Round && !inst.hasVotedZero && inst.numBvalZero[inst.round] > inst.f {
			inst.hasVotedZero = true
			collection := serialCollection(inst.bvalZeroSigns[msg.Round])
			m := &message.ConsMessage{
				Type:       message.BVAL_ZERO_COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: collection,
			}
			inst.tp.Broadcast(m)
		}
		if inst.round == msg.Round && !inst.zeroEndorsed && inst.numBvalZero[inst.round] >= inst.thld {
			inst.zeroEndorsed = true
			if !inst.hasSentAux {
				inst.hasSentAux = true
				m := &message.ConsMessage{
					Type:     message.AUX,
					Proposer: msg.Proposer,
					Round:    inst.round,
					Sequence: msg.Sequence,
					Content:  []byte{0}, // aux 0
				}
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
			b = true
		}
		if inst.round == msg.Round && !inst.hasVotedOne && inst.numBvalOne[inst.round] > inst.f {
			inst.hasVotedOne = true
			collection := serialCollection(inst.bvalOneSigns[msg.Round])
			m := &message.ConsMessage{
				Type:       message.BVAL_ONE_COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: collection,
			}
			inst.tp.Broadcast(m)
		}
		if inst.round == msg.Round && !inst.oneEndorsed && inst.numBvalOne[inst.round] >= inst.thld {
			inst.oneEndorsed = true
			if !inst.hasSentAux {
				inst.hasSentAux = true
				m := &message.ConsMessage{
					Type:     message.AUX,
					Proposer: msg.Proposer,
					Round:    inst.round,
					Sequence: msg.Sequence,
					Content:  []byte{1}} // aux 1
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
		}
		if b {
			return inst.isFastDecided()
		}
	case message.AUX:
		verify := VerifySign(msg.Signature, msg.Content, inst.priv)
		if !verify {
			return false, false
		}
		switch msg.Content[0] {
		case 0:
			inst.numAuxZero[msg.Round]++
			inst.auxZeroSigns[msg.Round][msg.From] = msg.Signature
		case 1:
			inst.numAuxOne[msg.Round]++
			inst.auxOneSigns[msg.Round][msg.From] = msg.Signature
		}
		if inst.round == msg.Round && msg.Content[0] == 0 && !inst.fastAuxZero && inst.numAuxZero[msg.Round] >= inst.fastgroup {
			inst.fastAuxZero = true
			collection := serialCollection(inst.auxZeroSigns[msg.Round])
			inst.tp.Broadcast(&message.ConsMessage{
				Type:       message.AUX_ZERO_COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: collection,
			})

			inst.zeroEndorsed = true
			if inst.canSkipCoin[inst.round] {
				inst.tp.Broadcast(&message.ConsMessage{
					Type:     message.SKIP,
					Proposer: msg.Proposer,
					Round:    inst.round,
					Sequence: msg.Sequence,
					Content:  msg.Content,
				})
			}
			return inst.isFastDecided()
		}
		if inst.round == msg.Round && msg.Content[0] == 1 && !inst.fastAuxOne && inst.numAuxOne[msg.Round] >= inst.fastgroup {
			inst.fastAuxOne = true
			collection := serialCollection(inst.auxOneSigns[msg.Round])
			inst.tp.Broadcast(&message.ConsMessage{
				Type:       message.AUX_ONE_COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: collection,
			})

			inst.oneEndorsed = true
			if inst.canSkipCoin[inst.round] {
				inst.tp.Broadcast(&message.ConsMessage{
					Type:     message.SKIP,
					Proposer: msg.Proposer,
					Round:    inst.round,
					Sequence: msg.Sequence,
					Content:  msg.Content,
				})
			}
			return inst.isFastDecided()
		}
		// return inst.isFastDecided()
	case message.SKIP:
		switch msg.Content[0] {
		case 0:
			inst.numZeroSkip++
		case 1:
			inst.numOneSkip++
		}
		if msg.Content[0] == 0 && inst.proposal != nil && !inst.isDecided && inst.numZeroSkip >= inst.fastgroup {
			return inst.isFastDecided()
		}
		if msg.Content[0] == 1 && inst.proposal != nil && !inst.isDecided && inst.numOneSkip >= inst.fastgroup {
			return inst.isFastDecided()
		}
		// return inst.isFastDecided()
	case message.ECHO_COLLECTION:
		if inst.fastRBC || inst.hasVotedZero || inst.hasVotedOne {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection {
			if inst.echoSigns[i] != nil || sign == nil || !VerifySign(sign, inst.hash, inst.priv) {
				continue
			}
			inst.numEcho++
			inst.echoSigns[i] = sign
		}
		inst.fastRBC = true
		// inst.tp.Broadcast(msg)
		inst.hasVotedOne = true
		m := &message.ConsMessage{
			Type:     message.BVAL,
			Proposer: msg.Proposer,
			Sequence: msg.Sequence,
			Content:  []byte{1}, // vote 1
		}
		GetSign(m, inst.priv)
		inst.tp.Broadcast(m)
		inst.isFastDecided()
	case message.BVAL_ZERO_COLLECTION:
		if inst.zeroEndorsed || inst.oneEndorsed || inst.hasSentAux {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection {
			if inst.bvalZeroSigns[msg.Round][i] != nil || sign == nil || !VerifySign(sign, []byte{0}, inst.priv) {
				continue
			}
			inst.numBvalZero[msg.Round]++
			inst.bvalZeroSigns[msg.Round][i] = sign
		}
		inst.zeroEndorsed = true
		// inst.tp.Broadcast(msg)
		if !inst.hasSentAux {
			inst.hasSentAux = true
			m := &message.ConsMessage{
				Type:     message.AUX,
				Proposer: msg.Proposer,
				Round:    inst.round,
				Sequence: msg.Sequence,
				Content:  []byte{0}} // aux 0
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
		}
		inst.isFastDecided()
	case message.BVAL_ONE_COLLECTION:
		if inst.oneEndorsed || inst.zeroEndorsed || inst.hasSentAux {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection {
			if inst.bvalOneSigns[msg.Round][i] != nil || sign == nil || !VerifySign(sign, []byte{1}, inst.priv) {
				continue
			}
			inst.numBvalOne[msg.Round]++
			inst.bvalOneSigns[msg.Round][i] = sign
		}
		inst.oneEndorsed = true
		// inst.tp.Broadcast(msg)
		if !inst.hasSentAux {
			inst.hasSentAux = true
			m := &message.ConsMessage{
				Type:     message.AUX,
				Proposer: msg.Proposer,
				Round:    inst.round,
				Sequence: msg.Sequence,
				Content:  []byte{1}} // aux 1
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
		}
		inst.isFastDecided()
	case message.AUX_ZERO_COLLECTION:
		if inst.fastAuxZero || inst.fastAuxOne || inst.round != msg.Round {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection {
			if inst.auxZeroSigns[msg.Round][i] != nil || sign == nil || !VerifySign(sign, []byte{0}, inst.priv) {
				continue
			}
			inst.numBvalZero[msg.Round]++
			inst.bvalZeroSigns[msg.Round][i] = sign
		}
		inst.fastAuxZero = true
		inst.zeroEndorsed = true
		if inst.canSkipCoin[inst.round] {
			inst.tp.Broadcast(&message.ConsMessage{
				Type:     message.SKIP,
				Proposer: msg.Proposer,
				Round:    inst.round,
				Sequence: msg.Sequence,
				Content:  []byte{0},
			})
		}
		return inst.isFastDecided()
	case message.AUX_ONE_COLLECTION:
		if inst.fastAuxZero || inst.fastAuxOne || inst.round != msg.Round {
			return false, false
		}
		collection := deserialCollection(msg.Collection)
		for i, sign := range collection {
			if inst.auxOneSigns[msg.Round][i] != nil || sign == nil || !VerifySign(sign, []byte{1}, inst.priv) {
				continue
			}
			inst.numBvalOne[msg.Round]++
			inst.bvalOneSigns[msg.Round][i] = sign
		}
		inst.fastAuxOne = true
		inst.oneEndorsed = true
		if inst.canSkipCoin[inst.round] {
			inst.tp.Broadcast(&message.ConsMessage{
				Type:     message.SKIP,
				Proposer: msg.Proposer,
				Round:    inst.round,
				Sequence: msg.Sequence,
				Content:  []byte{1},
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

func (inst *instance) getProposal() *message.ConsMessage {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	return inst.proposal
}

func serialCollection(collection [][]byte) []byte {
	mar_collection, err := json.Marshal(collection)
	if err != nil {
		panic("Marshal collection failed")
	}

	return mar_collection
}

func deserialCollection(data []byte) [][]byte {
	var collection [][]byte
	err := json.Unmarshal(data, &collection)
	if err != nil {
		panic("Unmarshal collection failed")
	}
	return collection
}

func GetSign(msg *message.ConsMessage, priv *ecdsa.PrivateKey) {
	hash, err := sha256.ComputeHash(msg.Content)
	if err != nil {
		panic("sha256 computeHash failed")
	}
	sig, err := myecdsa.SignECDSA(priv, hash)
	if err != nil {
		panic("myecdsa signECDSA failed")
	}

	msg.Signature = sig
}

func VerifySign(sign, content []byte, priv *ecdsa.PrivateKey) bool {
	hash, _ := sha256.ComputeHash(content)
	b, err := myecdsa.VerifyECDSA(&priv.PublicKey, sign, hash)
	if err != nil {
		fmt.Println("Failed to verify a message: ", err)
	}
	return b
}
