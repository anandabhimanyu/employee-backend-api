CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    job_title TEXT NOT NULL,
    country TEXT NOT NULL,
    salary NUMERIC NOT NULL CHECK (salary >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
