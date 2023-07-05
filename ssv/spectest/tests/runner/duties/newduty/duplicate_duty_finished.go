package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateDutyFinished is a test that runs the following scenario:
// - Runner is assigned a duty
// - Runner finishes the duty
// - Runner is assigned the same duty again
// TODO - Does it make sense that the runner starts? Does it make sense that duties with precon phase don't return an error?
func DuplicateDutyFinished() tests.SpecTest {
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
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:           &testingutils.TestingSyncCommitteeDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  "can't start new duty runner instance for duty: could not start new QBFT instance: instance already running",
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:   &testingutils.TestingAggregatorDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix)),
				Duty:   testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:           &testingutils.TestingAttesterDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  "can't start new duty runner instance for duty: could not start new QBFT instance: instance already running",
			},
		},
	}
}
