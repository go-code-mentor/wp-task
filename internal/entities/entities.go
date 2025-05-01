package entities

const (
	UserLoginKey = "user"
)

type Task struct {
	ID          uint64
	Name        string
	Description string
	Owner       string
}
