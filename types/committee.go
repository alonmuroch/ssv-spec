package types

import spec "github.com/attestantio/go-eth2-client/spec/phase0"

type CommitteeDuty struct {
	Slot         spec.Slot
	BeaconDuties []*Duty
}

func IsCommitteeDuty(role BeaconRole) bool {
	return role == BNRoleAttester || role == BNRoleSyncCommittee
}
