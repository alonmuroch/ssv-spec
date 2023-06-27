package proposer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ProposeRegularBlockDecidedBlinded tests proposing a regular block but the decided block is a blinded block. Full flow
func ProposeRegularBlockDecidedBlinded() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MsgProcessingSpecTest{
		Name:   "propose regular decide blinded",
		Runner: testingutils.ProposerRunner(ks),
		Duty:   testingutils.TestingProposerDutyV(spec.DataVersionCapella),
		Messages: []*types.SSVMessage{
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionCapella)),
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionCapella)),
			testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionCapella)),

			testingutils.SSVMsgProposer(
				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*bls.SecretKey{
						ks.Shares[1], ks.Shares[2], ks.Shares[3],
					},
					[]types.OperatorID{1, 2, 3},
					qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionCapella)),
					testingutils.ProposerMsgID,
					testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionCapella),
				), nil),

			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionCapella)),
			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionCapella)),
			testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionCapella)),
		},
		PostDutyRunnerStateRoot: "8069eb6e56029dbda522ae465f738a0cd09d5876c9ff6735d751010ad715e29c",
		OutputMessages: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionCapella),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionCapella),
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, spec.DataVersionCapella)),
		},
	}
}
