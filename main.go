package main

import(
	"log"
	"github.com/gin-gonic/gin"
)

func getPort () string {
port := "2020"
return ":" + port
}

func main(){
	router:= gin.Default()
	port:=getPort()
	log("port is" + port)
	router.Run(port)
}