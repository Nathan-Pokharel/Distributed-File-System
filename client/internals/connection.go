package internals

import (
	"fmt"
	"net"
)

var Server net.Conn
var err error

func HandleErrors(err error, message string) {
	if err != nil {
		fmt.Println(message, err)
		panic(message)
	}

}
func ConnectToServer() {
	Server, err = net.Dial("tcp", ":8080")
	HandleErrors(err, "Error Connecting To The Server")
	defer Server.Close()
	fmt.Println("Connected To The Server")
    
    ReadFromServer()
}
