package monitor

import (
	"github.com/joaosoft/errors"
)

type IStorageDB interface {
	GetProcess(idProcess string) (*Process, *errors.Err)
	GetProcesses(values map[string][]string) (ListProcess, errors.IErr)
	CreateProcess(newProcess *Process) errors.IErr
	UpdateProcess(updProcess *Process) errors.IErr
	DeleteProcess(idProcess string) errors.IErr
	DeleteProcesses() errors.IErr
}

type Interactor struct {
	storageDB IStorageDB
}

func NewInteractor(storageDB IStorageDB) *Interactor {
	return &Interactor{
		storageDB: storageDB,
	}
}

func (interactor *Interactor) GetProcesses(values map[string][]string) (ListProcess, errors.IErr) {
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

func (interactor *Interactor) GetProcess(idProcess string) (*Process, errors.IErr) {
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

func (interactor *Interactor) CreateProcess(newProcess *Process) errors.IErr {
	log.WithFields(map[string]interface{}{"method": "CreateProcess"})

	log.Infof("creating process with id %s", newProcess.IdProcess)
	if err := interactor.storageDB.CreateProcess(newProcess); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating process %s on storage database", newProcess.IdProcess, err).ToErr(err)
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) UpdateProcess(updProcess *Process) errors.IErr {
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

func (interactor *Interactor) DeleteProcess(idProcess string) errors.IErr {
	log.WithFields(map[string]interface{}{"method": "DeleteProcess"})
	log.Infof("deleting process %s", idProcess)
	if err := interactor.storageDB.DeleteProcess(idProcess); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting process %s on storage database %s", idProcess, err).ToErr(err)
		return err
	}
	return nil
}

func (interactor *Interactor) DeleteProcesses() errors.IErr {
	log.WithFields(map[string]interface{}{"method": "DeleteProcesses"})
	log.Info("deleting processes")
	if err := interactor.storageDB.DeleteProcesses(); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting processes on storage database %s", err).ToErr(err)
		return err
	}
	return nil
}
