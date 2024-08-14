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

func (ur UsersRepository) CreateUser(user *models.User) (*models.User, *models.ResponseError) {
	ok := ur.checkDuplicateUser(user)
	if !ok {
		return nil, &models.ResponseError{
			Message: "User already exists",
			Status:  http.StatusNotAcceptable,
		}
	}
	query := `insert into users(username, user_password, user_role)
				values($1,crypt($2,gen_salt('bf')),$3)
				returning id`
	rows, err := ur.dbHandler.Query(query, user.Username, user.Password, "user")
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var userId string
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: "Error while reading rows",
			Status:  http.StatusInternalServerError,
		}
	}
	return &models.User{
		ID:       userId,
		Username: user.Username,
		Role:     "user",
	}, nil
}

func (ur UsersRepository) checkDuplicateUser(user *models.User) bool {
	query := `select id from users
				where username=$1`
	rows, err := ur.dbHandler.Query(query, user.Username)
	if err != nil {
		return false
	}
	defer rows.Close()
	var userId string
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return false
		}
	}
	if rows.Err() != nil {
		return false
	}
	if userId == "" {
		return true
	} else {
		return false
	}
}
