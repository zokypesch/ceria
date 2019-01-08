package repository

import (
	"fmt"
	"strconv"

	hlp "github.com/zokypesch/ceria/helper"
	util "github.com/zokypesch/ceria/util"

	"github.com/gin-gonic/gin"
)

// RouteHandler struct for handler gin
type RouteHandler struct {
	rt       *gin.Engine
	repo     *MasterRepository
	httpUtil *hlp.HTTPHelper
	qryProps *QueryProps
}

// RouteHAndlerInterface for interfacing routeHandler
type RouteHAndlerInterface interface {
	PathRegister()
}

// NewServiceRouteHandler for new service handler
func NewServiceRouteHandler(gn *gin.Engine, rp *MasterRepository, qp *QueryProps) *RouteHandler {
	return &RouteHandler{
		rt:       gn,
		repo:     rp,
		httpUtil: hlp.NewServiceHTTPHelper(),
		qryProps: qp,
	}
}

// PathRegister for register path
func (routes *RouteHandler) PathRegister() {
	routes.rt.GET(routes.repo.Tablename, routes.GetAllHandler)
	routes.rt.POST(routes.repo.Tablename, routes.CreateHandler)
	routes.rt.PUT(routes.repo.Tablename+"/:id", routes.UpdateHandler)
	routes.rt.DELETE(routes.repo.Tablename+"/:id", routes.DeleteHandler)
	// additional
	routes.rt.POST(routes.repo.Tablename+"/find", routes.GetDataByfieldHandler)
	routes.rt.POST(routes.repo.Tablename+"/bulkcreate", routes.BulkCreateHandler)
	routes.rt.POST(routes.repo.Tablename+"/bulkdelete", routes.BulkDeleteHandler)
}

// PathRegisterWithMiddleware for register path
func (routes *RouteHandler) PathRegisterWithMiddleware(groupName string, middle func(c *gin.Context)) {
	grp := routes.rt.Group(groupName)
	grp.Use(middle)
	{
		routes.rt.GET(routes.repo.Tablename, routes.GetAllHandler)
		routes.rt.POST(routes.repo.Tablename, routes.CreateHandler)
		routes.rt.PUT(routes.repo.Tablename+"/:id", routes.UpdateHandler)
		routes.rt.DELETE(routes.repo.Tablename+"/:id", routes.DeleteHandler)
		// additional
		routes.rt.POST(routes.repo.Tablename+"/find", routes.GetDataByfieldHandler)
		routes.rt.POST(routes.repo.Tablename+"/bulkcreate", routes.BulkCreateHandler)
		routes.rt.POST(routes.repo.Tablename+"/bulkdelete", routes.BulkDeleteHandler)
	}
}

// GetAllHandler for function Handler get
func (routes *RouteHandler) GetAllHandler(ctx *gin.Context) {
	var page, limit int = 0, 0

	if routes.qryProps.WithPagination {

		pageSetting := ctx.DefaultQuery("page", "1")
		limitSetting := ctx.DefaultQuery("limit", "30")

		pages, errPage := strconv.Atoi(pageSetting)
		limits, errLimit := strconv.Atoi(limitSetting)

		if errPage == nil {
			page = pages
		}

		if errLimit == nil {
			limit = limits
		}

		routes.qryProps.WithPagination = true
		routes.qryProps.Limit = limit
		routes.qryProps.Offset = (page - 1) * limit
	}
	data, err := routes.repo.GetAllFromStruct(routes.qryProps)

	if err != nil {
		routes.httpUtil.EchoResponse(ctx, 400, false, "failed to get data", err.Error(), nil)
		return
	}

	if routes.qryProps.WithPagination {
		routes.httpUtil.EchoResponseWithPagination(ctx, data, strconv.Itoa(page), strconv.Itoa(len(data)))
		return
	}
	routes.httpUtil.EchoResponseSuccess(ctx, data)
}

// CreateHandler for create handler
func (routes *RouteHandler) CreateHandler(ctx *gin.Context) {
	// release new model
	newModel := util.NewServiceStructValue().SetNilValue(routes.repo.Model)

	errGin := ctx.ShouldBindJSON(newModel)
	if errGin != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed create data", errGin.Error())
		return
	}

	_, err := routes.repo.Create(newModel)

	if err != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed create data", err.Error())
		return
	}

	routes.httpUtil.EchoResponseCreated(ctx, newModel)
}

// UpdateHandler for update handler
func (routes *RouteHandler) UpdateHandler(ctx *gin.Context) {

	id := ctx.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed to update", "params must a valid number")
		return
	}

	type paramsMustHave struct {
		Data map[string]interface{} `validate:"required" form:"data" json:"data" binding:"required"`
	}

	var params paramsMustHave

	errGin := ctx.ShouldBindJSON(&params)
	if errGin != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", errGin.Error())
		return
	}

	if len(params.Data) == 0 {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", fmt.Errorf("your data length 0").Error())
		return
	}

	err := routes.repo.Update(map[string]interface{}{"id": id}, params.Data)
	if err != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", err.Error())
		return
	}

	routes.httpUtil.EchoResponseSuccess(ctx, params)

}

// DeleteHandler for delete handler
func (routes *RouteHandler) DeleteHandler(ctx *gin.Context) {

	id := ctx.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed to update", "params must a valid number")
		return
	}

	err := routes.repo.Delete(map[string]interface{}{"id": id})

	if err != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed delete data", err.Error())
		return
	}

	routes.httpUtil.EchoResponseSuccess(ctx, nil)
}

// BulkCreateHandler for bulk create handler
func (routes *RouteHandler) BulkCreateHandler(ctx *gin.Context) {

	utility := util.NewUtilConvertToMap()
	str := utility.RebuildToSlice(routes.repo.Model)

	errGin := ctx.ShouldBindJSON(str)

	if errGin != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", errGin.Error())
		return
	}
	_, err := routes.repo.BulkCreate(str)

	if len(err) > 0 {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", err[0].Error())
		return
	}

	routes.httpUtil.EchoResponseSuccess(ctx, str)
}

// BulkDeleteHandler from bulk delete handler
func (routes *RouteHandler) BulkDeleteHandler(ctx *gin.Context) {

	type paramsMustHave struct {
		// ID int `json:"data" binding:"required"`
		Data []map[string]interface{} `validate:"required" form:"data" json:"data" binding:"required"`
	}

	var params paramsMustHave

	errGin := ctx.ShouldBindJSON(&params)
	if errGin != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", errGin.Error())
		return
	}

	if len(params.Data) == 0 {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", fmt.Errorf("your data length 0").Error())
		return
	}

	err := routes.repo.BulkDelete(params.Data)
	if len(err) > 0 {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", err[0].Error())
		return
	}

	routes.httpUtil.EchoResponseSuccess(ctx, params)

}

// GetDataByfieldHandler from bulk delete handler
func (routes *RouteHandler) GetDataByfieldHandler(ctx *gin.Context) {

	type paramsMustHave struct {
		Condition map[string]interface{} `validate:"required" form:"condition" json:"condition" binding:"required"`
	}

	var params paramsMustHave

	errGin := ctx.ShouldBindJSON(&params)
	if errGin != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", errGin.Error())
		return
	}

	if len(params.Condition) == 0 {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", fmt.Errorf("your data length 0").Error())
		return
	}

	data, err := routes.repo.GetDataByfield(params.Condition)

	if err != nil {
		routes.httpUtil.EchoResponseBadRequest(ctx, "failed update data", err.Error())
		return
	}

	routes.httpUtil.EchoResponseSuccess(ctx, data)
}
