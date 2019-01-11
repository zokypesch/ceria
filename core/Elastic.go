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

// BulkConfigElastic for configuration elastic
type BulkConfigElastic struct {
	Index string
	ID    string
	Doc   interface{}
}

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

	var ch = make(chan *BulkConfigElastic)

	// var sf []reflect.StructField
	// for i := 0; i < field.NumField(); i++ {
	// 	fieldValue := field.Field(i)
	// 	fieldName := field.Type().Field(i).Name

	// 	if fieldName == "Model" && fieldValue.Kind() == reflect.Struct {
	// 		for j := 0; j < fieldValue.NumField(); j++ {
	// 			sf = append(sf, reflect.StructField{
	// 				Name: fieldValue.Type().Field(j).Name,
	// 				Type: fieldValue.Type().Field(j).Type,
	// 			})
	// 		}
	// 		continue
	// 	}

	// 	switch fieldValue.Kind() {
	// 	case reflect.Struct, reflect.Slice, reflect.Ptr:
	// 		continue
	// 	}

	// 	sf = append(sf, reflect.StructField{
	// 		Name: fieldName,
	// 		Type: field.Type().Field(i).Type,
	// 	})
	// }
	// sTi := reflect.StructOf(sf)
	// sTn := reflect.New(sTi).Elem()
	var ignore []reflect.Type

	for i := 0; i < field.NumField(); i++ {
		fieldValue := field.Field(i)
		fieldName := field.Type().Field(i).Name

		if fieldValue.Kind() == reflect.Invalid || (fieldName == "Model" && fieldValue.Kind() == reflect.Struct) {
			// for j := 0; j < fieldValue.NumField(); j++ {
			// 	sTn.FieldByName(fieldValue.Type().Field(j).Name).Set(fieldValue.Field(j))
			// }
			ignore = append(ignore, fieldValue.Type())
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {
				st := fieldValue.Index(j)
				ignore = append(ignore, fieldValue.Type())

				go elasticCore.sendExecuteBackgroud(st, ch)
				elasticCore.doExecuteBackhround(ch)
			}
		case reflect.Ptr, reflect.Struct:
			if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
				continue
			}
			subValue := fieldValue

			if fieldValue.Kind() == reflect.Ptr {
				subValue = fieldValue.Elem()
			}

			ignore = append(ignore, subValue.Type())

			go elasticCore.sendExecuteBackgroud(subValue, ch)
			elasticCore.doExecuteBackhround(ch)

		default:
			// sTn.FieldByName(fieldName).Set(fieldValue)
			continue
		}
	}

	rebuildProps := &util.RebuildProperty{
		IgnoreFieldString: []string{},
		IgnoreFieldType:   ignore,
		MoveToMember:      []string{"Model"},
	}

	sTn, _ := util.NewServiceStructValue().RebuilToNewStruct(field.Interface(), rebuildProps, true)

	go func() {
		ch <- &BulkConfigElastic{Index: elasticCore.Index, ID: IDParams, Doc: sTn}
	}()

	elasticCore.doExecuteBackhround(ch)

	go elasticCore.Client.Index().Refresh("wait_for").Do(context.Background())

	return nil
}

// SendExecuteBackgroud for execute create struct in golang
func (elasticCore *ElasticCore) sendExecuteBackgroud(rfl reflect.Value, ch chan<- *BulkConfigElastic) {
	idx := fmt.Sprintf("%ss", strings.ToLower(rfl.Type().Name()))
	var newID string
	newSubValue := rfl.Field(0)
	newID = newSubValue.String()

	if newSubValue.Type().Name() == "Model" {
		newID = util.NewUtilConvertToMap().ConvertDataToString(newSubValue.FieldByName("ID"))
	}

	ch <- &BulkConfigElastic{Index: idx, ID: newID, Doc: rfl.Interface()}
}

// DoExecuteBackground for execute create document as background
func (elasticCore *ElasticCore) doExecuteBackhround(ref <-chan *BulkConfigElastic) {

	newCfg := <-ref

	elasticCore.Client.Index().
		Index(newCfg.Index).
		Type("doc").
		Id(newCfg.ID).
		BodyJson(newCfg.Doc).
		Refresh("").
		Do(context.Background())
}
