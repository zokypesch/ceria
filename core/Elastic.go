package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/zokypesch/ceria/util"

	"github.com/olivere/elastic"
)

// ElasticCore for type elastic
type ElasticCore struct {
	Model  interface{}
	Index  string
	Client *elastic.Client
}

// ElasticCoreInter for interfacing function elastic
type ElasticCoreInter interface{}

// NewServiceElasticCore for new service elastic core
func NewServiceElasticCore(model interface{}, serverURL string) (*ElasticCore, error) {
	structName := util.NewServiceStructValue().GetNameOfStruct(model)
	idx := fmt.Sprintf("%ss", strings.ToLower(structName))
	client, err := Register(model, serverURL, idx)

	if err != nil {
		return nil, err
	}

	return &ElasticCore{
		Model:  model,
		Index:  idx,
		Client: client,
	}, nil
}

// Register for registration elastic function
func Register(model interface{}, serverURL string, index string) (*elastic.Client, error) {

	modelValue := reflect.Indirect(reflect.ValueOf(model)).Interface()
	ctx := context.Background()

	if reflect.TypeOf(modelValue).Kind() != reflect.Struct || serverURL == "" {
		return nil, fmt.Errorf("model must be a struct, connection cannot be empty")
	}

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(serverURL))

	if err != nil {
		return nil, err
	}

	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		return nil, err
	}

	if !exists {
		// Create an index
		client.CreateIndex(index).Do(ctx)
		// _, err = client.CreateIndex(index).Do(ctx)
		// if err != nil {
		// 	// Handle error
		// 	return nil, err
		// }

	}

	return client, nil
}

// AddDocument for add document and indexing elastic
func (elasticCore *ElasticCore) AddDocument(ID string, bodyJSON interface{}) error {

	maps := util.NewUtilConvertToMap().ConvertStructToSingeMap(bodyJSON)

	if len(maps) == 0 || ID == "" {
		return fmt.Errorf("value or ID cannot be null")
	}

	_, err := elasticCore.Client.Index().
		Index(elasticCore.Index).
		Type("doc").
		Id(ID).
		BodyJson(bodyJSON).
		Refresh("wait_for").
		Do(context.Background())

	if err != nil {
		// Handle error
		return err
	}

	return nil

}

// EditDocument for add document and indexing elastic
func (elasticCore *ElasticCore) EditDocument(ID string, bodyJSON interface{}) error {

	maps := util.NewUtilConvertToMap().ConvertStructToSingeMap(bodyJSON)
	if len(maps) == 0 || ID == "" {
		return fmt.Errorf("value or ID cannot be null")
	}

	_, err := elasticCore.Client.Index().
		Index(elasticCore.Index).
		Type("doc").
		Id(ID).
		BodyJson(bodyJSON).
		Refresh("wait_for").
		Do(context.Background())

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// DeleteDocument for delete th elastic search
func (elasticCore *ElasticCore) DeleteDocument(ID string) error {

	if ID == "" {
		return fmt.Errorf("value or ID cannot be null")
	}

	_, err := elasticCore.Client.Delete().Index(elasticCore.Index).Type("doc").Id(ID).Do(context.TODO())

	return err

}

// DeleteIndex for delete th elastic search
func (elasticCore *ElasticCore) DeleteIndex() error {
	// Delete an index.
	_, err := elasticCore.Client.DeleteIndex(elasticCore.Index).Do(context.TODO())

	if err != nil {
		// Handle error
		return err
	}

	return nil
}

// MultipleinsertDocumentByStruct for struct inheretence
func (elasticCore *ElasticCore) MultipleinsertDocumentByStruct(IDParams string, str interface{}) error {
	var field reflect.Value

	switch reflect.ValueOf(str).Kind() {
	case reflect.Ptr:
		field = reflect.ValueOf(str).Elem()
	case reflect.Struct:
		field = reflect.ValueOf(str)
	default:
		return nil
	}

	for i := 0; i < field.NumField(); i++ {
		fieldValue := field.Field(i)
		fieldName := util.NewUtilConvertToMap().ConvertDataToString(field.Type().Field(i).Name)

		if fieldValue.Kind() == reflect.Invalid || fieldName == "Model" {
			newValueOfInvalid := reflect.New(fieldValue.Type()).Elem()
			fieldValue.Set(newValueOfInvalid)
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {
				st := fieldValue.Index(j)

				idx := fmt.Sprintf("%ss", strings.ToLower(st.Type().Name()))

				var newID string
				newSubValue := st.Field(0)
				newID = newSubValue.String()

				if newSubValue.Type().Name() == "Model" {
					newID = util.NewUtilConvertToMap().ConvertDataToString(newSubValue.FieldByName("ID"))
				}

				elasticCore.Client.Index().
					Index(idx).
					Type("doc").
					Id(newID).
					BodyJson(st.Interface()).
					Refresh("wait_for").
					Do(context.Background())
			}
			newValueOfInvalid := reflect.New(fieldValue.Type()).Elem()
			fieldValue.Set(newValueOfInvalid)

		case reflect.Ptr, reflect.Struct:
			if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
				continue
			}
			idx := fmt.Sprintf("%ss", strings.ToLower(fieldName))
			subValue := fieldValue

			var newID string

			if fieldValue.Kind() == reflect.Ptr {
				subValue = fieldValue.Elem()
			}

			newSubValue := subValue.Field(0)
			newID = newSubValue.String()

			if newSubValue.Type().Name() == "Model" {
				newID = util.NewUtilConvertToMap().ConvertDataToString(newSubValue.FieldByName("ID"))
			}

			elasticCore.Client.Index().
				Index(idx).
				Type("doc").
				Id(newID).
				BodyJson(fieldValue.Interface()).
				Refresh("wait_for").
				Do(context.Background())

			newValueOfInvalid := reflect.New(fieldValue.Type()).Elem()
			fieldValue.Set(newValueOfInvalid)

		default:
			continue
		}
	}

	elasticCore.Client.Index().
		Index(elasticCore.Index).
		Type("doc").
		Id(IDParams).
		BodyJson(field.Interface()).
		Refresh("wait_for").
		Do(context.Background())

	return nil
}
