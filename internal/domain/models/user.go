package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey" example:"1"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email" example:"user@example.com"`
	FirstName string         `json:"first_name" gorm:"not null" validate:"required,min=2,max=50" example:"John"`
	LastName  string         `json:"last_name" gorm:"not null" validate:"required,min=2,max=50" example:"Doe"`
	Age       int            `json:"age" gorm:"not null" validate:"required,min=1,max=120" example:"30"`
	Active    bool           `json:"active" gorm:"default:true" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50" example:"John"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50" example:"Doe"`
	Age       int    `json:"age" binding:"required,min=1,max=120" example:"30"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50" example:"John"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50" example:"Doe"`
	Age       int    `json:"age" binding:"required,min=1,max=120" example:"30"`
	Active    bool   `json:"active" example:"true"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Age       int       `json:"age" example:"30"`
	Active    bool      `json:"active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ToResponse converts a User model to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Age:       u.Age,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromCreateRequest populates a User from CreateUserRequest
func (u *User) FromCreateRequest(req *CreateUserRequest) {
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Age = req.Age
	u.Active = true // Default to active when creating
}

// FromUpdateRequest populates a User from UpdateUserRequest
func (u *User) FromUpdateRequest(req *UpdateUserRequest) {
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Age = req.Age
	u.Active = req.Active
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdult returns true if the user is 18 or older
func (u *User) IsAdult() bool {
	return u.Age >= 18
}
