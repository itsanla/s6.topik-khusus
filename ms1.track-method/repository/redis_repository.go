package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"track-method/domain"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "track:event:"

type redisTrackRepository struct {
	client *redis.Client
}

func NewRedisTrackRepository(client *redis.Client) domain.TrackRepository {
	return &redisTrackRepository{client: client}
}

func (r *redisTrackRepository) IncrementEvent(ctx context.Context, eventName string) (int64, error) {
	key := fmt.Sprintf("%s%s", keyPrefix, eventName)
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("gagal increment event '%s': %w", eventName, err)
	}
	return count, nil
}

func (r *redisTrackRepository) GetEventCount(ctx context.Context, eventName string) (int64, error) {
	key := fmt.Sprintf("%s%s", keyPrefix, eventName)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil count event '%s': %w", eventName, err)
	}
	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("gagal parse nilai '%s': %w", val, err)
	}
	return count, nil
}

func (r *redisTrackRepository) GetAllEvents(ctx context.Context) ([]domain.TrackEvent, error) {
	pattern := fmt.Sprintf("%s*", keyPrefix)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil semua event keys: %w", err)
	}

	events := make([]domain.TrackEvent, 0, len(keys))
	for _, key := range keys {
		val, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		count, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			continue
		}
		eventName := strings.TrimPrefix(key, keyPrefix)
		events = append(events, domain.TrackEvent{
			EventName: eventName,
			Count:     count,
		})
	}
	return events, nil
}

func (r *redisTrackRepository) ResetEvent(ctx context.Context, eventName string) error {
	key := fmt.Sprintf("%s%s", keyPrefix, eventName)
	deleted, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("gagal menghapus event '%s': %w", eventName, err)
	}
	if deleted == 0 {
		return fmt.Errorf("event '%s' tidak ditemukan", eventName)
	}
	return nil
}
