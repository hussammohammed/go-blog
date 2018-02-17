package controllers
import(
	//"fmt"
	"net/http"
	"github.com/hussammohammed/go-blog/services"
	"github.com/hussammohammed/go-blog/utilities"
	"github.com/hussammohammed/go-blog/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct{
userService services.IUserService
cryptUtil   utilities.ICryptUtil
}

func NewUserController(cryptUtil utilities.ICryptUtil, userService services.IUserService) *UserController {
	controller := UserController{}
	controller.cryptUtil = cryptUtil
	controller.userService = userService
	return &controller
}

func (controller UserController) Root(c *gin.Context) {
	c.JSON(http.StatusOK, true)
}

func (controller UserController) CreateUser(c *gin.Context) {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// e-mail exist
	err, _ = controller.userService.FindOne(&bson.M{"email": user.Email})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-mail address '" + user.Email + "' is already exits"})
		return
	}
	// name exist
	err, _ = controller.userService.FindOne(&bson.M{"username": user.UserName})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username '" + user.UserName + "' is already exits"})
		return
	}
	passwordLength := len(user.Password)
	if passwordLength < 3 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Passord length should be 3 charachters atleast"})
		return
	}
	err = controller.userService.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user.UserName, "email": user.Email})
}
