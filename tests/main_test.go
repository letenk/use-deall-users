package tests

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/letenk/use_deal_user/config"
	"github.com/letenk/use_deal_user/router"
	"go.mongodb.org/mongo-driver/mongo"
)

var ConnTest *mongo.Database
var RouteTest *gin.Engine

func TestMain(m *testing.M) {
	// Open connection
	db := config.SetupDB()
	ConnTest = db

	// Setup router
	RouteTest = router.SetupRouter(db)

	m.Run()
}

/*
export MONGO_USER=root
export MONGO_PASSWORD=root
export MONGO_HOST=localhost
export MONGO_PORT=27017
export MONGO_DBNAME=usedeall
*/
