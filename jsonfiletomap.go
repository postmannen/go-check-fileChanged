//Package mapfile Checks if JSON file is updated.
//If the file is updated then decode the JSON,
//and put the content in the map.
//If it at some point fails, the current working map
//will be kept.
package mapfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

//Update holds the info about how the file watching are going
type Update struct {
	Time time.Time
	Err  error
}

//FileWatcher will
type FileWatcher struct {
	name      string
	Updates   chan Update
	fsWatcher *fsnotify.Watcher
}

//New create a new FileWatcher struct
func New(fileName string, updates chan Update) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.New("unable to start fsnotify watcher")
	}

	return &FileWatcher{fileName, updates, watcher}, nil
}

//Watch will check the file
func (fw *FileWatcher) Watch() error {
	err := fw.fsWatcher.Add(fw.name)
	if err != nil {
		return fmt.Errorf("unable to start watching %s", fw.name)
	}

	go func() {
		for {
			select {
			case event := <-fw.fsWatcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fw.Updates <- Update{Time: time.Now(), Err: nil}
				}
			case err := <-fw.fsWatcher.Errors:
				fw.Updates <- Update{Time: time.Now(), Err: err}
			}
		}
	}()

	return nil
}

//Close for closing down the file watcher
func (fw *FileWatcher) Close() error {
	return fw.fsWatcher.Close()
}

//------------------------------------------------------
//done is used for stopping the services.
var done = make(chan bool)

//StartFileWatcher starts the filewatcher.
func StartFileWatcher(d Data) {
	//************
	go checkFileUpdated(d)
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
func NewData(fileName string) Data {
	return Data{
		FileUpdated: make(chan bool),
		FileError:   make(chan error),
		FileName:    fileName,
		AMap:        make(map[string]string),
	}

}

//checkFileUpdated , this is basically the same code as given as example
//in the fsnotify doc.......with some minor changes.
func checkFileUpdated(d Data) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed fsnotify.NewWatcher")
		return
	}
	defer watcher.Close()

	go func() {
		//Give a true value to updated so it reads the file the first time.
		d.FileUpdated <- true

		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
					//testing with an update chan to get updates
					//instead of just logs
					d.FileUpdated <- true
				}
			case err := <-watcher.Errors:
				log.Println("pkg:jsonfiletomap checkFileUpdated:", err)
				d.FileError <- err
			}
		}
	}()

	err = watcher.Add(d.FileName)
	if err != nil {
		d.FileError <- err
	}
	<-done

	return
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
	defer f.Close()

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
