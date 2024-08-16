# SPMCB
**S**ingle **P**roducer **M**ultiple **C**onsumer **B**uffer!

It's a ring buffer!



### Usage

```go
package main

import "githuh.com/xslook/spmcb"

func main() {
    ctx := context.Background()

    rb, _ := spmcb.New[string](100)

    c1, _ := spmcb.NewConsumer(rb)
    c2, _ := spmcb.NewConsumer(rb)

    rb.Produce("hello")
    rb.Produce("world")

    v11, _ := c1.Consume(ctx) // v11 == "hello"
    v21, _ := c2.Consume(ctx) // v21 == "hello"
    v12, _ := c2.Consume(ctx) // v12 == "world"
}
```


### LICENSE
MIT LICENSE

