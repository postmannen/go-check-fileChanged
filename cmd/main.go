/*
Example code of how to use the package
*/
package main

import (
	"fmt"
	"log"

	"github.com/postmannen/mapfile"
)

func main() {
	updates := make(chan mapfile.Update)
	fw, err := mapfile.New("commandToTemplate.json", updates)
	if err != nil {
		log.Println("Main: Failed to create new FileWatcher struct: ", err)
	}

	err = fw.Watch()
	if err != nil {
		log.Println("Main: Failed to Watch: ", err)
	}

	for {
		select {
		case u := <-fw.Updates:
			fmt.Println(u)
		}
	}

	/*
		//Create data structure who holds the map and channels
		d := mapfile.NewData("./commandToTemplate.json")

		//Start the file watcher
		mapfile.StartFileWatcher(d)
		defer mapfile.StopFileWatcher()

		for {
			select {
			case <-d.FileUpdated:
				var err error
				//Convert the file content, and insert into map.
				//Will return current map if the new one fails.
				d.AMap, err = mapfile.Convert(d.FileName, d.AMap)
				if err != nil {
					log.Println("Main :", err)
				}

				printMap(d)

			//Catch the errors from the functions that are Go routines
			case errF := <-d.FileError:
				log.Println("Error: Main: Received on error channel.", errF)
			}
		}
	*/
}

func printMap(d mapfile.Data) {

	//Print out all the values for testing
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Content of the map unmarshaled from fileContent :")
	for key, value := range d.AMap {
		fmt.Println("key = ", key, "value = ", value)
	}
	fmt.Println("----------------------------------------------------------------")

}
