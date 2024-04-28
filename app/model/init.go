package model

import (
	database "github.com/ZiyanK/service-catalog-api/app/db"
	"github.com/ZiyanK/service-catalog-api/app/logger"
)

var (
	log = logger.CreateLogger()
	db  = database.GetDBInstance()
)
