package domain

type Permission string

const (
	PermissionRead     Permission = "read"
	PermissionWrite    Permission = "write"
	PermissionPublish  Permission = "publish"
	PermissionModerate Permission = "moderate"
	PermissionDelete   Permission = "delete"
)

func PermissionsFor(role string) []Permission {
	if role == "administrator" {
		return []Permission{PermissionRead, PermissionWrite, PermissionPublish, PermissionModerate, PermissionDelete}
	}
	if role == "visitor" {
		return []Permission{PermissionRead}
	}
	return []Permission{PermissionRead, PermissionWrite}
}
func HasPermission(role string, p Permission) bool {
	for _, v := range PermissionsFor(role) {
		if v == p {
			return true
		}
	}
	return false
}
func RequirePermission(role string, p Permission) error {
	if HasPermission(role, p) {
		return nil
	}
	return permissionError{role: role, permission: p}
}

type permissionError struct {
	role       string
	permission Permission
}

func (e permissionError) Error() string {
	return "permission denied: " + e.role + " cannot " + string(e.permission)
}
func IsAdmin(role string) bool { return role == "administrator" }
func IsRegistered(role string) bool {
	return role == "visitor" || role == "member" || role == "administrator"
}
func CanManageContent(role string) bool { return IsAdmin(role) }
func CanSubmitMessage(role string) bool { return IsRegistered(role) }
func CanCreateBooking(role string) bool { return IsRegistered(role) }
func CanFavorite(role string) bool      { return role == "member" || role == "administrator" }
func ActionsFor(role string) []string {
	out := []string{}
	for _, p := range PermissionsFor(role) {
		out = append(out, string(p))
	}
	return out
}
