package ssv

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// committeeDuty holds the tuple <validator, duty>
type committeeDuty struct {
	Validator *Validator
	Duty      *types.Duty
}

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
	committeeDuties := make([]*committeeDuty, 0)
	for _, d := range duty.BeaconDuties {
		v := c.getValidatorByPubkey(d.PubKey[:])
		if v == nil {
			return errors.New("validator not found")
		}

		if types.IsCommitteeDuty(d.Type) {
			committeeDuties = append(committeeDuties, &committeeDuty{
				Validator: v,
				Duty:      d,
			})
		} else {
			if err := v.StartDuty(d); err != nil {
				return err
			}
		}
	}

	c.CommitteeRunners[duty.Slot%32] = NewCommitteeRunner(committeeDuties)

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

			if runner := c.CommitteeRunners[signedMsg.Message.Height%32]; runner != nil {
				return runner.ProcessConsensus(signedSSVMessage)
			}
			return errors.New("could not find runner")
		case types.SSVPartialSignatureMsgType:
			// Decode
			signedMsg := &types.SignedPartialSignatureMessage{}
			if err := signedMsg.Decode(msg.GetData()); err != nil {
				return errors.Wrap(err, "could not get post consensus Message from network Message")
			}

			// Process
			if runner := c.CommitteeRunners[signedMsg.Message.Slot%32]; runner != nil {
				return runner.ProcessPostConsensus(pk, signedSSVMessage)
			}
			return errors.New("could not find runner")
		default:
			return errors.New("msg type not supported")
		}
	} else { // message to validator
		if v := c.getValidatorByPubkey(pk); v != nil {
			return v.ProcessMessage(signedSSVMessage)
		}
		return errors.New("validator not found")
	}
}
