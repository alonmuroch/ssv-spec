package validation

import (
	"context"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type MessageValidation struct {
	Schedualer *MessageSchedule
}

func (validation *MessageValidation) Validate(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	ssvMsg := validation.decodedPubsubMsg(msg.GetData())
	if ssvMsg == nil {
		return pubsub.ValidationReject
	}
	if !validation.validateSSVMessage(ssvMsg) {
		return pubsub.ValidationReject
	}

	switch ssvMsg.GetType() {
	case types.SSVConsensusMsgType:
		return validation.validateConsensusMsg()
	case types.SSVPartialSignatureMsgType:
		return validation.validatePartialSigMsg()
	default:
		return pubsub.ValidationReject
	}
}

func (validation *MessageValidation) decodedPubsubMsg(data []byte) *types.SSVMessage {
	panic("implement")
}

func (validation *MessageValidation) validateSSVMessage(msg *types.SSVMessage) bool {
	// find validator and runner
	// msg.validate()
	// validate active validator
	panic("implement")
}

// verifySignatureForShare returns nil if message is signed correctly by someone from the share's committee
func (validation *MessageValidation) verifySignatureForShare(msg types.MessageSignature, share *types.Share) error {
	panic("implement")
}

// decodePartialSignatureMessage returns the decoded SignedPartialSignatureMessage or error
func (validation *MessageValidation) decodePartialSignatureMessage(data []byte) (*ssv.SignedPartialSignatureMessage, error) {
	panic("implement")
}

//// MsgValidatorFunc represents a message validator
//type MsgValidatorFunc = func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult
//
//func MsgValidation(runner ssv.Runner) MsgValidatorFunc {
//	return func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
//		ssvMsg, err := DecodePubsubMsg(msg)
//		if err != nil {
//			return pubsub.ValidationReject
//		}
//		if validateSSVMessage(runner, ssvMsg) != nil {
//			return pubsub.ValidationReject
//		}
//
//		switch ssvMsg.GetType() {
//		case types.SSVConsensusMsgType:
//			if validateConsensusMsg(runner, ssvMsg.Data) != nil {
//				return pubsub.ValidationReject
//			}
//		case types.SSVPartialSignatureMsgType:
//			if validatePartialSigMsg(runner, ssvMsg.Data) != nil {
//				return pubsub.ValidationReject
//			}
//		default:
//			return pubsub.ValidationReject
//		}
//
//		return pubsub.ValidationAccept
//	}
//}
//
//func DecodePubsubMsg(msg *pubsub.Message) (*types.SSVMessage, error) {
//	byts := msg.GetData()
//	ret := &types.SSVMessage{}
//	if err := ret.Decode(byts); err != nil {
//		return nil, err
//	}
//	return ret, nil
//}
//
//func validateSSVMessage(runner ssv.Runner, msg *types.SSVMessage) error {
//	if !runner.GetBaseRunner().Share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
//		return errors.New("msg ID doesn't match validator ID")
//	}
//
//	if len(msg.GetData()) == 0 {
//		return errors.New("msg data is invalid")
//	}
//
//	return nil
//}
//
//func validateConsensusMsg(runner ssv.Runner, data []byte) error {
//	signedMsg := &qbft.SignedMessage{}
//	if err := signedMsg.Decode(data); err != nil {
//		return err
//	}
//
//	contr := runner.GetBaseRunner().QBFTController
//
//	if err := contr.BaseMsgValidation(signedMsg); err != nil {
//		return err
//	}
//
//	/**
//	Main controller processing flow
//	_______________________________
//	All decided msgs are processed the same, out of instance
//	All valid future msgs are saved in a container and can trigger highest decided futuremsg
//	All other msgs (not future or decided) are processed normally by an existing instance (if found)
//	*/
//	if qbft.IsDecidedMsg(contr.Share, signedMsg) {
//		return qbft.ValidateDecided(contr.GetConfig(), signedMsg, contr.Share)
//	} else if signedMsg.Message.Height > contr.Height {
//		return qbft.ValidateFutureMsg(contr.GetConfig(), signedMsg, contr.Share.Committee)
//	} else {
//		if inst := contr.StoredInstances.FindInstance(signedMsg.Message.Height); inst != nil {
//			return inst.BaseMsgValidation(signedMsg)
//		}
//		return errors.New("unknown instance")
//	}
//}
//
//func validatePartialSigMsg(runner ssv.Runner, data []byte) error {
//	signedMsg := &ssv.SignedPartialSignatureMessage{}
//	if err := signedMsg.Decode(data); err != nil {
//		return err
//	}
//
//	if signedMsg.Message.Type == ssv.PostConsensusPartialSig {
//		return runner.GetBaseRunner().ValidatePostConsensusMsg(runner, signedMsg)
//	}
//	return runner.GetBaseRunner().ValidatePreConsensusMsg(runner, signedMsg)
//}
