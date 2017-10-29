package api

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb-links/mdb/models"
	"github.com/Bnei-Baruch/mdb-links/utils"
)

type FileBackendRequest struct {
	SHA1     string `json:"sha1"`
	Name     string `json:"name"`
	ClientIP string `json:"clientip,omitempty"`
}

type FileBackendResponse struct {
	Url string `json:"url"`
}

var filerClient = &http.Client{
	Timeout: time.Second,
}

func FilesHandler(c *gin.Context) {
	uid := c.Param("uid")
	uid = strings.Split(uid, ".")[0] // ignore file extension
	if len(uid) != 8 {
		utils.NewBadRequestError(errors.New("Invalid UID")).Abort(c)
		return
	}

	db := c.MustGet("MDB_DB").(*sql.DB)
	urls := c.MustGet("BACKEND_URLS").([]string)
	resp, err := handleFile(db, urls, uid, c.ClientIP())
	if err != nil {
		err.Abort(c)
		return
	}

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusFound, resp.Url)
	}
}

func handleFile(db *sql.DB, urls []string, uid string, clientIP string) (*FileBackendResponse, *utils.HttpError) {
	body, ex := createRequestBody(db, uid, clientIP)
	if ex != nil {
		return nil, ex
	}

	var err error
	var res *http.Response
	for i, url := range urls {
		log.Infof("Calling backend number %d", i+1)
		res, err = callBackend(url, body)
		if err != nil || res.StatusCode >= http.StatusMultipleChoices {
			continue
		}
		break
	}

	if err != nil {
		return nil, utils.NewInternalError(errors.Wrapf(err, "Communication error"))
	}

	return processResponse(res)
}

func createRequestBody(db *sql.DB, uid string, clientIP string) (*bytes.Buffer, *utils.HttpError) {
	file, err := mdbmodels.Files(db,
		qm.Select("sha1", "content_unit_id", "name"),
		qm.Where("uid = ?", uid)).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError()
		} else {
			return nil, utils.NewInternalError(errors.Wrap(err, "Lookup file in MDB"))
		}
	}

	if !file.Sha1.Valid {
		return nil, utils.NewBadRequestError(errors.New("Not a physical file"))
	}

	data := FileBackendRequest{
		SHA1:     hex.EncodeToString(file.Sha1.Bytes),
		Name:     file.Name,
		ClientIP: clientIP,
	}

	log.Infof("File exists in MDB: %s %s", uid, data.SHA1)

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(data)
	if err != nil {
		return nil, utils.NewInternalError(errors.Wrap(err, "json.Encode request"))
	}

	return b, nil
}

func callBackend(url string, b *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return filerClient.Do(req)
}

func processResponse(res *http.Response) (*FileBackendResponse, *utils.HttpError) {

	// physical file doesn't exists
	if res.StatusCode == http.StatusNoContent {
		log.Infof("Files backend no-content")
		return nil, utils.NewNotFoundError()
	}

	// Files backend crushed
	if res.StatusCode >= http.StatusInternalServerError {
		return nil, utils.NewHttpError(
			http.StatusFailedDependency,
			errors.Errorf("Files backend crashed: [%d - %s] %s",
				res.StatusCode, http.StatusText(res.StatusCode), res.Status),
			gin.ErrorTypePrivate,
		)
	}

	defer res.Body.Close()

	// Physical file exists
	if res.StatusCode == http.StatusOK {
		var body FileBackendResponse
		err := json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, utils.NewInternalError(errors.Wrap(err, "json.Decode response"))
		}
		return &body, nil
	}

	// Unexpected response (maybe some 400's ?)
	// Anyway, we shouldn't be here...
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, utils.NewInternalError(errors.Wrap(err, "ioutil.ReadAll response"))
	}

	msg := fmt.Sprintf("Unexpected response [%d - %s] %s",
		res.StatusCode, http.StatusText(res.StatusCode), res.Status)
	log.Error(msg)
	log.Errorf("res.Body: %s", b)

	return nil, utils.NewInternalError(errors.Errorf(msg))
}
