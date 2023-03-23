package db

import (
	"encoding/json"
	"log"
	"mib-go/internal/pkg"
	"mib-go/internal/pkg/thash"
	"os"
	"sync"
)

type User struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsBlocked bool   `json:"is_blocked"`
}

func (u *User) SetPassword(password string) {
	u.Password = pkg.HashPassword(password)
}

func (u *User) PasswordEqual(password string) bool {
	return pkg.PasswordEqual(password, u.Password)
}

type Image struct {
	Users map[string]User `json:"users"`
}

type DB struct {
	Filename string
	Password string

	mu    sync.Mutex
	Image Image
}

func NewDB(filename string, password string) *DB {
	var db *DB

	// check is file exists
	_, err := os.Stat(filename)
	if err == nil {
		db, err = Open(filename, password)
		if err != nil {
			log.Fatal(err)
		}

		bakData, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		// create backup
		err = os.WriteFile(filename+"~", bakData, 0666)
		if err != nil {
			log.Fatal(err)
		}

		return db
	}

	err = db.Create(User{Login: "admin"})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (d *DB) Create(user User) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.Image.Users[user.Login]
	if ok {
		return ErrLoginAlreadyExists
	}

	d.Image.Users[user.Login] = user

	return d.Save()
}

func (d *DB) GetByLogin(login string) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	user, ok := d.Image.Users[login]
	if !ok {
		return nil, ErrNotFound
	}

	return &user, nil
}

func (d *DB) Update(login string, user User) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.Image.Users[login]
	if !ok {
		return nil, ErrNotFound
	}

	if login == "admin" && user.Login != login {
		return nil, ErrChangeAdminLogin
	}

	if login != user.Login {
		delete(d.Image.Users, login)
	}

	d.Image.Users[user.Login] = user

	return &user, d.Save()
}

func (d *DB) Delete(login string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.Image.Users, login)

	return d.Save()
}

func Open(filename, password string) (*DB, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	jsonData, err := thash.DecryptWithPass(string(data), password)
	if err != nil {
		return nil, err
	}

	db := &DB{
		Filename: filename,
		Image:    Image{},
		Password: password,
	}

	err = json.Unmarshal([]byte(jsonData), &db.Image)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Save() error {
	data, err := json.Marshal(&d.Image)
	if err != nil {
		return err
	}

	encrypted, err := thash.EncryptWithPass(string(data), d.Password)
	if err != nil {
		return err
	}

	err = os.WriteFile(d.Filename, []byte(encrypted), 0666)
	if err != nil {
		return err
	}

	return nil
}
