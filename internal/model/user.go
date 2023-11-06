package model

type User struct {
	Login          string `json:"login"`
	HashedPassword []byte `json:"-"`

	Email string `json:"email"`
}
