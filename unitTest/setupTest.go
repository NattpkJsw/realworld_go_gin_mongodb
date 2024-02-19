package unittest

import (
	"encoding/json"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/servers"
	"github.com/NattpkJsw/real-world-api-go/pkg/databases"
)

func SetupTest() servers.IModulefactory {
	cfg := config.LoadConfig("../.env.dev")

	db := databases.DbConnect(cfg.Db())

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
