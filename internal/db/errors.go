package db

import "fmt"

var (
	ErrNotFound           = fmt.Errorf("not found")
	ErrLoginAlreadyExists = fmt.Errorf("login already exists")
	ErrChangeAdminLogin   = fmt.Errorf("can't change admin login")
)
