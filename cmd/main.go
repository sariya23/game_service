package main

import (
	"fmt"

	"github.com/sariya23/game_service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
