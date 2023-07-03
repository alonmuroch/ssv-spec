package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateDutyFinished is a test that runs the following scenario:
// - Runner is assigned a duty
// - Runner finishes the duty
// - Runner is assigned the same duty again
// TODO - Does it make sense that the runner starts? Does it make sense that duties with precon phase don't return an error?
func DuplicateDutyFinished() *MultiStartNewRunnerDutySpecTest {
	ks := testingutils.Testing4SharesSet()

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
		Name: "duplicate duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:   &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
			},
			{
				Name:          "sync committee",
				Runner:        finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:          &testingutils.TestingSyncCommitteeDuty,
				ExpectedError: "can't start new duty runner instance for duty: could not start new QBFT instance: instance already running",
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:   &testingutils.TestingAggregatorDutyNextEpoch,
			},
			{
				Name:   "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix)),
				Duty:   testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
			},
			{
				Name:          "attester",
				Runner:        finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:          &testingutils.TestingAttesterDuty,
				ExpectedError: "can't start new duty runner instance for duty: could not start new QBFT instance: instance already running",
			},
		},
	}
}
