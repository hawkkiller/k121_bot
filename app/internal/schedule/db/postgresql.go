package db

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/hawkkiller/k121_bot/internal/schedule"
	"github.com/hawkkiller/k121_bot/pkg/client/postgresql"
	"github.com/hawkkiller/k121_bot/pkg/logging"
)

type db struct {
	client postgresql.Client
	logger logging.Logger
}

func NewStorage(c postgresql.Client, l logging.Logger) schedule.Storage {
	return &db{
		client: c,
		logger: l,
	}
}

func (db *db) Create(ctx context.Context, schedule schedule.Schedule) error {
	sq.Insert("schedule").Columns("")
	return nil
}

func (db *db) FindOne(ctx context.Context) schedule.Schedule {
	return schedule.Schedule{}
}
