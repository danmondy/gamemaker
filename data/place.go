package data

import (
	//"fmt"
	"time"
)

const CREATE_PLACE_TABLE = `
CREATE TABLE IF NOT EXISTS place (
id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
user_id INT UNSIGNED,
title VARCHAR(50) NOT NULL,
subtitle VARCHAR(255),
description VARCHAR(900),
lon float(10,6) NOT NULL,
lat float(10,6) NOT NULL,
date_created TIMESTAMP
);
`

type Place struct {
	Id          int64
	UserId      int64
	Title       string
	Subtitle    string
	Description string
	Lat, Lon    float64
	DateCreated time.Time
}

func NewPlace(u User) *Place {
	return &Place{
		UserId:      u.Id,
		Title:       "New Place",
		Subtitle:    "",
		Description: "",
		Lat:         32.80502926732618,
		Lon:         -79.94039009094328,
		DateCreated: time.Now(),
	}
}

func GetUserPlaces(userId int64) ([]Place, error) {
	places := []Place{}
	rows, err := DB.Query("SELECT * FROM place where user_id = ?", userId)
	if err != nil {
		return places, err
	}
	for rows.Next() {
		p := Place{}
		err := rows.Scan(p.Id, p.UserId, p.Title, p.Description, p.Lat, p.Lon, p.DateCreated)
		if err != nil {
			return places, err
		}
		places = append(places, p)
	}
	return places, nil
}

func InsertPlace(p *Place) error {
	r, err := DB.Exec("INSERT INTO place (user_id, title, subtitle, description,  lat, lon, date_created) values (?,?,?,?,?,?,?)", p.UserId, p.Title, p.Subtitle, p.Description, p.Lat, p.Lon, p.DateCreated)
	if err == nil {
		if _, err := r.RowsAffected(); err != nil {
			return err
		}
		val, err := r.LastInsertId()
		if err != nil {
			return err
		}
		p.Id = val
		return nil
	}
	return err
}
