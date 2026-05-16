package usecase

import (
	"context"
	"fmt"
	"strings"

	"track-method/domain"
)

type trackUsecase struct {
	repo domain.TrackRepository
}

func NewTrackUsecase(repo domain.TrackRepository) domain.TrackUsecase {
	return &trackUsecase{repo: repo}
}

func (u *trackUsecase) RecordEvent(ctx context.Context, eventName string) (*domain.TrackEvent, error) {
	eventName = sanitizeEventName(eventName)
	if eventName == "" {
		return nil, fmt.Errorf("nama event tidak boleh kosong")
	}

	count, err := u.repo.IncrementEvent(ctx, eventName)
	if err != nil {
		return nil, err
	}

	return &domain.TrackEvent{
		EventName: eventName,
		Count:     count,
	}, nil
}

func (u *trackUsecase) GetEventStats(ctx context.Context, eventName string) (*domain.TrackEvent, error) {
	eventName = sanitizeEventName(eventName)
	if eventName == "" {
		return nil, fmt.Errorf("nama event tidak boleh kosong")
	}

	count, err := u.repo.GetEventCount(ctx, eventName)
	if err != nil {
		return nil, err
	}

	return &domain.TrackEvent{
		EventName: eventName,
		Count:     count,
	}, nil
}

func (u *trackUsecase) GetAllStats(ctx context.Context) ([]domain.TrackEvent, error) {
	return u.repo.GetAllEvents(ctx)
}

func (u *trackUsecase) ResetEventStats(ctx context.Context, eventName string) error {
	eventName = sanitizeEventName(eventName)
	if eventName == "" {
		return fmt.Errorf("nama event tidak boleh kosong")
	}
	return u.repo.ResetEvent(ctx, eventName)
}

// sanitizeEventName membersihkan nama event dari spasi dan karakter tidak valid
func sanitizeEventName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
