package internals

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)
func ReadFromClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
    message = message[:len(message)-1]
	if err != nil {
		fmt.Println("Error Reading From Connection: ", conn.RemoteAddr())
		Mutex.Lock()
		delete(Connections, conn)
		Mutex.Unlock()
		return
	}
    err = json.Unmarshal([]byte(message),&conndata)
    if err !=nil{
        HandleErrors(err,"Error Unmarshalling Data") 
    }
	Mutex.Lock()
	Connections[conn] = conndata
    Mutex.Unlock()
}
func WriteToClient(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error Writing To Client: ", conn.RemoteAddr())
		Mutex.Lock()
		delete(Connections, conn)
		Mutex.Unlock()
		return
	}
    go func(conn net.Conn){
        ReadFromClient(conn)
    }(conn)
}

