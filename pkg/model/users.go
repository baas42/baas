package model

import "gorm.io/gorm"

// UserRole is an enum which stores the roles a user can have.
type UserRole string

const (
	// User can just use images and change their own image
	User UserRole = "user"
	// Moderator can change or upload system images
	Moderator = "moderator"
	// Admin can do anything on the system
	Admin = "admin"
)

// UserModel (noun) one who uses, not necessarily a single person
type UserModel struct {
	gorm.Model `json:"-"`

	// Name is a human-readable identifier for a user (or entity) of the system
	Username string `gorm:"unique;not null;primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Role     UserRole

	// Images is a list of ImageModel of this user
	Images []ImageModel
}
