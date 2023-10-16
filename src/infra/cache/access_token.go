package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kits/api/src/common/fault"
	"kits/api/src/core/domains"
)

type AccessTokenCache struct {
	client redis.UniversalClient
}

func NewAccessTokenCache(c redis.UniversalClient) domains.AccessTokenCache {
	return &AccessTokenCache{client: c}
}

func (a AccessTokenCache) buildKey(userId int64) string {
	return fmt.Sprintf("access_token:%d", userId)
}

func (a AccessTokenCache) Save(ctx context.Context, userId int64, token *domains.AccessToken) error {
	caller := "AccessTokenCache.Save"
	k := a.buildKey(userId)

	bytes, err := json.Marshal(token)
	if err != nil {
		return fault.Wrapf(err, "[%v] failed to marshal %+v", caller, *token)
	}

	err = a.client.Set(ctx, k, bytes, token.ExpireTime).Err()
	if err != nil {
		return fault.DBWrapf(err, "[%v] failed to save access token", caller)
	}
	return nil
}

func (a AccessTokenCache) FindByUserId(ctx context.Context, userId int64) (*domains.AccessToken, error) {
	caller := "AccessTokenCache.FindByUserId"
	k := a.buildKey(userId)

	var token *domains.AccessToken
	val, err := a.client.Get(ctx, k).Result()

	if err != nil {
		return nil, fault.DBWrapf(err, "[%v] failed to get val", caller)
	}

	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return nil, fault.Wrapf(err, "[%v] failed to unmarshal %s", caller, val)
	}

	return token, nil
}
