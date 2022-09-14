package installer

import (
	"fmt"
	"testing"
)

// This is for taking backups on an edge device

var testAppName = "flow-framework"

func Test_readZip(t *testing.T) {
	readZip("/home/aidan/backup/full/testDevice1234-full-backup-2022-08-05 06:30:05.zip")
}

func Test_BackupApp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	deviceName := "testDevice1234"
	back, err := app.BackupApp(testAppName, &deviceName)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(back)
}

func Test_FullBackUp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	deviceName := "testDevice1234"
	back, err := app.FullBackUp(&deviceName)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(back)
}

func Test_ListFullBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	backs, err := app.ListFullBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func Test_ListAppsBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	backs, err := app.ListAppsBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func Test_ListAppBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	backs, err := app.ListAppBackups(testAppName)
	fmt.Println(err)
	fmt.Println(backs)
}

func Test_DeleteAllBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	backs, err := app.WipeBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func Test_DeleteAllAppBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode, BackupsDir: "/home/aidan/backup"})
	backs, err := app.DeleteAllAppBackups(testAppName)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(backs)
}
