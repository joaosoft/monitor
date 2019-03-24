
-- migrate up
CREATE SCHEMA IF NOT EXISTS "monitor";

-- GLOBAL
CREATE OR REPLACE FUNCTION monitor.function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';


-- PROCESS
CREATE TABLE monitor.process (
  id_process              TEXT NOT NULL,
  "type"                  TEXT NOT NULL,
  name                    TEXT NOT NULL,
  description             TEXT,
  time_from               TIME,
  time_to                 TIME,
  date_from               DATE,
  date_to                 DATE,
  days_off                TEXT ARRAY,
  monitor                 TEXT,
  status                  TEXT,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  CONSTRAINT monitor_id_process_pkey PRIMARY KEY (id_process)
);

CREATE TRIGGER trigger_process_updated_at BEFORE UPDATE
  ON monitor.process FOR EACH ROW EXECUTE PROCEDURE monitor.function_updated_at();


-- HISTORY
CREATE TABLE monitor.process_history (LIKE monitor.process);
ALTER TABLE monitor.process_history ADD COLUMN _operation TEXT NOT NULL;
ALTER TABLE monitor.process_history ADD COLUMN "_user" TEXT NOT NULL;
ALTER TABLE monitor.process_history ADD COLUMN _operation_at TIMESTAMP DEFAULT NOW();

CREATE OR REPLACE FUNCTION function_process_history() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO monitor.process_history VALUES(OLD.*, 'D', user, now());
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO monitor.process_history VALUES(NEW.*, 'U', user, now());
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO monitor.process_history VALUES(NEW.*, 'I', user, now());
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_process_history
AFTER INSERT OR UPDATE OR DELETE ON monitor.process
    FOR EACH ROW EXECUTE PROCEDURE function_process_history();


-- migrate down
DROP TRIGGER trigger_process_history ON monitor.process;
DROP TRIGGER trigger_process_updated_at on monitor.process;
DROP FUNCTION monitor.function_updated_at();

DROP TABLE monitor.process_HISTORY;
DROP TABLE monitor.process;