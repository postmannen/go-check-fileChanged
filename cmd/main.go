package main

import (
	"fmt"
	"log"

	"github.com/postmannen/jsonfiletomap"
)

/*
The package should:
	- Be started from main
	- Set the file to watch
	- Export the map created when the file changed
	- Export the channel for when the file changed
*/

func main() {
	fileUpdated := make(chan bool)
	fileName := "commandToTemplate.json"
	jsonfiletomap.Run(fileName, fileUpdated)

	var cmdToTplMap map[string]string
	//var fileName = "./commandToTemplate.json"

	for {
		select {
		case <-fileUpdated:
			//load file, read it's content, parse JSON,
			//and return map with parsed values
			cmdToTplMap, err := jsonfiletomap.ReadJSONFileToMap(fileName, cmdToTplMap)
			if err != nil {
				log.Println("file to JSON to map problem : ", err)
			}

			if cmdToTplMap != nil {
				fmt.Println("\nContent of the map unmarshaled from fileContent :")
				for key, value := range cmdToTplMap {
					log.Println("key = ", key, "value = ", value)
				}
			}
		}
	}
}
