package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ConsensusNotStarted tests starting duty after prev already started but for some duties' consensus didn't start because pre-consensus didnt get quorum (different duties will enable starting a new duty)
func ConsensusNotStarted() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	startRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty consensus not started",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: startRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         startRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:           &testingutils.TestingSyncCommitteeDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "aggregator",
				Runner: startRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:   &testingutils.TestingAggregatorDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:   testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         startRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:           &testingutils.TestingAttesterDutyNextEpoch,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
