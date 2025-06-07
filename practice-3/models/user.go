package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) Save() error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	result, err := db.DB.Exec(query, user.Email, hashedPassword)
	if err != nil {
		panic("Could not execute statement: " + err.Error())
	}
	id, err := result.LastInsertId()
	user.ID = id

	return err
}

func (user *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, user.Email)

	var storedPassword string
	err := row.Scan(&user.ID, &storedPassword)
	if err != nil {
		return err
	}

	isValid := utils.CheckPasswordHash(user.Password, storedPassword)

	if !isValid {
		return errors.New("invalid credentials")
	}

	return nil
}
