//Package jsonfiletomap Check if file is updated.
//If the file is updated decode the JSON,
//and put the content in the map.
package jsonfiletomap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

//fileName file to follow
//var fileName = "./commandToTemplate.json"

//FileUpdated is a channel that will give
//a value if file is updated.
//var FileUpdated = make(chan bool)

/*
TODO:
Make the fileName exportet
Make the package return the map from a function to be used from main
*/

//Run starts the filewatcher.
func Run(fileName string, fileUpdated chan bool) {
	//************
	go checkFileUpdated(fileName, fileUpdated)
}

//NewMap creates a map to hold all the parsed file values
func NewMap() map[string]string {
	m := make(map[string]string)
	return m
}

//ReadJSONFileToMap Load file, read it's content, parse JSON,
//and return map with parsed values.
//If it fails at some point then return the current map.
func ReadJSONFileToMap(fileName string, currentMap map[string]string) (map[string]string, error) {
	cmdToTplMap := make(map[string]string)

	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Failed to open file %v\n", err)
		return currentMap, err
	}

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Failed reading file %v\n", err)
		return currentMap, err
	}

	fmt.Println("Content read from file : \n", string(fileContent))

	err = json.Unmarshal(fileContent, &cmdToTplMap)
	if err != nil {
		log.Printf("Failed unmarshaling %v\n", err)
		return currentMap, err
	}

	return cmdToTplMap, nil
}

//checkFileUpdated , this is basically the same code as given as example
//in the fsnotify doc.......with some minor changes.
func checkFileUpdated(fileName string, fileUpdated chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed fsnotify.NewWatcher")
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		//Give a true value to updated so it reads the file the first time.
		fileUpdated <- true

		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
					//testing with an update chan to get updates
					//instead of just logs
					fileUpdated <- true
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(fileName)
	if err != nil {
		log.Fatal(err)
	}
	<-done

	return
}
