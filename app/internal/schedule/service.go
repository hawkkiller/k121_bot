package schedule

import (
	"context"
	"github.com/hawkkiller/k121_bot/pkg/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	CreateSchedule(schedule Schedule) error
	GetSchedule(ctx context.Context, chatId int64) (Schedule, error)
}

func (s service) CreateSchedule(schedule Schedule) error {
	err := s.storage.Create(context.Background(), schedule)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetSchedule(ctx context.Context, chatId int64) (Schedule, error) {
	schedule, err := s.storage.FindOne(ctx, chatId)
	if err != nil {
		s.logger.Error(err)
		return Schedule{}, err
	}
	return schedule, nil
}
