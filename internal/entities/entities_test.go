package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	id := uint64(1)
	name := "task name"
	desc := "task desc"

	newTask := NewTask(id, name, desc)

	assert.Equal(t, id, newTask.ID)
	assert.Equal(t, name, newTask.Name)
	assert.Equal(t, desc, newTask.Description)
}
