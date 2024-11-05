package main

import (
	"jx3-osm/cmd"
)

func main() {
	go cmd.Execute()
	select {}
}
