package helper

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// HTTPHelper struct of helper Gin
type HTTPHelper struct{}

// HTTPRepoHelper repo for mock Helper
type HTTPRepoHelper interface {
	TestHttpResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) (bool, error)) error
	TestAPI(r http.Handler, method, path string) *httptest.ResponseRecorder
}

// NewServiceHTTPHelper for wrapping the http helper
func NewServiceHTTPHelper() *HTTPHelper {
	return &HTTPHelper{}
}

// TestHTTPResponse for wrapping a errro
func (helper *HTTPHelper) TestHTTPResponse(
	t *testing.T,
	r *gin.Engine,
	req *http.Request,
	f func(w *httptest.ResponseRecorder) (bool, error),
) error {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if ok, err := f(w); !ok {
		return err
	}

	return nil

}

// TestAPI for wrapping manual and return httptest.ResponseRecorder
func (helper *HTTPHelper) TestAPI(
	r http.Handler, method, path string,
	params []byte,
	header map[string]string,
) *httptest.ResponseRecorder {

	if method == "" || path == "" {
		return nil
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	for k, v := range header {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// EchoResponse to echo the response http
func (helper *HTTPHelper) EchoResponse(
	c *gin.Context,
	code int,
	status bool,
	message string,
	err string,
	data interface{},
) {
	c.JSON(
		code,
		gin.H{
			"status":  status,
			"message": message,
			"error":   err,
			"data":    data,
		},
	)
}

// EchoResponseSuccess for Response Success
func (helper *HTTPHelper) EchoResponseSuccess(c *gin.Context, data interface{}) {
	helper.EchoResponse(c, 200, true, "", "", data)
}

// EchoResponseBadRequest for Response Failed
func (helper *HTTPHelper) EchoResponseBadRequest(c *gin.Context, message string, err string) {
	helper.EchoResponse(c, 400, false, message, err, nil)
}

// EchoResponseCreated for Response Success
func (helper *HTTPHelper) EchoResponseCreated(c *gin.Context, data interface{}) {
	helper.EchoResponse(c, 201, true, "", "", data)
}

// EchoResponseWithPagination for Response Success
func (helper *HTTPHelper) EchoResponseWithPagination(c *gin.Context, data interface{}, page string, totalData string) {
	c.JSON(
		200,
		gin.H{
			"status":     true,
			"message":    "",
			"error":      "",
			"data":       data,
			"page":       page,
			"total_data": totalData,
		},
	)
}
