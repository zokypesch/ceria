package ceria

import (
	"fmt"

	"github.com/zokypesch/ceria/core"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	repo "github.com/zokypesch/ceria/repository"
	routeService "github.com/zokypesch/ceria/route"
)

// ConfigurationStruct struct for configuration
type ConfigurationStruct struct {
	MasterGin *gin.Engine
	MasterDB  *gorm.DB
}

// GroupConfiguration for model setup
type GroupConfiguration struct {
	Name       string
	Middleware func(c *gin.Context)
}

// RegisterModel for registration  your model
func RegisterModel(
	initRouter *routeService.GinCfg,
	conn *core.Connection,
	elastic *repo.ElasticProperties,
	model interface{},
	grp *GroupConfiguration,
	qyp *repo.QueryProps,
) (*gorm.DB, error) {

	if model == nil {
		return nil, fmt.Errorf("model cannot nil")
	}

	r, errRouting := initRouter.Register(false)
	if errRouting != nil {
		return nil, errRouting
	}

	db, errDB := conn.GetConn()

	if errDB != nil {
		return nil, errDB
	}

	rp := repo.NewMasterRepository(model, db, elastic)

	handl := repo.NewServiceRouteHandler(r, rp, qyp)

	if grp.Name != "" {
		handl.PathRegisterWithMiddleware(grp.Name, grp.Middleware)
		return db, nil
	}

	handl.PathRegister()

	return db, nil
}
