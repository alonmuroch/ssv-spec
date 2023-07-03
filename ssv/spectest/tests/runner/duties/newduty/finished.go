package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Finished tests a valid start duty after finished prev
func Finished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)
		r.GetBaseRunner().State.Finished = true
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name: "sync committee aggregator",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty),
				Duty: &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:           &testingutils.TestingSyncCommitteeDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:   &testingutils.TestingAggregatorDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name: "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks),
					testingutils.TestingProposerDutyV(spec.DataVersionCapella)),
				Duty: testingutils.TestingProposerDutyNextEpochV(spec.DataVersionCapella),
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionCapella), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:           &testingutils.TestingAttesterDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
