package config

import "github.com/turistikrota/service.shared/base_roles"

type cdnRoles struct {
	Upload string
}

type ownerRoles struct {
	Super         string
	UploadAvatar  string
	UploadCover   string
	ListingUpdate string
	ListingCreate string
}

type roles struct {
	base_roles.Roles
	Cdn      cdnRoles
	Business ownerRoles
}

var Roles = roles{
	Roles: base_roles.BaseRoles,
	Cdn: cdnRoles{
		Upload: "cdn.upload",
	},
	Business: ownerRoles{
		Super:         "business.super",
		UploadAvatar:  "business.upload.avatar",
		UploadCover:   "business.upload.cover",
		ListingUpdate: "listing.update",
		ListingCreate: "listing.delete",
	},
}
