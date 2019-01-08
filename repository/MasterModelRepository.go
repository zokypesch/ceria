package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	util "github.com/zokypesch/ceria/util"

	"github.com/jinzhu/gorm"
	core "github.com/zokypesch/ceria/core"
)

// QueryProps struct for property getall
type QueryProps struct {
	WithPagination bool
	Offset         int `default:"0"`
	Limit          int `default:"-1"`
	Condition      map[string]interface{}
	PreloadStatus  bool
	Preload        []string
}

//ElasticProperties for elastic configuration
type ElasticProperties struct {
	Status bool
	Host   string
	Port   string
}

// MasterRepository for all model
type MasterRepository struct {
	Model       interface{}
	Conn        *gorm.DB
	Tablename   string
	trx         *gorm.DB
	WithElastic bool
	coreElastic *core.ElasticCore
}

// MasterRepositoryInterface for interfacing repository
type MasterRepositoryInterface interface {
	GetAllFromStruct(str interface{}, props *QueryProps) ([]map[string]interface{}, error)
	GetAll(selectQRY string) ([]map[string]string, error)
	Create(override interface{}) (int64, error)
	Update(condition map[string]interface{}, data map[string]interface{}) error
	BulkCreate(param interface{}) ([]int64, []error)
	Delete(condition map[string]interface{}) error
	BulkDelete(condition []map[string]interface{}) []error
	GetDataByfield(str interface{}, list map[string]interface{}) ([]map[string]interface{}, error)
}

// NewMasterRepository new master repository
func NewMasterRepository(model interface{}, db *gorm.DB, el *ElasticProperties) *MasterRepository {
	// return repo

	paramElastic := false

	var err error
	var coreEl *core.ElasticCore

	if el.Status {
		coreEl, err = core.NewServiceElasticCore(model, fmt.Sprintf("http://%s:%s", el.Host, el.Port))
	}

	if err == nil && el.Status {
		paramElastic = true
	}

	structName := util.NewServiceStructValue().GetNameOfStruct(model)
	return &MasterRepository{
		Model:       model,
		Conn:        db,
		Tablename:   fmt.Sprintf("%ss", strings.ToLower(structName)),
		trx:         nil,
		WithElastic: paramElastic,
		coreElastic: coreEl,
	}
}

// GetAllFromStruct for get using struct it self
func (repo *MasterRepository) GetAllFromStruct(props *QueryProps) ([]map[string]interface{}, error) {
	var listArr []map[string]interface{}
	var res *gorm.DB
	var strI util.StructValueInterface

	strI = util.NewServiceStructValue()

	utility := util.NewUtilConvertToMap()
	str := utility.RebuildToSlice(repo.Model)

	newDB := repo.Conn

	if props.PreloadStatus && len(props.Preload) > 0 {
		newDB = repo.PreloadSetup(newDB, props.Preload)
	}

	if props.WithPagination {
		strI.SetDefaultValueStruct(props)
		res = newDB.Where(props.Condition).Limit(props.Limit).Offset(props.Offset).Find(str)
	} else {
		res = newDB.Find(str)
	}

	listArr = utility.ConvertMultiStructToMap(str)

	return listArr, res.Error
}

// GetAll for getall data
func (repo *MasterRepository) GetAll(selectQRY string) ([]map[string]string, error) {

	var res []map[string]string

	var rows *sql.Rows
	var err error

	rows, err = repo.Conn.Table(repo.Tablename).Select(selectQRY).Rows()

	if err != nil {
		return nil, err
	}

	// define utility
	utility := util.NewUtilConvertToMap()
	utilityGeneral := util.GeneralUtilService()

	// defined a column
	columns := utility.ConvertInterfaceToKeyStr(repo.Model)
	ok, _ := utilityGeneral.InArray("Model", columns)
	count := len(columns)
	if ok {
		columns = columns[1:count]
		count = count - 1
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := utility.ConvertToDynamicMap(columns, values)
		res = append(res, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, rows.Err()

}

// Create for create data
func (repo *MasterRepository) Create(override interface{}) (int64, error) {
	var ID int64
	ID = 0
	var realModel interface{}
	var values map[string]interface{}

	realTrans := repo.Conn
	isOver := false
	realModel = repo.Model
	typeOver := reflect.TypeOf(override)

	if typeOver != nil && (reflect.TypeOf(override).Kind() == reflect.Struct || reflect.TypeOf(override).Kind() == reflect.Ptr) {
		// change the type struct from interface to real interfalce
		realModel = override
		isOver = true
	}

	if isOver && repo.trx != nil {
		realTrans = repo.trx
	}

	utils := util.NewUtilService(realModel)
	err := utils.Validate()

	if err != nil {
		return ID, err
	}

	cr := realTrans.Create(realModel)

	typModel := reflect.TypeOf(realModel)

	if typModel == nil || (typModel.Kind() != reflect.Struct && typModel.Kind() != reflect.Ptr) {
		return 0, fmt.Errorf("failed")
	}

	cvrt := util.NewUtilConvertToMap()

	values = cvrt.ConvertStructToSingeMap(realModel)

	if _, ok := values["Model"]; !ok {
		return ID, cr.Error
	}

	modelData, _ := cvrt.ConvertInterfaceMaptoMap(values["Model"])
	if val, ok := modelData["ID"]; ok {
		ID, err = strconv.ParseInt(val, 0, 64)
	}

	// create in elastic
	if repo.WithElastic {
		repo.coreElastic.AddDocument(modelData["ID"], realModel)
	}

	return ID, cr.Error
}

// Update for update data to database
func (repo *MasterRepository) Update(condition map[string]interface{}, data map[string]interface{}) error {
	if len(condition) == 0 || len(data) == 0 {
		return fmt.Errorf("cannot empty param")
	}

	newModel := util.NewServiceStructValue().SetNilValue(repo.Model)

	newUtil := util.NewUtilConvertToMap()
	sliceModel := newUtil.RebuildToSlice(repo.Model)

	cr := repo.Conn.Model(newModel).Where(condition).Updates(data)

	var errSlice []error
	// update in elastic to be continued
	if repo.WithElastic && cr.Error == nil {
		repo.Conn.Where(condition).Find(sliceModel)
		valueSlice := reflect.ValueOf(sliceModel).Elem()

		for i := 0; i < valueSlice.Len(); i++ {
			st := valueSlice.Index(i)

			stValue := st.Elem()

			mdl := stValue.FieldByName("Model")

			errElastic := repo.coreElastic.EditDocument(newUtil.ConvertDataToString(mdl.Field(0).Interface()), st.Interface())
			errSlice = append(errSlice, errElastic)
		}
	}

	if len(errSlice) > 0 {
		return errSlice[0]
	}

	return cr.Error
}

// BulkCreate for bulk insert into database
func (repo *MasterRepository) BulkCreate(param interface{}) ([]int64, []error) {
	var err []error
	var res []int64

	switch reflect.TypeOf(param).Kind() {
	case reflect.Slice, reflect.Ptr:
		var s reflect.Value

		if reflect.TypeOf(param).Kind() == reflect.Ptr {
			s = reflect.ValueOf(param).Elem()
		} else if reflect.TypeOf(param).Kind() == reflect.Slice {
			s = reflect.ValueOf(param)
		}

		tx := repo.Conn.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		repo.trx = tx

		for i := 0; i < s.Len(); i++ {
			st := s.Index(i)
			realModel := st.Interface()

			if st.Kind() == reflect.Struct {
				// condition force access from model and the member value is not pointer / addressing
				conValues := reflect.ValueOf(st).Interface() // get interface value of st
				stAddr := conValues.(reflect.Value).Addr()   // st.Addr() // get address like &Model{}
				realModel = stAddr.Interface()               // get interface of address st
			}
			// ada masalah disini cuy periksa ya stuck waktu pake elastic khusus test aja sih

			rs, er := repo.Create(realModel)
			if rs > 0 {
				res = append(res, rs)
			}
			if er != nil {
				err = append(err, er)
			}
		}
		if len(err) > 0 {
			tx.Rollback()
		}
		tx.Commit()

		repo.trx = nil

	default:
		err = []error{
			fmt.Errorf("failed value"),
		}
	}

	return res, err
}

// Delete for delete data
func (repo *MasterRepository) Delete(condition map[string]interface{}) error {
	// validate
	if len(condition) == 0 {
		return fmt.Errorf("Empty data")
	}

	realTrans := repo.Conn

	if repo.trx != nil {
		realTrans = repo.trx
	}

	newModel := util.NewServiceStructValue().SetNilValue(repo.Model)

	newUtil := util.NewUtilConvertToMap()
	sliceModel := newUtil.RebuildToSlice(newModel)
	repo.Conn.Where(condition).Find(sliceModel)

	exec := realTrans.Where(condition).Delete(newModel)

	var errSlice []error
	// update in elastic to be continued
	if repo.WithElastic && exec.Error == nil {
		valueSlice := reflect.ValueOf(sliceModel).Elem()
		for i := 0; i < valueSlice.Len(); i++ {
			st := valueSlice.Index(i)

			stValue := st.Elem()
			mdl := stValue.FieldByName("Model")
			newID := newUtil.ConvertDataToString(mdl.Field(0).Interface())

			repo.coreElastic.DeleteDocument(newID)

			// if errElastic != nil && errElastic.Error() != "elastic: Error 404 (Not Found)" {
			// 	errSlice = append(errSlice, errElastic)
			// }

			// fmt.Println(errElastic.Error() != "elastic: Error 404 (Not Found)")

		}
	}

	if len(errSlice) > 0 {
		return errSlice[0]
	}

	return exec.Error
}

// BulkDelete for bulk delete
func (repo *MasterRepository) BulkDelete(condition []map[string]interface{}) []error {
	var err []error

	// validate
	if len(condition) == 0 {
		err = append(err, fmt.Errorf("data empty"))
		return err
	}

	tx := repo.Conn.Begin()
	repo.trx = tx
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, v := range condition {
		er := repo.Delete(v)
		if er != nil {
			err = append(err, er)
		}
	}

	if len(err) > 0 {
		tx.Rollback()
	}

	tx.Commit()
	repo.trx = nil
	return err
}

// GetDataByfield get data by field
func (repo *MasterRepository) GetDataByfield(list map[string]interface{}) ([]map[string]interface{}, error) {
	var res []map[string]interface{}

	// validate
	if len(list) == 0 {
		return res, fmt.Errorf("Empty params")
	}

	str := util.NewUtilConvertToMap().RebuildToSlice(repo.Model)

	tx := repo.Conn.Where(list).Find(str)
	utils := util.NewUtilConvertToMap()

	res = utils.ConvertMultiStructToMap(str)

	return res, tx.Error
}

// PreloadSetup for adding preload configuration
func (repo *MasterRepository) PreloadSetup(db *gorm.DB, params []string) *gorm.DB {
	newDB := db
	for _, v := range params {
		newDB = newDB.Preload(v)
	}
	return newDB
}
