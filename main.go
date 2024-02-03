package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/andrei0427/go-blockchain/core"
	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/network"
	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("localhost")
	trRemoteA := network.NewLocalTransport("remote_a")
	trRemoteB := network.NewLocalTransport("remote_b")
	trRemoteC := network.NewLocalTransport("remote_c")

	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)

	trRemoteA.Connect(trLocal)

	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	go func() {
		for {
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}

			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)

		trLate := network.NewLocalTransport("late")
		trRemoteC.Connect(trLate)
		sLate := makeServer(string(trLate.Addr()), trLate, nil)

		go sLate.Start()
	}()

	pk := crypto.NewPrivateKey()
	localServ := makeServer("localhost", trLocal, &pk)
	localServ.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("remote_%d", i+1)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		ID:         id,
		PrivateKey: pk,
		Transports: []network.Transport{
			tr,
		},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s

}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	pk := crypto.NewPrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(pk)

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Bytes())
}
