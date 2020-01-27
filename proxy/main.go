package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
	"os"
)
type Config struct {
	Target string `json:"target"`
	LocalPort string 	`json:"local_port"`
}
type Proxy struct {

}
var (
	config = &Config{}
)
func (P *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL, _ = url.Parse(config.Target + r.RequestURI)
	res, _ := http.DefaultTransport.RoundTrip(r)
	//copy header
	for key, v := range res.Header{
		for _, vv := range v {
			w.Header().Add(key, vv)
		}
	}
	defer res.Body.Close()
	//copy body
	io.Copy(w, res.Body)
}

func main() {

	f, err := os.Open("./config.json")
	defer f.Close()
	if err != nil {
		panic("Open config file error.")
	}

	bf, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Read config content error.")
	}

	json.Unmarshal(bf, config)
	if config.Target == "" {
		panic("Target is required.")
	}
	//SET DEFAULT PORT
	if config.LocalPort == "" {
		config.LocalPort = ":80"
	}
	if config.LocalPort != "" && config.LocalPort[0:1] != ":" {
		config.LocalPort = ":" + config.LocalPort
	}

	var P = new(Proxy)
	http.ListenAndServe(config.LocalPort, P)
}