package types

import "github.com/attestantio/go-eth2-client/spec/phase0"

type ExitMessage struct {
	ValidatorPubKey ValidatorPK `ssz-size:"48"`
	Message         *phase0.VoluntaryExit
}

func (msg *ExitMessage) Encode() ([]byte, error) {
	return msg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msg *ExitMessage) Decode(data []byte) error {
	return msg.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (msg *ExitMessage) GetRoot() ([32]byte, error) {
	return msg.HashTreeRoot()
}

type SignedExitMessage struct {
	// Signature using the ethereum address that registered the validator
	Signature [65]byte `ssz-size:"65"`
	Message   ExitMessage
}

func (msg *SignedExitMessage) Encode() ([]byte, error) {
	return msg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msg *SignedExitMessage) Decode(data []byte) error {
	return msg.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (msg *SignedExitMessage) GetRoot() ([32]byte, error) {
	return msg.HashTreeRoot()
}
