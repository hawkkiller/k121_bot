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
	CreateSchedule(ctx context.Context, schedule Schedule) error
	GetSchedule(ctx context.Context, chatId int64) (Schedule, error)
	DeleteSchedule(ctx context.Context, chatId int64) error
}

func (s service) CreateSchedule(ctx context.Context, schedule Schedule) error {
	err := s.storage.Create(ctx, schedule)
	if err != nil {
		return err
	}
	return nil
}

func (s service) DeleteSchedule(ctx context.Context, chatId int64) error {
	err := s.storage.Delete(ctx, chatId)
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
