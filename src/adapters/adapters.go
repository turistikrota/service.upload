package adapters

import (
	"api.turistikrota.com/upload/src/adapters/mongo"
	"api.turistikrota.com/upload/src/adapters/mysql"
)

var (
	MySQL = mysql.New()
	Mongo = mongo.New()
)
