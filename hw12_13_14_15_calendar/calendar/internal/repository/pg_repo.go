package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/model"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func NewPg(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (s *pgRepository) AddEvent(ctx context.Context, params AddEventParams) (id model.EventID, err error) {
	if len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 || params.OwnerID == 0 {
		err = fmt.Errorf("%w: event title, date, ownerUserId must be provided", ErrInvalidParams)
		return
	}
	query := `insert into events(
		id, 
		title, 
		date,
		duration, 
		description,
		owner_id,
		notify_for
	) values ($1, $2, $3, $4, $5, $6, $7)`
	id = GenerateEventID()
	var tx pgx.Tx
	if tx, err = s.db.Begin(ctx); err != nil {
		return
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queryParams := []interface{}{
		id,
		params.Title,
		params.Date,
		params.Duration,
		params.Description,
		params.OwnerID,
		params.NotifyFor,
	}
	if _, err = tx.Exec(ctx, query, queryParams...); err != nil {
		return
	}
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}

func (s *pgRepository) UpdateEvent(ctx context.Context, params UpdateEventParams) (err error) {
	if len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 {
		return fmt.Errorf("%w: params id, title, date must be provided", ErrInvalidParams)
	}
	query := `update events set 
		title = $1, 
		date = $2,
		duration = $3, 
		description = $4,
		notify_for = $5,
	where id = $6`
	var tx pgx.Tx
	if tx, err = s.db.Begin(ctx); err != nil {
		return
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queryParams := []interface{}{
		params.Title,
		params.Date,
		params.Duration,
		params.Description,
		params.NotifyFor,
		params.ID,
	}
	if _, err = tx.Exec(ctx, query, queryParams...); err != nil {
		return
	}
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}

func (s *pgRepository) DeleteEvent(ctx context.Context, id model.EventID) (err error) {
	query := "delete from events where id = $1"
	var tx pgx.Tx
	if tx, err = s.db.Begin(ctx); err != nil {
		return
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if _, err = tx.Exec(ctx, query, id); err != nil {
		return
	}
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}

func (s *pgRepository) GetEvents(ctx context.Context, params GetEventParams) (events []model.Event, err error) {
	if params.BeginDate.IsZero() || params.EndDate.IsZero() {
		err = fmt.Errorf("%w: beginDate and endDate must be provided", ErrInvalidParams)
		return
	}
	query := `select 
		e.id,
		e.title,
		e.date,
		e.duration,
		e.description,
		e.owner_id,
		e.notify_for
	from events e`
	beginDate := params.BeginDate.Format(DateFormat)
	endDate := params.EndDate.Format(DateFormat)
	queryParams := []interface{}{beginDate, endDate}
	if beginDate == endDate {
		query += "where event_date::date between = $1"
	} else {
		query += "where event_date::date between $1 and $2"
	}
	var rows pgx.Rows
	if rows, err = s.db.Query(ctx, query, queryParams...); err != nil {
		return
	}
	defer rows.Close()
	events = make([]model.Event, 0)
	for rows.Next() {
		var e model.Event
		if err = rows.Scan(
			&e.ID,
			&e.Title,
			&e.Date,
			&e.Duration,
			&e.Description,
			&e.OwnerID,
			&e.NotifyFor,
		); err != nil {
			return
		}
		events = append(events, e)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
