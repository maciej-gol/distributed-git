package main

import (
	"fmt"
	"github.com/hashicorp/raft"
	raftbolt "github.com/hashicorp/raft-boltdb"
	"io"
	"net"
	"os"
	"time"
)

type GitFSM struct{}

func (self *GitFSM) Apply(log *raft.Log) interface{} {
	fmt.Printf("Would apply %v\n", log)
	return nil
}

func (self *GitFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (self *GitFSM) Restore(io.ReadCloser) error {
	return nil
}

func main() {
	store, err := raftbolt.NewBoltStore(fmt.Sprintf("/tmp/bolt/%s", os.Args[1]))
	fmt.Println(err)
	snapshot, err := raft.NewFileSnapshotStore(fmt.Sprintf("/tmp/snapshots/%s", os.Args[1]), 2, os.Stdout)
	fmt.Println(err)
	a, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%s", os.Args[1]))
	fmt.Println(err)
	transport, err := raft.NewTCPTransport(a.String(), a, 5, 10*time.Second, os.Stdout)
	fmt.Println(err)
	fmt.Println("Listening on", a)
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(fmt.Sprintf(":%s", os.Args[1]))
	r, err := raft.NewRaft(config, &GitFSM{}, store, store, snapshot, transport)
	r.BootstrapCluster(raft.Configuration{[]raft.Server{{Suffrage: raft.Voter, ID: config.LocalID, Address: raft.ServerAddress(a.String())}}})
	fmt.Println(err)
	time.Sleep(10 * time.Second)
	if len(os.Args) > 2 {
		for i := 2; i < len(os.Args); i++ {
			fmt.Println(os.Args[i])
			r.AddVoter(
				raft.ServerID(fmt.Sprintf(":%s", os.Args[i])),
				raft.ServerAddress(fmt.Sprintf("127.0.0.1:%s", os.Args[i])),
				0,
				10*time.Second,
			)
		}
	}
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if r.Leader() == raft.ServerAddress(a.String()) {
					r.Apply([]byte("SOME CMD"), 1*time.Second)
				}
			}
		}
	}()
	time.Sleep(100 * time.Second)
	ticker.Stop()
	done <- true
}
