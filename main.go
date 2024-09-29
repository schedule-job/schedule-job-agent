package main

import (
	"fmt"

	parser "github.com/Sotaneum/go-args-parser"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/schedule-job/schedule-job-agent/internal/job"
	"github.com/schedule-job/schedule-job-agent/internal/pg"
)

type Options struct {
	Port           string
	PostgresSqlDsn string
}

var DEFAULT_OPTIONS = map[string]string{
	"PORT":             "8080",
	"POSTGRES_SQL_DSN": "",
}

func getOptions() *Options {
	rawOptions := parser.ArgsJoinEnv(DEFAULT_OPTIONS)

	options := new(Options)
	options.Port = rawOptions["PORT"]
	options.PostgresSqlDsn = rawOptions["POSTGRES_SQL_DSN"]

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
		panic("not found 'POSTGRES_SQL_DSN' options")
	}
	if len(options.Port) == 0 {
		panic("not found 'PORT' options")
	}

	database := pg.New(options.PostgresSqlDsn)

	router := gin.Default()
	router.Use(ginsession.New())

	router.POST("/api/v1/request", func(ctx *gin.Context) {
		var jobs = []job.Job{}
		bindDataErr := ctx.ShouldBindJSON(&jobs)

		if bindDataErr != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": "잘못된 파라미터 입니다. (" + bindDataErr.Error() + ")"})
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

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"code": 404, "message": "접근 할 수 없는 페이지입니다!"})
	})

	router.Run(":" + options.Port)
}
