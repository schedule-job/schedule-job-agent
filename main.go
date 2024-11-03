package main

import (
	"fmt"
	"strconv"
	"strings"

	parser "github.com/Sotaneum/go-args-parser"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/schedule-job/schedule-job-agent/internal/job"
	"github.com/schedule-job/schedule-job-database/pg"
	schedule_errors "github.com/schedule-job/schedule-job-errors"
)

type Options struct {
	Port           string
	PostgresSqlDsn string
	TrustedProxies string
}

var DEFAULT_OPTIONS = map[string]string{
	"PORT":             "8080",
	"POSTGRES_SQL_DSN": "",
	"TRUSTED_PROXIES":  "",
}

func getOptions() *Options {
	rawOptions := parser.ArgsJoinEnv(DEFAULT_OPTIONS)

	options := new(Options)
	options.Port = rawOptions["PORT"]
	options.PostgresSqlDsn = rawOptions["POSTGRES_SQL_DSN"]
	options.TrustedProxies = rawOptions["TRUSTED_PROXIES"]

	return options
}

func safeGo(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		f()
	}()
}

func main() {
	options := getOptions()
	if len(options.PostgresSqlDsn) == 0 {
		err := schedule_errors.InvalidArgumentError{Param: "POSTGRES_SQL_DSN"}
		panic(err)
	}
	if len(options.Port) == 0 {
		err := schedule_errors.InvalidArgumentError{Param: "PORT"}
		panic(err)
	}

	database := pg.New(options.PostgresSqlDsn)

	router := gin.Default()
	router.Use(ginsession.New())

	if options.TrustedProxies != "" {
		trustedProxies := strings.Split(options.TrustedProxies, ",")
		router.SetTrustedProxies(trustedProxies)
	}

	router.POST("/api/v1/request", func(ctx *gin.Context) {
		var jobs = []job.Job{}
		bindDataErr := ctx.ShouldBindJSON(&jobs)

		if bindDataErr != nil {
			err := schedule_errors.InvalidArgumentError{Param: "body", Message: bindDataErr.Error()}
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		for _, jobObj := range jobs {
			jobObj.SetDatabase(database)
			safeGo(func() {
				jobObj.Run()
			})
		}

		ctx.JSON(200, gin.H{"code": 200, "data": "ok"})
	})

	router.GET("/api/v1/request/:jobId/logs", func(ctx *gin.Context) {
		jobId := ctx.Param("jobId")
		lastId := ctx.Query("lastId")
		limit, cnvErr := strconv.Atoi(ctx.Query("limit"))
		if cnvErr != nil {
			limit = 20
		}
		logs, dbErr := database.SelectRequestLogs(jobId, lastId, limit)

		if dbErr != nil {
			err := schedule_errors.QueryError{Err: dbErr}
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": logs})
	})

	router.GET("/api/v1/request/:jobId/log/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		jobId := ctx.Param("jobId")

		log, dbErr := database.SelectRequestLog(id, jobId)

		if dbErr != nil {
			err := schedule_errors.QueryError{Err: dbErr}
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"code": 200, "data": &log})
	})

	router.NoRoute(func(ctx *gin.Context) {
		err := schedule_errors.NotFoundError{}
		ctx.JSON(404, gin.H{"code": 404, "message": err.Error()})
	})

	fmt.Println("Started Agent! on " + options.Port)

	router.Run(":" + options.Port)
}
