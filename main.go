package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	config "github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var codec Codec
var text = `
  >===>                >=>     >=> >======>     >=>          >=>>=>                                 >=>                                        
>>    >=>              >=>     >=> >=>    >=>   >=>        >=>    >=> >=>                           >=>                                        
>=>           >=>      >=>     >=> >=>    >=>   >=>         >=>       >=>         >=>     >> >==> >=>>==>   >==>    >==>>==>    >==>    >> >==>
>=>         >=>  >=>   >=>     >=> >> >==>      >=>           >=>     >=>>=>    >=>  >=>   >=>      >=>   >>   >=>   >=>  >=> >>   >=>   >=>   
>=>   >===> >=>  >=>   >=>     >=> >=>  >=>     >=>              >=>  >=>  >=> >=>    >=>  >=>      >=>   >>===>>=>  >=>  >=> >>===>>=>  >=>   
>=>    >>   >=>  >=>   >=>     >=> >=>    >=>   >=>        >=>    >=> >>   >=>  >=>  >=>   >=>      >=>   >>         >=>  >=> >>         >=>   
  >====>      >=>        >====>    >=>      >=> >=======>    >=>>=>   >=>  >=>    >=>     >==>       >=>   >====>   >==>  >=>  >====>   >==>

Running on port : %s`

func init() {
	codec = newBase64Codec()

	// initialize config service //
	config.SetConfigName("App")
	config.AddConfigPath("config")
	err := config.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conn := OpenDB()
	conn.AutoMigrate(&Shortener{})
	defer conn.Close()
}

func OpenDB() *gorm.DB {
	var (
		host     = config.GetString("database.host")
		port     = config.GetString("database.port")
		username = config.GetString("database.username")
		password = config.GetString("database.password")
		database = config.GetString("database.database")
	)

	address := fmt.Sprintf("tcp(%s:%s)", host, port)

	// open a db connection //
	conn, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, password, address, database))
	if err != nil {
		fmt.Println("Failed to connect database : ", err.Error())
	}

	return conn
}

type Codec interface {
	Encode(string) string
	Decode(string) (string, error)
}

type Base64Codec struct {
	e *base64.Encoding
}

func (b Base64Codec) Encode(s string) string {
	str := base64.URLEncoding.EncodeToString([]byte(s))
	return strings.Replace(str, "=", "", -1)
}

func (b Base64Codec) Decode(s string) (string, error) {
	if l := len(s) % 4; l != 0 {
		s += strings.Repeat("=", 4-l)
	}
	str, err := base64.URLEncoding.DecodeString(s)
	return string(str), err
}

func newBase64Codec() Base64Codec {
	return Base64Codec{base64.URLEncoding}
}

func main() {
	// show text //
	fmt.Println(fmt.Sprintf(text, config.GetString("app.port")))

	router := mux.NewRouter()
	router.HandleFunc("/shortener/create", handleCreate).Methods("POST")
	router.HandleFunc("/shortener/{id}", handleFind).Methods("GET")
	s := &http.Server{
		Addr:         ":" + config.GetString("app.port"),                              // set our http listener port
		Handler:      router,                                                          // set our request handler
		ReadTimeout:  time.Duration(config.GetInt("app.read_timeout")) * time.Second,  // set our service read timeout (5s default)
		WriteTimeout: time.Duration(config.GetInt("app.write_timeout")) * time.Second, // set our service write timeout (5s default)
	}

	log.Fatal(s.ListenAndServe())
}

type Shortener struct {
	Id  int    `gorm:"column:id;not null" json:"id"`
	Url string `gorm:"column:url;not null" json:"url"`
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var jsonReq Shortener

	// First, decode json data from post body request //
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jsonReq)

	if len(jsonReq.Url) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "missing url param")
		return
	}

	conn := OpenDB()

	err = conn.Save(&jsonReq).Error
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(201)
	fmt.Fprintf(w, "Location : http://"+hostname+":"+config.GetString("app.port")+"/shortener/"+codec.Encode(strconv.Itoa(jsonReq.Id)))
}

func handleFind(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	url := param["id"]

	result, err := getRecord(url)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "not found")
		return
	}

	w.WriteHeader(302)
	fmt.Fprintf(w, "Location : "+result)
}

func getRecord(id string) (url string, err error) {
	var data Shortener
	rId, err := codec.Decode(id)
	if err != nil {
		return "", err
	}

	conn := OpenDB()
	err = conn.Raw("SELECT id, url from shorteners where id = ?", rId).Scan(&data).Error
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	defer conn.Close()
	return data.Url, nil
}
