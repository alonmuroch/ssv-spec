package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFutureDecided tests starting duty after a future decided
// This can happen if we receive a future decided message from the network and we are behind.
func PostFutureDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	futureDecide := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)

		futureDecidedInstance := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot+50))
		futureDecidedInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, futureDecidedInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot + 50)
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post future decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: futureDecide(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:   &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         futureDecide(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:           &testingutils.TestingSyncCommitteeDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  "new duty post future decided: can't start new duty runner instance for duty: could not start new QBFT instance: attempting to start an instace with a past or current height",
			},
			{
				Name:   "aggregator",
				Runner: futureDecide(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:   &testingutils.TestingAggregatorDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: futureDecide(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix)),
				Duty:   testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         futureDecide(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:           &testingutils.TestingAttesterDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  "new duty post future decided: can't start new duty runner instance for duty: could not start new QBFT instance: attempting to start an instace with a past or current height",
			},
		},
	}
}
