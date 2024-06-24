package ssv

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// CommitteeRunner holds the validators for a specific committee duties run and their respective beacon duties
type CommitteeRunner struct {
	Validators []*Validator
	Duties     []*types.Duty

	LeadingAttesterValidator      *committeeDuty
	LeadingSyncCommitteeValidator *committeeDuty
}

func NewCommitteeRunner(cd []*committeeDuty) *CommitteeRunner {
	ret := &CommitteeRunner{
		Validators: make([]*Validator, 0),
		Duties:     make([]*types.Duty, 0),
	}

	for _, d := range cd {
		ret.Validators = append(ret.Validators, d.Validator)
		ret.Duties = append(ret.Duties, d.Duty)

		if d.Duty.Type == types.BNRoleAttester && ret.LeadingAttesterValidator == nil {
			ret.LeadingSyncCommitteeValidator = d
		}
		if d.Duty.Type == types.BNRoleSyncCommittee && ret.LeadingSyncCommitteeValidator == nil {
			ret.LeadingSyncCommitteeValidator = d
		}
	}
	return ret
}

func (cr CommitteeRunner) Start() error {
	if err := cr.LeadingAttesterValidator.Validator.StartDuty(cr.LeadingAttesterValidator.Duty); err != nil {
		return err
	}
	return cr.LeadingSyncCommitteeValidator.Validator.StartDuty(cr.LeadingSyncCommitteeValidator.Duty)
}

//func (cr CommitteeRunner) StartNewDuty(duties []*types.Duty, validators []*Validator) error {
//
//}

// ProcessPreConsensus processes all pre-consensus msgs, returns error if can't process
func (cr CommitteeRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	return errors.New("no pre consensus phase for committee runner")
}

// ProcessConsensus processes all consensus msgs, returns error if can't process
func (cr CommitteeRunner) ProcessConsensus(signedSSVMessage *types.SignedSSVMessage) error {
	// Decode the nested SSVMessage
	msg := &types.SSVMessage{}
	if err := msg.Decode(signedSSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode data into an SSVMessage")
	}

	id := msg.GetID()
	role := id.GetRoleType()

	switch role {
	case types.BNRoleAttester:
		if err := cr.LeadingAttesterValidator.Validator.ProcessMessage(signedSSVMessage); err != nil {
			return err
		}

		runner := cr.LeadingAttesterValidator.Validator.DutyRunners[types.BNRoleAttester]
		if runner.GetBaseRunner().State.DecidedValue != nil {
			for i, d := range cr.Duties {
				if d.Type != types.BNRoleAttester {
					continue
				}

				runner2 := cr.Validators[i].DutyRunners[types.BNRoleAttester]
				baseRunner2 := runner2.GetBaseRunner()
				baseRunner2.State.DecidedValue = runner.GetBaseRunner().State.DecidedValue
				baseRunner2.highestDecidedSlot = baseRunner2.State.DecidedValue.Duty.Slot

				return runner2.UponDecided()
			}
		}
	case types.BNRoleSyncCommittee:
		if err := cr.LeadingSyncCommitteeValidator.Validator.ProcessMessage(signedSSVMessage); err != nil {
			return err
		}
		// TODO simulate decided
	default:
		return errors.New("role not supported")
	}

}

// ProcessPostConsensus processes all post-consensus msgs, returns error if can't process
func (cr CommitteeRunner) ProcessPostConsensus(pk []byte, signedSSVMessage *types.SignedSSVMessage) error {
	v := cr.getValidatorByPubkey(pk)
	if v == nil {
		return errors.New("validator not found")
	}
	return v.ProcessMessage(signedSSVMessage)
}

func (cr CommitteeRunner) getValidatorByPubkey(pk []byte) *Validator {
	for _, v := range cr.Validators {
		if bytes.Equal(pk, v.Share.ValidatorPubKey[:]) {
			return v
		}
	}
	return nil
}
