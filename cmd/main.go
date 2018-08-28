package main

import "github.com/postmannen/jsonfiletomap"

/*
The package should:
	- Be started from main
	- Set the file to watch
	- Export the map created when the file changed
	- Export the channel for when the file changed
*/

func main() {
	jsonfiletomap.FileName = "./commandToTemplate2.json"
	jsonfiletomap.Run()
}
