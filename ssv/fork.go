package ssv

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// forkBasedOnLatestDecided will return fork digest based on b.Height instance that was previously decided
func (b *BaseRunner) forkBasedOnLatestDecided() (types.ForkDigest, error) {
	inst := b.QBFTController.InstanceForHeight(b.QBFTController.Height)
	if inst == nil {
		return b.SSVNetworkChain.DefaultForkDigest(), nil
	}

	_, decidedValue := inst.IsDecided()
	cd := &types.ConsensusData{}
	if err := cd.Decode(decidedValue); err != nil {
		return types.ForkDigest{}, errors.Wrap(err, "could not decoded consensus data")
	}

	currentForkDigest := b.SSVNetworkChain.DefaultForkDigest()
	for _, forkData := range b.SSVNetworkChain.GetForksData() {
		if b.BeaconNetwork.EstimatedEpochAtSlot(cd.Duty.Slot) >= forkData.Epoch {
			currentForkDigest = forkData.CalculateForkDigest()
		}
	}
	return currentForkDigest, nil
}
