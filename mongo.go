package gomongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"
)

var mdb *mongodb
var once sync.Once

type mongodb struct {
	database *mongo.Database
}

// Instance mongodb 的单例
var Instance *mongodb

// MustConnect 连接数据库，并创建单例
func MustConnect(dbHOST string, dbName string) {
	once.Do(func() {
		mdb = &mongodb{}
		if err := mdb._connect(dbHOST, dbName); err != nil {
			panic("数据库连接失败")
		}
		Instance = mdb
	})
}

// 链接数据库
func (db *mongodb) _connect(dbHOST string, dbName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbHOST))
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	db.database = client.Database(dbName)
	return nil
}

// 插入一条数据
func (db *mongodb) InsertOne(collectionName string, doc interface{}) bool {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return false
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.InsertOne(ctx, doc)
	if result == nil || err != nil {
		return false
	}
	return true
}

// 插入多条数据
func (db *mongodb) InsertMany(collectionName string, doc []interface{}) int64 {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return 0
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.InsertMany(ctx, doc)
	if result == nil || err != nil {
		return 0
	}
	return int64(len(result.InsertedIDs))
}

// 查询单条数据
func (db *mongodb) FindOne(collectionName string, filter interface{}) bson.M {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return bson.M(nil)
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result bson.M
	err := table.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return bson.M(nil)
	}
	return result
}

// 查询多少数据
func (db *mongodb) FindMany(collectionName string, filter interface{}, opts ...*options.FindOptions) []bson.M {
	resultArr := make([]bson.M, 0) //[]bson.M{}
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return resultArr
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := table.Find(ctx, filter, opts...)
	if err != nil {
		Error.Printf("mongo:%v\n", err)
		return resultArr
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		err3 := cur.Decode(&result)
		if err3 != nil {
			return resultArr
		}
		resultArr = append(resultArr, result)
	}
	return resultArr
}

// 修改一条数据
func (db *mongodb) UpdateOne(collectionName string, filter interface{}, update interface{}) bool {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return false
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.UpdateOne(ctx, filter, update)
	if err != nil {
		return false
	}
	return result.ModifiedCount == 1
}

// 修改多条数据
func (db *mongodb) UpdateMany(collectionName string, filter interface{}, update interface{}) int64 {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return 0
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0
	}
	return result.ModifiedCount
}

// 删除一条数据
func (db *mongodb) DeleteOne(collectionName string, filter interface{}) bool {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return false
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.DeleteOne(ctx, filter)
	if err != nil {
		return false
	}
	return result.DeletedCount == 1
}

// 删除多条数据
func (db *mongodb) DeleteMany(collectionName string, filter interface{}) int64 {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return 0
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.DeleteMany(ctx, filter)
	if err != nil {
		return 0
	}
	return result.DeletedCount
}

// 查询数据的条数
func (db *mongodb) Count(collectionName string, filter interface{}) int64 {
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return 0
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := table.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}
	return result
}

// 聚合查询
func (db *mongodb) Aggregate(collectionName string, filter interface{}, opts ...*options.AggregateOptions) []bson.M {
	resultArr := make([]bson.M, 0)
	if db.database == nil {
		Error.Printf(":%v\n", "没有初始化连接和数据库信息！")
		return resultArr
	}
	table := db.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := table.Aggregate(ctx, filter, opts...)
	if err != nil {
		Error.Printf("mongo:%v\n", err)
		return resultArr
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		err3 := cur.Decode(&result)
		if err3 != nil {
			return resultArr
		}
		resultArr = append(resultArr, result)
	}
	return resultArr
}

func (db *mongodb) ToStruct(bsonValue interface{}, objectValue interface{}) error {
	data, err := bson.Marshal(bsonValue)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(data, objectValue)
	if err != nil {
		return err
	}
	return nil
}
func (db *mongodb) ToObjectID(id string) primitive.ObjectID {
	oID, _ := primitive.ObjectIDFromHex(id)

	return oID
}
