package types

import (
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type SSVNetworkChain []byte

var (
	MainnetSSVNetworkChain      SSVNetworkChain = []byte("mainnet")
	ShifuTestnetSSVNetworkChain SSVNetworkChain = []byte("shifu_testnet")
)

func (chain SSVNetworkChain) GetForksData() []ForkData {
	return []ForkData{
		{
			Epoch:             0,
			GenesisIdentifier: chain,
		},
	}
}

func (chain SSVNetworkChain) DefaultForkDigest() ForkDigest {
	return chain.GetForksData()[0].CalculateForkDigest()
}

type ForkData struct {
	// Epoch for which the fork is triggered (for all messages >= Epoch)
	Epoch spec.Epoch
	// GenesisIdentifier is a unique constant identifier per chain
	GenesisIdentifier []byte
}

func (dd ForkData) GetRoot() ([]byte, error) {
	byts, err := json.Marshal(dd)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal ForkData")
	}
	ret := sha256.Sum256(byts)
	return ret[:], nil
}

// CalculateForkDigest returns the fork digest for the fork data
func (dd ForkData) CalculateForkDigest() ForkDigest {
	r, err := dd.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	ret := ForkDigest{}
	copy(ret[:], r[:4])

	return ret
}

// ForkDigest is a 4 byte identifier for a specific fork, calculated from ForkData
type ForkDigest [4]byte
