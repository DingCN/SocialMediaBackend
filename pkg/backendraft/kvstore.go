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

package backendraft

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"

	"go.etcd.io/etcd/etcdserver/api/snap"
)

// a key-value store backed by raft
type kvstore struct {
	proposeC    chan<- string // channel for proposing updates
	mu          sync.RWMutex
	Store       storage // current committed key-value pairs
	snapshotter *snap.Snapshotter
}

type kv struct {
	RPCfunctionNum int32
	Data           []byte
}

func newKVStore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *string, errorC <-chan error) *kvstore {
	s := &kvstore{proposeC: proposeC,
		Store: storage{
			UserList:         userlist{Users: map[string]*User{}},
			CentralTweetList: centraltweetlist{Tweets: []Tweet{}},
		},
		snapshotter: snapshotter}

	// replay log into key-value map
	s.readCommits(commitC, errorC)
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	return s
}

// func (s *kvstore) Lookup(key string) (string, bool) {
// 	s.mu.RLock()
// 	v, ok := s.kvStore[key]
// 	s.mu.RUnlock()
// 	return v, ok
// }

func (s *kvstore) GetUser(username string) (*User, error) {
	pUser, err := s.Store.GetUser(username)
	if err != nil {
		return nil, err
	}
	return pUser, nil
}

func (s *kvstore) GetUserProfile(username string) (*User, error) {
	pUser, err := s.Store.GetUserProfile(username)
	if err != nil {
		return nil, err
	}
	return pUser, nil
}
func (s *kvstore) GetTweetByUsername(username string) ([]Tweet, error) {
	tweetlist, err := s.Store.GetTweetByUsername(username)
	if err != nil {
		return nil, err
	}
	return tweetlist, nil
}

func (s *kvstore) MomentRandomFeeds() []Tweet {
	tweets := s.Store.MomentRandomFeeds()
	return tweets
}

func (s *kvstore) GetFollowingTweets(username string) ([]Tweet, error) {
	tweetlist, err := s.Store.GetFollowingTweets(username)
	if err != nil {
		return nil, err
	}
	return tweetlist, nil
}

func (s *kvstore) GetAllFollowing(username string) ([]string, error) {
	followinglist, err := s.Store.GetAllFollowing(username)
	if err != nil {
		return nil, err
	}
	return followinglist, nil
}

func (s *kvstore) CheckIfFollowing(username string, targetname string) (bool, error) {
	success, err := s.Store.CheckIfFollowing(username, targetname)
	if err != nil {
		return false, err
	}
	return success, nil
}

func (s *kvstore) Propose(RPCfunctionNum int32, data []byte) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv{RPCfunctionNum, data}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("buf: %+v", buf)
	s.proposeC <- buf.String()
	// sleep for operations to complete
	// TODO change to block
	time.Sleep(500 * time.Millisecond)
}

func (s *kvstore) readCommits(commitC <-chan *string, errorC <-chan error) {
	for data := range commitC {
		if data == nil {
			// done replaying log; new data incoming
			// OR signaled to load snapshot
			snapshot, err := s.snapshotter.Load()
			if err == snap.ErrNoSnapshot {
				return
			}
			if err != nil {
				log.Panic(err)
			}
			log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
			if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
				log.Panic(err)
			}
			continue
		}
		fmt.Println("data: %+v", data)
		var dataKv kv
		dec := gob.NewDecoder(bytes.NewBufferString(*data))
		if err := dec.Decode(&dataKv); err != nil {
			continue
			log.Fatalf("raftexample: could not decode message (%v)", err)
		}
		// if len(dataKv.data) == 0 {
		// 	continue
		// }
		if dataKv.RPCfunctionNum == protocol.Functions_FunctionName_value["SignupRPC"] {
			type st struct {
				Username string
				Password string
			}
			var store st
			json.Unmarshal(dataKv.Data, &store)
			s.mu.Lock()
			s.Store.AddUser(store.Username, store.Password)
			s.mu.Unlock()
		} else if dataKv.RPCfunctionNum == protocol.Functions_FunctionName_value["FollowUnFollowRPC"] {
			type st struct {
				Username   string
				Targetname string
			}
			var store st
			json.Unmarshal(dataKv.Data, &store)
			s.mu.Lock()

			s.Store.FollowUnFollow(store.Username, store.Targetname)
			s.mu.Unlock()
		} else if dataKv.RPCfunctionNum == protocol.Functions_FunctionName_value["AddTweetRPC"] {
			type st struct {
				Username  string
				Timestamp protocol.Timestamp
				Post      string
			}
			var store st
			json.Unmarshal(dataKv.Data, &store)
			s.mu.Lock()
			fmt.Println("timestamp commited %+v\n", store.Timestamp)
			s.Store.AddTweet(store.Username, store.Timestamp, store.Post)
			fmt.Printf("readcommitsCentralTweetList length: %d", len(s.Store.CentralTweetList.Tweets))
			s.mu.Unlock()
		} else {
			// s.mu.Lock()
			// s.kvStore[dataKv.Key] = dataKv.Val
			// s.mu.Unlock()
		}
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *kvstore) getSnapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.Store)
}

func (s *kvstore) recoverFromSnapshot(snapshot []byte) error {
	var store storage
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.mu.Lock()
	s.Store = store
	s.mu.Unlock()
	return nil
}
