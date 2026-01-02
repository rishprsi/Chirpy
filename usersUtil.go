package main

import (
	"Chirpy/internal/database"
)

func DBUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.CreatedAt,
		Email:     dbUser.Email,
	}
}
