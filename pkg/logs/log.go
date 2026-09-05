package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	mu          sync.Mutex
	database    *gorm.DB
	databaseDir string
	databaseDay string
)

type Level string

const (
	Debug Level = "debug"
	Error Level = "error"
)

type Event struct {
	Method   string
	Path     string
	Body     string
	Response string
	Status   int
	Level    Level
	Message  string
}

func Initialize(databaseConfig config.Database) error {
	if databaseConfig.Driver != "sqlite" {
		return fmt.Errorf("unsupported database driver: %s", databaseConfig.Driver)
	}

	cachePath := config.ExpandPath(databaseConfig.Path)
	databaseDir = filepath.Join(filepath.Dir(cachePath), strings.TrimSuffix(filepath.Base(cachePath), filepath.Ext(cachePath)))
	return os.MkdirAll(databaseDir, 0755)
}

func Record(event Event) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UTC()
	err := openDatabase(now)
	if err == nil {
		_ = database.Create(&models.Request{
			Method:   event.Method,
			Path:     event.Path,
			Body:     event.Body,
			Response: event.response(),
			Status:   event.Status,
		}).Error
	}

	switch event.Level {
	case Debug:
		console.Debugf("%s\n", event.Message)
	case Error:
		console.Errorf("%s\n", event.Message)
	}
}

func openDatabase(now time.Time) error {
	day := now.Format("2006-01-02")
	if database != nil && databaseDay == day {
		return nil
	}
	if database != nil {
		sqlDB, err := database.DB()
		if err != nil {
			return err
		}
		err = sqlDB.Close()
		if err != nil {
			return err
		}
		database = nil
	}

	path := filepath.Join(databaseDir, day+".db")
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return err
	}
	err = db.Exec("PRAGMA journal_mode=WAL;").Error
	if err != nil {
		return err
	}
	err = db.Exec("PRAGMA synchronous=NORMAL;").Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Request{})
	if err != nil {
		return err
	}
	database = db
	databaseDay = day
	return nil
}

func (event Event) response() string {
	if event.Response != "" {
		return event.Response
	}
	return event.Message
}
