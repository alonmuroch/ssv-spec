package validation

import (
	"github.com/bloxapp/ssv-spec/qbft"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (validation *MessageValidation) validateConsensusMsg() pubsub.ValidationResult {
	signedMsg, err := validation.decodeConsensusMessage(nil)
	if err != nil {
		return pubsub.ValidationReject
	}

	if signedMsg.Validate() != nil {
		return pubsub.ValidationReject
	}

	contr := validation.findController()
	if contr == nil {
		return pubsub.ValidationReject
	}

	/**
	Main controller processing flow
	_______________________________
	All decided msgs are processed the same, out of instance
	All valid future msgs are saved in a container and can trigger the highest decided future msg
	All other msgs (not future or decided) are processed normally by an existing instance (if found)
	*/
	if qbft.IsDecidedMsg(contr.Share, signedMsg) {
		if qbft.ValidateDecided(contr.GetConfig(), signedMsg, contr.Share) != nil {
			return pubsub.ValidationReject
		}
		if !validation.isTimelyDecidedMsg() {
			return pubsub.ValidationReject
		}
		if !validation.isBetterDecidedMsg(contr, signedMsg) {
			return pubsub.ValidationIgnore
		}
	} else if signedMsg.Message.Height > contr.Height {
		if qbft.ValidateFutureMsg(contr.GetConfig(), signedMsg, contr.Share.Committee) != nil {
			return pubsub.ValidationReject
		}
		if !validation.isTimelyAndUniqueFutureMsg() {
			return pubsub.ValidationIgnore
		}
	} else {
		inst := contr.StoredInstances.FindInstance(signedMsg.Message.Height)
		if inst == nil {
			return pubsub.ValidationIgnore
		}
		isDecided, _ := inst.IsDecided()
		if isDecided {
			if !validation.isCommitMsgAggregatable(inst, signedMsg) {
				return pubsub.ValidationReject
			}
		} else if inst.BaseMsgValidation(signedMsg) != nil {
			return pubsub.ValidationReject
		}
	}

	return pubsub.ValidationAccept
}

// isCommitMsgAggregatable will return true if the signed message can be aggregated to the decided message
func (validation *MessageValidation) isCommitMsgAggregatable(inst *qbft.Instance, msg *qbft.SignedMessage) bool {
	panic("implement")
}

// isTimelyAndUniqueFutureMsg returns if future message is timely and unique
// A timely and unique message is a unique message for height-round-signer, round is increment and the message is timely (round 4 message comes X seconds after round 3 message)
func (validation *MessageValidation) isTimelyAndUniqueFutureMsg() bool {
	panic("implement")
}

// isTimelyDecidedMsg returns true if decided message is timely (both for future and past decided messages)
// FUTURE: when a valid decided msg is received, the next duty for the runner is marked. The next decided message will not be validated before that time.
// Aims to prevent a byzantine committee rapidly broadcasting decided messages
// PAST: a decided message which is "too" old will be rejected as well
func (validation *MessageValidation) isTimelyDecidedMsg() bool {
	panic("implement")
}

// isBetterDecidedMsg returns true if the decided message is better than the best known decided
func (validation *MessageValidation) isBetterDecidedMsg(contr *qbft.Controller, msg *qbft.SignedMessage) bool {
	panic("implement")
}

func (validation *MessageValidation) findController() *qbft.Controller {
	panic("implement")
}

// decodePartialSignatureMessage returns the decoded signed message or error
func (validation *MessageValidation) decodeConsensusMessage(data []byte) (*qbft.SignedMessage, error) {
	panic("implement")
}
