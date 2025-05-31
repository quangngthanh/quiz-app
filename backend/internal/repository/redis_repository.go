package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"quiz-app/internal/model"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisRepository interface {
	SetQuizSession(quizID string, session *model.QuizSession) error
	GetQuizSession(quizID string) (*model.QuizSession, error)
	UpdateLeaderboard(quizID string, participants []model.Participant) error
	GetLeaderboard(quizID string) ([]model.LeaderboardEntry, error)
	// PublishLeaderboardUpdate(quizID string, leaderboard []model.LeaderboardEntry) error
	// SubscribeToLeaderboardUpdates(quizID string) *redis.PubSub
	// New cache management methods
	DeleteKey(key string) error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
	KeyExists(key string) (bool, error)
	FlushCache() error
}

type redisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *redisRepository) SetQuizSession(quizID string, session *model.QuizSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, fmt.Sprintf("quiz:%s", quizID), data, time.Hour).Err()
}

func (r *redisRepository) GetQuizSession(quizID string) (*model.QuizSession, error) {
	data, err := r.client.Get(r.ctx, fmt.Sprintf("quiz:%s", quizID)).Result()
	if err != nil {
		return nil, err
	}

	var session model.QuizSession
	err = json.Unmarshal([]byte(data), &session)
	return &session, err
}

func (r *redisRepository) UpdateLeaderboard(quizID string, participants []model.Participant) error {
	key := fmt.Sprintf("leaderboard:%s", quizID)

	// Clear existing leaderboard
	r.client.Del(r.ctx, key)

	// Add participants to sorted set
	for _, p := range participants {
		r.client.ZAdd(r.ctx, key, &redis.Z{
			Score:  float64(p.Score),
			Member: fmt.Sprintf("%s:%s", p.UserID, p.Username),
		})
	}

	// Set expiration
	r.client.Expire(r.ctx, key, time.Hour)
	return nil
}

func (r *redisRepository) GetLeaderboard(quizID string) ([]model.LeaderboardEntry, error) {
	key := fmt.Sprintf("leaderboard:%s", quizID)
	results, err := r.client.ZRevRangeWithScores(r.ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var leaderboard []model.LeaderboardEntry
	for i, result := range results {
		// Parse member: "userID:username"
		parts := strings.Split(result.Member.(string), ":")
		if len(parts) != 2 {
			continue
		}

		userID, _ := uuid.Parse(parts[0])
		leaderboard = append(leaderboard, model.LeaderboardEntry{
			UserID:   userID,
			Username: parts[1],
			Score:    int(result.Score),
			Rank:     i + 1,
		})
	}

	return leaderboard, nil
}

func (r *redisRepository) DeleteKey(key string) error {
	result := r.client.Del(r.ctx, key)
	if result.Err() != nil {
		return result.Err()
	}
	log.Printf("🗑️ Deleted cache key: %s", key)
	return nil
}

func (r *redisRepository) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, ttl).Err()
}

func (r *redisRepository) GetTTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

func (r *redisRepository) KeyExists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	return result > 0, err
}

func (r *redisRepository) FlushCache() error {
	return r.client.FlushDB(r.ctx).Err()
}
