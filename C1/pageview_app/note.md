# PageView APP with Redis

### start an Redis server
```bash
sudo docker run --name db -d redis
```

### main.go

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v7"
)

var dbClient *redis.Client
var key = "pageview"

func init() {
	dbClient = redis.NewClient(&redis.Options{
		Addr: "db:6379",
	})
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ping from %s", r.RemoteAddr)
	pageView, err := dbClient.Incr(key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Hello, you're visitor #%v.\n", pageView)
}

```

### build app and build image

`CGO_ENABLE=0` is needed or it gives `exec user process caused "no such file or directory"`

check `CGO_ENABLE` later at: https://johng.cn/cgo-enabled-affect-go-static-compile/

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o .
sudo docker build -t pageview:v0.0.1 .
```

### run container

```bash
$ sudo docker run -p 8080:8080 --link db:db -d  pageview:v0.0.1
e384b3e1326269f4a9a26d782701c27fcebbbca9836e26634ad2ad5a8f4b5bf6

```

### test

```bash
$ curl http://localhost:8080
Hello, you're visitor #1.

$ curl http://localhost:8080
Hello, you're visitor #2.

```

