package common

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	ListenAddress      string
	GinMode            string
	RollbarToken       string
	RollbarEnvironment string
	BaseUrl            string
	FilerUrls          []string
	PublicOnly         bool
	MDBUrl             string
	MDBMaxIdleConns    int
	MDBMaxOpenConns    int
	MDBMaxIdleTime     time.Duration
	MDBConnMaxLifetime time.Duration
}

func newConfig() *config {
	return &config{
		ListenAddress:      ":8081",
		GinMode:            "debug",
		RollbarToken:       "",
		RollbarEnvironment: "development",
		BaseUrl:            "http://localhost:8081/",
		FilerUrls:          []string{"http://files.kabbalahmedia.info/api/v1/get"},
		PublicOnly:         true,
		MDBUrl:             "postgres://user:password@localhost/mdb?sslmode=disable",
		MDBMaxIdleConns:    2,
		MDBMaxOpenConns:    5,
		MDBMaxIdleTime:     5 * time.Minute,
		MDBConnMaxLifetime: 60 * time.Minute,
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
	if val := os.Getenv("FILER_URLS"); val != "" {
		Config.FilerUrls = strings.Split(val, ",")
	}
	if val := os.Getenv("PUBLIC_ONLY"); val != "" {
		Config.PublicOnly = val == "true"
	}
	if val := os.Getenv("MDB_URL"); val != "" {
		Config.MDBUrl = val
	}
	if val := os.Getenv("MDB_MAX_IDLE_CONNS"); val != "" {
		x, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Malformed MDB_MAX_IDLE_CONNS [%s]: %v", val, err)
		} else {
			Config.MDBMaxIdleConns = x
		}
	}
	if val := os.Getenv("MDB_MAX_OPEN_CONNS"); val != "" {
		x, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Malformed MDB_MAX_OPEN_CONNS [%s]: %v", val, err)
		} else {
			Config.MDBMaxOpenConns = x
		}
	}
	if val := os.Getenv("MDB_MAX_IDLE_TIME"); val != "" {
		x, err := time.ParseDuration(val)
		if err != nil {
			log.Fatalf("Malformed MDB_MAX_IDLE_TIME [%s]: %v", val, err)
		} else {
			Config.MDBMaxIdleTime = x
		}
	}
	if val := os.Getenv("MDB_CONN_MAX_LIFETIME"); val != "" {
		x, err := time.ParseDuration(val)
		if err != nil {
			log.Fatalf("Malformed MDB_CONN_MAX_LIFETIME [%s]: %v", val, err)
		} else {
			Config.MDBConnMaxLifetime = x
		}
	}
}
