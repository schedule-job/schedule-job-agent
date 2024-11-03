package job

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/schedule-job/schedule-job-database/core"
	schedule_errors "github.com/schedule-job/schedule-job-errors"
)

type Response struct {
	ID    string
	Log   *core.RequestTypePayload
	Res   *http.Response
	Error error
}

type Job struct {
	ID       string              `json:"id"`
	Url      string              `json:"url"`
	Method   string              `json:"method"`
	Body     string              `json:"body"`
	Headers  map[string][]string `json:"headers"`
	database core.Database
}

func (j *Job) SetDatabase(database core.Database) {
	j.database = database
}

func (j Job) getDefaultLog() core.RequestTypePayload {
	var log = core.RequestTypePayload{}
	log.Status = "progress"
	log.RequestUrl = j.Url
	log.RequestMethod = j.Method
	log.RequestHeaders = j.Headers
	log.RequestBody = j.Body

	return log
}

func (j Job) done(res *http.Response, err error) {
	if res == nil {
		j.failed(err)
		return
	}

	log := j.getDefaultLog()

	if err != nil {
		j.requestFailed(res, err)
		return
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		j.requestFailed(res, err)
		return
	}

	body, bodyErr := io.ReadAll(res.Body)
	if bodyErr != nil {
		j.failed(bodyErr)
		return
	}

	log.Status = "succeed"
	log.ResponseBody = string(body)
	log.ResponseHeaders = res.Header
	log.ResponseStatusCode = res.StatusCode

	err = j.database.InsertRequestLog(j.ID, log)
	if err != nil {
		write := schedule_errors.LogWriteError{Reason: err.Error()}
		fmt.Println(write.Error())
	}
}

func (j Job) failed(err error) {
	log := j.getDefaultLog()
	log.Status = "failed"
	log.ResponseBody = err.Error()
	err = j.database.InsertRequestLog(j.ID, log)
	if err != nil {
		write := schedule_errors.LogWriteError{Reason: err.Error()}
		fmt.Println(write.Error())
	}
}

func (j Job) requestFailed(res *http.Response, err error) {
	log := j.getDefaultLog()
	log.Status = "failed"
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, bodyErr := io.ReadAll(res.Body)
		if bodyErr != nil {
			j.failed(bodyErr)
			return
		}
		log.ResponseBody = string(body)
	} else {
		log.ResponseBody = err.Error()
	}
	log.ResponseHeaders = res.Header
	log.ResponseStatusCode = res.StatusCode
	err = j.database.InsertRequestLog(j.ID, log)
	if err != nil {
		write := schedule_errors.LogWriteError{Reason: err.Error()}
		fmt.Println(write.Error())
	}
}

// Run : Job를 실행합니다.
func (j Job) Run() {
	j.database.InsertRequestLog(j.ID, j.getDefaultLog())

	method := strings.ToUpper(j.Method)
	req, reqErr := http.NewRequest(method, j.Url, strings.NewReader(j.Body))

	if reqErr != nil {
		j.done(nil, reqErr)
		return
	}

	req.Header = j.Headers

	client := &http.Client{}
	res, resErr := client.Do(req)

	j.done(res, resErr)
}
