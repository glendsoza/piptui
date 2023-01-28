package main

import (
	"log"

	"github.com/glendsoza/piptui/pkg/widgets"
)

func main() {
	ui := widgets.NewUI()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
