/*
Example code of how to use the package
*/
package main

import (
	"fmt"
	"log"

	"github.com/postmannen/jsonfiletomap"
)

func main() {
	fileUpdated := make(chan bool)
	fileName := "./commandToTemplate.json"

	//Start the file watcher
	jsonfiletomap.StartFileWatcher(fileName, fileUpdated)
	defer jsonfiletomap.Stop()

	cmdToTplMap := jsonfiletomap.NewMap()

	for {
		select {
		case <-fileUpdated:
			var err error
			//Convert the file content, and insert into map.
			//Will return current map if the new one fails.
			cmdToTplMap, err = jsonfiletomap.Convert(fileName, cmdToTplMap)
			if err != nil {
				log.Println("file to JSON to map problem : ", err)
			}

			fmt.Println("\nContent of the map unmarshaled from fileContent :")
			for key, value := range cmdToTplMap {
				log.Println("key = ", key, "value = ", value)
			}

		}
	}
}
