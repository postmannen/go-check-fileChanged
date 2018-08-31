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

//New creates a new FileWatcher struct.
//The "updates chan Update" is taken as input so the user
//can choose to make a buffered or unbuffered channel.
func New(fileName string, updates chan Update) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.New("unable to start fsnotify watcher")
	}

	return &FileWatcher{fileName, updates, watcher}, nil
}

//Watch will check the file, and send an update over the
//methods Updates channel if file is changed.
func (fw *FileWatcher) Watch() error {
	err := fw.fsWatcher.Add(fw.name)
	if err != nil {
		return fmt.Errorf("unable to start watching %s", fw.name)
	}

	go func() {
		//Trigger the first reading of the file by sending an update to the Updates channel
		fw.Updates <- Update{
			Time: time.Now(),
			Err:  nil,
		}

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

//Convert loads the file,
//reads it's content, parse the JSON
//and returns a new map with the parsed values.
//If it fails at some point then return the current map.
func (fw *FileWatcher) Convert(currentMap map[string]string) (map[string]string, error) {
	theMap := make(map[string]string)

	f, err := os.Open(fw.name)
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
