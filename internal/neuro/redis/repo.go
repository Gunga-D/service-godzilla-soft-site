package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/neuro"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	neuroKey = "neuro:%s"
)

type repo struct {
	db redis.Redis
}

func NewRepo(db redis.Redis) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) SetTaskResult(ctx context.Context, id string, taskResult neuro.TaskResult) error {
	raw, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}

	err = r.db.Set(ctx, fmt.Sprintf(neuroKey, id), raw, pointer.ToInt64(1800))
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetTaskResult(ctx context.Context, id string) (*neuro.TaskResult, error) {
	raw, err := redigo.Bytes(r.db.Get(ctx, fmt.Sprintf(neuroKey, id)))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, nil
		}
	}
	var res neuro.TaskResult
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
