package repository

import (
	"fmt"
	"testing"

	"github.com/zokypesch/ceria/core"
	"github.com/zokypesch/ceria/helper"

	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

type user struct {
	gorm.Model
	Name string
	Age  string `validate:"required"`
}

type newuser struct {
	gorm.Model
	Name    string
	Age     string   `validate:"required"`
	Credits []credit `gorm:"foreignkey:UserRefer"`
}

type credit struct {
	gorm.Model
	UserRefer int
	Number    string
}

func TestMasterRepository(t *testing.T) {
	db := core.GetTestConnection()

	// set withElastic
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

	// define struct
	newUsr := user{}

	repo := NewMasterRepository(newUsr, db, withElastic)
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	t.Run("GetAll function normal", func(t *testing.T) {

		commonReply := []map[string]interface{}{{"name": "FirstLast", "age": "30"}}
		mocket.Catcher.NewMock().OneTime().WithQuery("SELECT * FROM \"users\"").WithReply(commonReply)
		result, err := repo.GetAll("*")

		assert.Len(t, result, 1)
		assert.NoError(t, err)
	})

	t.Run("GetAll return 0 rows", func(t *testing.T) {

		commonReply := []map[string]interface{}{}
		mocket.Catcher.NewMock().OneTime().WithQuery("SELECT * FROM \"users\"").WithReply(commonReply)
		result, err := repo.GetAll("*")

		assert.Len(t, result, 0)
		assert.NoError(t, err)
	})

	t.Run("GetAll with error", func(t *testing.T) {
		mocket.Catcher.NewMock().OneTime().WithError(fmt.Errorf("Error When select"))

		result, err := repo.GetAll("*")
		assert.Error(t, err)
		assert.Empty(t, result)

	})

	t.Run("Error when Scan", func(t *testing.T) {
		commonReply := []map[string]interface{}{{"age": "30"}}
		mocket.Catcher.NewMock().OneTime().WithQuery(`SELECT * FROM "users"`).WithReply(commonReply)
		result, err := repo.GetAll("*")

		assert.Len(t, result, 0)
		assert.Error(t, err)
	})

	t.Run("Test orm go", func(t *testing.T) {

		commonReply := []map[string]interface{}{
			{
				"name": "FirstLast",
				"age":  "30",
				"credits": []map[string]string{
					map[string]string{"number": "455666678", "user_refer": "1"},
				},
			},
		}
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "newusers"`).WithReply(commonReply)
		newRepo := NewMasterRepository(&newuser{}, db, withElastic)

		// var bindUser []user
		page := &QueryProps{
			WithPagination: true,
			Offset:         0,
			Limit:          10,
			Condition: []map[string]interface{}{
				map[string]interface{}{"value": "ajoi", "field": "name", "operator": "LIKE"},
			},
			PreloadStatus: true,
			Preload: []string{
				"Credits",
			},
		}

		// res, _ := repo.GetAllFromStruct(&bindUser, page)
		res, _ := newRepo.GetAllFromStruct(page)
		assert.Len(t, res, 1)
	})

	t.Run("Test orm passing the empty map params", func(t *testing.T) {
		commonReply := []map[string]interface{}{{"name": "FirstLast", "age": "30"}}
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(commonReply)

		// var bindUser []user

		page := &QueryProps{}

		// res, _ := repo.GetAllFromStruct2(&bindUser, page)
		res, _ := repo.GetAllFromStruct(page)
		assert.Len(t, res, 1)
	})

	t.Run("Test orm set using pagination true but is valid", func(t *testing.T) {
		commonReply := []map[string]interface{}{
			{"name": "FirstLast", "age": "30"},
			{"name": "Ajoi", "age": "75"},
		}
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(commonReply)

		// var bindUser []user

		page := &QueryProps{WithPagination: true}

		// res, _ := repo.GetAllFromStruct(&bindUser, page)
		res, _ := repo.GetAllFromStruct(page)

		assert.Len(t, res, 2)
	})

	t.Run("Test Get All using where failed setting", func(t *testing.T) {
		commonReply := []map[string]interface{}{
			{"name": "FirstLast", "age": "30"},
			{"name": "Ajoi", "age": "75"},
		}
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(commonReply)
		page := &QueryProps{WithPagination: true,
			Condition: []map[string]interface{}{
				map[string]interface{}{"name": "Ajo", "age": "30"},
			}}
		res, _ := repo.GetAllFromStruct(page)

		assert.Len(t, res, 2)
	})
}

func helperlocal(name string, age string, menu int, dbParam *gorm.DB) (*gorm.DB, *MasterRepository) {
	var newUsr user
	var db *gorm.DB
	var repo *MasterRepository
	newUsr = user{}

	// set withElastic
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

	if name != "" && age != "" {
		newUsr = user{
			Name: name,
			Age:  age,
		}
	}

	switch menu {
	case 1:
		db = core.GetTestConnection()
		return db, nil
	case 2:
		repo = NewMasterRepository(newUsr, dbParam, withElastic)
		return nil, repo
	case 3:
		db = core.GetTestConnection()
		repo = NewMasterRepository(newUsr, db, withElastic)
		return db, repo
	}
	return db, repo
}

func TestCreateData(t *testing.T) {
	var err error
	var res int64

	// set withElastic
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

	db, _ := helperlocal("", "", 3, nil)
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	type UserCtm struct {
		gorm.Model
		Name string `json:"name"`
		Age  string `validate:"required" json:"age"`
	}

	t.Run("Tes create data and passing the empty value must be validate", func(t *testing.T) {
		var idRet int64
		idRet = 0

		repo := NewMasterRepository(&UserCtm{Name: "", Age: ""}, db, withElastic)

		res, err = repo.Create(nil)
		assert.Error(t, err)
		assert.Equal(t, idRet, res)
	})

	t.Run("Tes create data positif value", func(t *testing.T) {
		// override
		positiveUsr := user{
			Name: "udin",
			Age:  "10",
		}

		repoNew := NewMasterRepository(&positiveUsr, db, withElastic)

		var mockedID int64
		mockedID = 64

		mocket.Catcher.Reset().NewMock().WithQuery("INSERT INTO \"users\"").WithID(mockedID)

		res, err = repoNew.Create(nil)
		assert.NoError(t, err)
		assert.Equal(t, mockedID, res)
	})

	t.Run("Tes Create when SQL Return error", func(t *testing.T) {
		db.LogMode(true)

		mocket.Catcher.Logging = true

		mocket.Catcher.Reset().NewMock().WithError(fmt.Errorf("Insert SQL Error"))

		positiveUsr := user{
			Name: "udin",
			Age:  "10",
		}
		repoNewFail := NewMasterRepository(&positiveUsr, db, withElastic)

		_, err = repoNewFail.Create(nil)
		assert.Error(t, err)
	})
}
func TestBulkCreate(t *testing.T) {

	// set withElastic
	config := helper.NewReadConfigService()
	config.Init()
	var withElastic *ElasticProperties
	withElastic = &ElasticProperties{}

	confStatus := config.GetByName("elastic.status")
	if confStatus == "true" {
		// elastic configuration
		withElastic = &ElasticProperties{
			Status: false,
			Host:   config.GetByName("elastic.host"),
			Port:   config.GetByName("elastic.port"),
		}
	}

	db := core.GetTestConnection()
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	t.Run("Negative case error when validate", func(t *testing.T) {
		var bulk []user

		bulk = []user{
			user{
				Name: "udin",
				Age:  "",
			},
			user{
				Name: "parjo",
				Age:  "",
			},
		}

		repo := NewMasterRepository(&user{}, db, withElastic)

		res, err := repo.BulkCreate(bulk)
		assert.Len(t, err, 2)
		assert.Len(t, res, 0)
	})

	t.Run("Positif case bulk create", func(t *testing.T) {
		var bulk []user

		bulk = []user{
			user{Model: gorm.Model{ID: 1}, Name: "titl 1", Age: "15"},
			user{
				Model: gorm.Model{
					ID: 2,
				},
				Name: "parjo",
				Age:  "50",
			},
		}

		repo := NewMasterRepository(&user{}, db, withElastic)

		var mockedID int64
		mockedID = 64

		mocket.Catcher.Reset().NewMock().WithQuery(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","age") VALUES ('0001-01-01 00:00:00','0001-01-01 00:00:00',NULL,'udin','40')`).WithID(mockedID)

		rs, err := repo.BulkCreate(&bulk)
		assert.Len(t, err, 0)
		assert.Len(t, rs, 2)
	})
}
func TestUpdateData(t *testing.T) {
	var err error

	// set withElastic
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

	db := core.GetTestConnection()
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	t.Run("Check data invalid update", func(t *testing.T) {
		repo := NewMasterRepository(&user{}, db, withElastic)

		condition := map[string]interface{}{}
		values := map[string]interface{}{}
		err = repo.Update(condition, values)
		// expect error
		assert.Error(t, err)
	})

	t.Run("Check data with positif update", func(t *testing.T) {
		repoNew := NewMasterRepository(&user{}, db, withElastic)

		commonReply := []map[string]interface{}{
			{"name": "Paijo", "age": "30"},
			{"name": "Udin", "age": "50"},
		}
		mocket.Catcher.Reset()
		mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`).WithReply(commonReply)

		condition := map[string]interface{}{
			"Name": "udin",
		}
		values := map[string]interface{}{
			"Age": "70",
		}

		err = repoNew.Update(condition, values)

		assert.NoError(t, err)
	})

	t.Run("Error when SQL error", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithError(fmt.Errorf("Update SQL Error"))

		repo := NewMasterRepository(&user{}, db, withElastic)

		condition := map[string]interface{}{
			"Name": "udin",
		}
		values := map[string]interface{}{
			"Age": "70",
		}

		err = repo.Update(condition, values)
		// expect error
		assert.Error(t, err)
	})

}
func TestGetByField(t *testing.T) {

	// set withElastic
	config := helper.NewReadConfigService()
	config.Init()
	var withElastic *ElasticProperties
	withElastic = &ElasticProperties{}

	confStatus := config.GetByName("elastic.status")
	if confStatus == "true" {
		// elastic configuration
		withElastic = &ElasticProperties{
			Status: false,
			Host:   config.GetByName("elastic.host"),
			Port:   config.GetByName("elastic.port"),
		}
	}

	var err error
	var res []map[string]interface{}

	db := core.GetTestConnection()
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	t.Run("Check Error validation when GetBYField", func(t *testing.T) {

		condition := map[string]interface{}{}
		// var params []user
		repo := NewMasterRepository(&user{}, db, withElastic)

		res, err = repo.GetDataByfield(condition)

		assert.Error(t, err)
		assert.Len(t, res, 0)
	})

	t.Run("Check Positif case GetBYField", func(t *testing.T) {
		condition := map[string]interface{}{
			"age": "30",
		}
		// var params []user
		commonReply := []map[string]interface{}{
			{"name": "Paijo", "age": "30"},
			{"name": "Udin", "age": "50"},
		}
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(commonReply)

		repo := NewMasterRepository(&user{}, db, withElastic)
		res, err = repo.GetDataByfield(condition)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

}
func TestDelete(t *testing.T) {
	var err error

	// set withElastic
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

	db := core.GetTestConnection()
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	t.Run("Check Error validation when Delete", func(t *testing.T) {
		repo := NewMasterRepository(&user{}, db, withElastic)
		mapDel := map[string]interface{}{}

		err = repo.Delete(mapDel)

		assert.Error(t, err)
	})

	t.Run("Check Positif when Delete", func(t *testing.T) {
		repo := NewMasterRepository(&user{}, db, withElastic)
		mapDel := map[string]interface{}{
			"id": 1,
		}

		mocket.Catcher.Reset().NewMock().WithQuery(`DELETE FROM users`)
		err = repo.Delete(mapDel)

		assert.NoError(t, err)
	})

	t.Run("Check Error SQL when Delete", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithError(fmt.Errorf("Delete SQL Error"))

		repo := NewMasterRepository(&user{}, db, withElastic)
		mapDel := map[string]interface{}{
			"id": 1,
		}

		err = repo.Delete(mapDel)

		assert.Error(t, err)
	})
}
func TestBulkDelete(t *testing.T) {
	var err []error
	// set withElastic
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

	db := core.GetTestConnection()
	db.LogMode(true)
	defer db.Close()
	mocket.Catcher.Logging = true

	t.Run("Check Error validation when Delete", func(t *testing.T) {
		repo := NewMasterRepository(&user{}, db, withElastic)
		mapDel := []map[string]interface{}{}

		err = repo.BulkDelete(mapDel)

		assert.Len(t, err, 1)
	})

	t.Run("Check Positif when Delete", func(t *testing.T) {
		repo := NewMasterRepository(&user{}, db, withElastic)
		mapDel := []map[string]interface{}{
			map[string]interface{}{"id": 1},
		}

		mocket.Catcher.Reset().NewMock().WithQuery(`DELETE FROM users`)
		err = repo.BulkDelete(mapDel)

		assert.Len(t, err, 0)
	})

	t.Run("Check Error SQL when Delete", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithError(fmt.Errorf("Delete SQL Error"))

		repo := NewMasterRepository(&user{}, db, withElastic)
		mapDel := []map[string]interface{}{
			map[string]interface{}{"id": 1},
		}

		err = repo.BulkDelete(mapDel)

		assert.Len(t, err, 1)
	})
}

func TestPreload(t *testing.T) {
	db := core.GetTestConnection()
	repo := NewMasterRepository(&newuser{}, db, &ElasticProperties{})

	dbNew := repo.PreloadSetup(db, []string{"Credits"})
	assert.NotEmpty(t, dbNew)
}

func TestParseCondition(t *testing.T) {
	db := core.GetTestConnection()
	repo := NewMasterRepository(&newuser{}, db, &ElasticProperties{})
	t.Run("Nil return when pass the empty map", func(t *testing.T) {
		data := []map[string]interface{}{}

		res, args := repo.ParseConditionToWhere(data)
		assert.Empty(t, res)
		assert.Nil(t, args)
	})

	t.Run("Nil return when pass the not exisiting field in map", func(t *testing.T) {
		data := []map[string]interface{}{
			map[string]interface{}{"field": "age", "value": "30"},
		}

		res, args := repo.ParseConditionToWhere(data)
		assert.Empty(t, res)
		assert.Nil(t, args)
	})

	t.Run("expected return when pass the wrong operator field in map", func(t *testing.T) {
		data := []map[string]interface{}{
			map[string]interface{}{"field": "age", "value": "30", "operator": "WHAT???"},
		}

		res, args := repo.ParseConditionToWhere(data)
		assert.NotEmpty(t, res)
		assert.NotNil(t, args)
	})

	t.Run("expected return when pass the right", func(t *testing.T) {
		data := []map[string]interface{}{
			map[string]interface{}{"field": "age", "value": "30", "operator": "LIKE"},
			map[string]interface{}{"field": "name", "value": "udin", "operator": "EQUAL"},
		}

		res, args := repo.ParseConditionToWhere(data)
		assert.NotEmpty(t, res)
		assert.NotNil(t, args)
	})

	t.Run("expected return when pass the one record not valid", func(t *testing.T) {
		data := []map[string]interface{}{
			map[string]interface{}{"field": "age", "value": "30", "operator": "LIKE"},
			map[string]interface{}{"field": "name", "operator": "EQUAL"},
		}

		res, args := repo.ParseConditionToWhere(data)
		assert.NotEmpty(t, res)
		assert.NotNil(t, args)
	})

	t.Run("Check operator function", func(t *testing.T) {
		assert.NotEmpty(t, repo.CheckOperator("WHAT ??"))
		assert.NotEmpty(t, repo.CheckOperator("equal"))
		assert.NotEmpty(t, repo.CheckOperator("LIKE"))
	})
}
