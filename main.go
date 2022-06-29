package main

import (
	"context"
	"fmt"
	"replace-url-gin/config"
	"replace-url-gin/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	config.Init()
	ctx := context.Background()
	fmt.Println("Start URL")
	server.NewApp(ctx).Start()
}
