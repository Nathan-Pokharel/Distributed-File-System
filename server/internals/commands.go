package internals

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"text/template"

)

var CurrentClient net.Conn

func FileSystem(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")

	data := struct {
		Files   []string
		Folders []string
		Cwd     string
	}{

    }

	for conn := range Connections {
		remoteAddr := conn.RemoteAddr().String()
		if remoteAddr == segments[2] {
			CurrentClient = conn
			connData := Connections[conn]
			data.Files = connData.Files
			data.Folders = connData.Folders
			data.Cwd = connData.Cwd
			break
		}
	}
	tmpl, err := template.ParseFiles("/home/nathan/Projects/Distributed_File_System/server/internals/templates/layout.html", "/home/nathan/Projects/Distributed_File_System/server/internals/templates/filesystem.html")
	HandleErrors(err, "Error Parsing Files")
	if err := tmpl.Execute(w, data); err != nil {
		HandleErrors(err, "Error Executing Template")
	}
}
func Index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Connections map[net.Conn]ConnectionData
	}{
		Connections,
	}
	tmpl, err := template.ParseFiles("/home/nathan/Projects/Distributed_File_System/server/internals/templates/layout.html", "/home/nathan/Projects/Distributed_File_System/server/internals/templates/connections.html")
	HandleErrors(err, "Error Parsing Files")
	if err := tmpl.Execute(w, data); err != nil {
		HandleErrors(err, "Error Executing Template")
	}
}

func UpdatePaths(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Define a struct to hold the JSON data
	var requestData struct {
		Cwd string `json:"cwd"`
	}

	// Parse the JSON data from the request body
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Access the cwd value
	cwd := requestData.Cwd

    
    WriteToClient(CurrentClient,cwd)

	// Respond with a success status
	w.WriteHeader(http.StatusOK)
}

func InitWebServer() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/fs/", FileSystem)
	http.HandleFunc("/ufs/", UpdatePaths)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Server Startup Failed: ", err)
	}
}
