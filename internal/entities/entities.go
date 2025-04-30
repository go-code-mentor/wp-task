package entities

type Task struct {
	ID          uint64
	Name        string
	Description string
	Owner       string
}

type User struct {
	ID    uint64
	Login string
}

type AccessToken struct {
	ID     uint64
	UserId uint64
	Token  string
}
