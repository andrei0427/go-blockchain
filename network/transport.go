package network

type NetAddr string

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(to NetAddr, bytes []byte) error
	Broadcast(bytes []byte) error
	Addr() NetAddr
}
