package main

import (
	"log"

	"github.com/puper/ppgo/helpers"
)

func main() {
	a := map[string]interface{}{
		"a": false,
	}
	cfg := map[string]string{}
	err := helpers.StructDecode(a, &cfg, "json")
	log.Println(cfg, err)
}
