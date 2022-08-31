package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Fajar-Islami/scrapping_test/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisTrackingRepoImpl struct {
	redisClient *redis.Client
}

type RedisTrackingRepository interface {
	GetTrackingByQueryCtx(ctx context.Context, key string) (*model.DataStruct, error)
	SetTrackingCtx(ctx context.Context, key string, time int, data *model.DataStruct) error
	DeleteTrackingCtx(ctx context.Context, key string) error
}

func NewRedisTrackingRepo(redisClient *redis.Client) RedisTrackingRepository {
	return &redisTrackingRepoImpl{
		redisClient: redisClient,
	}
}

func (ra *redisTrackingRepoImpl) GetTrackingByQueryCtx(ctx context.Context, key string) (*model.DataStruct, error) {
	fmt.Printf("Get keys %s from redis\n", key)

	newBase := &model.DataStruct{}
	newBytes, err := ra.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "trackingRedisRepo.GetTrackingByQueryCtx.redisClient.Get")
	}

	if err = json.Unmarshal(newBytes, newBase); err != nil {
		return nil, errors.Wrap(err, "trackingRedisRepo.GetTrackingByQueryCtx.json.Unmarshal")
	}

	return newBase, nil
}

func (ra *redisTrackingRepoImpl) SetTrackingCtx(ctx context.Context, key string, times int, data *model.DataStruct) error {
	newBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "trackingRedisRepo.SetTrackingCtx.json.Marshal")
	}

	if err = ra.redisClient.Set(ctx, key, newBytes, time.Minute*time.Duration(times)).Err(); err != nil {
		return errors.Wrap(err, "trackingRedisRepo.SetTrackingCtx.redisClient.set")
	}
	fmt.Printf("Set keys %s to redis\n", key)
	return nil
}

func (ra *redisTrackingRepoImpl) DeleteTrackingCtx(ctx context.Context, key string) error {
	if err := ra.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "trackingRedisRepo.DeleteTrackingCtx.redisClient.Del")
	}
	fmt.Printf("Delete keys %s from redis\n", key)
	return nil
}
