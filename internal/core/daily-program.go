package core

import (
	"context"

	"github.com/sousair/gocore/pkg/database/entity"
	"github.com/sousair/gocore/pkg/database/repository"
)

type DailyProgram struct {
	entity.BaseEntity

	Objective        string `json:"objective"`
	CigarretesPerDay int    `json:"cigarretes_per_day"`
}

type DailyProgramRepository interface {
	repository.Repository[DailyProgram]
	FindProgramByPerDay(ctx context.Context, cigarretesPerDay int) (*DailyProgram, error)
}

type (
	FailedDailyProgram struct {
		Identifier string `json:"user_identifier"`
	}
)
