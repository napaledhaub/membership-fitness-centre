package models

import "time"

type Member struct {
	ID                int
	Username          string
	Email             string
	Password          string
	IsVerified        bool
	VerificationToken string
	TokenCreatedAt    time.Time
	PackageID         int
	ExpireDate        time.Time
}
