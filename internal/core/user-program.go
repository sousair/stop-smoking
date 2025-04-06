package core

import "github.com/sousair/gocore/pkg/database/entity"

type UserProgram struct {
	entity.BaseEntity

	UserID         uint  `json:"user_id"`
	UserIdentifier int64 `json:"user_identifier"`
	ProgramID      uint  `json:"program_id"`

	Smoked  int64 `json:"smoked_count"`
	Skipped int64 `json:"skipped_count"`
	Failed  int64 `json:"failed_count"`

	HasFailed     bool  `json:"has_failed"`
	LastSmokeTime int64 `json:"last_smoke_time"`
}
