ALTER TABLE widgets ADD is_recurring BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE widgets ADD plan_id VARCHAR(255) NOT NULL DEFAULT '';