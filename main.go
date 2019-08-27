package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mux "github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type httpService struct {
	server *http.Server
}

type jsonUser struct {
	Status        int    `json:"status"`
	Uri           string `json:"uri"`
	Location      string `json:"location"`
	PostBody      string `json:"postBody"`
	Param         string `json:"param[qq]"`
	Hash_sha256   string `json:"hash_sha256"`
	Hash_md5      string `json:"hash_md5"`
	Guid          string `json:"guid"`
	Guid_location string `json:"guid_location"`
	Hash_sha512   string `json:"hash_sha512"`
	H_sha512_gen  string `json:"h_sha512_gen"`
	Hash_h        string `json:"hash_h"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(3000)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func (this *httpService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var param string = "null"
	buf := string(req.URL.String())
	params, _ := url.Parse(buf)
	params_ := params.Query().Get("qq")
	if params_ != "" {
		param = params_
	}
	PostBody, _ := ioutil.ReadAll(req.Body)
	Hash_sha256 := sha256.Sum256([]byte(buf))
	Hash_md5 := md5.Sum([]byte(buf))
	Guid, _ := uuid.NewV4()
	Guid_location := uuid.NewV5(uuid.Nil, buf)
	Hash_sha512 := sha512.Sum512([]byte(buf))
	H_sha512_gen := sha512.Sum512([]byte(RandStringRunes()))
	Hash_h := sha256.Sum256([]byte(buf + hex.EncodeToString(Hash_md5[:]) + strconv.FormatInt(time.Now().Unix(), 10)))
	jsonA := &jsonUser{
		Status:        1,
		Uri:           buf,
		Location:      req.URL.EscapedPath(),
		PostBody:      string(PostBody),
		Param:         param,
		Hash_sha256:   hex.EncodeToString(Hash_sha256[:]),
		Hash_md5:      hex.EncodeToString(Hash_md5[:]),
		Guid:          Guid.String(),
		Guid_location: Guid_location.String(),
		Hash_sha512:   hex.EncodeToString(Hash_sha512[:]),
		H_sha512_gen:  hex.EncodeToString(H_sha512_gen[:]),
		Hash_h:        hex.EncodeToString(Hash_h[:])}
	jsonB, _ := json.Marshal(jsonA)
	fmt.Fprint(resp, string(jsonB))
}
func (this *httpService) Start() {
	if this.server != nil {
		panic("Server already started.")
	}

	router := mux.NewRouter()

	router.HandleFunc("/test", this.ServeHTTP).Methods("GET", "POST")

	this.server = &http.Server{
		Addr:    "localhost:10800",
		Handler: router,
	}

	go func() {
		if err := this.server.ListenAndServe(); err != nil {
			panic("Failed to listen.")
		}
	}()
}
func main() {
	msgHandler := httpService{}
	fmt.Println("Server is listening...")
	msgHandler.Start()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C for quit.")

	<-c
}
