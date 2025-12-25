package background_workers

import "context"

type BackgroundWorker interface {
	Run(ctx context.Context) error
}
