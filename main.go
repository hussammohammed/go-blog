package main

import(
	"log"
	"github.com/hussammohammed/go-blog/controllers"
	"github.com/hussammohammed/go-blog/services"
	"github.com/hussammohammed/go-blog/utilities"
	"github.com/gin-gonic/gin"
)

func getPort () string {
	confiUtil := utilities.NewConfigUtil()
	port := confiUtil.GetConfig("port")
	if port == "" {
		port = "2020"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}

func main(){
router:= gin.Default()

cryptUtil := utilities.NewCryptUtil()
configUtil := utilities.NewConfigUtil()

userService := services.NewUserService(configUtil, cryptUtil)
usersController := controllers.NewUserController(cryptUtil, userService)
router.GET("/", usersController.Root)
router.POST("/api/signup", usersController.CreateUser)

port:= getPort()
log.Println("Port is " + port)
router.Run(port)
}