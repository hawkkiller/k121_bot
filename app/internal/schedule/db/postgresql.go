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
	sql, args, err := sq.Insert("schedules").
		Columns("chat_id").
		Values(schedule.ChatId).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		db.logger.Error(err)
		return err
	}

	db.logger.Infof(sql)

	err = db.client.QueryRow(ctx, sql, args...).Scan(&schedule.ID)
	if err != nil {
		db.logger.Error(err)
		return err
	}
	for _, day := range schedule.Days {
		sql, args, err = sq.Insert("days").
			Columns("schedule_id", "caption").
			Values(schedule.ID, day.Caption).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		err = db.client.QueryRow(ctx, sql, args...).Scan(&day.ID)
		if err != nil {
			db.logger.Error(err)
			return err
		}

		for _, pair := range day.Pairs {
			sql, args, err = sq.Insert("pairs").
				Columns("day_id", "information", "title").
				Values(day.ID, pair.AdditionalInfo, pair.Title).
				Suffix("RETURNING id").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			err = db.client.QueryRow(ctx, sql, args...).Scan(&pair.ID)
			if err != nil {
				db.logger.Error(err)
				return err
			}
		}
	}

	return nil
}

func (db *db) FindOne(ctx context.Context, chatId int64) (schedule.Schedule, error) {

	// defer recover
	if r := recover(); r != nil {
		db.logger.Error(r)
		return schedule.Schedule{}, r.(error)
	}

	s := new(schedule.Schedule)
	//get schedule
	sql, args, _ := sq.Select("id", "chat_id").
		From("schedules").
		Where(sq.Eq{"chat_id": chatId}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_ = db.client.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.ChatId)

	days := make([]schedule.Day, 0)
	//get days
	sql, args, _ = sq.Select("id", "caption").
		From("days").
		Where(sq.Eq{"schedule_id": s.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, _ := db.client.Query(ctx, sql, args...)

	for rows.Next() {
		day := new(schedule.Day)
		_ = rows.Scan(&day.ID, &day.Caption)
		pairs := make([]schedule.Pair, 0)
		// get pairs from certain day
		sql, args, _ = sq.Select("id", "information", "title").
			From("pairs").
			Where(sq.Eq{"day_id": day.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		rows, _ := db.client.Query(ctx, sql, args...)
		for rows.Next() {
			pair := new(schedule.Pair)
			_ = rows.Scan(&pair.ID, &pair.AdditionalInfo, &pair.Title)
			pairs = append(pairs, *pair)
		}
		day.Pairs = pairs
		days = append(days, *day)
	}
	s.Days = days
	return *s, nil
}

func (db *db) Delete(ctx context.Context, chatId int64) error {
	sql, args, err := sq.Delete("schedules").
		Where(sq.Eq{"chat_id": chatId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		db.logger.Error(err)
		return err
	}

	_, err = db.client.Exec(ctx, sql, args...)
	if err != nil {
		db.logger.Error(err)
		return err
	}

	return nil
}
