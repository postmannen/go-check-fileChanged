//Package jsonfiletomap Check if file is updated.
//If the file is updated decode the JSON,
//and put the content in the map.
package jsonfiletomap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

//done is used for stopping the services.
var done = make(chan bool)

//StartFileWatcher starts the filewatcher.
func StartFileWatcher(fileName string, fileUpdated chan bool, fileError chan error) {
	//************
	go checkFileUpdated(fileName, fileUpdated, fileError)
}

//StopFileWatcher is used to stop all running Go routines
func StopFileWatcher() {
	done <- true
}

//Data holds all the variable types needed for the service
type Data struct {
	FileUpdated chan bool
	FileError   chan error
	FileName    string
	AMap        map[string]string
}

//NewData creates a data structure for
//the variables used in the package
func NewData(fileName string) *Data {
	return &Data{
		FileUpdated: make(chan bool),
		FileError:   make(chan error),
		FileName:    fileName,
		AMap:        make(map[string]string),
	}

}

//Convert loads the file,
//reads it's content, parse the JSON
//and returns a new map with the parsed values.
//If it fails at some point then return the current map.
func Convert(fileName string, currentMap map[string]string) (map[string]string, error) {
	theMap := make(map[string]string)

	f, err := os.Open(fileName)
	if err != nil {
		e := fmt.Sprintln("Convert: Keeping current map, file open failed :", err.Error())
		return currentMap, errors.New(e)
	}

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		e := fmt.Sprintln("Convert: Keeping current map, ReadAll failed :", err.Error())
		return currentMap, errors.New(e)
	}

	err = json.Unmarshal(fileContent, &theMap)
	if err != nil {
		e := fmt.Sprintln("Convert: Keeping current map, Unmarshal failed :", err.Error())
		return currentMap, errors.New(e)
	}

	//If no failures, return the new map.
	return theMap, nil
}

//checkFileUpdated , this is basically the same code as given as example
//in the fsnotify doc.......with some minor changes.
func checkFileUpdated(fileName string, fileUpdated chan bool, fileError chan error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed fsnotify.NewWatcher")
		return
	}
	defer watcher.Close()

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
				log.Println("pkg:jsonfiletomap checkFileUpdated:", err)
				fileError <- err
			}
		}
	}()

	err = watcher.Add(fileName)
	if err != nil {
		fileError <- err
	}
	<-done

	return
}
