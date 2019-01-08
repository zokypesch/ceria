package route

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GinCfg struct config
type GinCfg struct {
	config   *gin.Engine
	path     string
	testMode bool
	InitGin  bool
}

// GinRepo Interface for mocking the router
type GinRepo interface {
	Register(withTemplate bool) (*gin.Engine, error)
	SetPath(newPath string) error
}

var cfg *GinCfg

// NewRouteService for run service a route
func NewRouteService(rebuild bool, path string, test bool) *GinCfg {
	if cfg == nil || rebuild {
		cfg = &GinCfg{
			path:     path,
			testMode: test,
		}
	}

	return cfg

}

// SetPath for Change Of path
func (route *GinCfg) SetPath(newPath string) error {
	if newPath == "" {
		route.path = newPath
		return fmt.Errorf("Failed for set")
	}
	return nil
}

// Register for register router
func (route *GinCfg) Register(withTemplate bool) (*gin.Engine, error) {
	var err error
	err = fmt.Errorf("Failed execute gin")

	if route.testMode {
		gin.SetMode(gin.TestMode)
	}

	if route.InitGin {
		return route.config, nil
	}

	if route.path == "" {
		return nil, err
	}

	route.InitGin = true
	err = nil

	route.config = gin.Default()
	if withTemplate {
		route.config.LoadHTMLGlob(route.path)
	}

	return route.config, err
}
