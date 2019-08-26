package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type jsonUser struct {
	Status        int    `json:"status"`
	Lock          string `json:"lock"`
	Param         string `json:"param[qq]"`
	Hash_sha256   string `json:"hash_sha256"`
	Hash_md5      string `json:"hash_md5"`
	Guid          string `json:"guid"`
	Guid_location string `json:"guid_location"`
	Hash_sha512   string `json:"hash_sha512"`
	H_sha512_gen  string `json:"h_sha512_gen"`
	Hash_h        string `json:"hash_h"`
}
type msg string

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
func (msg) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var param string = "null"
	buf := string(req.URL.String())
	params, _ := url.Parse(buf)
	params_ := params.Query().Get("qq")
	if params_ != "" {
		param = params_
	}
	Hash_sha256 := sha256.Sum256([]byte(buf))
	Hash_md5 := md5.Sum([]byte(buf))
	Guid, _ := uuid.NewV4()
	Guid_location := uuid.NewV5(uuid.Nil, buf)
	Hash_sha512 := sha512.Sum512([]byte(buf))
	H_sha512_gen := sha512.Sum512([]byte(RandStringRunes()))
	Hash_h := sha256.Sum256([]byte(buf + hex.EncodeToString(Hash_md5[:]) + strconv.FormatInt(time.Now().Unix(), 10)))
	jsonA := &jsonUser{
		Status:        1,
		Lock:          buf,
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
func main() {
	msgHandler := msg("L")
	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:10800", msgHandler)
}
