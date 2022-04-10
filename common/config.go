package common

import (
	"os"
	"strings"
)

type config struct {
	ListenAddress      string
	GinMode            string
	RollbarToken       string
	RollbarEnvironment string
	BaseUrl            string
	MDBUrl             string
	FilerUrls          []string
	PublicOnly         bool
}

func newConfig() *config {
	return &config{
		ListenAddress:      ":8081",
		GinMode:            "debug",
		RollbarToken:       "",
		RollbarEnvironment: "development",
		BaseUrl:            "http://localhost:8081/",
		MDBUrl:             "postgres://user:password@localhost/mdb?sslmode=disable",
		FilerUrls:          []string{"http://files.kabbalahmedia.info/api/v1/get"},
		PublicOnly:         true,
	}
}

var Config *config

func Init() {
	Config = newConfig()

	if val := os.Getenv("LISTEN_ADDRESS"); val != "" {
		Config.ListenAddress = val
	}
	if val := os.Getenv("GIN_MODE"); val != "" {
		Config.GinMode = val
	}
	if val := os.Getenv("ROLLBAR_TOKEN"); val != "" {
		Config.RollbarToken = val
	}
	if val := os.Getenv("ROLLBAR_ENVIRONMENT"); val != "" {
		Config.RollbarEnvironment = val
	}
	if val := os.Getenv("BASE_URL"); val != "" {
		Config.BaseUrl = val
	}
	if val := os.Getenv("MDB_URL"); val != "" {
		Config.MDBUrl = val
	}
	if val := os.Getenv("FILER_URLS"); val != "" {
		Config.FilerUrls = strings.Split(val, ",")
	}
	if val := os.Getenv("PUBLIC_ONLY"); val != "" {
		Config.PublicOnly = val == "true"
	}
}
