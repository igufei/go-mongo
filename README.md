# Proxy Middleware for Gin Framework

实现对官方 [mongodb](https://github.com/mongodb/mongo-go-driver) 的golang驱动的单例封装

## Usage

Download and install it:

```sh
$ go get github.com/igufei/go-mongo
```

Import it in your code:

```go
import "github.com/igufei/go-mongo"
```

## Example

```go
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
```