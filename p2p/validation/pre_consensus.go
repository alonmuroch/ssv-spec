package validation

import (
	"github.com/bloxapp/ssv-spec/ssv"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (validation *MessageValidation) validatePreConsensus() pubsub.ValidationResult {
	signedMsg, err := validation.decodePartialSignatureMessage(nil)
	if err != nil {
		return pubsub.ValidationReject
	}

	if validation.verifySignatureForShare(signedMsg, nil) != nil {
		return pubsub.ValidationReject
	}

	if !validation.preConsensusInTime(signedMsg) {
		return pubsub.ValidationReject
	}

	if !validation.isUniquePreConsensusMessageForSigner(signedMsg) {
		return pubsub.ValidationReject
	}

	for _, msg := range signedMsg.Message.Messages {

	}

	return pubsub.ValidationAccept
}

// isUniquePreConsensusMessageForSigner returns true if the message is a unique first time pre-consensus msg for signer
func (validation *MessageValidation) isUniquePreConsensusMessageForSigner(msg *ssv.SignedPartialSignatureMessage) bool {
	panic("implement")
}

// decodePartialSignatureMessage returns the decoded SignedPartialSignatureMessage or error
func (validation *MessageValidation) decodePartialSignatureMessage(data []byte) (*ssv.SignedPartialSignatureMessage, error) {
	panic("implement")
}

// validBeaconSignature returns nil if beacon signature is valid
func (validation *MessageValidation) validBeaconSignature(message *ssv.PartialSignatureMessage) error {
	panic("implement")
}

// preConsensusInTime returns true if pre-consensus message is received in time
func (validation *MessageValidation) preConsensusInTime(msg *ssv.SignedPartialSignatureMessage) bool {
	// Attester - during the whole slot
	// Aggregator - 4-12 seconds in slot
	// Proposer - during the whole slot
	// Sync committee - 4-12 seconds in slot
	// Sync committee aggregator - 4-12 seconds in slot
	panic("implement")
}
