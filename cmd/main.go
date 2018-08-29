package main

import (
	"fmt"
	"log"

	"github.com/postmannen/jsonfiletomap"
)

/*
The package should:
	Have a channel "fileupdated"telling if the file was updated
*/

func main() {
	fileUpdated := make(chan bool)
	fileName := "commandToTemplate.json"
	//Start the file watcher
	jsonfiletomap.Run(fileName, fileUpdated)

	cmdToTplMap := jsonfiletomap.NewMap()

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
