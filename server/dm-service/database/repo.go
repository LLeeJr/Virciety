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

type UniqueUser struct {
	UserName string `json:"userName"`
}

type Repository interface {
	CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error)
	GetDms(ctx context.Context) ([]*model.Dm, error)
	GetChat(ctx context.Context, user1 string, user2 string) ([]*model.Dm, error)
	GetOpenChats(ctx context.Context, userName string) ([]*model.Chat, error)
}

type repo struct {
	DB          *sql.DB
	CurrentDmId string
	CurrentDms  []*model.Dm
}

func (r *repo) GetOpenChats(ctx context.Context, userName string) ([]*model.Chat, error) {

	query := `SELECT DISTINCT dms."toUser" FROM dms WHERE dms."fromUser" = $1`

	uniqueUsers := make([]UniqueUser, 0)
	rows, _ := r.DB.QueryContext(ctx, query, userName)
	for rows.Next() {
		var uniqueUser UniqueUser
		err := rows.Scan(&uniqueUser.UserName)
		if err != nil {
			continue
		}
		uniqueUsers = append(uniqueUsers, uniqueUser)
	}

	chats := make([]*model.Chat, 0)
	for _, user := range uniqueUsers {
		previewQuery := `SELECT dms."toUser", dms.msg FROM dms WHERE dms."fromUser"=$1 AND dms."toUser"=$2
                         AND dms."atTime" = (SELECT MAX(dms."atTime") FROM dms WHERE dms."fromUser"=$1 AND dms."toUser"=$2);`

		previewRows, _ := r.DB.QueryContext(ctx, previewQuery, userName, user.UserName)
		for previewRows.Next() {
			var chat model.Chat
			err := previewRows.Scan(&chat.WithUser, &chat.Preview)
			if err != nil {
				continue
			}
			chats = append(chats, &chat)
		}
	}

	if len(chats) == 0 {
		return nil, errors.New("no opened chats available")
	}

	return chats, nil
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{
		DB:          db,
		CurrentDmId: "",
		CurrentDms:  make([]*model.Dm, 0),
	}, nil
}

func (r *repo) InsertDmEvent(ctx context.Context, dmEvent DmEvent) (err error) {
	query := `INSERT INTO dms ("eventTime", "eventType", id, "fromUser", "toUser", "atTime", msg)
              VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`

	id := ""

	row := r.DB.QueryRowContext(ctx, query, dmEvent.EventTime, dmEvent.EventType,
		dmEvent.DmID, dmEvent.From, dmEvent.To, dmEvent.Time, dmEvent.Msg)

	err = row.Scan(&id)

	r.CurrentDmId = id

	return
}

func (r *repo) CreateDm(ctx context.Context, dmEvent DmEvent) (*model.Dm, error) {
	_ = r.InsertDmEvent(ctx, dmEvent)

	id := fmt.Sprintf("%s__%s__%s", dmEvent.From, dmEvent.Time, dmEvent.To)
	dm := &model.Dm{
		ID:  id,
		Msg: dmEvent.Msg,
	}

	r.CurrentDms = append(r.CurrentDms, dm)

	return dm, nil
}

func (r *repo) GetChat(ctx context.Context, user1 string, user2 string) ([]*model.Dm, error) {
	dms := make([]*model.Dm, 0)

	query := `SELECT * FROM dms WHERE dms."fromUser" = $1 AND "toUser" = $2
              OR dms."fromUser" = $2 AND "toUser" = $1`

	rows, _ := r.DB.QueryContext(ctx, query, user1, user2)
	for rows.Next() {
		var dmEvent DmEvent
		err := rows.Scan(&dmEvent.EventTime, &dmEvent.EventType, &dmEvent.DmID,
			&dmEvent.From, &dmEvent.To, &dmEvent.Time, &dmEvent.Msg)
		if err != nil {
			continue
		}
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
