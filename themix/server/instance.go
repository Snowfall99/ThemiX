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
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"go.themix.io/crypto/bls"
	myecdsa "go.themix.io/crypto/ecdsa"
	"go.themix.io/crypto/sha256"
	"go.themix.io/transport"
	"go.themix.io/transport/info"
	"go.themix.io/transport/message"
	"go.uber.org/zap"
)

// the maximum expected round that terminates consensus, P = 1 - pow(0.5, maxround)
var maxround = 30

type instance struct {
	tp           transport.Transport
	blsSig       *bls.BlsSig
	fastRBC      bool
	hasEcho      bool
	hasVotedZero bool
	hasVotedOne  bool
	hasSentAux   bool
	hasSentCoin  bool
	zeroEndorsed bool
	oneEndorsed  bool
	isDecided    bool
	isFinished   bool
	sequence     uint64
	n            uint64
	thld         uint64
	f            uint64
	fastgroup    uint64
	round        uint8
	numEcho      uint64
	numReady     uint64
	numOneSkip   uint64
	numZeroSkip  uint64
	binVals      uint8
	lastCoin     uint8
	sin          []bool
	canSkipCoin  []bool
	numBvalZero  []uint64
	numBvalOne   []uint64
	numAuxZero   []uint64
	numAuxOne    []uint64
	numCon       []uint64
	numCoin      []uint64
	proposal     *message.ConsMessage
	valMsgs      []*message.ConsMessage
	echoMsgs     []*message.ConsMessage
	readyMsgs    []*message.ConsMessage
	bvalZeroMsgs [][]*message.ConsMessage
	bvalOneMsgs  [][]*message.ConsMessage
	auxZeroMsgs  [][]*message.ConsMessage
	auxOneMsgs   [][]*message.ConsMessage
	coinMsgs     [][]*message.ConsMessage
	startR       bool
	startB       bool
	startS       bool
	tmrR         time.Timer
	tmrB         time.Timer
	tmrS         time.Timer
	priv         *ecdsa.PrivateKey
	lg           *zap.Logger
	lock         sync.Mutex
}

func initInstance(lg *zap.Logger, tp transport.Transport, blsSig *bls.BlsSig, pkPath string, sequence uint64, n uint64, thld uint64) *instance {
	inst := &instance{
		lg:           lg,
		tp:           tp,
		blsSig:       blsSig,
		sequence:     sequence,
		n:            n,
		thld:         thld,
		f:            n / 2,
		sin:          make([]bool, maxround),
		canSkipCoin:  make([]bool, maxround),
		valMsgs:      make([]*message.ConsMessage, n),
		echoMsgs:     make([]*message.ConsMessage, n),
		readyMsgs:    make([]*message.ConsMessage, n),
		bvalZeroMsgs: make([][]*message.ConsMessage, maxround),
		bvalOneMsgs:  make([][]*message.ConsMessage, maxround),
		auxZeroMsgs:  make([][]*message.ConsMessage, maxround),
		auxOneMsgs:   make([][]*message.ConsMessage, maxround),
		coinMsgs:     make([][]*message.ConsMessage, maxround),
		numBvalZero:  make([]uint64, maxround),
		numBvalOne:   make([]uint64, maxround),
		numAuxZero:   make([]uint64, maxround),
		numAuxOne:    make([]uint64, maxround),
		numCon:       make([]uint64, maxround),
		numCoin:      make([]uint64, maxround),
		lock:         sync.Mutex{}}
	inst.fastgroup = uint64(math.Ceil(3*float64(inst.f)/2)) + 1
	inst.priv, _ = myecdsa.LoadKey(pkPath)
	for i := 0; i < maxround; i++ {
		inst.bvalZeroMsgs[i] = make([]*message.ConsMessage, n)
		inst.bvalOneMsgs[i] = make([]*message.ConsMessage, n)
		inst.auxZeroMsgs[i] = make([]*message.ConsMessage, n)
		inst.auxOneMsgs[i] = make([]*message.ConsMessage, n)
		inst.coinMsgs[i] = make([]*message.ConsMessage, n)
		inst.canSkipCoin[i] = true
	}
	return inst
}

// return true if the instance is decided or finished at the first time
func (inst *instance) insertMsg(msg *message.ConsMessage) (bool, bool) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	// // Just for test
	if msg.Round > 0 {
		return false, false
	}

	if inst.isFinished {
		return false, false
	}

	if inst.fastRBC && msg.Type == message.READY {
		return false, false
	}

	if len(msg.Content) > 0 {
		inst.lg.Info("receive msg",
			zap.String("type", msg.Type.GetName()),
			zap.Int("proposer", int(msg.Proposer)),
			zap.Int("seq", int(msg.Sequence)),
			zap.Int("round", int(msg.Round)),
			zap.Int("from", int(msg.From)),
			zap.Int("content", int(msg.Content[0])))
	} else {
		inst.lg.Info("receive msg",
			zap.String("type", msg.Type.GetName()),
			zap.Int("proposer", int(msg.Proposer)),
			zap.Int("seq", int(msg.Sequence)),
			zap.Int("round", int(msg.Round)),
			zap.Int("from", int(msg.From)))
	}

	switch msg.Type {

	/*
	 * upon receiving VAL(v)src
	 * if have not sent any ECHO(*) then
	 * broadcast VAL(v)src, ECHO(v)i
	 * start timer tmrR <- 2*delta
	 */
	case message.VAL:
		VerifySign(*msg, inst.priv)
		inst.proposal = msg
		hash, _ := sha256.ComputeHash(msg.Content)
		inst.valMsgs[msg.From] = msg
		if !inst.hasEcho {
			// broadcast VAL(v)src, ECHO(v)i
			m := &message.ConsMessage{
				Type:     message.ECHO,
				Proposer: msg.Proposer,
				Sequence: msg.Sequence,
				Content:  hash}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
			inst.hasEcho = true
			m = &message.ConsMessage{
				Type:      message.VAL,
				Proposer:  msg.Proposer,
				Sequence:  msg.Sequence,
				Content:   hash,
				Signature: msg.Signature,
			}
			inst.tp.Broadcast(m)
		}
		inst.isReadyToSendCoin()
		return inst.isReadyToEnterNewRound()

	/*
	 * upon receiving f+1 ECHO(v), and tmrR expires
	 * if have not received any VAL(v')src (v' != v) then
	 * broadcast READY(v)i
	 */
	case message.ECHO:
		verify := VerifySign(*msg, inst.priv)
		if !verify {
			return false, false
		}
		inst.numEcho++
		inst.echoMsgs[msg.From] = msg
		if inst.numEcho == inst.f+1 && !inst.startR {
			// start tmrR
			inst.startR = true
			inst.tmrR = *time.NewTimer(2 * time.Second)
			go func() {
				<-inst.tmrR.C
				// upon receiving f+1 ECHO(v)* and tmrR expires
				// if have not received any VAL(v') (v' != v) then
				// broadcast READY(v)i
				if inst.numEcho >= inst.f+1 {
					var content []byte
					for _, msg := range inst.echoMsgs {
						if msg != nil && content == nil {
							content = msg.Content
						} else if msg != nil && !bytes.Equal(content, msg.Content) {
							return
						}
					}
					m := &message.ConsMessage{
						Type:     message.READY,
						Proposer: msg.Proposer,
						Sequence: msg.Sequence,
						Content:  msg.Content,
					}
					GetSign(m, inst.priv)
					inst.tp.Broadcast(m)
				}
			}()
			inst.isReadyToSendCoin()
			return inst.isReadyToEnterNewRound()
		}
		/*
		 * upon receiving ECHO(v) from fast group
		 * broadcast ECHO(v) sent by fast group
		 * deliver(v)
		 */
		if inst.numEcho == inst.fastgroup && inst.round == 0 {
			inst.fastRBC = true
			if inst.startR {
				inst.tmrR.Stop()
				inst.startR = false
			}
			var collection []*message.ConsMessage
			for _, m := range inst.echoMsgs {
				if m != nil {
					collection = append(collection, m)
				}
			}
			data := serialCollection(collection)
			m := &message.ConsMessage{
				Type:       message.COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: data,
			}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
			if !inst.hasVotedZero && !inst.hasVotedOne {
				inst.hasVotedOne = true
				m := &message.ConsMessage{
					Type:     message.BVAL,
					Proposer: msg.Proposer,
					Sequence: msg.Sequence,
					Content:  []byte{1}} // vote 1
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
			return inst.isReadyToEnterNewRound()
		}

	/*
	 * upon receiving f+1 READY(v)
	 * broadcast f+1 READY(v)
	 * deliver(v)
	 */
	case message.READY:
		verify := VerifySign(*msg, inst.priv)
		if !verify {
			return false, false
		}
		inst.numReady++
		inst.readyMsgs[msg.From] = msg
		if inst.numReady >= inst.f+1 && inst.round == 0 && !inst.fastRBC {
			var collection []*message.ConsMessage
			for _, m := range inst.readyMsgs {
				if m != nil {
					collection = append(collection, m)
				}
			}
			data := serialCollection(collection)
			m := &message.ConsMessage{
				Type:       message.COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: data,
			}
			GetSign(m, inst.priv)
			inst.tp.Broadcast(m)
			if !inst.hasVotedZero && !inst.hasVotedOne {
				inst.hasVotedOne = true
				m := &message.ConsMessage{
					Type:     message.BVAL,
					Proposer: msg.Proposer,
					Sequence: msg.Sequence,
					Content:  []byte{1}} // vote 1
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
			return inst.isReadyToEnterNewRound()
		}

	/*
	 * upon receiving f+1 BVAL(b, r)
	 * insert b into vals
	 * broadcast f+1 BVAL(b, r)
	 * if have not sent AUX(*, r)i then
	 * broadcast AUX(b, r)i
	 * start tmrB <- 2*delta
	 */
	case message.BVAL:
		verify := VerifySign(*msg, inst.priv)
		if !verify {
			return false, false
		}
		var b bool
		switch msg.Content[0] {
		case 0:
			inst.numBvalZero[msg.Round]++
			inst.bvalZeroMsgs[msg.Round][msg.From] = msg
		case 1:
			inst.numBvalOne[msg.Round]++
			inst.bvalOneMsgs[msg.Round][msg.From] = msg
		}
		if inst.round == msg.Round && !inst.hasVotedZero && inst.numBvalZero[inst.round] > inst.f {
			inst.hasVotedZero = true
			var collection []*message.ConsMessage
			for _, m := range inst.bvalZeroMsgs[msg.Round] {
				if m != nil {
					collection = append(collection, m)
				}
			}
			data := serialCollection(collection)
			m := &message.ConsMessage{
				Type:       message.COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: data,
			}
			GetSign(m, inst.priv)
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
					Content:  []byte{0}} // aux 0
				GetSign(m, inst.priv)
				inst.tp.Broadcast(m)
			}
			inst.isReadyToSendCoin()
			b = true
		}
		if inst.round == msg.Round && !inst.hasVotedOne && inst.numBvalOne[inst.round] > inst.f {
			inst.hasVotedOne = true
			var collection []*message.ConsMessage
			for _, m := range inst.bvalOneMsgs[msg.Round] {
				if m != nil {
					collection = append(collection, m)
				}
			}
			data := serialCollection(collection)
			m := &message.ConsMessage{
				Type:       message.COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: data,
			}
			GetSign(m, inst.priv)
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
			inst.isReadyToSendCoin()
			b = true
		}
		if b && inst.round == msg.Round && !inst.startS {
			// start tmrS
			inst.startS = true
			inst.tmrS = *time.NewTimer(2 * time.Second)
			go func() {
				<-inst.tmrS.C
				// if tmrS expires, canSkipCoin = false
				inst.canSkipCoin[msg.Round] = false
			}()
		}
		/*
		 * upon receiving f+1 BVAL(b, r-1) and coin(r-1) = b
		 * if have not sent BVAL(b, r) or AUX(*, r) then
		 * if sin(r-1) = false then
		 * broadcast BVAL(b, r)
		 */
		if inst.round == msg.Round+1 {
			switch msg.Content[0] {
			case 0:
				if inst.numBvalZero[msg.Round] == inst.f+1 && inst.lastCoin == 0 &&
					(!inst.hasVotedZero || !inst.hasSentAux) &&
					!inst.sin[msg.Round] {
					m := &message.ConsMessage{
						Type:     message.BVAL,
						Proposer: msg.Proposer,
						Round:    inst.round,
						Sequence: msg.Sequence,
						Content:  []byte{0},
					}
					GetSign(m, inst.priv)
					inst.tp.Broadcast(m)
				}
			case 1:
				if inst.numBvalOne[msg.Round] == inst.f+1 && inst.lastCoin == 1 &&
					(!inst.hasVotedOne || !inst.hasSentAux) &&
					!inst.sin[msg.Round] {
					m := &message.ConsMessage{
						Type:     message.BVAL,
						Proposer: msg.Proposer,
						Round:    inst.round,
						Sequence: msg.Sequence,
						Content:  []byte{1},
					}
					GetSign(m, inst.priv)
					inst.tp.Broadcast(m)
				}
			}
		}
		if b {
			return inst.isReadyToEnterNewRound()
		}

	/*
	 * upon receiving f+1 AUX(*, r) and f+1 BVAL(b, r)
	 * for each b in AUX(*, r) and tmrB expires
	 * broadcast CON(vals, r)i, COIN(r)i
	 */
	case message.AUX:
		verify := VerifySign(*msg, inst.priv)
		if !verify {
			return false, false
		}
		switch msg.Content[0] {
		case 0:
			inst.numAuxZero[msg.Round]++
			inst.auxZeroMsgs[msg.Round][msg.From] = msg
		case 1:
			inst.numAuxOne[msg.Round]++
			inst.auxOneMsgs[msg.Round][msg.From] = msg
		}

		if inst.numAuxZero[msg.Round]+inst.numAuxOne[msg.Round] == inst.f+1 && !inst.startB {
			// start tmrB
			inst.startB = true
			inst.tmrB = *time.NewTimer(4 * time.Second)
			go func() {
				<-inst.tmrB.C
				// upon receiving f+1 AUX(*, r) and f+1 BVAL(b, r)
				// for each b in AUX(*, r) and tmrB expires
				// broadcast CON(vals, r)i, COIN(r)i
				if inst.numAuxZero[msg.Round]+inst.numAuxOne[msg.Round] > inst.f &&
					((inst.numAuxZero[msg.Round] != 0 && inst.numBvalZero[msg.Round] > inst.f) ||
						(inst.numAuxOne[msg.Round] != 0 && inst.numBvalOne[msg.Round] > inst.f)) {
					if inst.hasVotedZero {
						inst.tp.Broadcast(&message.ConsMessage{
							Type:     message.CON,
							Proposer: msg.Proposer,
							Round:    inst.round,
							Sequence: msg.Sequence,
							Content:  []byte{0},
						})
					}
					if inst.hasVotedOne {
						inst.tp.Broadcast(&message.ConsMessage{
							Type:     message.CON,
							Proposer: msg.Proposer,
							Round:    inst.round,
							Sequence: msg.Sequence,
							Content:  []byte{1},
						})
					}
					inst.isReadyToSendCoin()
				}
			}()
		}

		/*
		 * upon receiving AUX(b, r) from fast group
		 * broadcast AUX(b, r), sent by fast group, COIN(r)
		 * vals <- {b}
		 * if canSkipCoin(r) == true then
		 * broadcast SKIP(b, r)
		 * NEWROUND()
		 */
		if inst.round == msg.Round && msg.Content[0] == 0 && inst.numAuxZero[msg.Round] == inst.fastgroup {
			var collection []*message.ConsMessage
			for _, m := range inst.auxZeroMsgs[msg.Round] {
				if m != nil {
					collection = append(collection, m)
				}
			}
			data := serialCollection(collection)
			inst.tp.Broadcast(&message.ConsMessage{
				Type:       message.COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: data,
			})

			inst.zeroEndorsed = true
			inst.sin[inst.round] = true
			if inst.canSkipCoin[inst.round] {
				inst.tp.Broadcast(&message.ConsMessage{
					Type:     message.SKIP,
					Proposer: msg.Proposer,
					Round:    inst.round,
					Sequence: msg.Sequence,
					Content:  msg.Content,
				})
			}
			inst.isReadyToSendCoin()
			return inst.isReadyToEnterNewRound()
		}
		if inst.round == msg.Round && msg.Content[0] == 1 && inst.numAuxOne[msg.Round] == inst.fastgroup {
			var collection []*message.ConsMessage
			for _, m := range inst.auxOneMsgs[msg.Round] {
				if m != nil {
					collection = append(collection, m)
				}
			}
			data := serialCollection(collection)
			inst.tp.Broadcast(&message.ConsMessage{
				Type:       message.COLLECTION,
				Proposer:   msg.Proposer,
				Round:      inst.round,
				Sequence:   msg.Sequence,
				Collection: data,
			})

			inst.oneEndorsed = true
			inst.sin[inst.round] = true
			if inst.canSkipCoin[inst.round] {
				inst.tp.Broadcast(&message.ConsMessage{
					Type:     message.SKIP,
					Proposer: msg.Proposer,
					Round:    inst.round,
					Sequence: msg.Sequence,
					Content:  msg.Content,
				})
			}
			inst.isReadyToSendCoin()
			return inst.isReadyToEnterNewRound()
		}

		/*
		 * upon receiving AUX(b, r-1) from fast group
		 * if have not sent BVAL(b, r) or AUX(*, r) then
		 * broadcast BVAL(b, r)
		 */
		if inst.round == msg.Round+1 {
			switch msg.Content[0] {
			case 0:
				if !inst.hasVotedZero || !inst.hasSentAux {
					m := &message.ConsMessage{
						Type:     message.BVAL,
						Proposer: msg.Proposer,
						Round:    inst.round,
						Sequence: msg.Sequence,
						Content:  []byte{0},
					}
					GetSign(m, inst.priv)
					inst.tp.Broadcast(m)
				}
			case 1:
				if !inst.hasVotedOne || !inst.hasSentAux {
					m := &message.ConsMessage{
						Type:     message.BVAL,
						Proposer: msg.Proposer,
						Round:    inst.round,
						Sequence: msg.Sequence,
						Content:  []byte{1},
					}
					GetSign(m, inst.priv)
					inst.tp.Broadcast(m)
				}
			}
		}

		if inst.round == msg.Round {
			inst.isReadyToSendCoin()
			return inst.isReadyToEnterNewRound()
		}

	/*
	 * upon receiving f+1 CON(*, r), COIN(r)
	 * NewRound()
	 */
	case message.CON:
		inst.numCon[msg.Round]++
		if inst.round == msg.Round && inst.numCon[msg.Round] >= inst.f+1 && inst.numCoin[msg.Round] >= inst.f+1 {
			return inst.isReadyToEnterNewRound()
		}
	case message.COIN:
		inst.coinMsgs[msg.Round][msg.From] = msg
		inst.numCoin[msg.Round]++
		if inst.round == msg.Round && inst.numCon[msg.Round] >= inst.f+1 && inst.numCoin[msg.Round] >= inst.f+1 {
			return inst.isReadyToEnterNewRound()
		}
	case message.SKIP:
		switch msg.Content[0] {
		case 0:
			inst.numZeroSkip++
		case 1:
			inst.numOneSkip++
		}
		// if receiving SKIP(b, r') (any r') from fastgroup
		// decide(b)
		if msg.Content[0] == 0 && inst.proposal != nil && inst.numZeroSkip == inst.fastgroup {
			return inst.fastDecide(0)
		}
		if msg.Content[0] == 1 && inst.proposal != nil && inst.numOneSkip == inst.fastgroup {
			return inst.fastDecide(1)
		}
	case message.COLLECTION:
		if msg.Collection == nil {
			fmt.Println("Nil collection")
		}
		collection := deserialCollection(msg.Collection)
		for _, m := range collection {
			switch m.Type {
			case message.READY:
				// if inst.readyMsgs[m.From] == nil {
				// 	inst.numReady++
				// }
				inst.readyMsgs[m.From] = &m
			case message.BVAL:
				switch m.Content[0] {
				case 0:
					// if inst.bvalZeroMsgs[msg.Round][m.From] == nil {
					// 	inst.numBvalZero[msg.Round]++
					// }
					inst.bvalZeroMsgs[msg.Round][m.From] = &m
				case 1:
					// if inst.bvalOneMsgs[msg.Round][m.From] == nil {
					// 	inst.numBvalOne[msg.Round]++
					// }
					inst.bvalOneMsgs[msg.Round][m.From] = &m
				}
			case message.ECHO:
				// if inst.echoMsgs[m.From] == nil {
				// 	inst.numEcho++
				// }
				inst.echoMsgs[m.From] = &m
			case message.AUX:
				switch m.Content[0] {
				case 0:
					inst.auxZeroMsgs[msg.Round][m.From] = &m
				case 1:
					inst.auxOneMsgs[msg.Round][m.From] = &m
				}
			}
		}

	default:
		return false, false
	}
	return false, false
}

// must be executed within inst.lock
func (inst *instance) isReadyToSendCoin() {
	if !inst.hasSentCoin && inst.proposal != nil {
		if inst.oneEndorsed && inst.numAuxOne[inst.round] >= inst.thld {
			inst.binVals = 1
		} else if inst.zeroEndorsed && inst.numAuxZero[inst.round] >= inst.thld {
			inst.binVals = 0
		} else if inst.oneEndorsed && inst.zeroEndorsed &&
			inst.numAuxOne[inst.round]+inst.numAuxZero[inst.round] >= inst.thld {
			inst.binVals = 2
		} else {
			return
		}

		inst.hasSentCoin = true
		inst.tp.Broadcast(&message.ConsMessage{
			Type:     message.COIN,
			Proposer: inst.proposal.Proposer,
			Round:    inst.round,
			Sequence: inst.sequence,
			Content:  inst.blsSig.Sign(inst.getCoinInfo())}) // threshold bls sig share
	}
}

// must be executed within inst.lock
// return true if the instance is decided or finished at the first time
func (inst *instance) isReadyToEnterNewRound() (bool, bool) {
	if inst.hasSentCoin &&
		inst.numCoin[inst.round] > inst.f &&
		inst.numCon[inst.round] > inst.f &&
		inst.proposal != nil &&
		(inst.numReady >= inst.thld || inst.fastRBC) &&
		inst.numAuxZero[inst.round]+inst.numAuxOne[inst.round] >= inst.thld &&
		((inst.oneEndorsed && inst.numAuxOne[inst.round] >= inst.thld) ||
			(inst.zeroEndorsed && inst.numAuxZero[inst.round] >= inst.thld) ||
			(inst.oneEndorsed && inst.zeroEndorsed)) {
		sigShares := make([][]byte, 0)
		sigNum := 0
		for _, m := range inst.coinMsgs[inst.round] {
			if m != nil {
				sigNum++
				sigShares = append(sigShares, m.Content)
			}
		}
		if sigNum <= int(inst.f) {
			return false, false
		}
		coin := inst.blsSig.Recover(inst.getCoinInfo(), sigShares, int(inst.f+1), int(inst.n))

		inst.lg.Info("coin result",
			zap.Int("proposer", int(inst.proposal.Proposer)),
			zap.Int("seq", int(inst.sequence)),
			zap.Int("round", int(inst.round)),
			zap.Int("coin", int(coin[0]%2)))

		var nextVote byte
		if coin[0]%2 == inst.binVals {
			if inst.isDecided {
				inst.isFinished = true
				return false, true
			}

			inst.lg.Info("decided",
				zap.Int("proposer", int(inst.proposal.Proposer)),
				zap.Int("seq", int(inst.sequence)),
				zap.Int("round", int(inst.round)),
				zap.Int("result", int(coin[0]%2)))

			inst.isDecided = true
			nextVote = inst.binVals
		} else if inst.binVals != 2 { // nextVote should insist the single value
			nextVote = inst.binVals
		} else {
			nextVote = coin[0] % 2
		}

		if nextVote == 0 {
			inst.hasVotedZero = true
			inst.hasVotedOne = false
		} else {
			inst.hasVotedZero = false
			inst.hasVotedOne = true
		}
		inst.hasSentAux = false
		inst.hasSentCoin = false
		inst.zeroEndorsed = false
		inst.oneEndorsed = false
		inst.lastCoin = coin[0]
		if inst.startB {
			inst.tmrB.Stop()
			inst.startB = false
		}
		if inst.startS {
			inst.tmrS.Stop()
			inst.startS = false
		}
		inst.round++

		m := &message.ConsMessage{
			Type:     message.BVAL,
			Proposer: inst.proposal.Proposer,
			Round:    inst.round,
			Sequence: inst.sequence,
			Content:  []byte{nextVote}}
		GetSign(m, inst.priv)
		inst.tp.Broadcast(m)

		if coin[0]%2 == inst.binVals && inst.isDecided {
			return true, false
		}
	}
	return false, false
}

func (inst *instance) fastDecide(value int) (bool, bool) {
	if inst.isDecided {
		inst.isFinished = true
		return false, true
	}

	inst.lg.Info("fast decide",
		zap.Int("proposer", int(inst.proposal.Proposer)),
		zap.Int("seq", int(inst.sequence)),
		zap.Int("round", int(inst.round)),
		zap.Int("result", value))

	inst.isDecided = true
	nextVote := inst.binVals
	if nextVote == 0 {
		inst.hasVotedZero = true
		inst.hasVotedOne = false
	} else {
		inst.hasVotedZero = false
		inst.hasVotedOne = true
	}
	inst.hasSentAux = false
	inst.hasSentCoin = false
	inst.zeroEndorsed = false
	inst.oneEndorsed = false
	if inst.startB {
		inst.tmrB.Stop()
		inst.startB = false
	}
	if inst.startS {
		inst.tmrS.Stop()
		inst.startS = false
	}
	inst.round++

	m := &message.ConsMessage{
		Type:     message.BVAL,
		Proposer: inst.proposal.Proposer,
		Round:    inst.round,
		Sequence: inst.sequence,
		Content:  []byte{nextVote},
	}
	GetSign(m, inst.priv)
	inst.tp.Broadcast(m)
	return true, false
}

func (inst *instance) getCoinInfo() []byte {
	bsender := make([]byte, 8)
	binary.LittleEndian.PutUint64(bsender, uint64(inst.proposal.Proposer))
	bseq := make([]byte, 8)
	binary.LittleEndian.PutUint64(bseq, inst.sequence)

	b := make([]byte, 17)
	b = append(b, bsender...)
	b = append(b, bseq...)
	b = append(b, inst.round)

	return b
}

func (inst *instance) canVoteZero(sender info.IDType, seq uint64) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	if inst.round == 0 && !inst.hasVotedZero && !inst.hasVotedOne {
		inst.hasVotedZero = true
		m := &message.ConsMessage{
			Type:     message.BVAL,
			Proposer: sender,
			Round:    inst.round,
			Sequence: seq,
			Content:  []byte{0}} // vote 0
		GetSign(m, inst.priv)
		inst.tp.Broadcast(m)
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

func serialCollection(collection []*message.ConsMessage) []byte {
	mar_collection, err := json.Marshal(collection)
	if err != nil {
		panic("Marshal collection failed")
	}

	return mar_collection
}

func deserialCollection(data []byte) []message.ConsMessage {
	var collection []message.ConsMessage
	err := json.Unmarshal(data, &collection)
	if err != nil {
		panic("Unmarshal collection failed")
	}
	return collection
}

func GetSign(msg *message.ConsMessage, priv *ecdsa.PrivateKey) {
	msg_type, _ := json.Marshal(msg.Type)
	msg_proposer, _ := json.Marshal(msg.Proposer)
	msg_round, _ := json.Marshal(msg.Round)
	msg_sequence, _ := json.Marshal(msg.Sequence)
	msg_content, _ := json.Marshal(msg.Content)
	mar_msg := append(msg_type, msg_proposer...)
	mar_msg = append(mar_msg, msg_round...)
	mar_msg = append(mar_msg, msg_sequence...)
	mar_msg = append(mar_msg, msg_content...)

	hash, err := sha256.ComputeHash(mar_msg)
	if err != nil {
		panic("sha256 computeHash failed")
	}

	sig, err := myecdsa.SignECDSA(priv, hash)
	if err != nil {
		panic("myecdsa signECDSA failed")
	}

	msg.Signature = sig
}

func VerifySign(msg message.ConsMessage, priv *ecdsa.PrivateKey) bool {
	msg_type, _ := json.Marshal(msg.Type)
	msg_proposer, _ := json.Marshal(msg.Proposer)
	msg_round, _ := json.Marshal(msg.Round)
	msg_sequence, _ := json.Marshal(msg.Sequence)
	msg_content, _ := json.Marshal(msg.Content)
	mar_msg := append(msg_type, msg_proposer...)
	mar_msg = append(mar_msg, msg_round...)
	mar_msg = append(mar_msg, msg_sequence...)
	mar_msg = append(mar_msg, msg_content...)

	hash, err := sha256.ComputeHash(mar_msg)
	if err != nil {
		panic("sha256 computeHash failed")
	}

	b, err := myecdsa.VerifyECDSA(&priv.PublicKey, msg.Signature, hash)

	if err != nil {
		fmt.Println("Failed to verify a message: ", err)
	}

	return b
}
