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

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb-links/common"
	"github.com/Bnei-Baruch/mdb-links/mdb/models"
	"github.com/Bnei-Baruch/mdb-links/utils"
)

type FileBackendRequest struct {
	SHA1     string `json:"sha1"`
	Name     string `json:"name"`
	ClientIP string `json:"clientip,omitempty"`
}

type FileBackendResponse struct {
	Url           string `json:"url"`
	IsAlternative bool   `json:"alternative"`
}

var filerClient = &http.Client{
	Timeout: 10 * time.Second,
}

func FilesHandler(c *gin.Context) {
	uid := c.Param("uid")

	if uid == "health_check" {
		HealthCheckHandler(c)
		return
	}

	resp, err := handleFile(c, uid, c.ClientIP())
	if err != nil {
		err.Abort(c)
		return
	}

	// client asked not to be redirected
	if nr, ok := c.GetQuery("no-redirect"); ok {
		cnr := strings.ToLower(strings.TrimSpace(nr))
		if cnr == "t" || cnr == "true" || cnr == "1" {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	// redirect type based on alternative status
	code := http.StatusFound
	if resp.IsAlternative {
		code = http.StatusMovedPermanently
	}

	c.Redirect(code, resp.Url)
}

func handleFile(cp utils.ContextProvider, uidParam string, clientIP string) (*FileBackendResponse, *utils.HttpError) {
	s := strings.Split(uidParam, ".") // strip file extension if provided
	uid := s[0]
	if len(uid) != 8 {
		return nil, utils.NewBadRequestError(errors.Errorf("Invalid UID: %s", uid))
	}

	db := cp.MustGet("MDB_DB").(*sql.DB)

	file, herr := lookupFile(db, uid, common.Config.PublicOnly)
	if herr != nil {
		return nil, herr
	}

	// are we redirecting to alternative file ?
	if file.UID != uid {
		resp := new(FileBackendResponse)
		ext := ""
		if len(s) > 1 {
			ext = fmt.Sprintf(".%s", strings.Join(s[1:], "."))
		}
		resp.Url = fmt.Sprintf("%s%s%s", common.Config.BaseUrl, file.UID, ext)
		resp.IsAlternative = true
		return resp, nil
	}

	// File seems reasonable. Proceed to filer backend
	sha1sum := hex.EncodeToString(file.Sha1.Bytes)
	filename := getFilename(file, db)
	body, herr := createRequestBody(sha1sum, filename, clientIP)
	if herr != nil {
		return nil, herr
	}

	var err error
	var res *http.Response
	for i, url := range common.Config.FilerUrls {
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

func lookupFile(db *sql.DB, uid string, publicOnly bool) (*mdbmodels.File, *utils.HttpError) {
	mods := []qm.QueryMod{
		qm.Select("id", "uid", "sha1", "content_unit_id", "name", "removed_at"),
		qm.Where("uid = ?", uid),
	}

	if publicOnly {
		mods = append(mods, qm.Where("secure = 0"))
	}

	file, err := mdbmodels.Files(mods...).One(db)
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

	if file.RemovedAt.Valid {
		log.Infof("File removed, look for alternative: %s", uid)
		file, err = lookupAlternative(file, db)
		if err != nil {
			return nil, utils.NewInternalError(errors.Wrap(err, "Lookup alternative file in MDB"))
		}
		if file == nil {
			log.Infof("No alternative file found")
			return nil, utils.NewNotFoundError()
		}
	}

	return file, nil
}

func lookupAlternative(file *mdbmodels.File, db *sql.DB) (*mdbmodels.File, error) {
	if err := file.Reload(db); err != nil {
		return nil, errors.Wrap(err, "reload file from MDB")
	}

	// alternative lookup makes sense only inside content units
	if !file.ContentUnitID.Valid {
		return nil, nil
	}

	mods := []qm.QueryMod{
		qm.Select("uid"),
		qm.Where("sha1 IS NOT NULL AND removed_at IS NULL"),     // physical, not removed file
		qm.And("id <> ?", file.ID),                              // not this file
		qm.And("secure <= ?", file.Secure),                      // at least secure as this one
		qm.And("published = ?", file.Published),                 // same published status
		qm.And("content_unit_id = ?", file.ContentUnitID.Int64), // in the same unit
		qm.And("type = ?", file.Type),                           // same type
		qm.And("language = ?", file.Language.String),            // same language
		qm.OrderBy("created_at desc"),                           // solves most mime_type / sub_type conflicts
	}

	alts, err := mdbmodels.Files(mods...).All(db)
	if err != nil {
		return nil, errors.Wrap(err, "fetch alternative files from MDB")
	}

	if len(alts) == 0 {
		return nil, nil
	}

	return alts[0], nil
}

// getFilename returns the final file name displayed to the end user.
// It takes whatever metadata needed from MDB to determine that.
// In case of error the original name of the file is returned.
func getFilename(file *mdbmodels.File, db *sql.DB) string {
	if file.ContentUnitID.IsZero() {
		return file.Name
	}

	cu, err := file.ContentUnit().One(db)
	if err != nil {
		log.Errorf("getFilename fetch CU from MDB [%d]: %v", file.ContentUnitID.Int64, err)
		return file.Name
	}

	var props map[string]interface{}
	if err := cu.Properties.Unmarshal(&props); err != nil {
		log.Errorf("getFilename unmarshal CU properties [%d]: %v", file.ContentUnitID.Int64, err)
		return file.Name
	}

	captureDate, cdOK := props["capture_date"]
	filmDate, fdOK := props["film_date"]

	// film_date and capture_date are not in dispute
	if !cdOK || !fdOK || captureDate == filmDate {
		return file.Name
	}

	return strings.Replace(file.Name, captureDate.(string), filmDate.(string), -1)
}

func createRequestBody(sha1sum, filename, clientIP string) ([]byte, *utils.HttpError) {
	data := FileBackendRequest{
		SHA1:     sha1sum,
		Name:     filename,
		ClientIP: clientIP,
	}

	jsonB, err := json.Marshal(data)
	if err != nil {
		return nil, utils.NewInternalError(errors.Wrap(err, "json.Marshal"))
	}

	return jsonB, nil
}

func callBackend(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return filerClient.Do(req)
}

func processResponse(res *http.Response) (*FileBackendResponse, *utils.HttpError) {

	// physical file doesn't exists
	if res.StatusCode == http.StatusNotFound {
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

	// physical file doesn't exists
	if res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusNotFound {
		return nil, utils.NewNotFoundError()
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
