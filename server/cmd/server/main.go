package main

import (
	"session-demo/internal/config"
	"session-demo/internal/routes"
	"strconv"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	_, err = config.ConnectDB(cfg.MongoURI)

	if err != nil {
		panic(err)
	}

	r := routes.SetUpRouter()

	r.Run(":" + strconv.Itoa(cfg.PORT))
}
