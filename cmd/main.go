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
	fileError := make(chan error)
	fileName := "./commandToTemplate.json"

	//Start the file watcher
	jsonfiletomap.StartFileWatcher(fileName, fileUpdated, fileError)
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

			fmt.Println("----------------------------------------------------------------")
			fmt.Println("Content of the map unmarshaled from fileContent :")
			for key, value := range cmdToTplMap {
				fmt.Println("key = ", key, "value = ", value)
			}
			fmt.Println("----------------------------------------------------------------")

		case errF := <-fileError:
			fmt.Println("---Main: Received on error channel.", errF)
		}
	}
}
