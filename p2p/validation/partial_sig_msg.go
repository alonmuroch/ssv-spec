package validation

import (
	"github.com/bloxapp/ssv-spec/ssv"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (validation *MessageValidation) validatePartialSigMsg() pubsub.ValidationResult {
	signedMsg, err := validation.decodePartialSignatureMessage(nil)
	if err != nil {
		return pubsub.ValidationReject
	}

	if signedMsg.Validate() != nil {
		return pubsub.ValidationReject
	}

	if validation.verifySignatureForShare(signedMsg, nil) != nil {
		return pubsub.ValidationReject
	}

	if !validation.partialSigInTime(signedMsg) {
		return pubsub.ValidationReject
	}

	if !validation.isUniquePartialSigMessageForSigner(signedMsg) {
		return pubsub.ValidationReject
	}

	for _, msg := range signedMsg.Message.Messages {
		if validation.validBeaconSignature(msg) != nil {
			return pubsub.ValidationReject
		}
	}

	return pubsub.ValidationAccept
}

// isUniquePartialSigMessageForSigner returns true if the message is a unique first time pre/post-consensus msg for signer (for epoch)
func (validation *MessageValidation) isUniquePartialSigMessageForSigner(msg *ssv.SignedPartialSignatureMessage) bool {
	panic("implement")
}

// validBeaconSignature returns nil if beacon signature is valid
func (validation *MessageValidation) validBeaconSignature(message *ssv.PartialSignatureMessage) error {
	panic("implement")
}

// partialSigInTime returns true if pre/post - consensus message is received in time
func (validation *MessageValidation) partialSigInTime(msg *ssv.SignedPartialSignatureMessage) bool {
	// Test msg in epoch
	panic("implement")
}
