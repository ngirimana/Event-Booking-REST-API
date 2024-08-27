package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64  `json:"id" `
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(name, email, password) VALUES(?, ?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	newPass, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Name, u.Email, newPass)
	if err != nil {
		return err
	}

	u.ID, err = result.LastInsertId()
	return err
}
func (u *User) Authenticate() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("Credentials invalid.")
	}

	isPasswordValid := utils.CheckPassword(u.Password, retrievedPassword)
	if !isPasswordValid {
		return errors.New("Invalid credentials.")
	}
	return nil

}
