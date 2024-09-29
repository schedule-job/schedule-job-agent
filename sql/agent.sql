/*
  sql for postgresql
*/

/* CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; */

CREATE TABLE request_logs
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    job_id uuid NOT NULL,
    status text NOT NULL DEFAULT 'progress',
    request_url text NOT NULL,
    request_method text NOT NULL,
    request_headers json,
    request_body text DEFAULT '',
    response_headers json,
    response_body text DEFAULT '',
    response_status_code integer,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    UNIQUE (id)
);

