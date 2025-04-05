package entities

type Task struct {
	ID          uint64
	Name        string
	Description string
}

func NewTask(id uint64, name string, desc string) *Task {
	return &Task{id, name, desc}
}
