package database

import (
	"context"
	"database/sql"
	"dm-service/graph/model"
	"errors"
	"fmt"
)

type DmEvent struct {
	EventTime string `json:"eventTime"`
	EventType string `json:"eventType"`
	DmID string `json:"id"`
	From string `json:"from"`
	To string `json:"to"`
	Time string `json:"time"`
	Msg string `json:"msg"`
}

type Repository interface {
	CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error)
	GetDms(ctx context.Context) ([]*model.Dm, error)
	GetDmsByFromTo(ctx context.Context, from string, to string) ([]*model.Dm, error)
}

type repo struct {
	DB *sql.DB
	DmEvents []*DmEvent
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:       db,
		DmEvents: make([]*DmEvent, 0),
	}, nil
}

func (r *repo) CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error) {
	r.DmEvents = append(r.DmEvents, &dmEvent)

	query := `INSERT INTO dms ("eventTime", "eventType", id, "fromUser", "toUser", "atTime", msg)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.DB.ExecContext(ctx, query, dmEvent.EventTime, dmEvent.EventType,
		dmEvent.DmID, dmEvent.From, dmEvent.To, dmEvent.Time, dmEvent.Msg)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s__%s__%s", dmEvent.From, dmEvent.Time, dmEvent.To)
	return &model.Dm{
		ID:  id,
		Msg: dmEvent.Msg,
	}, nil
}

func (r *repo) GetDmsByFromTo(ctx context.Context, from string, to string) ([]*model.Dm, error) {
	dms := make([]*model.Dm, 0)

	query := `SELECT * FROM dms WHERE dms."fromUser" = $1 AND "toUser" = $2`

	rows, _ := r.DB.QueryContext(ctx, query, from, to)
	for rows.Next() {
		var dmEvent DmEvent
		err := rows.Scan(&dmEvent.EventTime, &dmEvent.EventType, &dmEvent.DmID,
			&dmEvent.From, &dmEvent.To, &dmEvent.Time, &dmEvent.Msg)
		if err != nil {
			continue
		}
		if dmEvent.From == from && dmEvent.To == to {
			dm := &model.Dm{
				ID:  dmEvent.DmID,
				Msg: dmEvent.Msg,
			}
			dms = append(dms, dm)

		}
	}

	if len(dms) == 0 {
		return nil, errors.New("no chat history available")
	}

	return dms, nil

}

func (r *repo) GetDms(ctx context.Context) ([]*model.Dm, error) {
	dms := make([]*model.Dm, 0)

	query := `SELECT * FROM dms`

	rows, _ := r.DB.QueryContext(ctx, query)
	for rows.Next() {
		var dmEvent DmEvent
		err := rows.Scan(&dmEvent.EventTime, &dmEvent.EventType, &dmEvent.DmID,
			&dmEvent.From, &dmEvent.To, &dmEvent.Time, &dmEvent.Msg)
		if err != nil {
			continue
		}
		fmt.Println(dmEvent)
		dm := &model.Dm{
			ID:  dmEvent.DmID,
			Msg: dmEvent.Msg,
		}
		dms = append(dms, dm)
	}

	if len(dms) == 0 {
		return nil, errors.New("no chat history available")
	}

	return dms, nil
}
