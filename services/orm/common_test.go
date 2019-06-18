package services

import (
	"reflect"
	"weassistant/conf"

	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

var (
	config conf.Config
)

func init() {
	config = conf.MustNewConfig()
	err := config.Load("../../config.json")
	if err != nil {
		panic(err)
	}
}

func checkIsFunc(t *testing.T, fnValue reflect.Value) {
	if fnValue.Kind() != reflect.Func {
		t.Fatalf("funcValue is not a func, is %v", fnValue.Kind())
	}
}

// 断言reflect.Value的Kind
func assertKind(value reflect.Value, kind reflect.Kind) {
	if value.Kind() != kind {
		panic(fmt.Sprintf("valud.Kind %v is not kind %v", value.Kind(), kind))
	}
}

// testService 通用测试
func testService(t *testing.T, createServFn interface{}, createModelFn interface{}) {
	var servValue reflect.Value
	// 获取待测服务的value
	crtServFnValue := reflect.ValueOf(createServFn)
	switch crtServFnValue.Kind() {
	case reflect.Func:
		db := config.GetMainDB()
		callRes := crtServFnValue.Call([]reflect.Value{reflect.ValueOf(db)})
		servValue = callRes[0].Elem()
	case reflect.Ptr:
		servValue = crtServFnValue
	default:
		panic(fmt.Sprintf("unknown createServFn Kind %v", crtServFnValue.Kind()))
	}
	// 获取生成新模型的FuncValue
	crtMdlFnValue := reflect.ValueOf(createModelFn)
	assertKind(crtMdlFnValue, reflect.Func)
	// 获取方法
	servSaveFnValue := servValue.MethodByName("Save")
	assertKind(servSaveFnValue, reflect.Func)
	servGetFnValue := servValue.MethodByName("Get")
	assertKind(servGetFnValue, reflect.Func)
	servDeleteFnValue := servValue.MethodByName("Delete")
	assertKind(servDeleteFnValue, reflect.Func)
	// 调用方法返回时返回结果的value储存
	var fnRetValues []reflect.Value
	// 单个测试
	{
		newDataPtrValue := crtMdlFnValue.Call(nil)[0]
		assertKind(newDataPtrValue, reflect.Ptr)
		errValue := servSaveFnValue.Call([]reflect.Value{newDataPtrValue})[0]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		if newDataPtrValue.Elem().FieldByName("ID").Uint() == 0 {
			t.Fatal("service.Save fail: merchant.ID is 0")
		}
		t.Log("new data ID is ", newDataPtrValue.Elem().FieldByName("ID").Uint())
		fnRetValues = servGetFnValue.Call([]reflect.Value{newDataPtrValue.Elem().FieldByName("ID")})
		errValue = fnRetValues[1]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		dbDataValue := fnRetValues[0]
		if !dbDataValue.MethodByName("Equal").Call([]reflect.Value{newDataPtrValue.Elem()})[0].Bool() {
			t.Fatalf("db data diff\n%+v\n%+v", newDataPtrValue.Elem().Interface(), dbDataValue.Interface())
		}
		// 尝试删除
		errValue = servDeleteFnValue.Call([]reflect.Value{newDataPtrValue})[0]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		fnRetValues = servGetFnValue.Call([]reflect.Value{dbDataValue.FieldByName("ID")})
		errValue = fnRetValues[1]
		if errValue.IsNil() { // 判断Error
			t.Fatalf("data still exist after delete: %v[%d]", dbDataValue.Type().Name(), dbDataValue.FieldByName("ID").Uint())
		} else if errValue.Interface().(error) != gorm.ErrRecordNotFound {
			t.Fatal(errValue.Interface().(error))
		}
	}
	// 获取方法
	servGetByWhereOptionsFnValue := servValue.MethodByName("GetByWhereOptions")
	assertKind(servGetByWhereOptionsFnValue, reflect.Func)
	servGetListByWhereOptionsFnValue := servValue.MethodByName("GetListByWhereOptions")
	assertKind(servGetListByWhereOptionsFnValue, reflect.Func)
	servGetCountByWhereOptionsFnValue := servValue.MethodByName("GetCountByWhereOptions")
	assertKind(servGetCountByWhereOptionsFnValue, reflect.Func)
	servDeleteByWhereOptionsFnValue := servValue.MethodByName("DeleteByWhereOptions")
	assertKind(servDeleteByWhereOptionsFnValue, reflect.Func)
	// WhereOption测试
	{
		fnRetValues = servGetCountByWhereOptionsFnValue.Call([]reflect.Value{reflect.ValueOf([]OrmWhereOption{})})
		errValue := fnRetValues[1]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		existDataCount := fnRetValues[0].Uint()
		t.Logf("data count: %d", existDataCount)
		var dataPtrValues []reflect.Value
		var dataIDs []uint64
		// 先制造测试数据
		for i := 0; i < 100; i++ {
			dataPtrValue := crtMdlFnValue.Call(nil)[0]
			assertKind(dataPtrValue, reflect.Ptr)
			errValue = servSaveFnValue.Call([]reflect.Value{dataPtrValue})[0]
			if !errValue.IsNil() { // 判断Error
				t.Fatal(errValue.Interface().(error))
			}
			if dataPtrValue.Elem().FieldByName("ID").Uint() == 0 {
				t.Fatal("service.Save fail: merchant.ID is 0")
			}
			dataPtrValues = append(dataPtrValues, dataPtrValue)
			dataIDs = append(dataIDs, dataPtrValue.Elem().FieldByName("ID").Uint())
		}
		whereOptions := []OrmWhereOption{
			OrmWhereOption{Query: "id in (?)", Item: []interface{}{dataIDs}},
		}
		// 测试获取全部
		fnRetValues = servGetListByWhereOptionsFnValue.Call([]reflect.Value{
			reflect.ValueOf(whereOptions),
			reflect.ValueOf([]string{}),
			reflect.ValueOf(int64(0)),
			reflect.ValueOf(int64(0)),
		})
		errValue = fnRetValues[1]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		dbDatasValue := fnRetValues[0]
		if len(dataPtrValues) != dbDatasValue.Len() {
			t.Fatalf("len(datas) %d different with len(dbDatas) %d", len(dataPtrValues), dbDatasValue.Len())
		}
		for k, dataPtrValue := range dataPtrValues {
			dbDataValue := dbDatasValue.Index(k)
			if !dbDataValue.MethodByName("Equal").Call([]reflect.Value{dataPtrValue.Elem()})[0].Bool() {
				t.Fatalf("db data diff\n%+v\n%+v", dataPtrValue.Elem().Interface(), dbDataValue.Interface())
			}
			// 尝试单个读取检验是否正确
			whereOptions := []OrmWhereOption{
				OrmWhereOption{Query: "id = ?", Item: []interface{}{dataPtrValue.Elem().FieldByName("ID").Uint()}},
			}
			fnRetValues = servGetByWhereOptionsFnValue.Call([]reflect.Value{
				reflect.ValueOf(whereOptions),
			})
			errValue = fnRetValues[1]
			if !errValue.IsNil() { // 判断Error
				t.Fatal(errValue.Interface().(error))
			}
			dbDatValue := fnRetValues[0]
			if !dbDatValue.MethodByName("Equal").Call([]reflect.Value{dbDataValue})[0].Bool() {
				t.Fatalf("data diff with dataList elem:\n%+v\n%+v", dbDatValue.Interface(), dbDataValue.Interface())
			}
		}
		// 测试Limit、Offset
		fnRetValues = servGetListByWhereOptionsFnValue.Call([]reflect.Value{
			reflect.ValueOf(whereOptions),
			reflect.ValueOf([]string{}),
			reflect.ValueOf(int64(10)),
			reflect.ValueOf(int64(10)),
		})
		errValue = fnRetValues[1]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		dbDatasValue = fnRetValues[0]
		if dbDatasValue.Len() != 10 {
			t.Fatalf("len(dbDatas) %d is not equal 10", dbDatasValue.Len())
		}
		for k, dataPtrValue := range dataPtrValues[10:20] {
			dbDataValue := dbDatasValue.Index(k)
			if !dbDataValue.MethodByName("Equal").Call([]reflect.Value{dataPtrValue.Elem()})[0].Bool() {
				t.Fatalf("db data diff\n%+v\n%+v", dataPtrValue.Elem().Interface(), dbDataValue.Interface())
			}
			// 尝试单个读取检验是否正确
			whereOptions := []OrmWhereOption{
				OrmWhereOption{Query: "id = ?", Item: []interface{}{dataPtrValue.Elem().FieldByName("ID").Uint()}},
			}
			fnRetValues = servGetByWhereOptionsFnValue.Call([]reflect.Value{
				reflect.ValueOf(whereOptions),
			})
			errValue = fnRetValues[1]
			if !errValue.IsNil() { // 判断Error
				t.Fatal(errValue.Interface().(error))
			}
			dbDatValue := fnRetValues[0]
			if !dbDatValue.MethodByName("Equal").Call([]reflect.Value{dbDataValue})[0].Bool() {
				t.Fatalf("data diff with dataList elem:\n%+v\n%+v", dbDatValue.Interface(), dbDataValue.Interface())
			}
		}
		// 删除这些不再有利用价值的数据
		errValue = servDeleteByWhereOptionsFnValue.Call([]reflect.Value{
			reflect.ValueOf(whereOptions),
		})[0]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		// 检查是否删干净了
		fnRetValues = servGetListByWhereOptionsFnValue.Call([]reflect.Value{
			reflect.ValueOf([]OrmWhereOption{}),
			reflect.ValueOf([]string{}),
			reflect.ValueOf(int64(0)),
			reflect.ValueOf(int64(0)),
		})
		errValue = fnRetValues[1]
		if !errValue.IsNil() { // 判断Error
			t.Fatal(errValue.Interface().(error))
		}
		dbDatasValue = fnRetValues[0]
		if dbDatasValue.Len() != int(existDataCount) {
			t.Fatalf("len(datas) %d is not existDataCount %d after delete", dbDatasValue.Len(), existDataCount)
		}
	}
}
