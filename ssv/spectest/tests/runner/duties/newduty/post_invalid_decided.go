package newduty

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostInvalidDecided tests starting a new duty after prev was decided with an invalid decided value
func PostInvalidDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	consensusDataByts := func(role types.BeaconRole) []byte {
		cd := &types.ConsensusData{
			Duty: types.Duty{
				Type:                    100, // invalid
				PubKey:                  testingutils.TestingValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
		}
		byts, _ := cd.Encode()
		return byts
	}

	// TODO: not sure what is the value of doing this. We initialize the runner with an impossible decided value.
	// Maybe we should ensure that `ValidateDecided()` doesn't let the runner enter this state and delete the test?
	decideWrong := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)

		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		decidedValue := sha256.Sum256(consensusDataByts(r.GetBaseRunner().BeaconRoleType))
		r.GetBaseRunner().State.RunningInstance.State.DecidedValue = decidedValue[:]

		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post invalid decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name: "sync committee aggregator",
				Runner: decideWrong(testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty),
				Duty: &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         decideWrong(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:           &testingutils.TestingSyncCommitteeDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "aggregator",
				Runner: decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:   &testingutils.TestingAggregatorDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:   testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         decideWrong(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:           &testingutils.TestingAttesterDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
