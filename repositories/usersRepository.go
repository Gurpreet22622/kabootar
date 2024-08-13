package repositories

import (
	"database/sql"
	"kabootar/models"
	"net/http"
)

type UsersRepository struct {
	dbHandler *sql.DB
}

func NewUsersRepository(dbHandler *sql.DB) *UsersRepository {
	return &UsersRepository{
		dbHandler: dbHandler,
	}
}

func (ur UsersRepository) LoginUser(username string, password string) (string, *models.ResponseError) {
	query := `
			select id from users
			where username=$1 and
			user_password=crypt($2, user_password)`
	rows, err := ur.dbHandler.Query(query, username, password)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
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
	return id, nil
}

func (ur UsersRepository) GetUserRole(accessToken string) (string, *models.ResponseError) {
	query := `
			select user_role from users
			where access_token=$1`
	rows, err := ur.dbHandler.Query(query, accessToken)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var role string
	for rows.Next() {
		err := rows.Scan(&role)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: "Error while reading users role",
			Status:  http.StatusInternalServerError,
		}
	}
	return role, nil
}

func (ur UsersRepository) SetAccessToken(accessToken string, id string) *models.ResponseError {
	query := `
			update users
			set access_token=$1
			where id=$2`
	_, err := ur.dbHandler.Exec(query, accessToken, id)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (ur UsersRepository) RemoveAccessToken(accessToken string) *models.ResponseError {
	query := `
			update users
			set access_token=''
			where access_token=$1`
	_, err := ur.dbHandler.Exec(query, accessToken)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}
