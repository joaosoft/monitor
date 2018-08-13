package monitor

import (
	"time"

	"strings"

	errors "github.com/joaosoft/errors"
	types "github.com/joaosoft/types"
)

type IStorageDB interface {
	GetProcess(idProcess string) (*Process, *errors.Err)
	GetProcesses(values map[string][]string) (ListProcess, *errors.Err)
	CreateProcess(newProcess *Process) *errors.Err
	UpdateProcess(updProcess *Process) *errors.Err
	UpdateProcessStatus(idProcess string, status Status) *errors.Err
	DeleteProcess(idProcess string) *errors.Err
	DeleteProcesses() *errors.Err
}

type Interactor struct {
	storageDB IStorageDB
}

func NewInteractor(storageDB IStorageDB) *Interactor {
	return &Interactor{
		storageDB: storageDB,
	}
}

func (interactor *Interactor) GetProcesses(values map[string][]string) (ListProcess, *errors.Err) {
	log.WithFields(map[string]interface{}{"method": "GetProcesses"})
	log.Info("getting processes")
	if categories, err := interactor.storageDB.GetProcesses(values); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting processes on storage database %s", err).ToErr(err)
		return nil, err
	} else {
		return categories, nil
	}
}

func (interactor *Interactor) GetProcess(idProcess string) (*Process, *errors.Err) {
	log.WithFields(map[string]interface{}{"method": "GetProcess"})
	log.Infof("getting process %s", idProcess)
	if category, err := interactor.storageDB.GetProcess(idProcess); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting process %S on storage database %s", idProcess, err).ToErr(err)
		return nil, err
	} else {
		return category, nil
	}
}

func (interactor *Interactor) CreateProcess(newProcess *Process) *errors.Err {
	log.WithFields(map[string]interface{}{"method": "CreateProcess"})

	log.Infof("creating process with id %s", newProcess.IdProcess)
	if err := interactor.storageDB.CreateProcess(newProcess); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating process %s on storage database %s", newProcess.IdProcess, err).ToErr(err)
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) UpdateProcess(updProcess *Process) *errors.Err {
	log.WithFields(map[string]interface{}{"method": "UpdateProcess"})
	log.Infof("updating process %s", updProcess.IdProcess)
	if err := interactor.storageDB.UpdateProcess(updProcess); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error updating process %s on storage database %s", updProcess.IdProcess, err).ToErr(err)
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) UpdateProcessStatus(idProcess string, status Status) errors.ListErr {
	log.WithFields(map[string]interface{}{"method": "UpdateProcessStatus"})
	log.Infof("updating process %s to status %s", idProcess, status)

	if canExecuite, errs := interactor.CanExecute(idProcess); canExecuite {

		if err := interactor.storageDB.UpdateProcessStatus(idProcess, status); err != nil {
			log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
				Errorf("error updating process %s to status %s on storage database %s", idProcess, status, err).ToErr(err)
			return []*errors.Err{err}
		} else {
			return nil
		}
	} else {
		return errs
	}
}

func (interactor *Interactor) UpdateProcessStatusCheck(idProcess string, status Status) errors.ListErr {
	log.WithFields(map[string]interface{}{"method": "UpdateProcessStatusCheck"})
	log.Infof("check updating process %s to status %s", idProcess, status)

	if canExecuite, errs := interactor.CanExecute(idProcess); canExecuite {
		return nil
	} else {
		return errs
	}
}

func (interactor *Interactor) DeleteProcess(idProcess string) *errors.Err {
	log.WithFields(map[string]interface{}{"method": "DeleteProcess"})
	log.Infof("deleting process %s", idProcess)
	if err := interactor.storageDB.DeleteProcess(idProcess); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting process %s on storage database %s", idProcess, err).ToErr(err)
		return err
	}
	return nil
}

func (interactor *Interactor) DeleteProcesses() *errors.Err {
	log.WithFields(map[string]interface{}{"method": "DeleteProcesses"})
	log.Info("deleting processes")
	if err := interactor.storageDB.DeleteProcesses(); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting processes on storage database %s", err).ToErr(err)
		return err
	}
	return nil
}

func (interactor *Interactor) CanExecute(idProcess string) (bool, errors.ListErr) {
	var errs errors.ListErr
	process, err := interactor.GetProcess(idProcess)
	if err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting process %s on storage database %s", idProcess, err).ToErr(err)
		return false, *errs.Add(err)
	}

	now := time.Now()
	if process.Status != nil && *process.Status == StatusRunning {
		errors.New("0", "the process is already running!")
	}
	if process.DaysOff != nil && process.DaysOff.Contains(types.Day(strings.ToLower(now.Weekday().String()))) {
		errs.Add(errors.New("1", "the process cannot the executed on %+v!", process.DaysOff))
	}
	if process.DateFrom != nil && now.Format("02-01-2006") < string(*process.DateFrom) {
		errs.Add(errors.New("2", "the process can just be started after %s", *process.DateFrom))
	}
	if process.DateTo != nil && now.Format("02-01-2006") > string(*process.DateTo) {
		errs.Add(errors.New("3", "the process could just be started before %s", *process.DateTo))
	}
	if process.TimeFrom != nil && now.Format("15:04:05") < string(*process.TimeFrom) {
		errs.Add(errors.New("4", "the process can just be started after %s", *process.TimeFrom))
	}
	if process.TimeTo != nil && now.Format("15:04:05") > string(*process.TimeTo) {
		errs.Add(errors.New("5", "the process could just be started before %s", *process.TimeTo))
	}

	return errs.IsEmpty(), errs
}
