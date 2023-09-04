package internals

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionData struct {
	Files   []string
	Folders []string
	Cwd     string
}

var conndata ConnectionData
var (
	Connections = make(map[net.Conn]ConnectionData)
	Mutex       sync.Mutex
)

func HandleErrors(err error, message string) {
	if err != nil {
		fmt.Println(message, err)
		panic(message)
	}
}

// Initialise Server And Start Accepting Connections
func AcceptConnections() {

	listner, err := net.Listen("tcp", ":8080")
	HandleErrors(err, "Error Listening")
	defer listner.Close()
	fmt.Println("Initialised Server Listening on :8080")
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection")
			continue
		}
		fmt.Println(conn, " Connected To Server")
			go func(conn net.Conn) {
				Mutex.Lock()
				Connections[conn] = ConnectionData{}
				WriteToClient(conn, "Ready")
				Mutex.Unlock()
			}(conn)
	}

}
