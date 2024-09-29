package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/schedule-job/schedule-job-agent/internal/job"
)

func (p *PostgresSQL) InsertRequestLog(jobID string, data job.Request) error {
	_, err := p.usePostgresSQL(func(client *pgx.Conn, ctx context.Context) (result interface{}, err error) {
		return client.Exec(ctx, "INSERT INTO request_logs (job_id, status, request_url, request_method, request_headers, request_body, response_headers, response_body, response_status_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", jobID, data.Status, data.RequestUrl, data.RequestMethod, data.RequestHeaders, data.RequestBody, data.ResponseHeaders, data.ResponseBody, data.ResponseStatusCode)
	})
	return err
}
