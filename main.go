package main

import (
	"time"

	"github.com/andrei0427/go-blockchain/network"
)

func main() {
	trLocal := network.NewLocalTransport("localhost")
	trRemote := network.NewLocalTransport("remote")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{
			trLocal,
			trRemote,
		},
	}

	s := network.NewServer(opts)
	s.Start()
}
