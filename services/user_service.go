package services
import(
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	//"time"
	"github.com/hussammohammed/go-blog/utilities"
	"github.com/hussammohammed/go-blog/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IUserService  interface {
	FindById(id string) (error, *models.User)
	FindOne(query *bson.M) (error, *models.User)
	Find(query *models.User) (users []*models.User)
	Insert(user *models.User) error
}
type UserService struct {
	dbName            string
	uri               string
	dialMongoWithInfo string
	collectionName    string
	cryptUtil         utilities.ICryptUtil
	configUtil        utilities.IConfigUtil
	//sessionService    *SessionService
}

func NewUserService(configUtil utilities.IConfigUtil, cryptUtil utilities.ICryptUtil) *UserService {
	r := UserService{}
	r.uri = configUtil.GetConfig("dbUri")
	r.dbName = configUtil.GetConfig("dbName")
	r.dialMongoWithInfo = configUtil.GetConfig("dialMongoWithInfo")
	r.collectionName = "users"
	r.cryptUtil = cryptUtil
	r.configUtil = configUtil
	//r.sessionService = NewSessionService(configUtil, cryptUtil)
	return &r
}

func (r UserService) populateRole(user *models.User, session *mgo.Session) {
	roleCollection := session.DB(r.dbName).C("roles")
	roleCollection.FindId(user.RoleId).One(&user.Role)
}

func (r UserService) FindById(id string) (error, *models.User) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	user := models.User{}
	err := collection.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return err, nil
	}
	r.populateRole(&user, session)
	return nil, &user
}

func (r UserService) newSession() (*mgo.Session, error) {
	fmt.Println("connection is ")
	fmt.Println(r.uri)
	fmt.Println("is dial with info")
	fmt.Println(r.dialMongoWithInfo)
	if r.dialMongoWithInfo == "true" {
		tlsConfig := &tls.Config{}
		roots := x509.NewCertPool()
		path := r.configUtil.GetConfig("path")
		if ca, err := ioutil.ReadFile(path + "/ssh/mongo.pem"); err == nil {
			roots.AppendCertsFromPEM(ca)
		}
		tlsConfig.RootCAs = roots

		dialInfo, _ := mgo.ParseURL(r.uri)
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			fmt.Println("try connect")
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		//Here is the session you are looking for. Up to you from here ;)
		return mgo.DialWithInfo(dialInfo)
	}
	fmt.Println("without ssl")
	return mgo.Dial(r.uri)
}

func (r UserService) FindOne(query *bson.M) (error, *models.User) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	user := models.User{}
	err := collection.Find(query).One(&user)
	return err, &user
}

func (r UserService) Find(query *models.User) (users []*models.User) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	bsonQuery, _ := bson.Marshal(query)
	find := collection.Find(bsonQuery).Iter()
	user := models.User{}
	for find.Next(&user) {
		r.populateRole(&user, session)
		users = append(users, &user)
	}
	return users
}

func (r UserService) Insert(user *models.User) error {
	user.Id = bson.NewObjectId()
	//user.Slug = r.slugUtil.GetSlug(user.UserName)
	user.Password = r.cryptUtil.Bcrypt(user.Password)
	user.EmailVerified = false
	user.VerifyToken = r.cryptUtil.NewEncryptedToken()

	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Insert(user)
	// send verify email
	//globalVars := map[string]interface{}{"FNAME": user.FirstName, "ACTIVATE_ACCOUNT": "google.com/verify/registration?token=" + user.VerifyToken + "&email=" + user.Email}

	/*_, mailErr := r.mailUtil.SendTemplate(user.Email, "john@curtisdigital.com", "John Curtis", "Verify your e-mail", globalVars, 606061)
	if mailErr != nil {
		fmt.Println("verify email send error")
		fmt.Println(mailErr)
	}*/
	return err
}