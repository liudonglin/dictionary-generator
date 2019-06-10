package core

// User ...
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
	Active    bool   `json:"active"`
	Avatar    string `json:"avatar"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	LastLogin string `json:"last_login"`
}

// UserStore ...
type UserStore interface {
	// Count returns a count of active users.
	Count() (int64, error)

	// Create persists a new user to the datastore.
	Create(*User) error

	FindLogin(login string) (*User, error)
}
