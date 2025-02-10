package models

import "time"

// Model of User
type User struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Email      string `json:"email" db:"email"`
	AvatarPath string `json:"avatar_path" db:"avatar_path"`
	//CreatedAt time.Time `json:"created_at" db:"created_at"`
	//ChangedAt time.Time `json:"changed_at" db:"changed_at"`
}

// Model of Container
type Container struct {
	ID          int       `json:"id" db:"id" validate:"numeric"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description" db:"description"`
	OwnerId     int       `json:"owner_id" db:"owner_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ChangedAt   time.Time `json:"changed_at" db:"changed_at"`
}

// Model of Container in request
type ContainerReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	OwnerId     int    `json:"owner_id" validate:"required"`
}

// TODO:
type name interface {
	getDbMap() map[string]interface{}
}
