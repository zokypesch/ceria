package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zokypesch/ceria/core"
	"github.com/zokypesch/ceria/helper"
	routeService "github.com/zokypesch/ceria/route"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	gorm.Model
	Title  string `validate:"required" form:"title" json:"title" binding:"required"`
	Author string `validate:"required" form:"author" json:"author" binding:"required"`
}

func TestRouteHandler(t *testing.T) {

	db := core.GetTestConnection()

	initRouter := routeService.NewRouteService(true, "../templates/*", true)

	newHelper := helper.NewServiceHTTPHelper()
	r, errRouting := initRouter.Register(true)
	assert.NoError(t, errRouting)

	var response map[string]interface{}

	// check is router is fine
	t.Run("Check general response is fine", func(t *testing.T) {
		r.GET("/hello", showIndexPageAPI)
		// Api Test
		w := newHelper.TestAPI(r, "GET", "/hello", nil, nil)
		body := gin.H{
			"title": "hello World",
		}
		// Assert we encoded correctly,
		// the request gives a 200
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["title"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, body["title"], value)

	})

	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true
	mocket.Catcher.Reset()

	// check status test withElastic
	config := helper.NewReadConfigService()
	config.Init()
	var withElastic *ElasticProperties
	withElastic = &ElasticProperties{}

	confStatus := config.GetByName("elastic.status")
	if confStatus == "true" {
		// elastic configuration
		withElastic = &ElasticProperties{
			Status: true,
			Host:   config.GetByName("elastic.host"),
			Port:   config.GetByName("elastic.port"),
		}
	}

	// myStruct := Example{Title: "Welcome to jungle", Author: "Maulana"}
	myStruct := Example{}
	repo := NewMasterRepository(&myStruct, db, withElastic)

	handl := NewServiceRouteHandler(r, repo, &QueryProps{WithPagination: true})

	handl.PathRegister()

	// fill the query result notification
	commonReply := []map[string]interface{}{
		{"id": 1, "title": "FirstLast", "author": "dodo"},
		{"id": 2, "title": "LastFirst", "author": "udin"},
		{"id": 3, "title": "check kehutanan indonesia", "author": "triadi"},
	}
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT * FROM \"examples\"").WithReply(commonReply)

	// Api Check all functional
	t.Run("Check GetAll Handler", func(t *testing.T) {
		response = nil
		var responseGetAll map[string]interface{}

		w := newHelper.TestAPI(r, "GET", "/examples?page=1&limit=10", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &responseGetAll)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Len(t, responseGetAll, 6)
		assert.True(t, responseGetAll["status"].(bool))
		assert.Len(t, responseGetAll["data"], 3)
	})

	t.Run("Check GetAll Handler with pagination", func(t *testing.T) {
		response = nil
		var responseGetAll map[string]interface{}

		w := newHelper.TestAPI(r, "GET", "/examples?page=test123&limit=cdef", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &responseGetAll)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Len(t, responseGetAll, 6)
		assert.True(t, responseGetAll["status"].(bool))
		assert.Len(t, responseGetAll["data"], 3)
	})

	t.Run("Check failure GetAll Handler but passed with pagination and condition", func(t *testing.T) {
		response = nil
		var responseGetAll map[string]interface{}

		w := newHelper.TestAPI(r, "GET", "/examples?page=1&limit=10&where=title:welcome|author_id:1:LIKE", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &responseGetAll)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Len(t, responseGetAll, 6)

		assert.True(t, responseGetAll["status"].(bool))
		assert.Len(t, responseGetAll["data"], 3)
	})

	t.Run("Check expected GetAll Handler but passed with pagination and condition", func(t *testing.T) {
		response = nil
		var responseGetAll map[string]interface{}

		w := newHelper.TestAPI(r, "GET", "/examples?page=1&limit=10&where=title:welcome:EQUAL|author_id:1:LIKE", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &responseGetAll)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Len(t, responseGetAll, 6)
		assert.True(t, responseGetAll["status"].(bool))
		assert.Len(t, responseGetAll["data"], 3)
	})

	t.Run("Check Failure GetAll Handler params hacker", func(t *testing.T) {
		response = nil
		var responseGetAll map[string]interface{}

		w := newHelper.TestAPI(r, "GET", "/examples?page=1&limit=10&where=title:=SQLBLABLA:EQUAL|futher:=?PSQL:LIKE", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &responseGetAll)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Check Create failure Handler", func(t *testing.T) {
		var responsed map[string]interface{}
		responsed = nil
		jsonParams := map[string]string{"author": "admin"}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &responsed)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.False(t, responsed["status"].(bool))
		assert.NotEmpty(t, responsed["message"].(string))

	})

	t.Run("Check Create Handler", func(t *testing.T) {
		response = nil
		jsonParams := map[string]string{"title": "hi, there this is title", "author": "admin"}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.True(t, response["status"].(bool))
	})

	t.Run("Check Update failure Handler string mode", func(t *testing.T) {
		response = nil
		jsonParams := map[string]interface{}{"data": map[string]string{"author": "maulana"}}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "PUT", "/examples/abcdefg", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))
	})

	t.Run("Check Update failure Handler empty params", func(t *testing.T) {
		response = nil
		jsonParams := map[string]string{}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "PUT", "/examples/1", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))
	})

	t.Run("Check Update failure length of params", func(t *testing.T) {
		response = nil
		jsonParams := map[string]interface{}{"data": map[string]interface{}{}}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "PUT", "/examples/1", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))
	})

	t.Run("Check Update Handler", func(t *testing.T) {
		response = nil
		jsonParams := map[string]interface{}{"data": map[string]string{"author": "maulana"}}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "PUT", "/examples/1", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, response["status"].(bool))
	})

	t.Run("check failure delete handler wrong parameter", func(t *testing.T) {
		w := newHelper.TestAPI(r, "DELETE", "/examples/abcdef", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))
	})

	t.Run("check failure delete handler not pass anything", func(t *testing.T) {
		w := newHelper.TestAPI(r, "DELETE", "/examples", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Check Delete Handler", func(t *testing.T) {
		w := newHelper.TestAPI(r, "DELETE", "/examples/1", nil, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, response["status"].(bool))
		assert.Empty(t, response["message"].(string))
	})

	t.Run("check failure Bulk Create Handler", func(t *testing.T) {
		jsonParams := map[string]interface{}{
			"test": "gile",
		}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples/bulkcreate", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))

	})

	t.Run("Check Bulk Create Handler", func(t *testing.T) {
		jsonParams := []Example{
			Example{Model: gorm.Model{ID: 1}, Title: "titl 1", Author: "triadi"},
			Example{Model: gorm.Model{ID: 2}, Title: "titl 2", Author: "triadi"},
			Example{Model: gorm.Model{ID: 3}, Title: "titl 3", Author: "triadi"},
		}
		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples/bulkcreate", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, response["status"].(bool))
	})

	t.Run("Check Failure Bulk Delete Handler", func(t *testing.T) {
		response = nil
		jsonParams := []map[string]int{}

		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples/bulkdelete", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))

	})

	t.Run("Check Bulk Delete Handler", func(t *testing.T) {
		response = nil

		jsonParams2 := map[string]interface{}{
			"data": []map[string]interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
			},
		}

		jsonValue2, _ := json.Marshal(jsonParams2)

		wd := newHelper.TestAPI(r, "POST", "/examples/bulkdelete", jsonValue2, nil)
		json.Unmarshal([]byte(wd.Body.String()), &response)

		assert.Equal(t, http.StatusOK, wd.Code)
		assert.True(t, response["status"].(bool))
	})

	t.Run("Find Data Failure", func(t *testing.T) {
		response = nil
		jsonParams := map[string]interface{}{
			"condition": map[string]string{},
		}

		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples/find", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.False(t, response["status"].(bool))
		assert.NotEmpty(t, response["message"].(string))
	})

	t.Run("Find Data", func(t *testing.T) {
		response = nil

		jsonParams := map[string]interface{}{
			"condition": map[string]string{
				"title":  "welcome to jungle",
				"author": "admin",
			},
		}

		jsonValue, _ := json.Marshal(jsonParams)

		w := newHelper.TestAPI(r, "POST", "/examples/find", jsonValue, nil)
		json.Unmarshal([]byte(w.Body.String()), &response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, response["status"].(bool))
	})
}

func TestWithMiddleWare(t *testing.T) {
	db := core.GetTestConnection()

	initRouter := routeService.NewRouteService(true, "../templates/*", true)

	r, errRouting := initRouter.Register(true)
	assert.NoError(t, errRouting)

	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true
	mocket.Catcher.Reset()

	// check status test withElastic
	config := helper.NewReadConfigService()
	config.Init()
	var withElastic *ElasticProperties
	withElastic = &ElasticProperties{}

	confStatus := config.GetByName("elastic.status")
	if confStatus == "true" {
		// elastic configuration
		withElastic = &ElasticProperties{
			Status: true,
			Host:   config.GetByName("elastic.host"),
			Port:   config.GetByName("elastic.port"),
		}
	}

	// myStruct := Example{Title: "Welcome to jungle", Author: "Maulana"}
	myStruct := Example{}
	repo := NewMasterRepository(&myStruct, db, withElastic)

	handl := NewServiceRouteHandler(r, repo, &QueryProps{})

	handl.PathRegisterWithMiddleware("test", func(c *gin.Context) {
		fmt.Println("Hii iam example middleware")
		c.Next()
	})

	t.Run("tes register all handler", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		list := handl.RegisterAllHandler()
		assert.NotEmpty(t, list)
	})

	t.Run("tes modified nil params ignore", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		err := handl.ModifiedListHandler(nil)
		assert.Error(t, err)
	})

	t.Run("tes modified success params ignore", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		handl.RegisterAllHandler()
		err := handl.ModifiedListHandler([]string{"getall", "create", "update"})

		assert.Len(t, handl.listHandler, 4)
		assert.NoError(t, err)
	})

	t.Run("tes register url with worng type", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		handl.RegisterAllHandler()

		randomHdlr := handl.listHandler[0]
		fn := randomHdlr["fn"]

		newFn := fn.(func(ctx *gin.Context))

		handl.RegisterURL("WHAT??", "/fbs", newFn)
	})

	t.Run("tes register url with real type", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		handl.RegisterAllHandler()

		randomHdlr := handl.listHandler[0]
		fn := randomHdlr["fn"]

		newFn := fn.(func(ctx *gin.Context))

		handl.RegisterURL("GET", "/str", newFn)
	})

	t.Run("tes register url with real type", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		handl.RegisterAllHandler()

		randomHdlr := handl.listHandler[0]
		fn := randomHdlr["fn"]

		newFn := fn.(func(ctx *gin.Context))

		handl.RegisterURL("GET", "", newFn)
	})

	t.Run("tes register url with group", func(t *testing.T) {
		repo := NewMasterRepository(&myStruct, db, withElastic)
		handl := NewServiceRouteHandler(r, repo, &QueryProps{})

		handl.RegisterAllHandler()

		randomHdlr := handl.listHandler[0]
		fn := randomHdlr["fn"]

		newFn := fn.(func(ctx *gin.Context))

		grp := handl.rt.Group("test123")
		grp.Use()
		{
			handl.RegisterURLFromGroup(grp, "GET", "/strgroup", newFn)
			handl.RegisterURLFromGroup(grp, "WHATS??", "", newFn)
		}

	})

}

// example api test
func showIndexPageAPI(context *gin.Context) {
	context.JSON(
		http.StatusOK,
		gin.H{
			"title": "hello World",
		},
	)
}
