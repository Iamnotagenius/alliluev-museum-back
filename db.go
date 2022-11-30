package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type MuseumDB struct {
	db *sql.DB
}

func NewDB(address string) MuseumDB {
	cfg := mysql.Config{
		User:   "root",
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Addr:   address,
		Net:    "tcp",
		DBName: "museum",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return MuseumDB{db}
}

func (db MuseumDB) queryExhibits(errHandler func(error) error, options string, args ...any) ([]Exhibit, error) {
	result := make([]Exhibit, 0)

	rows, err := db.db.Query("SELECT * FROM exhibits "+options, args...)
	if err != nil {
		return result, errHandler(err)
	}

	for rows.Next() {
		var current Exhibit
		var paths string
		err := rows.Scan(&current.ID, &current.Name, &paths, &current.Description, &current.Room)
		if err != nil {
			return result, fmt.Errorf("Error while scanning rows: %v", err)
		}
		current.Pictures = strings.Split(paths, ",")
		result = append(result, current)
	}

	return result, nil
}

func (db MuseumDB) ExhibitByID(id int) (Exhibit, error) {
	var result Exhibit
	var paths string

	row := db.db.QueryRow("SELECT * FROM exhibits WHERE id = ?", id)

	err := row.Scan(&result.ID, &result.Name, &paths, &result.Description, &result.Room)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("Exhibit with id %d not found.", id)
		}
		return result, fmt.Errorf("Error querying exhibit No. %d: %v", id, err)
	}

	result.Pictures = strings.Split(paths, ",")

	return result, nil
}

func (db MuseumDB) ExhibitsByRoomID(id int) ([]Exhibit, error) {
	return db.queryExhibits(func(err error) error {
		if err == sql.ErrNoRows {
			return fmt.Errorf("There is no exhibits found in room No. %d", id)
		}
		return fmt.Errorf("Error when querying exhibits by room No. %d: %v", id, err)
	}, "WHERE room = ?", id)
}

func (db MuseumDB) RoomByID(id int) (Room, error) {
	var result Room
	var paths string

	row := db.db.QueryRow("SELECT * FROM rooms WHERE id = ?", id)

	err := row.Scan(&result.ID, &result.Name, &paths)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("Room with id %d not found.", id)
		}
		return result, fmt.Errorf("Error querying room No. %d: %v", id, err)
	}

	result.Pictures = strings.Split(paths, ",")

	return result, nil
}

func (db MuseumDB) GetAllExhibits() ([]Exhibit, error) {
	return db.queryExhibits(func(err error) error {
		if err == sql.ErrNoRows {
			return fmt.Errorf("There is no exhibits in database")
		}
		return fmt.Errorf("Error when querying for all exhibits: %v", err)
	}, "")
}

func (db MuseumDB) GetAllRooms() ([]Room, error) {
	result := make([]Room, 0)

	rows, err := db.db.Query("SELECT * FROM rooms")
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("There is no exhibits in database")
		}
		return result, fmt.Errorf("Error when querying for all exhibits: %v", err)
	}

	for rows.Next() {
		var current Room
		var paths string
		err := rows.Scan(&current.ID, &current.Name, &paths)
		if err != nil {
			return result, fmt.Errorf("Error while scanning rows: %v", err)
		}
		current.Pictures = strings.Split(paths, ",")
		result = append(result, current)
	}

	return result, nil
}
