package bundle

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type yamlFile struct {
	Name        string   `yaml:"name"`
	ContentType string   `yaml:"contentType"`
	Paths       []string `yaml:"paths"`
	Description string   `yaml:"description"`
	DirTypes    []string `yaml:"dirTypes"`
}

func init() {
	var files []yamlFile
	err := yaml.Unmarshal([]byte(filesYAML), &files)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fileType, err := convert(file)
		if err != nil {
			panic(err)
		}
		RegisterFileType(fileType)
	}
}

func convert(y yamlFile) (fileType FileType, err error) {
	fileType.Name = FileTypeName(y.Name)
	fileType.Description = y.Description
	switch y.ContentType {
	case string(CTJson):
		fileType.ContentType = CTJson
	case "journal":
		fileType.ContentType = CTJournal
	case "dmesg":
		fileType.ContentType = CTDmesg
	case "output":
		fileType.ContentType = CTOutput
	case "other":
		fileType.ContentType = CTOther
	default:
		err = fmt.Errorf("FileType '%v' has unknown ContentType '%v'", fileType.Name, y.ContentType)
	}
	for _, s := range y.DirTypes {
		var dirType DirType
		dirType, err = convertDirType(s)
		if err != nil {
			return
		}
		fileType.DirTypes = append(fileType.DirTypes, dirType)
	}
	fileType.Paths = y.Paths
	return
}

func convertDirType(s string) (d DirType, err error) {
	switch s {
	case "root":
		d = DTRoot
	case "master":
		d = DTMaster
	case "agent":
		d = DTAgent
	case "public agent":
		d = DTPublicAgent
	default:
		err = fmt.Errorf("unknown DirType: %v", s)
	}
	return
}
