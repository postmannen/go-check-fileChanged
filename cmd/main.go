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
	defer close(updates)

	fw, err := mapfile.New("commandToTemplate.json", updates)
	if err != nil {
		log.Println("Main: Failed to create new FileWatcher struct: ", err)
	}

	err = fw.Watch()
	if err != nil {
		log.Println("Main: Failed to Watch: ", err)
	}

	defer fw.Close()

	aMap := make(map[string]string)

	for {
		select {
		case u := <-fw.Updates:
			if u.Err == nil {
				fmt.Println("No FileWatch Error: ", u)

				aMap, err = fw.Convert(aMap)
				if err != nil {
					log.Println("Error: ", err)
				}

				printMap(aMap)
			} else {
				fmt.Println("FileWatch Error: ", u)
			}
		}
	}
}

func printMap(m map[string]string) {

	//Print out all the values for testing
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Content of the map unmarshaled from fileContent :")
	for key, value := range m {
		fmt.Println("key = ", key, "value = ", value)
	}
	fmt.Println("----------------------------------------------------------------")

}
