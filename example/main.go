package example

import (
	"fmt"
	mongo "github.com/igufei/go-mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	mongo.MustConnect("mongodb://localhost:27017", "usermanager")
}

func main() {
	value := mongo.Instance.FindOne("agent", bson.M{
		"username": "gufeng",
	})
	fmt.Printf("%v", value)
}
