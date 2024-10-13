package log

import "time"

type Log struct {
	Id                 string    `json:"id"`
	JobId              string    `json:"jobId"`
	Status             string    `json:"status"`
	RequestUrl         string    `json:"requestUrl"`
	RequestMethod      string    `json:"requestMethod"`
	ResponseStatusCode int       `json:"responseStatusCode"`
	CreatedAt          time.Time `json:"createdAt"`
}

type DetailLog struct {
	Log
	RequestHeaders  map[string][]string `json:"requestHeaders"`
	RequestBody     string              `json:"requestBody"`
	ResponseHeaders map[string][]string `json:"responseHeaders"`
	ResponseBody    string              `json:"responseBody"`
}
