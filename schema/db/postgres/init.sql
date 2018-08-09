CREATE SCHEMA monitor;

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
  days                    TEXT ARRAY,
  status                  TEXT,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  CONSTRAINT monitor_id_process_pkey PRIMARY KEY (id_process),
);

CREATE TRIGGER trigger_process_updated_at BEFORE UPDATE
  ON monitor.process FOR EACH ROW EXECUTE PROCEDURE monitor.function_updated_at();


-- HISTORY
CREATE TABLE monitor.process_history (LIKE monitor.process);
