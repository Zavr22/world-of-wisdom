package main

import (
	"context"
	"fmt"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/cache"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/clock"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/config"
	"github.com/Zavr22/world-of-wisdom/internal/server"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("start server")

	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)
	ctx = context.WithValue(ctx, "clock", &clock.SystemClock{})

	cacheInst := cache.InitInMemoryCache(&clock.SystemClock{})
	if cacheInst == nil {
		fmt.Println("error creating cache instance")
		return
	}
	ctx = context.WithValue(ctx, "cache", cacheInst)

	rand.Seed(time.Now().UnixNano())

	// run server
	serverAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
	err = server.Run(ctx, serverAddress)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
