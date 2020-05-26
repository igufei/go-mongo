# Proxy Middleware for Gin Framework

This is a middleware for [Gin](https://github.com/gin-gonic/gin) framework.

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