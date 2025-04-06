package core

import "github.com/sousair/gocore/pkg/event"

const (
	StartProgramQueueEventType         event.EventType = "queue.start-program"
	ScheduleSmokingTimeQueueEventType  event.EventType = "queue.smoking-time"
	CalculateNextProgramQueueEventType event.EventType = "queue.calculate-next-program"

	NotifySmokingTimeEventType event.EventType = "notify.smoking-time"
)
