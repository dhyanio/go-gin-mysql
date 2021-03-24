package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)
type Version struct {
	TemplateID  string `json:"templateID"`
	ID          int16  `json:"id"`
	Provider    string `json:"provider"`
	CreatedTime string `json:"createdTime"`
	Creator     string `json:"creator"`
	VersionTag  string `json:"versionTag"`
}
func errHandler(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
// updateTag API
func latestVersion(c *gin.Context) {
	c.JSON(http.StatusOK, GetTag())
}
// updateTag api
func updateVersion(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	errHandler(err)
	var postVersion Version
	json.Unmarshal(jsonData, &postVersion)
	c.JSON(http.StatusOK, PostTag(postVersion.TemplateID, postVersion.Provider, postVersion.Creator, postVersion.VersionTag, postVersion.CreatedTime, postVersion.ID))
}
// indexTag api
func indexVersion(c *gin.Context) {
	c.JSON(http.StatusOK, GetTag())
}
func PostTag(templateID, provider, creator, versionTag, createdTime string, id int16) string {
	host := "
	database := ""
	user := ""
	password := ""
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	errHandler(err)
	defer db.Close()
	prepStatement, err := db.Prepare("INSERT INTO version(id, creator, templateID, provider, versionTag, createdTime) VALUES(?,?,?,?,?,?)")
	errHandler(err)
	prepStatement.Exec(id, creator, templateID, provider, versionTag, createdTime)
	defer prepStatement.Close()
	return "Inserted"
}
func GetTag() []Version {
	host := "
	database := ""
	user := ""
	password := ""
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	errHandler(err)
	defer db.Close()
	selDB, err := db.Query()
	errHandler(err)
	emp := Version{}
	res := []Version{}
	for selDB.Next() {
		var id int16
		var creator, versionTag, templateID, provider, createdTime string
		err = selDB.Scan(&id, &creator, &templateID, &provider, &versionTag, &createdTime)
		errHandler(err)
		emp.TemplateID = templateID
		emp.ID = id
		emp.CreatedTime = createdTime
		emp.Creator = creator
		emp.VersionTag = versionTag
		emp.Provider = provider
		res = append(res, emp)
	}
	return res
}
func main() {
	router := gin.Default()
	router.GET("/abc", latestVersion)
	router.POST("/abc", updateVersion)
	router.Run()
}
