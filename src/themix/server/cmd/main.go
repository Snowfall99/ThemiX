package main

import (
	"flag"
	"fmt"
	"strings"

	"go.themix.io/crypto/bls"
	myecdsa "go.themix.io/crypto/ecdsa"
	"go.themix.io/themix/config"
	"go.themix.io/themix/logger"
	"go.themix.io/themix/server"
	"go.themix.io/transport/http"
	"go.uber.org/zap"
)

var DEBUG = flag.Bool("debug", true, "enable debug logging")
var SIGN = flag.Bool("sign", false, "enable client request signature verification")
var BATCH = flag.Int("batch", 1, "the number of client requests in a batch")

func main() {
	flag.Parse()

	config, err := config.ReadConfig("node.json")
	if err != nil {
		panic(err)
	}

	lg, err := logger.NewLogger(int(config.Id), *DEBUG)
	if err != nil {
		panic(err)
	}
	defer lg.Sync()

	addrs := strings.Split(config.Cluster, ",")
	lg.Info("Init",
		zap.Int("ID", int(config.Id)),
		zap.Strings("Addresses", addrs),
		zap.Int("Nodes", len(addrs)))
	bls, err := bls.InitBLS(config.Key, len(addrs), int(len(addrs)/2+1), int(config.Id))
	if err != nil {
		panic(fmt.Sprint("bls.InitBLS: ", err))
	}
	pk, _ := myecdsa.LoadKey(config.Pk)
	var peers []http.Peer
	for i := 0; i < len(addrs); i++ {
		peer := http.Peer{
			PeerID:    uint32(i),
			PublicKey: pk,
			Addr:      addrs[i],
		}
		peers = append(peers, peer)
	}
	ckPath := config.Ck
	ck, err := myecdsa.LoadKey(ckPath)
	if err != nil {
		panic(fmt.Sprintf("ecdsa.LoadKey: %v", err))
	}
	server.InitNode(lg, bls, config.Pk, uint32(config.Id), uint64(len(addrs)), config.Port, peers, *BATCH, ck, *SIGN)
}
