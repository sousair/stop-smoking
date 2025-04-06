package core

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/sousair/gocore/pkg/database/entity"
	"github.com/sousair/gocore/pkg/database/repository"
)

type DailyProgram struct {
	entity.BaseEntity

	Objective        string `json:"objective"`
	CigarretesPerDay int    `json:"cigarretes_per_day"`
	JumpOffset       uint   `json:"jump_offset"`
}

type DailyProgramRepository interface {
	repository.Repository[DailyProgram]
	FindProgramByPerDay(ctx context.Context, cigarretesPerDay int) (*DailyProgram, error)
}

type SmokedRequest struct {
	Identifier int64 `json:"identifier"`
	Failed     bool  `json:"failed"`
}

type ScheduleSmokingPayload struct {
	Identifier int64     `json:"identifier"`
	SmokeTime  time.Time `json:"smoke_time"`
}

const (
	SleepSmokeIntervalInMinutes     = 30
	NeedSleepSmokeIntervalInMinutes = 10
	AwakeSmokeIntervalInMinutes     = 30
)

var NoSmokingTimeAvailableErr = errors.New("No smoking time available")

type SmokingSchedule struct {
	Identifier          int64         `json:"identifier"`
	UserProgramID       uint          `json:"user_program_id"`
	NextSmoke           time.Time     `json:"next_smoke"`
	IntervalInMs        time.Duration `json:"interval_in_ms"`
	RemainingCigarretes int           `json:"remaining_cigarretes"`
}

func (s *SmokingSchedule) Next() error {
	if s.RemainingCigarretes == 0 {
		return NoSmokingTimeAvailableErr
	}

	s.NextSmoke = s.NextSmoke.Add(s.IntervalInMs)
	s.RemainingCigarretes--

	log.Printf("Next smoke time for user [%d] is [%s]", s.Identifier, s.NextSmoke)

	return nil
}

func (dp DailyProgram) CalculateUserSchedule(user *User, userProgramID uint, referenceTime time.Time) (*SmokingSchedule, error) {
	var metadata *UserMetadata
	if err := json.Unmarshal(user.Metadata, &metadata); err != nil {
		return nil, err
	}

	sleepMin := metadata.SleepHour * 60
	wakeUpMin := metadata.WakeUpHour * 60

	if sleepMin < wakeUpMin {
		sleepMin += 24 * 60
	}

	awakeDuration := sleepMin - wakeUpMin
	availableSmokeTime := awakeDuration - AwakeSmokeIntervalInMinutes

	if metadata.SmokeToSleep {
		availableSmokeTime -= NeedSleepSmokeIntervalInMinutes
	} else {
		availableSmokeTime -= SleepSmokeIntervalInMinutes
	}

	remainingCigarretes := dp.CigarretesPerDay - 1
	smokingIntervalInMs :=
		time.Duration((availableSmokeTime*60*1000)/remainingCigarretes) *
			time.Millisecond

	firstSmokeTime := getDateTimeFromMinutes(referenceTime, wakeUpMin+AwakeSmokeIntervalInMinutes)

	res := &SmokingSchedule{
		Identifier:          user.Identifier,
		UserProgramID:       userProgramID,
		NextSmoke:           firstSmokeTime,
		IntervalInMs:        smokingIntervalInMs,
		RemainingCigarretes: remainingCigarretes,
	}

	if res.NextSmoke.After(referenceTime) {
		log.Printf("Next smoke time for user [%d] is [%s]", user.Identifier, res.NextSmoke)
		return res, nil
	}

	for res.RemainingCigarretes > 0 {
		res.RemainingCigarretes--
		res.NextSmoke = res.NextSmoke.Add(smokingIntervalInMs)

		if res.NextSmoke.After(referenceTime) {
			log.Printf("Next smoke time for user [%d] is [%s]", user.Identifier, res.NextSmoke)
			return res, nil
		}
	}

	return nil, NoSmokingTimeAvailableErr
}

func getDateTimeFromMinutes(baseDate time.Time, minutes int) time.Time {
	return time.Date(
		baseDate.Year(),
		baseDate.Month(),
		baseDate.Day(),
		minutes/60,
		minutes%60,
		0,
		0,
		time.Local,
	)
}
