package network

import (
	"fmt"
	"time"

	"github.com/andrei0427/go-blockchain/core"
	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	RPCHandler RPCHandler
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {

	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	s := &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	if s.RPCHandler == nil {
		s.RPCHandler = NewDefaultRPCHandler(s)
	}

	return s

}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
				logrus.Error(err)
			}

		case <-s.quitCh:
			break free

		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}

	}

	fmt.Println("Server shutting down")
}

func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already inserted in pool")

		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeenOn(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":     hash,
		"poolSize": s.memPool.Len(),
	}).Info("adding new tx to pool")

	// TODO: Broadcast to other peers

	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}

		}(tr)

	}
}
