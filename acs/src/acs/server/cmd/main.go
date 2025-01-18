package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.themix.io/acs/server"
	"go.themix.io/crypto/bls"
	"go.themix.io/transport/info"
)

const CONFIG_FILE = "node.json"

func newLogger(id int) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"log/server" + strconv.Itoa(id),
	}
	cfg.Sampling = nil
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	return cfg.Build()
}

type Configuration struct {
	Id      uint64 `json:"id"`
	Port    int    `json:"port"`
	Key     string `json:"key_path"`
	Cluster string `json:"cluster"`
}

func main() {
	jsonFile, err := os.Open(CONFIG_FILE)
	if err != nil {
		panic(fmt.Sprint("os.Open: ", err))
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Sprint("io.ReadAll: ", err))
	}
	var config Configuration
	json.Unmarshal([]byte(data), &config)

	lg, err := newLogger(int(config.Id))
	if err != nil {
		panic(fmt.Sprintf("newLogger: %v", err))
	}
	defer lg.Sync()

	addrs := strings.Split(config.Cluster, ",")
	fmt.Printf("%d %s %d\n", config.Id, addrs, len(addrs))

	bls, err := bls.InitBLS(config.Key, len(addrs), int(len(addrs)/3+1), int(config.Id))
	if err != nil {
		panic(fmt.Sprint("bls.InitBLS: ", err))
	}

	server.InitNode(lg, bls, info.IDType(config.Id), uint64(len(addrs)), config.Port, addrs)
}
