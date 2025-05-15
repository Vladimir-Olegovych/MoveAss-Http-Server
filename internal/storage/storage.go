package storage

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID           int
	UserToken    string
	UserName     string
	UserPassword string
}

type Stats struct {
	ID        int
	UserToken string
	UserMoney int64
}

type DataBaseService interface {
	Close()

	AddUser(name, password, token string) error
	FindUserByName(name string) (*User, error)
	FindUserByToken(token string) (*User, error)

	AddStats(token string, money int64) error
	FindStatsByToken(token string) (*Stats, error)
}

type SQLService struct {
	DataBase *sql.DB
}

func (s *SQLService) AddUser(name, password, token string) error {
	query := `INSERT INTO users(user_token, user_name, user_password) VALUES (?, ?, ?)`
	_, err := s.DataBase.Exec(query, token, name, password)
	return err
}

func (s *SQLService) FindUserByName(name string) (*User, error) {
	query := `SELECT id, user_token, user_name, user_password FROM users WHERE user_name = ?`
	row := s.DataBase.QueryRow(query, name)

	var user User
	err := row.Scan(&user.ID, &user.UserToken, &user.UserName, &user.UserPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *SQLService) FindUserByToken(token string) (*User, error) {
	query := `SELECT id, user_token, user_name, user_password FROM users WHERE user_token = ?`
	row := s.DataBase.QueryRow(query, token)

	var user User
	err := row.Scan(&user.ID, &user.UserToken, &user.UserName, &user.UserPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *SQLService) AddStats(token string, money int64) error {
	query := `INSERT INTO stats(user_token, user_money) VALUES (?, ?)`
	_, err := s.DataBase.Exec(query, token, money)
	return err
}

func (s *SQLService) FindStatsByToken(token string) (*Stats, error) {
	query := `SELECT id, user_token, user_money FROM stats WHERE user_token = ?`
	row := s.DataBase.QueryRow(query, token)

	var stats Stats
	err := row.Scan(&stats.ID, &stats.UserToken, &stats.UserMoney)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &stats, nil
}

func (s *SQLService) Close() {
	s.DataBase.Close()
}

func Open(outfilename string, sqlfilenames ...string) (DataBaseService, error) {

	db, sqlError := sql.Open("sqlite3", outfilename)
	if sqlError == nil {
		for _, filename := range sqlfilenames {
			content, fileError := os.ReadFile(filename)
			if fileError != nil {
				return nil, fileError
			}
			db.Exec(string(content))
		}
	} else {
		return nil, sqlError
	}

	var databaseService DataBaseService = &SQLService{
		DataBase: db,
	}

	return databaseService, nil
}
