// Copyright 2015 The etcd Authors
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

package raftnode

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"

	"go.etcd.io/etcd/raft/raftpb"
)

// Handler for a http based key-value store backed by raft
type httpKVAPI struct {
	store       *Store
	confChangeC chan<- raftpb.ConfChange
}

func handleConnection(conn net.Conn, kv *Store) {
	for {
		dec := gob.NewDecoder(conn)
		message := &Message{}
		dec.Decode(message)
		kv.Propose(message.RPCfunctionNum, message.data)
		fmt.Printf("Received : %+v", message)
	}
	conn.Close()
}

// https://stackoverflow.com/questions/11202058/unable-to-send-gob-data-over-tcp-in-go-programming/11202252#11202252
func Serve(kv *Store, port int, confChangeC chan<- raftpb.ConfChange, errorC <-chan error) {
	fmt.Println("start")
	Addr := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", Addr)
	if err != nil {
		// handle error
		panic(err)
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn, kv) // a goroutine handles conn so that the loop can accept other connections
	}
}

///////////////////////////////
// 	srv := http.Server{
// 		Addr: ":" + strconv.Itoa(port),
// 		Handler: &httpKVAPI{
// 			store:       kv,
// 			confChangeC: confChangeC,
// 		},
// 	}
// 	go func() {
// 		if err := srv.ListenAndServe(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	// exit when raft goes down
// 	if err, ok := <-errorC; ok {
// 		log.Fatal(err)
// 	}
// }
