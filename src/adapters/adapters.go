package adapters

import (
	"github.com/turistikrota/service.upload/src/adapters/mongo"
	"github.com/turistikrota/service.upload/src/adapters/mysql"
)

var (
	MySQL = mysql.New()
	Mongo = mongo.New()
)
