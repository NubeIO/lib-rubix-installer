package installer

import "os"

// MakeStoreAll  => /data, /data/store, /data/store/apps
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
	return os.MkdirAll(inst.StoreDir, os.FileMode(inst.FileMode))
}

// MakeStoreApps  => /data/store/apps
func (inst *App) MakeStoreApps() error {
	return os.MkdirAll(inst.GetStoreAppsDir(), os.FileMode(inst.FileMode))
}
