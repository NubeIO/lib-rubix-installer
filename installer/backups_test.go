package installer

import (
	"fmt"
	"testing"
)

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
