package domain

import "context"

// TrackEvent merepresentasikan data sebuah event yang dilacak
type TrackEvent struct {
	EventName string `json:"event_name"`
	Count     int64  `json:"count"`
}

// TrackRepository adalah kontrak untuk interaksi dengan penyimpanan data (Redis)
type TrackRepository interface {
	IncrementEvent(ctx context.Context, eventName string) (int64, error)
	GetEventCount(ctx context.Context, eventName string) (int64, error)
	GetAllEvents(ctx context.Context) ([]TrackEvent, error)
	ResetEvent(ctx context.Context, eventName string) error
}

// TrackUsecase adalah kontrak untuk logika bisnis tracking
type TrackUsecase interface {
	RecordEvent(ctx context.Context, eventName string) (*TrackEvent, error)
	GetEventStats(ctx context.Context, eventName string) (*TrackEvent, error)
	GetAllStats(ctx context.Context) ([]TrackEvent, error)
	ResetEventStats(ctx context.Context, eventName string) error
}
