package installer

import (
	"github.com/NubeIO/lib-command/command"
	"github.com/NubeIO/lib-command/unixcmd"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

type BuildDetails struct {
	MatchedName    string `json:"matched_name"`
	MatchedVersion string `json:"matched_version"`
	MatchedArch    string `json:"matched_arch"`
	ZipName        string `json:"zip_name"`
}

func (inst *App) GetZipBuildDetails(zipName string) *BuildDetails {
	parts := strings.Split(zipName, "-")
	count := 0
	name := ""
	version := ""
	arch := ""
	for i, part := range parts { // match version
		p := strings.Split(part, ".")
		// if len is 3 eg, 0.0.1
		isNum := 0
		if len(p) == 3 || len(p) == 4 {
			// check if they are numbers
			for _, s := range p {
				if _, err := strconv.Atoi(s); err == nil {
					isNum++
				}
			}
			if isNum == 3 {
				count = i
				version = part
				version = strings.Trim(version, ".zip")
			}
			for _, s := range p { // match arch
				if s == "amd64" {
					arch = "amd64"
				}
				if s == "armv7" {
					arch = "armv7"
				}
			}
		}
	}
	name = strings.Join(parts[0:count], "-")
	return &BuildDetails{
		MatchedName:    name,
		MatchedVersion: version,
		MatchedArch:    arch,
		ZipName:        zipName,
	}
}

// GetBuildZipNameByArch // get a build by its arch
func (inst *App) GetBuildZipNameByArch(path, arch string) (*BuildDetails, error) {
	var out *BuildDetails
	details, err := getFileDetails(path)
	if err != nil {
		return out, err
	}
	for _, name := range details {
		app := inst.GetZipBuildDetails(name.Name)
		if app.MatchedArch == arch {
			out = app
		}
	}
	return out, nil
}

// GetBuildZipNames // get all the builds zips from a path
func (inst *App) GetBuildZipNames(path string) ([]BuildDetails, error) {
	var out []BuildDetails
	details, err := getFileDetails(path)
	if err != nil {
		return out, err
	}
	for _, name := range details {
		app := inst.GetZipBuildDetails(name.Name)
		out = append(out, *app)
	}
	return out, nil
}

type fileDetails struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	IsDir     bool   `json:"is_dir"`
}

func getFileDetails(dir string) ([]fileDetails, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []fileDetails
	var f fileDetails
	for _, file := range files {
		var extension = filepath.Ext(file.Name())
		f.Extension = extension
		f.Name = file.Name()
		f.IsDir = file.IsDir()
		out = append(out, f)
	}
	return out, nil
}

// matchRepoName get the tag name from the zip eg, wires-builds-0.5.5-1575cf89.amd64.zip => wires-builds
// 	returns
// 	- true if is a match if it is a match
// 	- match count
// 	- string version name
// 	- arch match`
// 	- arch type
func matchRepoName(zipName, repoName string) (bool, int, string, bool, string) {
	parts := strings.Split(zipName, "-")
	repoNameParts := strings.Split(repoName, "-")
	count := 0
	version := ""
	arch := ""
	archMatch := false
	repoMatch := false
	for i, part := range parts {
		p := strings.Split(part, ".")
		// if len is 3 eg, 0.0.1
		isNum := 0
		if len(p) == 3 || len(p) == 4 {
			// check if they are numbers
			for _, s := range p {
				if _, err := strconv.Atoi(s); err == nil {
					isNum++
				}
			}
			if isNum == 3 {
				count = i
				version = part
				version = strings.Trim(version, ".zip")
			}
		}
	}
	match := 0
	for i := 0; i < count; i++ {
		if isMatch(parts, repoNameParts[i]) {
			match++
		}
	}
	if match == count {
		repoMatch = true
	}
	if repoName != "wires-builds" { // wires can run on any os
		arch, _ = getArch()
		if contains(parts, arch) {
			if repoName == "wires-builds" {

			}
			archMatch = true
		}
	} else {
		archMatch = true
	}
	return repoMatch, count, version, archMatch, arch
}

var cmd = unixcmd.New(&command.Command{})

func getArch() (string, error) {
	arch, err := cmd.DetectArch()
	if err != nil {
		return "", err
	}
	return arch.ArchModel, err
}

func isMatch(s []string, term string) bool {
	count := 0
	for _, item := range s {
		if item == term {
			count++
			return true
		}
	}
	return false
}

func contains(s []string, term string) bool {
	count := 0
	for _, item := range s {
		if strings.Contains(item, term) {
			count++
			return true
		}
	}
	return false
}
