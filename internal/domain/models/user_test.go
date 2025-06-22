package models

import (
	"testing"
	"time"
)

func TestUser_TableName(t *testing.T) {
	var user User
	if user.TableName() != "users" {
		t.Errorf("expected table name 'users', got '%s'", user.TableName())
	}
}

func TestUser_ToResponse(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := user.ToResponse()

	if resp.ID != user.ID || resp.Email != user.Email || resp.FirstName != user.FirstName || resp.LastName != user.LastName || resp.Age != user.Age || resp.Active != user.Active || resp.CreatedAt != user.CreatedAt || resp.UpdatedAt != user.UpdatedAt {
		t.Errorf("ToResponse did not map fields correctly")
	}
}

func TestUser_FromCreateRequest(t *testing.T) {
	req := &CreateUserRequest{
		Email:     "new@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
		Age:       25,
	}
	user := &User{}
	user.FromCreateRequest(req)

	if user.Email != req.Email || user.FirstName != req.FirstName || user.LastName != req.LastName || user.Age != req.Age {
		t.Error("FromCreateRequest did not map fields correctly")
	}
	if !user.Active {
		t.Error("Expected user to be active by default")
	}
}

func TestUser_FromUpdateRequest(t *testing.T) {
	req := &UpdateUserRequest{
		Email:     "updated@example.com",
		FirstName: "John",
		LastName:  "Smith",
		Age:       35,
		Active:    false,
	}
	user := &User{}
	user.FromUpdateRequest(req)

	if user.Email != req.Email || user.FirstName != req.FirstName || user.LastName != req.LastName || user.Age != req.Age || user.Active != req.Active {
		t.Error("FromUpdateRequest did not map fields correctly")
	}
}

func TestUser_GetFullName(t *testing.T) {
	user := &User{FirstName: "John", LastName: "Doe"}
	expected := "John Doe"
	if fullName := user.GetFullName(); fullName != expected {
		t.Errorf("expected full name '%s', got '%s'", expected, fullName)
	}
}

func TestUser_IsAdult(t *testing.T) {
	t.Run("is adult", func(t *testing.T) {
		user := &User{Age: 18}
		if !user.IsAdult() {
			t.Error("expected user to be an adult")
		}
	})

	t.Run("is not adult", func(t *testing.T) {
		user := &User{Age: 17}
		if user.IsAdult() {
			t.Error("expected user not to be an adult")
		}
	})

	t.Run("is adult above 18", func(t *testing.T) {
		user := &User{Age: 50}
		if !user.IsAdult() {
			t.Error("expected user to be an adult")
		}
	})
}
