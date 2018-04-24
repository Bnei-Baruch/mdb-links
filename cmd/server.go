package cmd

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"github.com/volatiletech/sqlboiler/boil"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"fmt"
	"github.com/Bnei-Baruch/mdb-links/api"
	"github.com/Bnei-Baruch/mdb-links/utils"
	"github.com/Bnei-Baruch/mdb-links/version"
	"net/url"
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
	mdbDB, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	defer mdbDB.Close()
	boil.DebugMode = viper.GetString("server.mode") == "debug"

	// Verify file service backend urls
	backendUrls := viper.GetStringSlice("file_service.urls")
	if len(backendUrls) == 0 {
		panic("No file service urls found")
	}
	for i, x := range backendUrls {
		if _, err := url.ParseRequestURI(x); err != nil {
			panic(fmt.Sprintf("Bad Url [%d]: %s", i, x))
		}
		log.Debug(x)
	}

	// Setup Rollbar
	rollbar.Token = viper.GetString("server.rollbar-token")
	rollbar.Environment = viper.GetString("server.rollbar-environment")
	rollbar.CodeVersion = version.Version

	publicOnly := viper.GetBool("permissions.public-only")

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router := gin.New()
	router.Use(
		utils.EnvironmentMiddleware(mdbDB, backendUrls, publicOnly),
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
		router.Run(viper.GetString("server.bind-address"))
	}
}
