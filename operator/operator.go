package operator

import (
	"github.com/DryginAlexander/notifier"
)

type Operator struct {
	noteChanals map[string]chan string
	storage     notifier.Storage
}

func NewOperator(stor notifier.Storage) Operator {
	return Operator{
		make(map[string]chan string),
		stor,
	}
}
