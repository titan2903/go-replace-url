package main

import (
	"context"
	"replace-url-gin/config"
	"replace-url-gin/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	config.Init()

	ctx := context.Background()
	server.NewApp(ctx).Start()
}
