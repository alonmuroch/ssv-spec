package ssv

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/p2p"
	"github.com/bloxapp/ssv-spec/types"
)

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[types.BeaconRole]Runner

// DutyRunnerForMsgID returns a Runner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID types.MessageID) Runner {
	role := msgID.GetRoleType()
	return ci[role]
}

// Network is the network interface for SSV
type Network interface {
	p2p.Broadcaster
}

// AttesterCalls interface has all attester duty specific calls
type AttesterCalls interface {
	// GetAttestationData returns attestation data by the given slot and committee index
	GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (interface{}, spec.DataVersion, error)
	// SubmitAttestation submit the attestation to the node
	SubmitAttestation(attestation interface{}, version spec.DataVersion) error
}

// ProposerCalls interface has all block proposer duty specific calls
type ProposerCalls interface {
	// GetBeaconBlock returns beacon block by the given slot and committee index
	GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (interface{}, spec.DataVersion, error)
	// GetBlindedBeaconBlock returns blinded beacon block by the given slot and committee index
	GetBlindedBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (interface{}, spec.DataVersion, error)
	// SubmitBeaconBlock submit the block to the node
	SubmitBeaconBlock(block interface{}, version spec.DataVersion) error
	// SubmitBlindedBeaconBlock submit the blinded block to the node
	SubmitBlindedBeaconBlock(block interface{}, version spec.DataVersion) error
}

// AggregatorCalls interface has all attestation aggregator duty specific calls
type AggregatorCalls interface {
	// SubmitAggregateSelectionProof returns an AggregateAndProof object
	SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, index phase0.ValidatorIndex, slotSig []byte) (interface{}, spec.DataVersion, error)
	// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
	SubmitSignedAggregateSelectionProof(msg interface{}, version spec.DataVersion) error
}

// SyncCommitteeCalls interface has all sync committee duty specific calls
type SyncCommitteeCalls interface {
	// GetSyncMessageBlockRoot returns beacon block root for sync committee
	GetSyncMessageBlockRoot(slot phase0.Slot) (phase0.Root, error)
	// SubmitSyncMessage submits a signed sync committee msg
	SubmitSyncMessage(msg interface{}, version spec.DataVersion) error
}

// SyncCommitteeContributionCalls interface has all sync committee contribution duty specific calls
type SyncCommitteeContributionCalls interface {
	// IsSyncCommitteeAggregator returns true if aggregator
	IsSyncCommitteeAggregator(proof []byte) (bool, error)
	// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
	SyncCommitteeSubnetID(index phase0.CommitteeIndex) (uint64, error)
	// GetSyncCommitteeContribution returns
	GetSyncCommitteeContribution(slot phase0.Slot, subnetID uint64) (interface{}, spec.DataVersion, error)
	// SubmitSignedContributionAndProof broadcasts to the network
	SubmitSignedContributionAndProof(contribution interface{}, version spec.DataVersion) error
}

type DomainCalls interface {
	DomainData(epoch phase0.Epoch, domain phase0.DomainType) (phase0.Domain, error)
}

type BeaconNode interface {
	// GetBeaconNetwork returns the beacon network the node is on
	GetBeaconNetwork() types.BeaconNetwork
	AttesterCalls
	ProposerCalls
	AggregatorCalls
	SyncCommitteeCalls
	SyncCommitteeContributionCalls
	DomainCalls
}
