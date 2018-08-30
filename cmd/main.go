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
	defer jsonfiletomap.StopFileWatcher()

	myMap := jsonfiletomap.NewMap()

	for {
		select {
		case <-fileUpdated:
			var err error
			//Convert the file content, and insert into map.
			//Will return current map if the new one fails.
			myMap, err = jsonfiletomap.Convert(fileName, myMap)
			if err != nil {
				log.Println("file to JSON to map problem : ", err)
			}

			fmt.Println("----------------------------------------------------------------")
			fmt.Println("Content of the map unmarshaled from fileContent :")
			for key, value := range myMap {
				fmt.Println("key = ", key, "value = ", value)
			}
			fmt.Println("----------------------------------------------------------------")

		//Catch the errors from the functions that are Go routines
		case errF := <-fileError:
			fmt.Println("---Main: Received on error channel.", errF)
		}
	}
}
