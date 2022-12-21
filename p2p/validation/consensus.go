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
		if !validation.isTimelyDecidedMsg(signedMsg) {
			return pubsub.ValidationReject
		}

		if !validation.isBetterDecidedMsg(contr, signedMsg) {
			return pubsub.ValidationIgnore
		}

		if qbft.ValidateDecided(contr.GetConfig(), signedMsg, contr.Share) != nil {
			return pubsub.ValidationReject
		}
		validation.Schedualer.MarkDecidedMessage(signedMsg.Message.Identifier, signedMsg.Signers, signedMsg.Message.Height)
		return pubsub.ValidationAccept
	} else { // non-decided messages
		if !validation.isTimelyMsg(signedMsg) {
			return pubsub.ValidationReject
		}

		inst := contr.StoredInstances.FindInstance(signedMsg.Message.Height)
		// No existing instance, make basic validation on message as it's timely
		if inst == nil {
			if qbft.ValidateBaseMsg(contr.GetConfig(), signedMsg, contr.Share.Committee) != nil {
				return pubsub.ValidationReject
			}
			validation.Schedualer.MarkConsensusMessage(signedMsg.Message.Identifier, signedMsg.Signers[0], signedMsg.Message.Round, signedMsg.Message.MsgType)
			return pubsub.ValidationAccept
		}

		isDecided, _ := inst.IsDecided()
		// If instance is decided we only accept aggregatable commit messages, other messages are not useful
		if isDecided {
			if qbft.BaseCommitValidation(contr.GetConfig(), signedMsg, signedMsg.Message.Height, contr.Share.Committee) != nil {
				return pubsub.ValidationReject
			}
			if !validation.isCommitMsgAggregatable(inst, signedMsg) {
				return pubsub.ValidationReject
			}
			validation.Schedualer.MarkConsensusMessage(signedMsg.Message.Identifier, signedMsg.Signers[0], signedMsg.Message.Round, signedMsg.Message.MsgType)
			return pubsub.ValidationAccept
		}

		// If instance exists, make a full stateful validation
		if inst.BaseMsgValidation(signedMsg) != nil {
			return pubsub.ValidationReject
		}
		validation.Schedualer.MarkConsensusMessage(signedMsg.Message.Identifier, signedMsg.Signers[0], signedMsg.Message.Round, signedMsg.Message.MsgType)
		return pubsub.ValidationAccept
	}
}

// isCommitMsgAggregatable will return true if the signed message can be aggregated to the decided message
func (validation *MessageValidation) isCommitMsgAggregatable(inst *qbft.Instance, msg *qbft.SignedMessage) bool {
	panic("implement")
}

// isTimelyMsg returns if future message is timely
func (validation *MessageValidation) isTimelyMsg(msg *qbft.SignedMessage) bool {
	return validation.Schedualer.IsConsensusMessageTimely(
		msg.Message.Identifier,
		msg.Signers[0],
		msg.Message.Round,
		msg.Message.MsgType,
	)
}

// isTimelyDecidedMsg returns true if decided message is timely (both for future and past decided messages)
// FUTURE: when a valid decided msg is received, the next duty for the runner is marked. The next decided message will not be validated before that time.
// Aims to prevent a byzantine committee rapidly broadcasting decided messages
// PAST: a decided message which is "too" old will be rejected as well
func (validation *MessageValidation) isTimelyDecidedMsg(msg *qbft.SignedMessage) bool {
	return validation.Schedualer.IsTimelyDecidedMessage(
		msg.Message.Identifier,
		msg.Signers,
		msg.Message.Height,
	)
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
