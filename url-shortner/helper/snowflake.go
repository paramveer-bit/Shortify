package helper

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Epoch is the starting timestamp (customize as needed)
	Epoch int64 = 1672531200000 // Jan 1, 2023 in milliseconds

	// Bit allocations
	TimestampBits = 41
	MachineIDBits = 10
	SequenceBits  = 12

	// Bit shifts
	MachineIDShift = SequenceBits
	TimestampShift = SequenceBits + MachineIDBits

	// Max values
	MaxMachineID = (1 << MachineIDBits) - 1
	MaxSequence  = (1 << SequenceBits) - 1
)

// Snowflake struct
type Snowflake struct {
	machineID     int64
	lastTimestamp int64
	sequence      int64
	mutex         sync.Mutex
}

// NewSnowflake creates a new Snowflake instance
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > MaxMachineID {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", MaxMachineID)
	}
	return &Snowflake{
		machineID: machineID,
	}, nil
}

// GenerateID generates a unique ID
func (s *Snowflake) GenerateID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentTimestamp := currentMillis()

	if currentTimestamp < s.lastTimestamp {
		panic("clock moved backwards. Refusing to generate ID")
	}

	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & MaxSequence
		if s.sequence == 0 {
			// Wait for the next millisecond
			for currentTimestamp <= s.lastTimestamp {
				currentTimestamp = currentMillis()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	id := ((currentTimestamp - Epoch) << TimestampShift) |
		(s.machineID << MachineIDShift) |
		s.sequence

	return id
}

// currentMillis returns the current timestamp in milliseconds
func currentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
