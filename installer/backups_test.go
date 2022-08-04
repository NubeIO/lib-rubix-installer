package installer

import (
	"fmt"
	pprint "github.com/NubeIO/lib-rubix-installer/print"
	"testing"
)

/*
This is for taking backups on an edge device
*/

func Test_readZip(t *testing.T) {
	readZip("/home/aidan/backup/full/testDevice1234-full-backup-2022-08-05 06:30:05.zip")

}

func Test_BackupApp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	back, err := app.BackupApp(appName, "testDevice1234")
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(back)

}

func Test_FullBackUp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	back, err := app.FullBackUp("testDevice1234")
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(back)
}

func TestApp_ListFullBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	backs, err := app.ListFullBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_ListAppBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	backs, err := app.ListAppBackupsDirs()
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_ListBackupsByApp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	backs, err := app.ListBackupsByApp(appName)
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_DeleteAllBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	backs, err := app.WipeBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_DeleteAppAllBackUpByName(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	backs, err := app.DeleteAppAllBackUpByName(appName)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(backs)

}

func TestApp_DeleteAppBackUp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, BackupsDir: "/home/aidan/backup"})
	byApp, err := app.ListBackupsByApp(appName)
	fmt.Println(err)
	if err != nil {
		return
	}
	if byApp == nil {

	}
	pprint.PrintJOSN(byApp)

}
