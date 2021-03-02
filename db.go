package gomongo

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
	"time"
)

var obj = Instance

type DB struct {
	colName string
	doc     interface{}
}

func New(colName string, doc interface{}) *DB {
	return &DB{
		colName: colName,
		doc:     doc,
	}
}

// Parse 转换数据
func (my *DB) Parse(data bson.M, schema interface{}) error {
	if data == nil {
		return errors.New("数据不存在")
	}
	if err := obj.ToStruct(data, schema); err != nil {
		return err
	}
	return nil
}

// AddOne 添加
func (my *DB) Save() error {
	var data bson.M
	bytes, err := bson.Marshal(my.doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	id, ok := data["_id"].(primitive.ObjectID)
	if !ok {
		createTime, ok := data["create_time"].(int64)
		if !ok || createTime == 0 {
			data["create_time"] = time.Now().UnixNano() / 1e6
		}
		success := obj.InsertOne(my.colName, data)
		if !success {
			return errors.New("保存失败")
		}
	} else {
		ok := obj.UpdateOne(my.colName, bson.M{"_id": id}, bson.M{"$set": my.doc})
		if !ok {
			return errors.New("保存失败")
		}
	}

	return nil
}

// DeleteOneByID 删除
func (my *DB) Delete() error {
	v := struct2Map(my.doc)
	id := v["_id"].(primitive.ObjectID)
	if id.Hex() == "000000000000000000000000" {
		return errors.New("删除失败")
	}
	b := obj.DeleteOne(my.colName, bson.M{"_id": id})
	if !b {
		return errors.New("删除失败")
	}
	return nil
}

// Load 加载文档
func (my *DB) LoadDoc(filter bson.M) error {
	userF := obj.FindOne(my.colName, filter)
	if userF == nil {
		return errors.New("数据不存在！")
	}
	err := my.Parse(userF, my.doc)
	if err != nil {
		return errors.New("数据解析失败！")
	}
	return nil
}

// FindMany 查询多个
func (my *DB) Find(pageNum int64, filter bson.M) ([]bson.M, error) {
	pageNum = pageNum - 1
	if pageNum < 0 {
		pageNum = 0
	}
	opts := new(options.FindOptions)
	list := obj.FindMany(my.colName, filter,
		opts.SetSkip(pageNum*20),
		opts.SetLimit(20),
		opts.SetSort(bson.M{
			"create_time": 1,
		}))

	return list, nil
}
func (my *DB) Aggregate(filter interface{}, opts ...*options.AggregateOptions) ([]bson.M, error) {
	list := obj.Aggregate(my.colName, filter)
	return list, nil
}

// FindMany 查询多个
func (my *DB) All(filter bson.M) ([]bson.M, error) {
	list := obj.FindMany(my.colName, filter)
	return list, nil
}

// Has 文档是否存在
func (my *DB) Exist(filter bson.M) bool {
	userF := obj.FindOne(my.colName, filter)
	if userF == nil {
		return false
	}
	return true
}

// Count 数量
func (my *DB) Count(filter bson.M) (int64, error) {
	return obj.Count(my.colName, filter), nil
}

func (my *DB) SetDoc(doc interface{}) {
	my.doc = doc
}
func struct2Map(obj interface{}) map[string]interface{} {
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	var data = make(map[string]interface{})

	for i := 0; i < relType.NumField(); i++ {
		name := strings.SplitN(relType.Field(i).Tag.Get("bson"), ",", 2)[0]
		data[name] = elem.Field(i).Interface()
	}
	return data
}
