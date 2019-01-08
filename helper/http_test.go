package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	rt "github.com/zokypesch/ceria/route"

	mock "github.com/zokypesch/ceria/helper/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHttpHelper(t *testing.T) {

	var req *http.Request

	httpHelper := NewServiceHTTPHelper()

	newRoute := rt.NewRouteService(true, "../templates/*", true)
	// i expected route is tested so, skip test for route
	r, _ := newRoute.Register(true)
	req, _ = http.NewRequest("GET", "/", nil)

	t.Run("Expected Error when fetch web", func(t *testing.T) {
		// newMock := new(mock.HTTPRepoHelper)

		// anythingType := "func(w *httptest.ResponseRecorder) (bool, error){ return false, fmt.Errorf(\"cannot fetch web\") }"
		// newMock.On("TestHTTPResponse", t, r, req, mockOri.AnythingOfType(anythingType)).Return(fmt.Errorf("cannot fetch web"))
		err := httpHelper.TestHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) (bool, error) {
			return false, fmt.Errorf("cannot fetch web")
		})

		assert.Error(t, err)
	})

	t.Run("Expected Success when fetch web", func(t *testing.T) {
		err := httpHelper.TestHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) (bool, error) {
			return true, nil
		})

		assert.NoError(t, err)
	})
}

func TestHttpHelperManual(t *testing.T) {
	// var req *http.Request

	httpHelper := NewServiceHTTPHelper()

	newRoute := rt.NewRouteService(true, "../templates/*", true)
	r, _ := newRoute.Register(true)
	// req, _ = http.NewRequest("GET", "/hello", nil)

	t.Run("Expected Nil when fetch web manually", func(t *testing.T) {
		newMock := new(mock.HTTPRepoHelper)

		newMock.On("TestHttpResponseManual", r, "", "/hello").Return(nil)
		w := httpHelper.TestAPI(r, "", "/hello", nil, nil)

		assert.Equal(t, newMock.TestHttpResponseManual(r, "", "/hello"), w)
	})

	t.Run("Expected Ok result", func(t *testing.T) {
		w := httpHelper.TestAPI(r, "GET", "/hello", nil, nil)

		assert.NotEmpty(t, w)
	})

}

func TestHttpResponse(t *testing.T) {

	httpHelper := NewServiceHTTPHelper()

	newRoute := rt.NewRouteService(true, "../templates/*", true)
	r, _ := newRoute.Register(true)

	r.GET("/", showIndexPageAPI)
	r.GET("/success", showIndexPageAPISuccess)
	r.GET("/failed", showIndexPageAPIFailed)
	r.POST("/created", showIndexPageAPICreated)
	r.GET("/pagination", showIndexPageAPIPagination)

	w := httpHelper.TestAPI(r, "GET", "/", nil, nil)
	var response map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, response["status"])

	// test http response success
	ws := httpHelper.TestAPI(r, "GET", "/success", nil, nil)

	errs := json.Unmarshal([]byte(ws.Body.String()), &response)

	assert.NoError(t, errs)
	assert.Equal(t, http.StatusOK, ws.Code)

	// test http response failed
	wf := httpHelper.TestAPI(r, "GET", "/failed", nil, nil)

	errf := json.Unmarshal([]byte(wf.Body.String()), &response)

	assert.NoError(t, errf)
	assert.Equal(t, http.StatusBadRequest, wf.Code)

	// test http response created
	wc := httpHelper.TestAPI(r, "POST", "/created", nil, nil)

	errc := json.Unmarshal([]byte(wc.Body.String()), &response)

	assert.NoError(t, errc)
	assert.Equal(t, http.StatusCreated, wc.Code)

	// test http response with pagination
	webPage := httpHelper.TestAPI(r, "GET", "/pagination", nil, nil)

	errps := json.Unmarshal([]byte(webPage.Body.String()), &response)

	assert.NoError(t, errps)
	assert.Equal(t, http.StatusOK, webPage.Code)
}

// example api test
func showIndexPageAPI(context *gin.Context) {
	httpHelper := NewServiceHTTPHelper()

	httpHelper.EchoResponse(context, 200, true, "success", "", nil)
}

func showIndexPageAPICreated(context *gin.Context) {
	httpHelper := NewServiceHTTPHelper()

	httpHelper.EchoResponseCreated(context, nil)
}

func showIndexPageAPISuccess(context *gin.Context) {
	httpHelper := NewServiceHTTPHelper()

	httpHelper.EchoResponseSuccess(context, nil)
}

func showIndexPageAPIFailed(context *gin.Context) {
	httpHelper := NewServiceHTTPHelper()

	httpHelper.EchoResponseBadRequest(context, "failed get", "error description")
}

func showIndexPageAPIPagination(context *gin.Context) {
	httpHelper := NewServiceHTTPHelper()

	httpHelper.EchoResponseWithPagination(context, nil, "1", "10")
}
