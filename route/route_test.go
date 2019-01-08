package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	mock "github.com/zokypesch/ceria/route/mocks"

	helper "github.com/zokypesch/ceria/helper"

	"github.com/stretchr/testify/assert"
)

var path string
var routecfg GinRepo

func TestRoute(t *testing.T) {

	path = "../templates/*"
	routecfg = NewRouteService(false, path, true)
	gin.SetMode(gin.TestMode)

	t.Run("Set Path both condition", func(t *testing.T) {

		err := routecfg.SetPath("../templates/*")

		assert.NoError(t, err)

		errNew := routecfg.SetPath("")
		assert.Error(t, errNew)
	})

	t.Run("Run negative case of router", func(t *testing.T) {

		newMock := new(mock.GinRepo)
		routecfg = NewRouteService(true, "", true)

		newMock.On("Register", false).Return(gin.Default(), fmt.Errorf("Failed execute gin"))

		_, err := routecfg.Register(false)
		_, errEcpect := newMock.Register(false)

		assert.Equal(t, errEcpect, err)
		assert.Error(t, err)

	})

	t.Run("Run positif case type of router same initial of route", func(t *testing.T) {

		newMock := new(mock.GinRepo)
		routecfg = NewRouteService(true, path, true)

		newMock.On("Register", false).Return(gin.Default(), nil)

		positifCfg, err := routecfg.Register(false)
		_, errEcpectPos := newMock.Register(false)

		assert.Equal(t, err, errEcpectPos)
		assert.NotNil(t, positifCfg)
	})

	t.Run("Run with Html", func(t *testing.T) {
		_, err := routecfg.Register(true)

		assert.NoError(t, err)

	})

}

func TestPageRouter(t *testing.T) {
	initRouter := NewRouteService(true, "../templates/*", true)
	newHelper := helper.NewServiceHTTPHelper()
	r, _ := initRouter.Register(true)

	r.GET("/", showIndexPage)
	r.GET("/hello", showIndexPageAPI)

	req, _ := http.NewRequest("GET", "/", nil)

	// Web test
	errWeb := newHelper.TestHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) (bool, error) {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>hello world !</title>") > 0

		if statusOK && pageOK {
			return true, nil
		}
		return false, fmt.Errorf("Error when fetch a title")
	})
	assert.NoError(t, errWeb)

	// Api Test
	w := newHelper.TestAPI(r, "GET", "/hello", nil, nil)
	body := gin.H{
		"title": "hello World",
	}
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["title"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["title"], value)

}

// Example test page
func showIndexPage(context *gin.Context) {
	context.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "hello world !",
		},
	)
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
