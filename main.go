package main

import (
	"log"
	"os"

	_ "github.com/MarcosRoch4/gofeed/matchers"
	"github.com/MarcosRoch4/gofeed/search"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("president")
}
