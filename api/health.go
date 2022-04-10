package api

import (
	"database/sql"
	"database/sql/driver"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb-links/common"
)

func HealthCheckHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// ping mdb
	g.Go(func() error {
		mdb := c.MustGet("MDB_DB").(*sql.DB)
		return PingDB(ctx, mdb)
	})

	// ping filer backends
	var operatingBackends int32
	for _, x := range common.Config.FilerUrls {
		p, _ := url.Parse(x)
		p.Path = "/"
		u := p.String()
		g.Go(func() error {
			// we're not returning an error in this goroutine
			// As that would cancel the entire errgroup
			// instead we log it and count the number of healthy backends

			err := PingHttp(ctx, u)
			if err != nil {
				log.Errorf("filer backend %s: %s", u, err.Error())
			} else {
				atomic.AddInt32(&operatingBackends, 1)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil || operatingBackends == 0 {
		reason := "No filer backend is alive"
		if err != nil {
			reason = err.Error()
		}

		c.JSON(http.StatusFailedDependency, gin.H{
			"status": "error",
			"error":  reason,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// PingDB Temporary implementation until lib/pq PR is merged.
// See https://github.com/lib/pq/pull/737
func PingDB(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "select 1")
	if err != nil {
		return driver.ErrBadConn // https://golang.org/pkg/database/sql/driver/#Pinger
	}
	defer rows.Close()
	return nil
}

func PingHttp(ctx context.Context, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "http error")
	}

	// all is good - return
	if resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	// try to log the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response body")
	}
	log.Warnf("Bad dependency status: %d", resp.StatusCode)
	log.Warn(string(body))

	return errors.Errorf("Bad dependency status: %d", resp.StatusCode)
}
