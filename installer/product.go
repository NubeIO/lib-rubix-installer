package installer

import (
	"errors"
	"github.com/NubeIO/lib-files/fileutils/jparse"
	log "github.com/sirupsen/logrus"
	"runtime"
)

type Product struct {
	EdgeVersion  string `json:"edge_version"`
	FlowVersion  string `json:"flow_version"`
	ImageVersion string `json:"image_version"`
	Product      string `json:"product"` // RubixCompute, RubixComputeIO, RubixCompute5, Server, Edge28, Nuc  see https://github.com/NubeIO/lib-command/blob/master/product/product.go#L7
	Arch         string `json:"arch"`    // armv7 amd64
	OS           OS     `json:"os"`      // Linux, Windows, Darwin
}

func (inst *App) GetProduct(fileAndPath ...string) (*Product, error) {
	edgeVersion := inst.GetAppVersion("rubix-edge")
	ffVersion := inst.GetAppVersion("flow-framework")
	product, err := read(fileAndPath...)
	if product != nil {
		product.EdgeVersion = edgeVersion
		product.FlowVersion = ffVersion
	}
	return product, err
}

const (
	FilePath = "/data/product.json"
)

func read(fileAndPath ...string) (*Product, error) {
	path := FilePath
	if len(fileAndPath) > 0 {
		path = fileAndPath[0]
		if path == "" {
			return nil, errors.New("path can not be nil")
		}
	}
	p := &Product{}
	j := jparse.New()
	var err error
	if readErr := j.ParseToData(path, p); readErr != nil {
		log.Errorln("read-product: read from json err", readErr.Error())
		err = readErr
		return nil, readErr
	}
	p.OS = GetOS()
	return p, err
}

type OS struct {
	Type    string `json:"type,omitempty"`
	Windows bool   `json:"windows"`
	Linux   bool   `json:"linux"`
	Darwin  bool   `json:"darwin"`
}

func GetOS() (arch OS) {
	s := runtime.GOOS
	switch s {
	case "linux":
		arch.Linux = true
	case "windows":
		arch.Windows = true
	case "darwin":
		arch.Darwin = true
	}
	arch.Type = s
	return arch
}

func (inst *App) GetOS() (arch OS) {
	s := runtime.GOOS
	switch s {
	case "linux":
		arch.Linux = true
	case "windows":
		arch.Windows = true
	case "darwin":
		arch.Darwin = true
	}
	arch.Type = s
	return arch
}
