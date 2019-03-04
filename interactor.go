package monitor

import (
	"github.com/joaosoft/logger"
	"time"

	"strings"

	errors "github.com/joaosoft/errors"
	types "github.com/joaosoft/types"
)

type IStorageDB interface {
	GetProcess(idProcess string) (*Process, error)
	GetProcesses(values map[string][]string) (ListProcess, error)
	CreateProcess(newProcess *Process) error
	UpdateProcess(updProcess *Process) error
	UpdateProcessStatus(idProcess string, status Status) error
	DeleteProcess(idProcess string) error
	DeleteProcesses() error
}

type Interactor struct {
	storageDB IStorageDB
	logger logger.ILogger
}

func (monitor *Monitor) NewInteractor(storageDB IStorageDB) *Interactor {
	return &Interactor{
		storageDB: storageDB,
		logger: monitor.logger,
	}
}

func (interactor *Interactor) GetProcesses(values map[string][]string) (ListProcess, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "GetProcesses"})
	interactor.logger.Info("getting processes")
	if categories, err := interactor.storageDB.GetProcesses(values); err != nil {
		err = interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting processes on storage database %s", err).ToError()
		return nil, err
	} else {
		return categories, nil
	}
}

func (interactor *Interactor) GetProcess(idProcess string) (*Process, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "CheckUser"})
	interactor.logger.Infof("getting process %s", idProcess)
	if category, err := interactor.storageDB.GetProcess(idProcess); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting process %S on storage database %s", idProcess, err).ToError()
		return nil, err
	} else {
		return category, nil
	}
}

func (interactor *Interactor) CreateProcess(newProcess *Process) error {
	interactor.logger.WithFields(map[string]interface{}{"method": "CreateProcess"})

	interactor.logger.Infof("creating process with id %s", newProcess.IdProcess)
	if err := interactor.storageDB.CreateProcess(newProcess); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error creating process %s on storage database %s", newProcess.IdProcess, err).ToError()
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) UpdateProcess(updProcess *Process) error {
	interactor.logger.WithFields(map[string]interface{}{"method": "UpdateProcess"})
	interactor.logger.Infof("updating process %s", updProcess.IdProcess)
	if err := interactor.storageDB.UpdateProcess(updProcess); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error updating process %s on storage database %s", updProcess.IdProcess, err).ToError()
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) UpdateProcessStatus(idProcess string, status Status) errors.ListErr {
	interactor.logger.WithFields(map[string]interface{}{"method": "UpdateProcessStatus"})
	interactor.logger.Infof("updating process %s to status %s", idProcess, status)

	if canExecuite, errs := interactor.CanExecute(idProcess); canExecuite {

		if err := interactor.storageDB.UpdateProcessStatus(idProcess, status); err != nil {
			err = interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
				Errorf("error updating process %s to status %s on storage database %s", idProcess, status, err).ToError()
			return []*errors.Err{errors.New(errors.ErrorLevel, 0, err)}
		} else {
			return nil
		}
	} else {
		return errs
	}
}

func (interactor *Interactor) UpdateProcessStatusCheck(idProcess string, status Status) errors.ListErr {
	interactor.logger.WithFields(map[string]interface{}{"method": "UpdateProcessStatusCheck"})
	interactor.logger.Infof("check updating process %s to status %s", idProcess, status)

	if canExecuite, errs := interactor.CanExecute(idProcess); canExecuite {
		return nil
	} else {
		return errs
	}
}

func (interactor *Interactor) DeleteProcess(idProcess string) error {
	interactor.logger.WithFields(map[string]interface{}{"method": "DeleteProcess"})
	interactor.logger.Infof("deleting process %s", idProcess)
	if err := interactor.storageDB.DeleteProcess(idProcess); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting process %s on storage database %s", idProcess, err).ToError()
		return err
	}
	return nil
}

func (interactor *Interactor) DeleteProcesses() error {
	interactor.logger.WithFields(map[string]interface{}{"method": "DeleteProcesses"})
	interactor.logger.Info("deleting processes")
	if err := interactor.storageDB.DeleteProcesses(); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting processes on storage database %s", err).ToError()
		return err
	}
	return nil
}

func (interactor *Interactor) CanExecute(idProcess string) (bool, errors.ListErr) {
	var errs errors.ListErr
	process, err := interactor.GetProcess(idProcess)
	if err != nil {
		err = interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting process %s on storage database %s", idProcess, err).ToError()
		return false, []*errors.Err{errors.New(errors.ErrorLevel, 0, err)}
	}

	now := time.Now()
	if process.Status != nil && *process.Status == StatusRunning {
		errors.New(errors.ErrorLevel, 0, "the process is already running!")
	}
	if process.DaysOff != nil && process.DaysOff.Contains(types.Day(strings.ToLower(now.Weekday().String()))) {
		errors.New(errors.ErrorLevel, 0, "the process cannot the executed on %+v!", process.DaysOff)
	}
	if process.DateFrom != nil && now.Format("02-01-2006") < string(*process.DateFrom) {
		errors.New(errors.ErrorLevel, 0, "the process can just be started after %s", *process.DateFrom)
	}
	if process.DateTo != nil && now.Format("02-01-2006") > string(*process.DateTo) {
		errors.New(errors.ErrorLevel, 0, "the process could just be started before %s", *process.DateTo)
	}
	if process.TimeFrom != nil && now.Format("15:04:05") < string(*process.TimeFrom) {
		errors.New(errors.ErrorLevel, 0, "the process can just be started after %s", *process.TimeFrom)
	}
	if process.TimeTo != nil && now.Format("15:04:05") > string(*process.TimeTo) {
		errors.New(errors.ErrorLevel, 0, "the process could just be started before %s", *process.TimeTo)
	}

	return errs.IsEmpty(), errs
}
