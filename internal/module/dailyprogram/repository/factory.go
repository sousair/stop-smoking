package repository

import (
	"context"

	"github.com/sousair/gocore/pkg/database/repository"
	"github.com/sousair/stop-smoking/internal/core"
	"gorm.io/gorm"
)

type DailyProgramRepository struct {
	db *gorm.DB
	repository.Repository[core.DailyProgram]
}

var _ core.DailyProgramRepository = (*DailyProgramRepository)(nil)

func NewDailyProgramRepository(db *gorm.DB) (*DailyProgramRepository, error) {
	repository, err := repository.NewRepository[core.DailyProgram](db)
	if err != nil {
		return nil, err
	}

	return &DailyProgramRepository{db, repository}, nil
}

func (r DailyProgramRepository) FindProgramByPerDay(
	ctx context.Context,
	cigarretesPerDay int,
) (*core.DailyProgram, error) {
	tx := r.db

	if dbTx, err := repository.FromContext(ctx); err == nil {
		tx = dbTx
	}

	query := &core.DailyProgram{
		CigarretesPerDay: cigarretesPerDay,
	}

	if err := tx.Where(query).First(query).Error; err != nil {
		return nil, err
	}

	return query, nil
}
