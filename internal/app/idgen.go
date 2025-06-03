package app

import (
	"context"

	"github.com/sony/sonyflake"
)

type SonyflakeIDGenerator struct {
	flake *sonyflake.Sonyflake
}

func NewSonyflakeIDGenerator() *SonyflakeIDGenerator {
	return &SonyflakeIDGenerator{
		flake: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

func (g *SonyflakeIDGenerator) NextID(ctx context.Context) (int64, error) {
	id, err := g.flake.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
