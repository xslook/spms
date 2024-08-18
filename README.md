# SPMS
**S**ingle **P**ublisher **M**ultiple **S**ubscriber!

It's a ring buffer!



### Usage

```go
package main

import "githuh.com/xslook/spms"

func main() {
    ctx := context.Background()

    rb, _ := spms.New[string](100)

    c1, _ := spms.NewSubscriber(rb)
    c2, _ := spms.NewSubscriber(rb)

    rb.Produce("hello")
    rb.Produce("world")

    v11, _ := c1.Read(ctx) // v11 == "hello"
    v21, _ := c2.Read(ctx) // v21 == "hello"
    v12, _ := c2.Read(ctx) // v12 == "world"
}
```


### LICENSE
MIT LICENSE

