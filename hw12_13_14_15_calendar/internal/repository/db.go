package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/logger"
	"github.com/kpechenenko/hw12_13_14_15_calendar/internal/model"
	"go.uber.org/zap"
	"time"
)

// PgEventRepository хранит события в БД postgres.
type PgEventRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewPgEventRepository(db *sql.DB) *PgEventRepository {
	const loggerName = "pgEventRepository"
	return &PgEventRepository{
		db:     db,
		logger: logger.GetNamed(loggerName),
	}
}

func (r *PgEventRepository) CreateEvent(ctx context.Context, params CreateEventParams) (eventID model.EventID, err error) {
	if len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 || params.OwnerUserID == 0 {
		err = fmt.Errorf("%w: event title, date, ownerUserId must be provided", ErrInvalidParam)
		return
	}
	query := `
	insert into events(
		title,
		event_date,
		duration,
		description,
		owner_user_id,
		notify_time
	) values ($1, $2, $3, $4, $5, $6)
	returning event_id
	`
	queryParams := []interface{}{
		params.Title,
		params.Date.Format(DateTimeFormat),
		params.Duration,
		params.Description,
		params.OwnerUserID,
		params.NotifyTime,
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("begin tx: %v", err)
		return
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, query, queryParams...)
	if err != nil {
		r.logger.Errorw("save event", "error", err, "query", query, "queryParams", queryParams)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		r.logger.Errorf("get last inserted id: %v", err)
		return
	}
	if err = tx.Commit(); err != nil {
		r.logger.Errorf("commit tx: %v", err)
		return
	}
	eventID = model.EventID(id)
	return
}

func (r *PgEventRepository) UpdateEvent(ctx context.Context, params UpdateEventParams) (err error) {
	if params.EventID == 0 || len(params.Title) == 0 || params.Date.IsZero() || params.Duration.Abs() == 0 {
		return fmt.Errorf("%w: params id, title, date must be provided", ErrInvalidParam)
	}
	query := `
	update events 
	set 
		title = $2,
		event_date = $3,
		duration = $4,
		description = $5,
		notify_time = $6.
	where event_id = $1'
	`
	queryParams := []interface{}{
		params.EventID,
		params.Title,
		params.Date,
		params.Duration,
		params.Description,
		params.NotifyTime,
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("begin tx: %v", err)
		return
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, query, queryParams...)
	if err != nil {
		r.logger.Errorw("update event", "error", err, "query", query, "queryParams", queryParams)
		return
	}
	if affectedCnt, err := res.RowsAffected(); err == nil {
		r.logger.Debugf("update event %d affected rows %d", params.EventID, affectedCnt)
	}
	if err = tx.Commit(); err != nil {
		r.logger.Errorf("commit tx: %v", err)
		return
	}
	return
}

func (r *PgEventRepository) DeleteEvent(ctx context.Context, eventID model.EventID) (err error) {
	if eventID == 0 {
		err = fmt.Errorf("%w: eventID must be provided", ErrInvalidParam)
		return
	}
	query := "delete from events where event_id = $1"
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("begin tx: %v", err)
		return
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, query, eventID)
	if err != nil {
		r.logger.Errorw("delete event", "error", err, "query", query, "eventID", eventID)
		return
	}
	if affectedCnt, err := res.RowsAffected(); err == nil {
		r.logger.Debugf("delete event %d affected rows %d", eventID, affectedCnt)
	}
	if err = tx.Commit(); err != nil {
		r.logger.Errorf("commit tx: %v", err)
		return
	}
	return
}

func (r *PgEventRepository) FindEvent(ctx context.Context, params FindEventParams) (events []model.Event, err error) {
	if params.BeginDate.IsZero() || params.EndDate.IsZero() {
		err = fmt.Errorf("%w: beginDate and endDate must be provided", ErrInvalidParam)
		return
	}
	query := `
	select
		event_id,
		title, 
		event_date,
		duration,
		description,
		owner_user_id,
		notify_time
	from
		events
	`
	beginDate := params.BeginDate.Format(DateFormat)
	endDate := params.EndDate.Format(DateFormat)
	queryParams := []interface{}{beginDate, endDate}
	if beginDate == endDate {
		query += "where event_date::date between = $1"
	} else {
		query += "where event_date::date between $1 and $2"
	}
	r.logger.Debugf("query: %s", query)
	rows, err := r.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		r.logger.Errorw("find event", "error", err, "query", query)
		return
	}
	defer rows.Close()
	events = make([]model.Event, 0, 10)
	for rows.Next() {
		var e model.Event
		var date string
		err = rows.Scan(
			&e.EventID,
			&e.Title,
			&date,
			&e.Duration,
			&e.Description,
			&e.OwnerUserID,
			&e.NotifyTime,
		)
		if e.Date, err = time.Parse(DateTimeFormat, date); err != nil {
			r.logger.Errorw("parse event date", "date", date, "layout", DateTimeFormat, "eventID", e.EventID)
			return
		}
		if err != nil {
			r.logger.Errorf("scan event: %v", err)
			return
		}
		events = append(events, e)
	}
	if err = rows.Err(); err != nil {
		r.logger.Errorf("iteration err: %v", err)
		return
	}
	return
}
