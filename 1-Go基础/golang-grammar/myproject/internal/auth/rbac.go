package auth

// Role 常量
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// HasPermission 检查角色是否拥有指定权限
func HasPermission(role, permission string) bool {
	permissions := rolePermissions[role]
	for _, p := range permissions {
		if p == permission || p == "*" {
			return true
		}
	}
	return false
}

var rolePermissions = map[string][]string{
	RoleUser: {
		"user:read:self",
		"order:create",
		"order:read:self",
		"order:pay:self",
	},
	RoleAdmin: {
		"*", // 管理员拥有所有权限
	},
}