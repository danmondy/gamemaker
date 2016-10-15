package data

import(
	//"fmt"
	"time"
)

const CREATE_MAP_TABLE = `
CREATE TABLE IF NOT EXISTS map (
id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
user_id INT UNSIGNED,
name VARCHAR(50) NOT NULL UNIQUE,
lon float(10,6) NOT NULL,
lat float(10,6) NOT NULL,
width int NOT NULL,
height int NOT NULL,
style VARCHAR(50),
date_created TIMESTAMP
);
`

var styles = []string{
		"mapbox/emerald-v8",
		"mapbox/satellite-v8",
		"mapbox/dark-v8",
		"mapbox/light-v8",
		"mapbox/satellite-hybrid-v8",
		"mapbox/streets-v8",
		"dmondy/dark",
		"dmondy/cikep085000cx9fm160lx34za",
		"dmondy.effe5575",
	        "examples.a4c252ab",
}

type Map struct{
	Id int64
	UserId int64
	Name string
	Lat, Lon float64
	Width, Height int
	Style string//May change this to a struct of it's own if neessary
	DateCreated time.Time
}

func NewMap(u User)*Map{
	return &Map{
		UserId: u.Id,
		Name: "Default Map",
		Lat:32.80502926732618,
		Lon:-79.94039009094328,
		Width:900,
		Height: 300,
		Style:"examples.a4c252ab",
	}
}


func InsertMap(m *Map) error {
	r, err := DB.Exec("INSERT INTO map (user_id, name, lat, lon, width, height, style) values (?,?,?,?,?,?,?,?)", m.UserId, m.Name, m.Lat, m.Lon, m.Width, m.Height, m.Style, m.DateCreated)
	if err == nil {
		if _, err := r.RowsAffected(); err != nil {
			return err
		}
		val, err := r.LastInsertId()
		if err != nil {
			return err
		}
		m.Id = val
		return nil
	}
	return err
}

func UpdateMap(m *Map) error {
	return nil
}
