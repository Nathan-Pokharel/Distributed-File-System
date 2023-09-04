package main
import (
    "github.com/Nathan-Pokharel/DistributedFS/server/internals"

)

func main(){
    go internals.InitWebServer()
    internals.AcceptConnections()

}
