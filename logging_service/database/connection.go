package database

import (
	"fmt"
	"logging_service/config"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateConnectionConfig() {
	var conf = config.GetConfig()
	var connectionString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", conf.Database.DatabaseUsername, conf.Database.DatabasePassword, conf.Database.DatabaseURL, conf.Database.DatabaseName)
	fmt.Println(options.Client().ApplyURI(connectionString))
	err := mgm.SetDefaultConfig(nil, "LoggingService", options.Client().ApplyURI(connectionString))
	if err != nil {
		panic("hi")
	}
}
