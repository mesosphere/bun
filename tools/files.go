package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
)

// FindFiles finds file types from the bundle directory.
func FindFiles(p string) ([]bundle.FileType, error) {
	info, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%v is not a directory", p)
	}
	fileTypes, errors := readDir(p)
	if len(errors) != 0 {
		// TODO: concat and return all errors
		return nil, fmt.Errorf("errors found")
	}
	fileTypes = squash(fileTypes)
	for i, t := range fileTypes {
		t.Name = bundle.FileTypeName(chooseName(t.Paths[0]))
		fileTypes[i] = t
	}
	sort.Slice(fileTypes, func(i, j int) bool {
		return fileTypes[i].Name < fileTypes[j].Name
	})
	return fileTypes, nil
}

// choose a name for the file type based on the file name
func chooseName(p string) string {
	p = path.Base(p)
	p = strings.Replace(p, "_", "-", -1)
	p = strings.Replace(p, "8080", "marathon", -1)
	p = strings.Replace(p, "8443", "marathon", -1)
	p = strings.Replace(p, "5050", "mesos-master", -1)
	p = strings.Replace(p, "5051", "mesos-agent", -1)
	p = strings.Replace(p, "slave", "agent", -1)
	p = regexp.MustCompile(`v[0-9]+-`).ReplaceAllLiteralString(p, "")
	p = regexp.MustCompile(`^[0-9]+-`).ReplaceAllLiteralString(p, "")
	p = strings.TrimPrefix(p, "dcos-")
	ext := path.Ext(p)
	if ext != ".socket" && ext != ".timer" {
		p = strings.TrimSuffix(p, ext)
	} else {
		p = strings.TrimSuffix(p, ext) + "-" + strings.TrimPrefix(ext, ".")
	}
	p = regexp.MustCompile(`-[0-9]$`).ReplaceAllLiteralString(p, "")
	p = strings.Replace(p, ".", "-", -1)
	p = toSnakeCase(p)
	p = regexp.MustCompile(` {1,}`).ReplaceAllLiteralString(p, "-")
	p = regexp.MustCompile(`-{2,}`).ReplaceAllLiteralString(p, "-")
	p = strings.TrimSuffix(p, "-json")
	pp := strings.Split(p, "-")
	pp = removeDuplicates(pp)
	p = strings.Join(pp, "-")
	p = strings.Trim(p, "- ")
	return p
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	str = matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	str = matchAllCap.ReplaceAllString(str, "${1}-${2}")
	return strings.ToLower(str)
}
func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// squash file types by filenames. Sort by Path.
func squash(types []bundle.FileType) []bundle.FileType {
	typesByPath := make(map[string]bundle.FileType)
	for _, t := range types {
		p := strings.TrimSuffix(t.Paths[0], ".gz")
		d := t.DirTypes[0]
		existing, ok := typesByPath[p]
		if !ok {
			typesByPath[p] = t
			continue
		}
		existing.DirTypes, ok = appendUnique(existing.DirTypes, d)
		if ok {
			typesByPath[p] = existing
		}
	}
	squashed := make([]bundle.FileType, 0, len(typesByPath))
	for _, t := range typesByPath {
		squashed = append(squashed, t)
	}
	return squashed
}

// returns true if appended
func appendUnique(slice []bundle.DirType, d bundle.DirType) ([]bundle.DirType, bool) {
	for _, dir := range slice {
		if dir == d {
			return slice, false
		}
	}
	return append(slice, d), true
}

func readDir(p string) ([]bundle.FileType, []error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, []error{err}
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, []error{err}
	}
	if info.IsDir() {
		infos, err := ioutil.ReadDir(p)
		if err != nil {
			return nil, []error{err}
		}
		fileTypes := []bundle.FileType{}
		errors := []error{}
		for _, info := range infos {
			f, e := readDir(path.Join(p, info.Name()))
			fileTypes = append(fileTypes, f...)
			errors = append(errors, e...)
		}
		return fileTypes, errors
	}

	name := strings.TrimSuffix(info.Name(), ".gz")
	f := bundle.FileType{}

	// ContentType
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".json":
		f.ContentType = bundle.CTJson
	case ".service":
		f.ContentType = bundle.CTJournal
	case ".socket":
		f.ContentType = bundle.CTJournal
	case ".timer":
		f.ContentType = bundle.CTJournal
	case ".output":
		if strings.HasPrefix(name, "dmesg") {
			f.ContentType = bundle.CTDmesg
		} else {
			f.ContentType = bundle.CTOutput
		}
	default:
		f.ContentType = bundle.CTOther
	}

	// Paths
	f.Paths = append(f.Paths, strings.TrimSuffix(trimBasePath(p), ".gz"))

	//DirTypes
	f.DirTypes = append(f.DirTypes, pathToDirType(p))

	return []bundle.FileType{f}, nil
}

var re = regexp.MustCompile(`((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))_(agent_public|agent|master)(.*)`)

// detect DirType by path
func pathToDirType(p string) bundle.DirType {
	groups := re.FindStringSubmatch(p)
	if groups == nil {
		return bundle.DTRoot
	}
	switch groups[5] {
	case "master":
		return bundle.DTMaster
	case "agent":
		return bundle.DTAgent
	case "agent_public":
		return bundle.DTPublicAgent
	default:
		panic(fmt.Sprintf("unknown directory type: %v", groups[5]))
	}
}

// trim the path to the base directory of the DirType
func trimBasePath(p string) string {
	groups := re.FindStringSubmatch(p)
	if groups == nil {
		// root directory
		return path.Base(p)
	}
	return strings.TrimPrefix(groups[6], "/")
}
