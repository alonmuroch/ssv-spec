package ssv

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type Committee struct {
	Validators       []*Validator
	CommitteeRunners [32]*CommitteeRunner
}

func NewCommittee(validators []*Validator, runners [32]*CommitteeRunner) *Committee {
	return &Committee{
		Validators:       validators,
		CommitteeRunners: runners,
	}
}

func (c *Committee) StartDuty(duty *types.CommitteeDuty) error {
	committeeDuties := map[types.BeaconRole][]*types.Duty{
		types.BNRoleAttester:      {},
		types.BNRoleSyncCommittee: {},
	}
	for _, d := range duty.BeaconDuties {
		if types.IsCommitteeDuty(d.Type) {
			committeeDuties[d.Type] = append(committeeDuties[d.Type], d)
		} else {
			if v := c.getValidatorByPubkey(d.PubKey[:]); v != nil {
				if err := v.StartDuty(d); err != nil {
					return err
				}
			} else {
				return errors.New("validator not found")
			}
		}
	}

	if err := c.CommitteeRunners[duty.Slot%32].StartNewDuty(committeeDuties[types.BNRoleAttester]); err != nil {
		return err
	}
	if err := c.CommitteeRunners[duty.Slot%32].StartNewDuty(committeeDuties[types.BNRoleSyncCommittee]); err != nil {
		return err
	}

	return nil
}

func (c *Committee) getValidatorByPubkey(pk []byte) *Validator {
	for _, v := range c.Validators {
		if bytes.Equal(pk, v.Share.ValidatorPubKey[:]) {
			return v
		}
	}
	return nil
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Decode the nested SSVMessage
	msg := &types.SSVMessage{}
	if err := msg.Decode(signedSSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode data into an SSVMessage")
	}

	id := msg.GetID()
	role := id.GetRoleType()
	pk := id.GetPubKey()
	if types.IsCommitteeDuty(role) {
		switch msg.GetType() {
		case types.SSVConsensusMsgType:
			// Decode
			signedMsg := &qbft.SignedMessage{}
			if err := signedMsg.Decode(msg.GetData()); err != nil {
				return errors.Wrap(err, "could not get consensus Message from network Message")
			}

			// Check signer consistency
			if !signedMsg.CommonSigners([]types.OperatorID{signedSSVMessage.OperatorID}) {
				return errors.New("SignedSSVMessage's signer not consistent with SignedMessage's signers")
			}

			if runner := c.CommitteeRunners[signedMsg.Message.Height%32]; runner != nil {
				return runner.ProcessConsensus(signedMsg)
			}
			return errors.New("could not find runner")
		case types.SSVPartialSignatureMsgType:
			// Decode
			signedMsg := &types.SignedPartialSignatureMessage{}
			if err := signedMsg.Decode(msg.GetData()); err != nil {
				return errors.Wrap(err, "could not get post consensus Message from network Message")
			}

			// Check signer consistency
			if signedMsg.Signer != signedSSVMessage.OperatorID {
				return errors.New("SignedSSVMessage's signer not consistent with SignedPartialSignatureMessage's signer")
			}

			// Process
			if runner := c.CommitteeRunners[signedMsg.Message.Slot%32]; runner != nil {
				if signedMsg.Message.Type == types.PostConsensusPartialSig {
					return runner.ProcessPostConsensus(signedMsg)
				}
				return runner.ProcessPreConsensus(signedMsg)
			}
			return errors.New("could not find runner")
		}
	} else { // message to validator
		if v := c.getValidatorByPubkey(pk); v != nil {
			return v.ProcessMessage(signedSSVMessage)
		}
		return errors.New("validator not found")
	}
}
