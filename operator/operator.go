package operator

import (
	"github.com/DryginAlexander/notifier"
)

type Client struct {
	channel chan string
	finish  func()
}

type Operator struct {
	clients map[string]Client
	storage     notifier.Storage
}

func NewOperator(stor notifier.Storage) Operator {
	return Operator{
		make(map[string]Client),
		stor,
	}
}
