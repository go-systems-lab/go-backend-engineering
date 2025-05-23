package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kuluruvineeth/social-go/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, id int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user:%v", id)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	user := &store.User{}
	if data != "" {
		if err := json.Unmarshal([]byte(data), user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user:%v", user.ID)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEx(ctx, cacheKey, jsonData, UserExpTime).Err()
}
