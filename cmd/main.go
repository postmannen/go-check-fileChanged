package main

import (
	"fmt"
	"log"

	jftm "github.com/postmannen/jsonfiletomap"
)

func main() {
	fileUpdated := make(chan bool)
	fileName := "commandToTemplate.json"
	//Start the file watcher
	jftm.Run(fileName, fileUpdated)

	cmdToTplMap := jftm.NewMap()

	for {
		select {
		case <-fileUpdated:
			cmdToTplMap, err := jftm.ReadJSONFileToMap(fileName, cmdToTplMap)
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
