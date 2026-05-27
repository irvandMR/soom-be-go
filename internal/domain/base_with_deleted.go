package domain

import "time"

type BaseModelWithDeleted struct {
	BaseModel
	DeletedAt *time.Time
	DeletedBy *string
}
