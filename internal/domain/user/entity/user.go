package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Password  string     `json:"-"` // Never expose password in JSON
	FullName  string     `json:"full_name"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewUser(email, username, password, fullName, role string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  password,
		FullName:  fullName,
		Role:      role,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) IsActive() bool {
	return u.Status == "active" && u.DeletedAt == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) MarkAsDeleted() {
	now := time.Now()
	u.DeletedAt = &now
	u.Status = "inactive"
	u.UpdatedAt = now
}

func (u *User) UpdateProfile(fullName string) {
	if fullName != "" {
		u.FullName = fullName
	}
	u.UpdatedAt = time.Now()
}

func (u *User) UpdatePassword(hashedPassword string) {
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
}

func (u *User) ChangeStatus(status string) {
	u.Status = status
	u.UpdatedAt = time.Now()
}
