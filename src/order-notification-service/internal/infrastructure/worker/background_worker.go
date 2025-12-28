package worker

import "context"

type BackgroundWorker interface {
	Run(ctx context.Context) error
}
