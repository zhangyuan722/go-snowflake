package snowflake

import (
	"hash/fnv"
	"math"
	"net"
	"time"
)

const (
	workerIdBits = 5

	centerIdBits = 5
	sequenceBits = 12

	maxCenterId = -1 ^ (-1 << centerIdBits)
	maxSequence = -1 ^ (-1 << sequenceBits)
)

var (
	startTimestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
)

// SnowflakeIdGenerator Snowflake id generator
type SnowflakeIdGenerator struct {
	CenterId      int64
	WorkerId      int64
	Sequence      int64
	LastTimestamp int64
}

// NewSnowflakeIdGenerator Create a new generator
func NewSnowflakeIdGenerator(centerId int64) *SnowflakeIdGenerator {
	if centerId < 0 || centerId > maxCenterId {
		panic("CenterId must be in the range [0, 31]")
	}
	return &SnowflakeIdGenerator{
		CenterId:      centerId,
		Sequence:      0,
		WorkerId:      getWorkerId(),
		LastTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

// getWorkerId Get worker unique id
func getWorkerId() int64 {
	ips, err := net.LookupHost("localhost")
	if err != nil {
		panic("failed to get ip address")
	}

	h := fnv.New64()
	for _, ip := range ips {
		_, err = h.Write([]byte(ip))
		if err != nil {
			return 1
		}
	}

	hashSum := h.Sum64()

	workerId := int64(math.Abs(float64(hashSum%31)) + 1)

	return workerId
}

// GenerateId Generate snowflake id
func (g *SnowflakeIdGenerator) GenerateId() int64 {
	currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)

	if currentTimestamp < g.LastTimestamp {
		panic("Clock moved backwards. Refusing to generate id")
	}

	// If they are generated at the same time, the sequence number will be incremented
	if currentTimestamp == g.LastTimestamp {
		g.Sequence = (g.Sequence + 1) & maxSequence
		// The sequence number overflows, waiting for the next timestamp
		if g.Sequence == 0 {
			for currentTimestamp <= g.LastTimestamp {
				currentTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		// The current timestamp is different from the previous timestamp, reset the sequence number
		g.Sequence = 0
	}

	g.LastTimestamp = currentTimestamp

	return (currentTimestamp-startTimestamp)<<(workerIdBits+centerIdBits+sequenceBits) |
		g.CenterId<<(workerIdBits+sequenceBits) |
		g.WorkerId<<sequenceBits |
		g.Sequence
}
