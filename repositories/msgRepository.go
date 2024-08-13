package repositories

import (
	"database/sql"
	"kabootar/models"
	"net/http"
)

type MsgRepository struct {
	dbHandler *sql.DB
}

func NewMsgRepository(dbHandler *sql.DB) *MsgRepository {
	return &MsgRepository{
		dbHandler: dbHandler,
	}
}

func (mr MsgRepository) SaveMessage(message string, username string) (string, *models.ResponseError) {
	query := `
			insert into chats(message,username)
			values
			($1,$2)
			returning message`
	rows, err := mr.dbHandler.Query(query, message, username)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var messg string
	for rows.Next() {
		err := rows.Scan(&messg)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: "Error while reading user id",
			Status:  http.StatusInternalServerError,
		}
	}
	return messg, nil
}

func (mr MsgRepository) PullMessage() ([][]string, *models.ResponseError) {
	query := `
			select message,timestamp,username
			from chats`
	rows, err := mr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var retData [][]string
	var mssg, timst, usrnm string
	for rows.Next() {
		var tmpData []string
		err := rows.Scan(&mssg, &timst, &usrnm)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		tmpData = append(tmpData, mssg, timst, usrnm)
		retData = append(retData, tmpData)
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: "Error while reading rows",
			Status:  http.StatusInternalServerError,
		}
	}
	return retData, nil
}
