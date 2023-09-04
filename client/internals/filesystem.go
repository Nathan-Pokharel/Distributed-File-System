package internals

import (
	"encoding/json"
	"os"
)

var response []byte

type filesystem struct {
	Files   []string
	Folders []string
	Cwd     string
}

var FileSystem filesystem

func Populate(cwd string){
    entries, err := os.ReadDir(FileSystem.Cwd)
    FileSystem.Folders = []string{}
    FileSystem.Files= []string{}
    HandleErrors(err, "Error Reading The File System")
    for _, entry := range entries {
        if entry.IsDir() {
            FileSystem.Folders = append(FileSystem.Folders, entry.Name())
        } else {
            FileSystem.Files = append(FileSystem.Files, entry.Name())
        }
}
}
func ReadFromServer() {
	for {
		response = make([]byte, 1024)
		n, err := Server.Read(response)
		HandleErrors(err, "Error Reading From the Server")
		if string(response[:n]) == "Ready" {
            FileSystem.Cwd, err = os.Getwd()
            HandleErrors(err, "Error reading Current Working Directory")
            Populate(FileSystem.Cwd)
		}else{
            FileSystem.Cwd = string(response[:n])
            Populate(FileSystem.Cwd)
        } 
		WriteToServer()
	}
}

func WriteToServer() {
	jsonData, err := json.Marshal(FileSystem)
	jsonData = append(jsonData, '\n')
	HandleErrors(err, "Error Marshalling The FileSystem")
	_, err = Server.Write([]byte(jsonData))
	HandleErrors(err, "Error Writing To The Server")

}
