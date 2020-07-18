package operator

import (
	"github.com/DryginAlexander/notifier"
)

type Operator struct {
	storage notifier.Storage
}

func NewOperator(stor notifier.Storage) Operator {
	return Operator{stor}
}
