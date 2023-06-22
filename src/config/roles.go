package config

import "github.com/turistikrota/service.shared/base_roles"

type cdnRoles struct {
	Upload string
}

type roles struct {
	base_roles.Roles
	Cdn cdnRoles
}

var Roles = roles{
	Roles: base_roles.BaseRoles,
	Cdn: cdnRoles{
		Upload: "cdn.upload",
	},
}
