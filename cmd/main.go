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
	//Create data structure who holds the map and channels
	d := jsonfiletomap.NewData("./commandToTemplate.json")

	//Start the file watcher
	jsonfiletomap.StartFileWatcher(d)
	defer jsonfiletomap.StopFileWatcher()

	for {
		select {
		case <-d.FileUpdated:
			var err error
			//Convert the file content, and insert into map.
			//Will return current map if the new one fails.
			d.AMap, err = jsonfiletomap.Convert(d.FileName, d.AMap)
			if err != nil {
				log.Println("Main :", err)
			}

			printMap(d)

		//Catch the errors from the functions that are Go routines
		case errF := <-d.FileError:
			fmt.Println("---Main: Received on error channel.", errF)
		}
	}
}

func printMap(d jsonfiletomap.Data) {

	//Print out all the values for testing
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Content of the map unmarshaled from fileContent :")
	for key, value := range d.AMap {
		fmt.Println("key = ", key, "value = ", value)
	}
	fmt.Println("----------------------------------------------------------------")

}
