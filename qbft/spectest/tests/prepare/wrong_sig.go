package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSignature tests a single prepare received with a wrong signature
func WrongSignature() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(2)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare wrong sig",
		Pre:           pre,
		PostRoot:      "69c049da1936e3727d09f976754cc7ee3a5cb7d85fa1e079f0465096b0a15cdb",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg signature invalid: failed to verify signature",
	}
}
