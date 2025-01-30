package models

import (
	"time"
)

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
)

type Task struct {
	UniqueId    string    `bson:"uniqueId,omitempty"`
	Title       string    `bson:"title,omitempty"`
	Description string    `bson:"description,omitempty"`
	Status      Status    `bson:"status,omitempty"`
	CreatedAt   time.Time `bson:"createdAt,omitempty"`
}

type TaskPayload struct {
	Title       string `bson:"title" validate:"required"`
	Description string `bson:"description" validate:"required"`
	Status      Status `bson:"status" validate:"required,oneof=todo in_progress done"`
}

type TaskStatusPayload struct {
	UniqueId string `bson:"uniqueId" validate:"required"`
	Status   Status `bson:"status" validate:"required,oneof=todo in_progress done"`
}
