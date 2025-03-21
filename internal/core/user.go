package core

import (
	"encoding/json"
	"time"

	"github.com/sousair/gocore/pkg/database/entity"
)

type User struct {
	entity.BaseEntity

	Idendifier int64  `json:"idendifier"`
	Name       string `json:"name"`

	StartDate       time.Time `json:"start_date"`
	QuitDate        time.Time `json:"quit_date"`
	LastSmokeDate   time.Time `json:"last_smoke_date"`
	TotalCigarettes int       `json:"total_cigarettes"`

	ProgramID uint `json:"program_id"`

	Metadata json.RawMessage `json:"metadata"`
}

type UserMetadata struct {
	WakeUpHour   int  `json:"wake_up_hour"`
	SleepHour    int  `json:"sleep_hour"`
	SmokeToSleep bool `json:"smoke_to_sleep"`
}

type (
	StartProgramRequest struct {
		Identifier int64 `json:"identifier"`

		CigarettesPerDay int  `json:"cigarettes_per_day"`
		WakeUpHour       int  `json:"wake_up_hour"`
		SleepHour        int  `json:"sleep_hour"`
		SmokeToSleep     bool `json:"smoke_to_sleep"`
	}
)
