package main

import (
	jsonparse "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type (
	// args hold argument passed to JSON RPC server
	Args struct {
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
	TwitterProfile struct {
		Id        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Username  string `json:"username,omitempty"`
		Followers string `json:"followers,omitempty"`
		Following string `json:"following,omitempty"`
	}
	JsonServer struct{}
)

func (s *JsonServer) GetProfileDetail(r *http.Request, args *Args, reply *TwitterProfile) error {
	var profiles []TwitterProfile
	raw, err := ioutil.ReadFile("./twitterProfile.json")
	if err != nil {
		log.Fatal("error: ", err)
	}
	err = jsonparse.Unmarshal(raw, &profiles)
	if err != nil {
		log.Fatal("unmarshal error: ", err)
	}
	for _, profile := range profiles {
		if profile.Id == args.Id {
			fmt.Println(profile)
			*reply = profile
			break
		}
	}
	return nil
}

func main() {
	// create a new RPC server
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	// register RPC server
	s.RegisterService(new(JsonServer), "")
	r := mux.NewRouter()
	// NOTE: can I register just a function, no method from object or receiver func?
	// short: kinda can't, ref: https://go.dev/src/net/rpc/server.go
	r.Handle("/rpc", s)
	http.ListenAndServe(":9000", r)
}
