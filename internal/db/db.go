package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/petara94/mib-go/internal/pkg"
	"github.com/petara94/mib-go/internal/pkg/thash"
)

type User struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CheckPass bool   `json:"check_pass"`
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

	mu     sync.Mutex
	Image  Image
	logins []string
}

func NewDB(filename string, password string) *DB {
	if filename == "" {
		log.Fatal("filename is empty")
	}

	var db = &DB{
		Filename: filename,
		Password: password,
		Image: Image{
			Users: make(map[string]User),
		},
	}

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

	d.logins = append(d.logins, user.Login)
	sort.Strings(d.logins)

	return d.Save()
}

// GetWithPagination returns users with pagination
func (d *DB) GetWithPagination(page int, limit int) ([]*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if page < 1 {
		return nil, fmt.Errorf("page must be greater than 0 (page: %d)", page)
	}

	if limit < 1 {
		return nil, fmt.Errorf("limit must be greater than 0 (limit: %d)", limit)
	}

	var users []*User

	from := (page - 1) * limit
	to := page * limit

	if from > len(d.logins) {
		return []*User{}, nil
	}

	if to > len(d.logins) {
		to = len(d.logins)
	}

	for _, login := range d.logins[from:to] {
		u := d.Image.Users[login]
		users = append(users, &u)
	}

	return users, nil
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

		d.logins = append(d.logins, user.Login)
		sort.Strings(d.logins)
	}

	d.Image.Users[user.Login] = user

	return &user, d.Save()
}

func (d *DB) Delete(login string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.Image.Users, login)

	for i, l := range d.logins {
		if l == login {
			d.logins = append(d.logins[:i], d.logins[i+1:]...)
			break
		}
	}
	sort.Strings(d.logins)

	return d.Save()
}

func Open(filename, password string) (*DB, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var jsonData = string(data)

	if password != "" {
		jsonData, err = thash.DecryptWithPass(string(data), password)
		if err != nil {
			return nil, fmt.Errorf("can't decrypt data: %w", err)
		}
	}

	db := &DB{
		Filename: filename,
		Image:    Image{},
		Password: password,
	}

	err = json.Unmarshal([]byte(jsonData), &db.Image)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %w", err)
	}

	// check is admin exists
	_, ok := db.Image.Users["admin"]
	if !ok {
		return nil, fmt.Errorf("db file is corrupted: admin user not found")
	}

	// create logins list
	db.logins = make([]string, 0, len(db.Image.Users))
	for login := range db.Image.Users {
		db.logins = append(db.logins, login)
	}

	sort.Strings(db.logins)

	return db, nil
}

func (d *DB) Save() error {
	data, err := json.Marshal(&d.Image)
	if err != nil {
		return err
	}

	var outData = string(data)

	if d.Password != "" {
		outData, err = thash.EncryptWithPass(string(data), d.Password)
		if err != nil {
			return fmt.Errorf("can't encrypt data: %w", err)
		}
	}

	err = os.WriteFile(d.Filename, []byte(outData), 0666)
	if err != nil {
		return fmt.Errorf("can't write file: %w", err)
	}

	return nil
}
