package main

import (
	"Chirpy/internal/database"
)

func DBUserToUser(dbUser database.User, token string, refreshToken string) User {
	return User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.CreatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
}
