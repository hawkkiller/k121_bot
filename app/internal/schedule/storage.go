package schedule

import "context"

type Storage interface {
	Create(ctx context.Context, schedule Schedule) error
	FindOne(ctx context.Context) Schedule
}
