package main

import (
	"dynamic-user-segmentation/pkg/config"
	"dynamic-user-segmentation/pkg/repository"
	"fmt"
	_ "github.com/jackc/pgx"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: init logger

	// TODO: init database
	db, err := repository.NewPostgresDB()
	_, _ = db, err

	// TODO: init router

}
