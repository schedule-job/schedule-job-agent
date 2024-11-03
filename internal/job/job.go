package job

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Status             string
	RequestUrl         string
	RequestMethod      string
	RequestHeaders     map[string][]string
	RequestBody        string
	ResponseHeaders    map[string][]string
	ResponseBody       string
	ResponseStatusCode int
}

type Response struct {
	ID    string
	Log   *Request
	Res   *http.Response
	Error error
}

type Database interface {
	InsertRequestLog(jobID string, data interface{}) error
}

type Job struct {
	ID       string              `json:"id"`
	Url      string              `json:"url"`
	Method   string              `json:"method"`
	Body     string              `json:"body"`
	Headers  map[string][]string `json:"headers"`
	database Database
}

func (j *Job) SetDatabase(database Database) {
	j.database = database
}

func (j Job) getDefaultLog() Request {
	var log = Request{}
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
		fmt.Println("로그 삽입 실패:", err)
	}
}

func (j Job) failed(err error) {
	log := j.getDefaultLog()
	log.Status = "failed"
	log.ResponseBody = err.Error()
	err = j.database.InsertRequestLog(j.ID, log)
	if err != nil {
		fmt.Println("로그 삽입 실패 :", err)
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
		fmt.Println("로그 삽입 실패 :", err)
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
