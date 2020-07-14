package bundle

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

const hostRegexp = `^((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))_(agent_public|agent|master)$`

// Host represents a host in a DC/OS cluster.
type Host struct {
	IP IP
	Directory
}

type IP string

// Bundle describes DC/OS diagnostics bundle.
type Bundle struct {
	Hosts []Host
	Directory
}

// New creates new Bundle
func New(path string) (Bundle, error) {
	b := Bundle{}
	var err error
	b.Type = DTRoot
	b.Path, err = filepath.Abs(path)
	if err != nil {
		log.Printf("bun.New: cannot determine absolute path: %v", err)
		return b, err
	}
	infos, err := ioutil.ReadDir(b.Path)
	if err != nil {
		return b, err
	}
	re := regexp.MustCompile(hostRegexp)
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		groups := re.FindStringSubmatch(info.Name())
		if groups == nil {
			continue
		}
		var host Host
		host.IP = IP(groups[1])
		host.Path = filepath.Join(b.Path, info.Name())
		switch groups[5] {
		case "master":
			host.Type = DTMaster
		case "agent":
			host.Type = DTAgent
		case "agent_public":
			host.Type = DTPublicAgent
		default:
			panic(fmt.Sprintf("Unknown directory type: %v", groups[5]))
		}
		b.Hosts = append(b.Hosts, host)
	}

	if len(b.Hosts) == 0 {
		_, err = b.OpenFile("summary-report")
		if err != nil {
			_, err = b.OpenFile("summary-errors-report")
			if err != nil {
				return b, errors.New("bundle not found in the given directory")
			}
		}
	}
	return b, nil
}

func (b Bundle) Masters() []Host {
	return b.filter(DTMaster)
}

func (b Bundle) Agents() []Host {
	return b.filter(DTAgent)
}

func (b Bundle) PublicAgents() []Host {
	return b.filter(DTPublicAgent)
}

func (b Bundle) filter(t DirType) (hosts []Host) {
	hosts = make([]Host, 0, len(b.Hosts))
	for _, host := range b.Hosts {
		if host.Type == t {
			hosts = append(hosts, host)
		}
	}
	return
}

// ForEachFile finds all the files of a given type and pass them one by one to the do function.
// It stops if the do function returns true.
func (b Bundle) ForEachFile(fileTypeName FileTypeName, do func(f File) (stop bool)) {
	f := func(d Directory) bool {
		if f, err := d.OpenFile(fileTypeName); err == nil {
			return do(f)
		}
		return false
	}
	b.ForEachDirectory(fileTypeName, f)
}

// ForEachDirectory finds all the bundle directories of a given type and pass them one by one to the do function.
// It stops if the do function returns true.
func (b Bundle) ForEachDirectory(fileTypeName FileTypeName, do func(d Directory) (stop bool)) {
	if t := GetFileType(fileTypeName); t.ExistsOn(b.Type) {
		if do(b.Directory) {
			return
		}
	}
	for _, host := range b.Hosts {
		if t := GetFileType(fileTypeName); t.ExistsOn(host.Type) {
			if do(host.Directory) {
				return
			}
		}
	}
}
