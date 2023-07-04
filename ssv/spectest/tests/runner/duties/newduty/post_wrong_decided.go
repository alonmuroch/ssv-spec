package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostWrongDecided tests starting a new duty after prev was decided wrongly (future decided)
// This can happen if we receive a future decided message from the network.
func PostWrongDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: not sure what is the value of doing this. We initialize the runner with an impossible decided value.
	// Maybe we should ensure that `ValidateDecided()` doesn't let the runner enter this state and delete the test?
	decideWrong := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		storedInstances := r.GetBaseRunner().QBFTController.StoredInstances
		storedInstances = append(storedInstances, nil)
		storedInstances = append(storedInstances, nil)

		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		storedInstances[1] = r.GetBaseRunner().State.RunningInstance

		higherDecided := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			50)
		higherDecided.State.Decided = true
		higherDecided.State.DecidedValue = []byte{1, 2, 3, 4}
		storedInstances[0] = higherDecided
		r.GetBaseRunner().QBFTController.Height = 50
		// TODO: hacky fix to a bug in the test.
		// You can't append a copied slice and expect the original to change in go. Since maybe we want to delete
		// the test I didn't do it nicer.
		r.GetBaseRunner().QBFTController.StoredInstances = storedInstances
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post wrong decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: decideWrong(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "sync committee",
				Runner:         decideWrong(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:           &testingutils.TestingSyncCommitteeDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  "can't start new duty runner instance for duty: could not start new QBFT instance: attempting to start an instace with a past or current height",
			},
			{
				Name:   "aggregator",
				Runner: decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:   &testingutils.TestingAggregatorDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "proposer",
				Runner: decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionCapella)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionCapella),
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionCapella), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         decideWrong(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:           &testingutils.TestingAttesterDuty,
				OutputMessages: []*types.SignedPartialSignatureMessage{},
				ExpectedError:  "can't start new duty runner instance for duty: could not start new QBFT instance: attempting to start an instace with a past or current height",
			},
		},
	}
}
