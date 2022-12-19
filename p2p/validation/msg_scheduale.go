package validation

import (
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"time"
)

const (
	MaxDutyTypePerEpoch = 1
	Slot                = time.Second * 12
	Epoch               = Slot * 32
)

// RoundTimeout returns timeout duration for round
var RoundTimeout = func(round qbft.Round) time.Duration {
	return time.Second * 2
}

type mark struct {
	HighestRound    qbft.Round
	FirstMsgInRound time.Time
	MsgTypesInRound map[qbft.MessageType]bool

	Last3DecidedTimes [3]time.Time
	HighestDecided    qbft.Height
	MarkedDecided     int
}

func markID(id []byte, signer types.OperatorID) string {
	return fmt.Sprintf("%x%d", id, signer)
}

// DurationForLast3Decided returns duration for the last 3 decided including a new decided added at time.Now()
func (mark *mark) DurationForLast3Decided() time.Duration {
	return time.Now().Sub(mark.Last3DecidedTimes[1])
}

func (mark *mark) AddDecidedMark() {
	copy(mark.Last3DecidedTimes[1:], mark.Last3DecidedTimes[:1]) // shift array one up to make room at index 0
	mark.Last3DecidedTimes[0] = time.Now()

	mark.MarkedDecided++
}

func (mark *mark) ResetForNewRound(round qbft.Round, msgType qbft.MessageType) {
	mark.HighestRound = round
	mark.FirstMsgInRound = time.Now()
	mark.MsgTypesInRound = map[qbft.MessageType]bool{
		msgType: true,
	}
}

// MessageSchedule keeps track of consensus msg schedules to determine timely receiving msgs
type MessageSchedule struct {
	Marks map[string]*mark
}

func NewMessageSchedule() *MessageSchedule {
	return &MessageSchedule{
		Marks: map[string]*mark{},
	}
}

// MarkConsensusMessage marks a msg
func (schedule *MessageSchedule) MarkConsensusMessage(id []byte, signer types.OperatorID, round qbft.Round, msgType qbft.MessageType) {
	idStr := markID(id, signer)
	var signerMark *mark
	signerMark, found := schedule.Marks[idStr]
	if !found {
		signerMark = &mark{}
		schedule.Marks[idStr] = signerMark
	}

	if signerMark.HighestRound < round {
		signerMark.ResetForNewRound(round, msgType)
	} else {
		signerMark.MsgTypesInRound[msgType] = true
	}
}

func (schedule *MessageSchedule) IsConsensusMessageTimely(id []byte, signer types.OperatorID, round qbft.Round, msgType qbft.MessageType) bool {
	idStr := markID(id, signer)
	var signerMark *mark
	signerMark, found := schedule.Marks[idStr]
	if !found {
		return true
	}

	if signerMark.HighestRound < round {
		// if new round msg, check at least round timeout has passed
		return time.Now().After(signerMark.FirstMsgInRound.Add(RoundTimeout(signerMark.HighestRound)))
	} else if signerMark.HighestRound == round {
		_, found := signerMark.MsgTypesInRound[msgType]
		return !found
	} else {
		// past rounds are not timely
		return false
	}
}

func (schedule *MessageSchedule) MarkDecidedMessage(id []byte, signers []types.OperatorID, height qbft.Height) {
	for _, signer := range signers {
		idStr := markID(id, signer)
		var signerMark *mark
		signerMark, found := schedule.Marks[idStr]
		if !found {
			signerMark = &mark{}
			schedule.Marks[idStr] = signerMark
		}

		if signerMark.HighestDecided < height {
			signerMark.HighestDecided = height
			signerMark.AddDecidedMark()
			signerMark.HighestRound = qbft.NoRound
		}
	}
}

func (schedule *MessageSchedule) IsTimelyDecidedMessage(id []byte, signers []types.OperatorID, height qbft.Height) bool {
	ret := false
	for _, signer := range signers {
		idStr := markID(id, signer)
		var signerMark *mark
		signerMark, found := schedule.Marks[idStr]
		if !found {
			ret = true
			continue
		}

		if signerMark.HighestDecided >= height {
			return false
		}

		if signerMark.MarkedDecided >= 2 {
			return signerMark.DurationForLast3Decided() >= time.Second*2
		} else {
			ret = true
		}
	}
	return ret
}
