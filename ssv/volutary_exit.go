package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
)

type VoluntaryExitMsgType uint64
type RequestID [24]byte

const (
	VoluntaryExitInitMsgType VoluntaryExitMsgType = iota
	VoluntaryExitPartialSigMsgType
)

type VoluntaryExitMessage struct {
	MsgType   VoluntaryExitMsgType
	RequestID RequestID
	Data      []byte
}

type SignedVoluntaryExitMessage struct {
	Message   *VoluntaryExitMessage
	Signer    []types.OperatorID
	Signature types.Signature
}

// VoluntaryExitInit is the first message sent to initiate a voluntary exit signature from operators
type VoluntaryExitInit struct {
	ValidatorPK types.ValidatorPK
	ExitMessage *spec.VoluntaryExit
}

// VoluntaryExitPartialSig is sent by each operator in response to VoluntaryExitInit
type VoluntaryExitPartialSig struct {
	Signature types.Signature
}
