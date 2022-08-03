package installer

import (
	"fmt"
	"testing"
)

/*
This is for taking backups on an edge device
*/

func Test_BackupApp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
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
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
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
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	backs, err := app.ListFullBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_ListAppBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	backs, err := app.ListAppBackupsDirs()
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_ListBackupsByApp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	backs, err := app.ListBackupsByApp(appName)
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_DeleteAllBackups(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	backs, err := app.WipeBackups()
	fmt.Println(err)
	fmt.Println(backs)
}

func TestApp_DeleteAppAllBackUpByName(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
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
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	byApp, err := app.ListBackupsByApp(appName)
	fmt.Println(err)
	if err != nil {
		return
	}
	for i, folder := range byApp {
		fmt.Println(i, folder)
		if i == 0 {
			backs, err := app.DeleteAppOneBackUpByName(appName, folder)
			fmt.Println(err)
			fmt.Println(backs)
		}
	}
}
