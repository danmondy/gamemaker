package data

import (
//"time"
)

type Design struct {
	Id           int           `json:"id"`
	H            int           `json:"h"`
	W            int           `json:"w"`
	Name         string        `json:"name"`
	Img          string        `json:"img"`
	ImgBack      string        `json:"imgBack"`
	FieldLayouts []FieldLayout `json:"fieldLayouts"`
}

type Card struct {
	DesignId int      `json:"designId"`
	Title    string   `json:"title"`
	Fields   []string `json:"fields"`
	Images   []Image  `json:"images"`
	Active   bool     `json:"active"`
	Quantity int      `json:"quantity"`
}

type FieldLayout struct {
	Font     string `json:"font"`
	FontSize int    `json:"fontSize"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	W        int    `json:"w"`
	H        int    `json:"h"`
}

type Image struct {
	File string `json:"file"`
	W    int    `json:"w"`
	H    int    `json:"h"`
	Y    int    `json:"y"`
	X    int    `json:"x"`
}

type Config struct {
	Host         string   `json:"host"`
	Port         string   `json:"port"`
	Db           DbConfig `json:"db"`
	Docroot      string   `json:"docroot"`
	FileLocation string   `json:"fileLocation"`

	//OtherDb data.DbConfig `json:"otherdb,omitempty"`
}

type DbConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
