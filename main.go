package main

import (
	"github.com/soider/schnur/targets"
	"github.com/soider/schnur/targets/yaml"
	"golang.org/x/crypto/ssh"

	"flag"
	"github.com/gorilla/mux"
	"github.com/soider/d"
	"github.com/soider/schnur/handlers"
	"github.com/soider/schnur/targets/manager"
	"net/http"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type clientBuildFunction func(cfg *ssh.ClientConfig) (*ssh.Client, error)

type Schnur struct {
	address string
	Targets targets.TargetList
	Loader  interface {
		Targets() (targets.TargetList, error)
	}

	mux *mux.Router
}

func (s *Schnur) Init() {
	tm := manager.New(s.Loader)
	s.mux = mux.NewRouter()
	s.mux.HandleFunc("/ws/{target:.+}", handlers.Proxy)
	s.mux.Handle("/api/preconnect/{target:.+}", handlers.NewPrepareHandler(tm))
	s.mux.Handle("/api/list", handlers.NewListHandler(tm))
	s.mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

}

func (s *Schnur) Run() {
	serv := http.Server{
		Addr:    s.address,
		Handler: s.mux,
	}
	d.D("Listeon on", s.address)
	panic(serv.ListenAndServe())
}

func main() {
	var targetsPath string
	var address string
	flag.StringVar(&targetsPath, "targets", "targets.yaml", "Path to yaml file with targets")
	flag.StringVar(&address, "address", "0.0.0.0:8080", "Listen on")
	flag.Parse()
	app := &Schnur{
		Loader:  yaml.New(targetsPath),
		address: address,
	}

	app.Init()
	app.Run()
	t, err := app.Loader.Targets()
	d.D(t, err)
}
