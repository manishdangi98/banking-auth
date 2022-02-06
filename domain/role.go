package domain

import "strings"

type RolePermission struct {
	rolePermission map[string][]string
}

func (p RolePermission) IsAuthorizedFor(role string, routeName string) bool {
	perms := p.rolePermission[role]
	for _, r := range perms {
		if r == strings.TrimSpace(routeName) {
			return true
		}
	}
	return false
}

func GetRolePermissions() RolePermission {
	return RolePermission{map[string][]string{
		"admin": {
			"GetAllCustomer", "GetCustomer", "NewAccount", "NewTransaction",
		},
		"user": {
			"GetCustomer", "NewTransaction",
		},
	}}
}
