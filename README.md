# What does this package does

* Continously check for json file change.
* If changed, read and decode json.
* Put the decoded content into a map.
* If any operation fails the current working map will be kept, and an error will be printed to console

Example of how to use can be found under ./cmd/

## Notes

The JSON file is structured as key/value pairs..like the map

    "addButton": "buttonTemplate1",
    "addHeader": "socketTemplate1",
    "addParagraph": "paragraphTemplate1"
