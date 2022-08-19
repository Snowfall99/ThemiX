package server

import (
	"crypto/ecdsa"
	"runtime"

	"go.themix.io/crypto/bls"
	"go.themix.io/transport"
	"go.themix.io/transport/http"
	"go.themix.io/transport/proto/consmsgpb"
	"go.uber.org/zap"
)

type Node struct {
	state *state
	msgc  chan *consmsgpb.WholeMessage
}

// InitNode initiate a node for processing messages
func InitNode(lg *zap.Logger, blsSig *bls.BlsSig, pk *ecdsa.PrivateKey, ck *ecdsa.PrivateKey, id uint32, n uint64, port int, peers []http.Peer, batch int, sign bool) {
	tp, msgc, reqc, repc := transport.InitTransport(lg, id, port, peers, ck, sign, batch)
	proposer := initProposer(lg, tp, id, reqc)
	state := initState(lg, tp, blsSig, pk, id, proposer, n, repc, batch)
	node := &Node{
		state: state,
		msgc:  msgc,
	}
	node.run()
}

func (n *Node) run() {
	for i := 0; i < runtime.NumCPU()-1; i++ {
		go n.insertMsg()
	}
	n.insertMsg()
}

func (n *Node) insertMsg() {
	for {
		msg := <-n.msgc
		n.state.insertMsg(msg)
	}
}
