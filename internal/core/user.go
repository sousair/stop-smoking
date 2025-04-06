package core

import (
	"encoding/json"
	"time"

	"github.com/sousair/gocore/pkg/database/entity"
)

type User struct {
	entity.BaseEntity

	Identifier int64  `json:"identifier"`
	Name       string `json:"name"`

	StartDate       time.Time `json:"start_date"`
	QuitDate        time.Time `json:"quit_date"`
	TotalCigarettes int       `json:"total_cigarettes"`

	Metadata json.RawMessage `json:"metadata"`
}

type UserMetadata struct {
	WakeUpHour             int  `json:"wake_up_hour"`
	SleepHour              int  `json:"sleep_hour"`
	SmokeToSleep           bool `json:"smoke_to_sleep"`
	RolledOverSleepRoutine bool `json:"round_up_routine"`
}

type (
	UserMetadataSentRequest struct {
		Identifier int64 `json:"identifier"`

		CigarettesPerDay int  `json:"cigarettes_per_day"`
		WakeUpHour       int  `json:"wake_up_hour"`
		SleepHour        int  `json:"sleep_hour"`
		SmokeToSleep     bool `json:"smoke_to_sleep"`
	}

	StartProgramScheduleRequest struct {
		Identifier       int64 `json:"identifier"`
		CigarettesPerDay int   `json:"cigarettes_per_day"`
	}
)
