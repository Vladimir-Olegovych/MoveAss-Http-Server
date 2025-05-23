package storage

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID           string
	UserName     string
	UserPassword string
}

type Stats struct {
	ID        string
	UserMoney int64
}

type DataBaseService interface {
	Close()

	CreateUser(id, name, password string) error
	FindUser(name, password string) (*User, error)
	FindStatById(id string) (*Stats, error)
}

type SQLService struct {
	DataBase *sql.DB
}

func (s *SQLService) CreateUser(id, name, password string) error {
	query := `INSERT INTO users(id, user_name, user_password) VALUES (?, ?, ?)`
	_, err := s.DataBase.Exec(query, id, name, password)

	if err != nil {
		return err
	}

	query = `INSERT INTO stats(id, user_money) VALUES (?, ?)`
	_, err = s.DataBase.Exec(query, id, 100)

	return err
}

func (s *SQLService) FindUser(name, password string) (*User, error) {
	query := `SELECT id, user_name, user_password FROM users WHERE user_name = ? AND user_password = ?`
	row := s.DataBase.QueryRow(query, name, password)

	var user User
	err := row.Scan(&user.ID, &user.UserName, &user.UserPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found or invalid password")
		}
		return nil, err
	}
	return &user, nil
}

func (s *SQLService) FindStatById(id string) (*Stats, error) {
	query := `SELECT id, user_money FROM stats WHERE id = ?`
	row := s.DataBase.QueryRow(query, id)

	var stat Stats
	err := row.Scan(&stat.ID, &stat.UserMoney)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("stat not found")
		}
		return nil, err
	}
	return &stat, nil
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
