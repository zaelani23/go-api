package main

import (
	"context"
	"fmt"

	_ "github.com/zaelani/go-api/config"
	_ "github.com/zaelani/go-api/counter"
	_ "github.com/zaelani/go-api/counter/task"
	_ "github.com/zaelani/go-api/httpc"
	_ "github.com/zaelani/go-api/method"
	_ "github.com/zaelani/go-api/sql"
	"github.com/zaelani/go-api/trace"
)

func main() {
	fmt.Println("Whatap Golang api")

	ctx, _ := trace.Start(context.Background(), "Test")
	trace.UpdateMtraceWithContext(ctx, make(map[string][]string))

}
