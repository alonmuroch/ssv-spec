package validation

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMessageSchedule_MarkConsensusMessage(t *testing.T) {
	s := NewMessageSchedule()
	s.MarkConsensusMessage([]byte{1, 2, 3, 4}, 0, qbft.FirstRound, qbft.PrepareMsgType)
	require.False(t, s.IsConsensusMessageTimely([]byte{1, 2, 3, 4}, 0, qbft.FirstRound, qbft.PrepareMsgType))
	require.True(t, s.IsConsensusMessageTimely([]byte{1, 2, 3, 4}, 0, qbft.FirstRound, qbft.CommitMsgType))

	require.False(t, s.IsConsensusMessageTimely([]byte{1, 2, 3, 4}, 0, 2, qbft.PrepareMsgType))
	<-time.After(time.Second * 2)
	require.True(t, s.IsConsensusMessageTimely([]byte{1, 2, 3, 4}, 0, 2, qbft.PrepareMsgType))

	s.MarkDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 1)
	require.True(t, s.IsConsensusMessageTimely([]byte{1, 2, 3, 4}, 0, qbft.FirstRound, qbft.PrepareMsgType))
}

func TestMessageSchedule_MarkDecidedMessage(t *testing.T) {
	s := NewMessageSchedule()
	require.True(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 1))
	s.MarkDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 1)
	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 1))

	require.True(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 2))
	s.MarkDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 2)
	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 2))

	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 3))
	<-time.After(time.Second * 2)
	require.True(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 3))
	s.MarkDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 3)
	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 3))

	require.True(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 4))
	s.MarkDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 4)
	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 4))

	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 5))
	<-time.After(time.Second * 2)
	require.True(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 5))
	s.MarkDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 5)
	require.False(t, s.IsTimelyDecidedMessage([]byte{1, 2, 3, 4}, []types.OperatorID{0, 1, 2}, 5))
}
