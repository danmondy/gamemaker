package data

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

//TODO: add create db script with user and permissions
const CREATE_DB = `
CREATE DATABASE IF NOT EXISTS mass;
GRANT USAGE ON *.* TO massuser@localhost IDENTIFIED BY 'h4h4h3h3h3';
GRANT ALL PRIVELEGES ON mass TO massuser@localhost;
FLUSH PRIVILEGES;
`

func InitDb(dbconf DbConfig) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name))
	if err != nil {
		return err
	}
	DB = db
	createDB()
	return nil
}

//data.Query allows you to pass your args as a separate value
//but will forces mysql to execute it unprepaired.
func Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return DB.Query(fmt.Sprintf(sql, args...))
}

func createDB() {
	/*fmt.Println("Creating db")
	result, err := DB.Exec(CREATE_DB)
	fmt.Println(result, err)*/

	fmt.Println("Creating User Table")
	result, err := DB.Exec("drop table user")
	fmt.Println(result, err)
	result, err = DB.Exec(CREATE_USER_TABLE)
	fmt.Println(result, err)

	fmt.Println("Creating Map Table")
	result, err = DB.Exec(CREATE_MAP_TABLE)
	fmt.Println(result, err)

	fmt.Println("Creating Place Table")
	result, err = DB.Exec(CREATE_PLACE_TABLE)
	fmt.Println(result, err)

	daniel := NewUser("danmondy@gmail.com", "boom")
	err = daniel.Insert()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pw match?:", daniel.PasswordIs("boom"))
	fmt.Println("pw match?:", daniel.PasswordIs("boomer"))

	//mymap := NewMap(*daniel)
	//err = InsertMap(mymap)
	//if err != nil {
	//	fmt.Println(err)
	//}
}
