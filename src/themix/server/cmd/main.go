package main

import (
	"flag"
	"strings"

	"go.themix.io/crypto/bls"
	"go.themix.io/crypto/ecdsa"
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
		panic(err)
	}
	pk, err := ecdsa.LoadKey(config.Pk)
	if err != nil {
		panic(err)
	}
	ck, err := ecdsa.LoadKey(config.Ck)
	if err != nil {
		panic(err)
	}
	var peers []http.Peer
	for i := 0; i < len(addrs); i++ {
		peer := http.Peer{
			PeerID:    uint32(i),
			PublicKey: pk,
			Addr:      addrs[i],
		}
		peers = append(peers, peer)
	}
	server.InitNode(lg, bls, pk, ck, uint32(config.Id), uint64(len(addrs)), config.Port, peers, *BATCH, *SIGN)
}
