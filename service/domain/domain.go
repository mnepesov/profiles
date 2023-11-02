package domain

const (
	IsAdminCtxKey = "IsAdmin"
)

type Profile struct {
	Id       string
	Username string
	Email    string
	Password string
	IsAdmin  bool
}
