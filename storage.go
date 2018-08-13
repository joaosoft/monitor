package monitor

import (
	"database/sql"

	"fmt"

	errors "github.com/joaosoft/errors"
	manager "github.com/joaosoft/manager"
	"github.com/lib/pq"
)

type StoragePostgres struct {
	conn manager.IDB
}

func NewStoragePostgres(connection manager.IDB) *StoragePostgres {
	return &StoragePostgres{
		conn: connection,
	}
}

func (storage *StoragePostgres) GetProcess(idProcess string) (*Process, *errors.Err) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
		    "type",
			"name",
			description,
			date_from,
			date_to,
			time_from,
			time_to,
			days_off,
			monitor,
			status,
			updated_at,
			created_at
		FROM monitor.process
		WHERE id_process = $1
	`, idProcess)

	process := &Process{IdProcess: idProcess}
	if err := row.Scan(
		&process.Type,
		&process.Name,
		&process.Description,
		&process.DateFrom,
		&process.DateTo,
		&process.TimeFrom,
		&process.TimeTo,
		&process.DaysOff,
		&process.Monitor,
		&process.Status,
		&process.UpdatedAt,
		&process.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, errors.New("0", err)
		}

		return nil, nil
	}

	return process, nil
}

func (storage *StoragePostgres) GetProcesses(values map[string][]string) (ListProcess, *errors.Err) {
	query := `
	    SELECT
			id_process,
		    "type",
			"name",
			description,
			date_from,
			date_to,
			time_from,
			time_to,
			days_off,
			monitor,
			status,
			updated_at,
			created_at
		FROM monitor.process
	`

	index := 1
	params := make([]interface{}, 0)

	if values != nil {
		for key, value := range values {
			if index == 1 {
				query += ` WHERE `
			} else {
				query += ` AND `
			}
			query += fmt.Sprintf(`%s = $%d`, key, index)
			index = index + 1

			params = append(params, value[0])
		}
	}

	rows, err := storage.conn.Get().Query(query, params...)
	defer rows.Close()
	if err != nil {
		return nil, errors.New("0", err)
	}

	processes := make(ListProcess, 0)
	for rows.Next() {
		process := &Process{}
		if err := rows.Scan(
			&process.IdProcess,
			&process.Type,
			&process.Name,
			&process.Description,
			&process.DateFrom,
			&process.DateTo,
			&process.TimeFrom,
			&process.TimeTo,
			&process.DaysOff,
			&process.Monitor,
			&process.Status,
			&process.UpdatedAt,
			&process.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, errors.New("0", err)
			}
			return nil, nil
		}
		processes = append(processes, process)
	}

	return processes, nil
}

func (storage *StoragePostgres) CreateProcess(newProcess *Process) *errors.Err {
	if _, err := storage.conn.Get().Exec(`
		INSERT INTO monitor.process(
			id_process, 
			"type",
			"name", 
			description,
			date_from,
			date_to,
			time_from,
			time_to,
			days_off,
			monitor,
			status)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		newProcess.IdProcess,
		newProcess.Type,
		newProcess.Name,
		newProcess.Description,
		newProcess.DateFrom,
		newProcess.DateTo,
		newProcess.TimeFrom,
		newProcess.TimeTo,
		newProcess.DaysOff,
		newProcess.Monitor,
		newProcess.Status); err != nil {
		return errors.New("0", err)
	}

	return nil
}

func (storage *StoragePostgres) UpdateProcess(updProcess *Process) *errors.Err {
	if _, err := storage.conn.Get().Exec(`
		UPDATE monitor.process SET 
			"type" = $1, 
			"name" = $2, 
			description = $3,
			date_from = $4,
			date_to = $5,
			time_from = $6,
			time_to = $7,
			days_off = $8,
			monitor = $9,
			status = $10,
			updated_at = $11
		WHERE id_process = $12
	`, updProcess.Type,
		updProcess.Name,
		updProcess.Description,
		updProcess.DateFrom,
		updProcess.DateTo,
		updProcess.TimeFrom,
		updProcess.TimeTo,
		pq.Array(updProcess.DaysOff),
		updProcess.Monitor,
		updProcess.Status,
		updProcess.UpdatedAt,
		updProcess.IdProcess); err != nil {
		return errors.New("0", err)
	}

	return nil
}

func (storage *StoragePostgres) UpdateProcessStatus(idProcess string, status Status) *errors.Err {
	if _, err := storage.conn.Get().Exec(`
		UPDATE monitor.process SET 
			status = $1
		WHERE id_process = $2
	`, status, idProcess); err != nil {
		return errors.New("0", err)
	}

	return nil
}

func (storage *StoragePostgres) DeleteProcess(idProcess string) *errors.Err {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM monitor.process
		WHERE id_process = $1
	`, idProcess); err != nil {
		return errors.New("0", err)
	}

	return nil
}

func (storage *StoragePostgres) DeleteProcesses() *errors.Err {
	if _, err := storage.conn.Get().Exec(`
	    DELETE FROM monitor.process`); err != nil {
		return errors.New("0", err)
	}

	return nil
}