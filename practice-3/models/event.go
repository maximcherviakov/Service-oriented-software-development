package models

import (
	"time"

	"example.com/rest-api/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
	UserID      int64     `json:"user_id"`
}

func (event *Event) Save() error {
	query := "INSERT INTO events (name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)"
	result, err := db.DB.Exec(query, event.Name, event.Description, event.Location, event.DateTime, event.UserID)
	if err != nil {
		panic("Could not execute statement: " + err.Error())
	}
	id, err := result.LastInsertId()
	event.ID = id

	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		panic("Could not execute statement: " + err.Error())
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			panic("Could not scan row: " + err.Error())
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"

	var event Event
	row := db.DB.QueryRow(query, id)
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := "UPDATE events SET name = ?, description = ?, location = ?, dateTime = ? WHERE id = ?"

	_, err := db.DB.Exec(query, event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	_, err := db.DB.Exec(query, event.ID)

	return err
}

func (event Event) RegisterForEvent(userId int64) error {
	query := "INSERT INTO registrations (event_id, user_id) VALUES (?, ?)"
	_, err := db.DB.Exec(query, event.ID, userId)

	return err
}

func (event Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	_, err := db.DB.Exec(query, event.ID, userId)

	return err
}
