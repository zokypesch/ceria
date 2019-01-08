package ceria

import (
	"strconv"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/zokypesch/ceria/core"
	helper "github.com/zokypesch/ceria/helper"
	repo "github.com/zokypesch/ceria/repository"
	routeService "github.com/zokypesch/ceria/route"
)

type Article struct {
	gorm.DB
	Title    string
	Body     string
	Comments []Comment `gorm:"foreignkey:ArticleID"`
}

type Comment struct {
	gorm.DB
	body      string
	ArticleID uint
}

func TestWrapper(t *testing.T) {

	type User struct {
		Name string
		Age  string
	}

	config := helper.NewReadConfigService()
	config.Init()

	port, _ := strconv.Atoi(config.GetByName("db.port"))

	srv := core.NewServiceConnection(
		config.GetByName("db.driver"),
		config.GetByName("db.host"),
		port,
		config.GetByName("db.user"),
		config.GetByName("db.password"),
		config.GetByName("db.name"),
	)

	initRouter := routeService.NewRouteService(true, "./templates", true)

	myModel := &User{}
	myElastic := &repo.ElasticProperties{
		Status: true,
		Host:   config.GetByName("elastic.host"),
		Port:   config.GetByName("elastic.port"),
	}

	t.Run("Tes Register Model Error", func(t *testing.T) {
		_, errNew := RegisterModel(initRouter, srv, myElastic, nil, &GroupConfiguration{}, &repo.QueryProps{})
		assert.Error(t, errNew)
	})

	t.Run("Tes Register Model non group", func(t *testing.T) {
		_, errNew := RegisterModel(initRouter, srv, myElastic, myModel, &GroupConfiguration{}, &repo.QueryProps{})
		assert.NoError(t, errNew)
	})

	t.Run("Tes register wrong host DB", func(t *testing.T) {
		srvNew := core.NewServiceConnection(
			config.GetByName("db.driver"),
			"192.168.0.1",
			port,
			config.GetByName("db.user"),
			config.GetByName("db.password"),
			config.GetByName("db.name"),
		)

		_, errNew := RegisterModel(initRouter, srvNew, myElastic, myModel, &GroupConfiguration{}, &repo.QueryProps{})
		assert.Error(t, errNew)

	})

	t.Run("Tes Register Model use group", func(t *testing.T) {
		db, errNew := RegisterModel(
			initRouter,
			srv,
			myElastic,
			&Article{},
			&GroupConfiguration{
				Name:       "testaja",
				Middleware: nil,
			},
			&repo.QueryProps{
				Preload:       []string{"Comments"},
				PreloadStatus: true,
			},
		)

		assert.NoError(t, errNew)
		if db != nil {
			db.Close()
		}
	})
}
