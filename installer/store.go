package installer

import (
	"errors"
	"fmt"
	"os"
)

// MakeStoreAll  => /data/store
func (inst *App) MakeStoreAll() error {
	if err := inst.MakeDataDir(); err != nil {
		return err
	}
	if err := inst.MakeStoreDir(); err != nil {
		return err
	}
	if err := inst.MakeStoreApps(); err != nil {
		return err
	}
	return nil
}

// MakeStoreDir  => /data/store
func (inst *App) MakeStoreDir() error {
	if err := checkDir(inst.DataDir); err != nil {
		return errors.New(fmt.Sprintf("store dir not exists %s", inst.DataDir))
	}
	return makeDirectoryIfNotExists(inst.StoreDir, os.FileMode(inst.FilePerm))
}

// MakeStoreApps  => /data/store/apps
func (inst *App) MakeStoreApps() error {
	if err := checkDir(inst.GetStoreDir()); err != nil {
		return errors.New(fmt.Sprintf("store/apps not exists %s", inst.GetStoreDir()))
	}
	return makeDirectoryIfNotExists(fmt.Sprintf("%s/apps", inst.StoreDir), os.FileMode(inst.FilePerm))
}
