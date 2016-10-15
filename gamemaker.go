package main

import (
	"github.com/danmondy/gamemaker/data"
	"github.com/danmondy/gamemaker/site"

	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	CONFIG_PATH = "/opt/gm/config.json"
)

var cfg data.Config

func main() {

	//site
	r := mux.NewRouter()

	site.SetRoutes(r)

	fs := http.FileServer(http.Dir(cfg.Docroot))

	r.PathPrefix("/img").Handler(fs)
	r.PathPrefix("/css").Handler(fs)
	r.PathPrefix("/js").Handler(fs)

	fmt.Printf("Listening on %s:%s", cfg.Host, cfg.Port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), r)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	readConfig()
	//INIT DB's
	/*err := data.InitDb(cfg.Db)
	if err != nil {
		panic(err)
	}*/

	//INIT Site
	site.InitSite(cfg)
	//site.InitSite("public")
	//fmt.Println("Connected to the", cfg.Db.Name, "database.")
}

func readConfig() {
	dat, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		panic(err)
	}
}
