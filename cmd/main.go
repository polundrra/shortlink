package main

import (
	"github.com/polundrra/shortlink/internal/api"
	"github.com/polundrra/shortlink/internal/repo"
	"github.com/polundrra/shortlink/internal/service"
	"github.com/valyala/fasthttp"
	"log"
	"github.com/BurntSushi/toml"
	"time"
)

type conf struct {
	ServerPort string
	ReadTimeout int
	IdleTimeout int
	WriteTimeout int
	RepoOpts repo.Opts
}

func main() {
	var conf conf
	if _, err := toml.DecodeFile("/etc/shortlink/conf.toml", &conf); err != nil {
		log.Fatal(err)
	}
	repo, err := repo.New(conf.RepoOpts)
	if err != nil {
		log.Fatal(err)
	}

	service := service.New(conf.RepoOpts, repo)

	api := api.New(service)

	server := fasthttp.Server{
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		Handler:      api.Router(),
	}

	server.ListenAndServe(conf.ServerPort)
}
