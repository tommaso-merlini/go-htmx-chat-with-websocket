package types

const UserContextKey = "user"

type AuthenticatedUser struct {
	AuthID     string
	Email      string
	IsLoggedIn bool
}
