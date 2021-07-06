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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"go.themix.io/crypto/bls"
	"go.themix.io/themix/server"
	"go.themix.io/transport/info"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// func initZapLog() *zap.Logger {
// 	config := zap.NewDevelopmentConfig()
// 	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
// 	config.EncoderConfig.TimeKey = "timestamp"
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	logger, _ := config.Build()
// 	return logger
// }

func newLogger(id int) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"log/server" + strconv.Itoa(id),
	}
	cfg.Sampling = nil
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	return cfg.Build()
}

func removeLastRune(s string) string {
	r := []rune(s)
	return string(r[:len(r)-1])
}

func main() {

	// loggerMgr := initZapLog()
	// zap.ReplaceGlobals(loggerMgr)
	// defer loggerMgr.Sync() // flushes buffer, if any
	// logger := loggerMgr.Sugar()
	// logger.Debug("START!")

	id := flag.Uint64("id", 0, "process ID")
	port := flag.Int("port", 11200, "port for themix server")
	keys := flag.String("keys", "keys", "the folder sotring keys")
	cluster := flag.String("cluster", "http://127.0.0.1:11200", "cluster members seperated by comma")
	clusterFile := flag.String("cluster-file", "address", "cluster members defined in the given file")
	// number := flag.Int("number", 10000, "number for benchmark test")
	// size := flag.Int("size", 10000, "content size for benchmark test")
	flag.Parse()

	addrs := []string{}
	lg, err := newLogger(int(*id))
	defer lg.Sync()

	if err != nil {
		fmt.Println("zap logger initialization failed: ", err)
		return
	}

	file, err := os.Open(*clusterFile)
	if err != nil {
		fmt.Println(err)
		addrs = strings.Split(*cluster, ",")
	} else {
		reader := bufio.NewReader(file)
		for {
			addr, err := reader.ReadString(' ')
			if err != nil && err == io.EOF {
				break
			}
			prt, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}
			addr = "http://" + removeLastRune(addr) + ":" + removeLastRune(prt)
			addrs = append(addrs, addr)
		}
	}
	defer file.Close()

	fmt.Printf("%d %s %d\n", *id, addrs, len(addrs))

	bls, err := bls.InitBLS(*keys, len(addrs), int(len(addrs)/3+1), int(*id))

	if err != nil {
		fmt.Println(err)
		return
	}

	server.InitNode(lg, bls, info.IDType(*id), uint64(len(addrs)), *port, addrs)

	// time.Sleep(5 * time.Second)

	// if *id == 0 {
	// 	msg := &transport.ConsMessage{To: 1, From: 0,
	// 		Content: make([]byte, *size)}
	// 	start := time.Now().UnixNano()
	// 	for i := 0; i < int(*number); i++ {
	// 		tp.Send(msg)
	// 	}
	// 	end := time.Now().UnixNano()
	// 	total := (int64)((*number) * (*size))
	// 	fmt.Printf("throughput is %d MB/s\n",
	// 		(total * 1000 / (end - start)))
	// }
}
