package schedule

import "context"

type Storage interface {
	Create(ctx context.Context, schedule Schedule) error
	Delete(ctx context.Context, chatId int64) error
	FindOne(ctx context.Context, chatId int64) (Schedule, error)
}
