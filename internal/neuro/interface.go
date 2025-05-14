package neuro

import "context"

type Repository interface {
	CreateFinishedNeuroTask(ctx context.Context, finishedTask Task) (int64, error)
}
