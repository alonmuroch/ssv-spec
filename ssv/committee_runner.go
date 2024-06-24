package ssv

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CommitteeRunner struct {
}

func NewCommitteeRunner() *CommitteeRunner {
	return &CommitteeRunner{}
}

func (cr CommitteeRunner) StartNewDuty(duties []*types.Duty) error {

}

// ProcessPreConsensus processes all pre-consensus msgs, returns error if can't process
func (cr CommitteeRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {

}

// ProcessConsensus processes all consensus msgs, returns error if can't process
func (cr CommitteeRunner) ProcessConsensus(msg *qbft.SignedMessage) error {

}

// ProcessPostConsensus processes all post-consensus msgs, returns error if can't process
func (cr CommitteeRunner) ProcessPostConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	
}
