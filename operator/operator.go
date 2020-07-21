package operator

import (
	"sync"

	"github.com/DryginAlexander/notifier"
)

type ChanMessage struct {
	text string
	id   uint
}

type Client struct {
	username string
	channel  chan ChanMessage
	finish   func()
}

type Operator struct {
	clients map[string]Client
	storage notifier.Storage
	mutex   *sync.Mutex
}

func NewOperator(stor notifier.Storage) Operator {
	mu := &sync.Mutex{}
	return Operator{
		make(map[string]Client),
		stor,
		mu,
	}
}
