package snowflake

import (
	"testing"
)

func TestGenerateSnowflakeId(t *testing.T) {
	snowflake := NewSnowflakeIdGenerator(1)
	uniqueIds := make(map[int64]struct{})
	total := 1000000

	for i := 0; i < total; i++ {
		id := snowflake.GenerateId()
		uniqueIds[id] = struct{}{}
	}

	if len(uniqueIds) != total {
		t.Errorf("The same Id is generated")
	}
}
