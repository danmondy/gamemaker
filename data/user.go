package data

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

const CREATE_USER_TABLE = `
CREATE TABLE IF NOT EXISTS user (
id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
email VARCHAR(50) NOT NULL UNIQUE,
hashword varchar(64) NOT NULL,
salt varchar(32) NOT NULL,
admin boolean NOT NULL default 0,
confirmed boolean NOT NULL default 0,
join_date TIMESTAMP
);
`

//DB.Exec(CREATE_USER_TABLE)

type User struct {
	Id        int64
	Email     string
	Hashword  string `json:"-"`
	Salt      string `json:"-"`
	Confirmed bool   `json:"-"`
	Admin     bool   `json:"-"`
	JoinDate  time.Time
}

func NewUser(email string, password string) *User {
	u := &User{Email: email}
	u.SetHashword(password)
	u.JoinDate = time.Now()
	return u
}

func (u *User) PasswordIs(p string) bool {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s%s", p, u.Salt)))
	if err != nil {
		fmt.Println(err)
		return false
	}
	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return hash == u.Hashword
}

func (u *User) SetHashword(p string) error {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	u.Salt = base64.StdEncoding.EncodeToString(salt)
	fmt.Println(u.Salt)

	hasher := sha1.New()
	_, err = hasher.Write([]byte(fmt.Sprintf("%s%s", p, u.Salt)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	u.Hashword = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return err
}

//REPO
func (u *User) Insert() error {
	insert := fmt.Sprintf(`insert into user (email, hashword, salt, admin, confirmed, join_date) values ("%s",  "%s", "%s","%s", %t, "%s")`, u.Email, u.Hashword, u.Salt, u.Admin, u.Confirmed, u.JoinDate)
	r, err := DB.Exec(insert)
	if err == nil {
		if _, err := r.RowsAffected(); err != nil {
			return err
		}
		val, err := r.LastInsertId()
		if err != nil {
			return err
		}
		u.Id = val
		return nil
	}
	return err
}

func (u *User) Update() error {
	update := fmt.Sprintf("pdate user set email = %[1]s, Hashword = %[2]s, Salt = %[3]s, u.Admin = %[5]s, where id = %[4]d", u.Email, u.Hashword, u.Salt, u.Id, u.Admin)
	_, err := DB.Exec(update)
	return err
}

func Authenticate(email string, pass string) (*User, bool) {
	if user, err := GetUserByEmail(email); err == nil {
		if user.PasswordIs(pass) {
			return user, true
		}
	} else {
		fmt.Println(err)
	}
	return nil, false
}

//SQL
func GetUserByEmail(email string) (*User, error) {
	return GetUser("select * from user where email = ?", email)
}

func GetUserById(id int64) (*User, error) {
	return GetUser("select * from user where userid = ?", id)
}

func GetUser(query string, args ...interface{}) (*User, error) {
	u := User{}
	err := DB.QueryRow(query, args...).Scan(&u.Id, &u.Email, &u.Hashword, &u.Salt, &u.Admin, &u.Confirmed, &u.JoinDate)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUsers(query string, args ...interface{}) ([]*User, error) {
	users := []*User{}
	rows, err := DB.Query(query, args)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.Id, &u.Email, &u.Hashword, &u.Salt, &u.Admin, &u.Confirmed, &u.JoinDate)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
