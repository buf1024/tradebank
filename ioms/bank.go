package ioms

import (
	"github.com/golang/protobuf/proto"
)

// ActionNext the next Action
type ActionNext int8

const (
	// ActionNextNone discard the message
	ActionNextNone ActionNext = iota
	// ActionNextBank proccess message the bank
	ActionNextBank
	// ActionNextExchange proccess message to exchange
	ActionNextExchange
)

// Action next proccess
type Action struct {
	// SID use for session manager, empty just ignore
	SID string
	// Next next proccess step
	Next ActionNext
	// Msg next proccess message
	Msg []byte
	// Session session data
	Session interface{}
}

// TimerHandler timer handler
type TimerHandler func() ([]Action, error)

// MyBank define the interface of bank plugin
type MyBank interface {
	Name() string
	ID() int64

	Init(s *Server) error

	ExchReq(command int64, msg proto.Message) ([]Action, error)
	ExchRsp(command int64, msg proto.Message) ([]Action, error)

	BankReq([]byte) ([]Action, error)
	BankRsp([]byte) ([]Action, error)
}

var mybank map[string]MyBank

// Register Register the bank
func Register(id string, bank MyBank) {
	if mybank == nil {
		mybank = make(map[string]MyBank)
	}
	if _, exists := mybank[id]; !exists {
		mybank[id] = bank
	}
}

// GetBank return the current bank
func GetBank(id string) MyBank {
	if bank, exists := mybank[id]; exists {
		return bank
	}
	return nil
}
