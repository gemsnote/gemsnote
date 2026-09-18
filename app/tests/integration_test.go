package tests

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/config"
	"github.com/revel/revel"
)

var integrationOnce sync.Once

func requireMongoIntegration(t *testing.T) {
	t.Helper()
	url := os.Getenv("GEMSNOTE_INTEGRATION_MONGO_URL")
	if url == "" {
		url = os.Getenv("LEANOTE_INTEGRATION_MONGO_URL")
	}
	if url == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_MONGO_URL to run database integration tests")
	}
	integrationOnce.Do(func() {
		basePath, err := filepath.Abs("../..")
		if err != nil {
			t.Fatal(err)
		}
		revel.BasePath = basePath
		revel.Config, err = config.LoadContext("app.conf", []string{filepath.Join(basePath, "conf")})
		if err != nil {
			t.Fatal(err)
		}
		revel.Config.SetOption("db.type", "mongodb")
		db.Init(url, "")
		service.InitService()
		service.ConfigS.InitGlobalConfigs()
	})
}
