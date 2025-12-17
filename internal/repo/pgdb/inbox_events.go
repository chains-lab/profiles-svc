package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const inboxEventsTable = "inbox_events"

type InboxEvent struct {
	ID           uuid.UUID `db:"id"`
	Topic        string    `db:"topic"`
	EventType    string    `db:"event_type"`
	EventVersion uint      `db:"event_version"`
	Key          string    `db:"key"`
	Payload      []byte    `db:"payload"`

	Status      string     `db:"status"`
	Attempts    uint       `db:"attempts"`
	NextRetryAt time.Time  `db:"next_retry_at"`
	CreatedAt   time.Time  `db:"created_at"`
	ProcessedAt *time.Time `db:"processed_at"`
}

type InboxEventsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewInboxEventsQ(db *sql.DB) InboxEventsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return InboxEventsQ{
		db:       db,
		selector: builder.Select("*").From(inboxEventsTable),
		inserter: builder.Insert(inboxEventsTable),
		updater:  builder.Update(inboxEventsTable),
		deleter:  builder.Delete(inboxEventsTable),
		counter:  builder.Select("COUNT(*) AS count").From(inboxEventsTable),
	}
}

func (q InboxEventsQ) New() InboxEventsQ {
	return NewInboxEventsQ(q.db)
}

func (q InboxEventsQ) Insert(ctx context.Context, input InboxEvent) error {
	values := map[string]interface{}{
		"id":            input.ID,
		"topic":         input.Topic,
		"event_type":    input.EventType,
		"event_version": input.EventVersion,
		"key":           input.Key,
		"payload":       input.Payload,
		"status":        input.Status,
		"attempts":      input.Attempts,
		"next_retry_at": input.NextRetryAt,
		"created_at":    input.CreatedAt,
		"processed_at":  input.ProcessedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", inboxEventsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q InboxEventsQ) Update(ctx context.Context) error {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", inboxEventsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q InboxEventsQ) UpdateStatus(status string) InboxEventsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q InboxEventsQ) AddAttempts() InboxEventsQ {
	q.updater = q.updater.Set("attempts", sq.Expr("attempts + 1"))
	return q
}

func (q InboxEventsQ) UpdateAttempts(attempts int) InboxEventsQ {
	q.updater = q.updater.Set("attempts", attempts)
	return q
}

func (q InboxEventsQ) UpdateNextRetryAt(t time.Time) InboxEventsQ {
	q.updater = q.updater.Set("next_retry_at", t)
	return q
}

func (q InboxEventsQ) UpdateProcessedAt(processedAt *time.Time) InboxEventsQ {
	q.updater = q.updater.Set("processed_at", processedAt)
	return q
}

func (q InboxEventsQ) Get(ctx context.Context) (InboxEvent, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return InboxEvent{}, fmt.Errorf("building get query for %s: %w", inboxEventsTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var e InboxEvent
	err = row.Scan(
		&e.ID,
		&e.Topic,
		&e.EventType,
		&e.EventVersion,
		&e.Key,
		&e.Payload,
		&e.Status,
		&e.Attempts,
		&e.NextRetryAt,
		&e.CreatedAt,
		&e.ProcessedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return InboxEvent{}, nil
		}
		return InboxEvent{}, err
	}

	return e, nil
}

func (q InboxEventsQ) Select(ctx context.Context) ([]InboxEvent, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", inboxEventsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []InboxEvent
	for rows.Next() {
		var e InboxEvent
		err = rows.Scan(
			&e.ID,
			&e.Topic,
			&e.EventType,
			&e.EventVersion,
			&e.Key,
			&e.Payload,
			&e.Status,
			&e.Attempts,
			&e.NextRetryAt,
			&e.CreatedAt,
			&e.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning inbox event: %w", err)
		}
		out = append(out, e)
	}

	return out, nil
}

func (q InboxEventsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", inboxEventsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q InboxEventsQ) FilterID(id ...uuid.UUID) InboxEventsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	return q
}

func (q InboxEventsQ) FilterStatus(status string) InboxEventsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	return q
}

func (q InboxEventsQ) FilterReadyToProcess(now time.Time) InboxEventsQ {
	q.selector = q.selector.Where(sq.LtOrEq{"next_retry_at": now})
	q.counter = q.counter.Where(sq.LtOrEq{"next_retry_at": now})
	return q
}

func (q InboxEventsQ) Page(limit, offset uint) InboxEventsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q InboxEventsQ) OrderCreatedAt(ascending bool) InboxEventsQ {
	if ascending {
		q.selector = q.selector.OrderBy("created_at ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC")
	}
	return q
}

func (q InboxEventsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", inboxEventsTable, err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q InboxEventsQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
