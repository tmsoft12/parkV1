package util

import (
	modelsuser "park/models/modelsUser"
)

var validRoles = []modelsuser.RoleType{
	modelsuser.AdminRole,
	modelsuser.OperatorRole,
	modelsuser.AccountantRole,
}

func IsValidRole(role modelsuser.RoleType) bool {
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}
