package cmd

import (
	"database/sql"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb-links/api"
	"github.com/Bnei-Baruch/mdb-links/common"
	"github.com/Bnei-Baruch/mdb-links/utils"
	"github.com/Bnei-Baruch/mdb-links/version"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A backend service for serving links to file backend",
	Run:   serverFn,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverFn(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.Info("Setting up connection to MDB")
	mdbDB, err := sql.Open("postgres", common.Config.MDBUrl)
	utils.Must(err)
	defer mdbDB.Close()
	//boil.DebugMode = true
	mdbDB.SetMaxIdleConns(common.Config.MDBMaxIdleConns)
	mdbDB.SetMaxOpenConns(common.Config.MDBMaxOpenConns)
	mdbDB.SetConnMaxIdleTime(common.Config.MDBMaxIdleTime)
	mdbDB.SetConnMaxLifetime(common.Config.MDBConnMaxLifetime)

	// read and validate config
	if len(common.Config.FilerUrls) == 0 {
		panic("No file service urls found")
	}
	for i, x := range common.Config.FilerUrls {
		if _, err := url.ParseRequestURI(x); err != nil {
			log.Fatalf("Malformed filer url [%d]: %s %s", i, x, err.Error())
		}
		log.Debug(x)
	}

	if _, err := url.ParseRequestURI(common.Config.BaseUrl); err != nil {
		log.Fatalf("Malformed base-url: %s", err.Error())
	}

	// Setup Rollbar
	rollbar.Token = common.Config.RollbarToken
	rollbar.Environment = common.Config.RollbarEnvironment
	rollbar.CodeVersion = version.Version

	// Setup gin
	gin.SetMode(common.Config.GinMode)
	router := gin.New()
	router.Use(
		utils.EnvironmentMiddleware(mdbDB),
		utils.ErrorHandlingMiddleware(),
		cors.New(cors.Config{
			AllowMethods:     []string{"GET"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
			AllowCredentials: false,
			AllowAllOrigins:  true,
			MaxAge:           12 * time.Hour,
		}),
		utils.RecoveryMiddleware())

	api.SetupRoutes(router)

	log.Infoln("Running mdb-links service")
	if cmd != nil {
		router.Run(common.Config.ListenAddress)
	}
}
