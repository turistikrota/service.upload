package config

import "github.com/turistikrota/service.shared/base_roles"

type cdnRoles struct {
	Upload string
}

type ownerRoles struct {
	Super        string
	UploadAvatar string
	UploadCover  string
}

type roles struct {
	base_roles.Roles
	Cdn   cdnRoles
	Owner ownerRoles
}

var Roles = roles{
	Roles: base_roles.BaseRoles,
	Cdn: cdnRoles{
		Upload: "cdn.upload",
	},
	Owner: ownerRoles{
		Super:        "owner.super",
		UploadAvatar: "owner.upload.avatar",
		UploadCover:  "owner.upload.cover",
	},
}
