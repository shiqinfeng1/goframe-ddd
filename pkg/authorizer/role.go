package authorizer

import "slices"

// Role defines a user access level.
type Role string

const (
	// RoleInvalid 无效的角色，所有操作将被拒绝
	RoleInvalid Role = ""
	// RoleObserver 观察者只能连接到设备并查看设备和会话信息.
	RoleObserver Role = "observer"
	// RoleOperator 具有设备相关的所有权限，除了不能删除设备.
	RoleOperator Role = "operator"
	// RoleAdministrator 除了不能删除命名空间，不能处理付费相关的接口，其他权限同owner
	RoleAdministrator Role = "administrator"
	// RoleOwner 具有命名空间内的所有权限.
	RoleOwner Role = "owner"
)

// RoleFromString 角色名称转换为角色类型
func RoleFromString(str string) Role {
	switch str {
	case "owner":
		return RoleOwner
	case "administrator":
		return RoleAdministrator
	case "operator":
		return RoleOperator
	case "observer":
		return RoleObserver
	default:
		return RoleInvalid
	}
}

// String converts the given role to its corresponding string.
// If the string is not a valid role, it returns "N/A".
func (r Role) String() string {
	switch r {
	case RoleOwner:
		return "owner"
	case RoleAdministrator:
		return "administrator"
	case RoleOperator:
		return "operator"
	case RoleObserver:
		return "observer"
	default:
		return ""
	}
}

// code converts the given role to its corresponding integer.
// If the role is not a valid one, it returns 0.
func (r Role) code() int {
	switch r {
	case RoleOwner:
		return 4
	case RoleAdministrator:
		return 3
	case RoleOperator:
		return 2
	case RoleObserver:
		return 1
	default:
		return 0
	}
}

// Permissions returns all permissions associated with the role r.
// If the role is [RoleInvalid], it returns an empty slice.
func (r Role) Permissions() []Permission {
	permissions := make([]Permission, 0)
	switch r {
	case RoleOwner:
		permissions = ownerPermissions
	case RoleAdministrator:
		permissions = adminPermissions
	case RoleOperator:
		permissions = operatorPermissions
	case RoleObserver:
		permissions = observerPermissions
	}

	return permissions
}

// HasPermission reports whether the role r has the specified permission.
func (r Role) HasPermission(permission Permission) bool {
	return slices.Contains(r.Permissions(), permission)
}

// HasAuthority reports whether the role r has greater or equal authority compared to the passive role.
// It always returns false if either role is invalid or if the passive role is [RoleOwner].
func (r Role) HasAuthority(passive Role) bool {
	return passive != RoleOwner && r.code() >= passive.code()
}
