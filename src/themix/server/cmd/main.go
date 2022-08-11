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
)

func removeLastRune(s string) string {
	r := []rune(s)
	return string(r[:len(r)-1])
}

func main() {
	sign := flag.Bool("sign", false, "whether to verify client sign or not")
	batchsize := flag.Int("batch", 1, "how many times for a client signature being verified")
	flag.Parse()

	config, err := config.ReadConfig("node.json")
	if err != nil {
		panic(err)
	}

	lg, err := logger.NewLogger(int(config.Id))
	if err != nil {
		panic(err)
	}
	defer lg.Sync()

	addrs := strings.Split(config.Cluster, ",")
	fmt.Printf("%d %s %d\n", config.Id, addrs, len(addrs))
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
	server.InitNode(lg, bls, config.Pk, uint32(config.Id), uint64(len(addrs)), config.Port, peers, *batchsize, ck, *sign)
}
