package user

import "../../core"

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.User) map[string]interface{} {
	return map[string]interface{}{
		"user_id":         u.ID,
		"user_login":      u.Login,
		"user_email":      u.Email,
		"user_admin":      u.Admin,
		"user_active":     u.Active,
		"user_avatar":     u.Avatar,
		"user_created":    u.Created,
		"user_updated":    u.Updated,
		"user_last_login": u.LastLogin,
	}
}
