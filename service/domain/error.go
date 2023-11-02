package domain

import "fmt"

var (
	AlreadyExistError = fmt.Errorf("already exist")
	NotFoundError     = fmt.Errorf("not found")
)
