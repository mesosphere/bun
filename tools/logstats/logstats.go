package logstats

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
)

func LogStats(b *bundle.Bundle, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	err := csvWriter.Write([]string{"IP", "Host.Type", "Service", "Lines"})
	if err != nil {
		return fmt.Errorf("cannot write CSV. Cause: %s", err.Error())
	}
	var serviceRegexp = regexp.MustCompile(`.*(\.service$|\.service.gz)`)
	for _, host := range b.Hosts {
		files, err := ioutil.ReadDir(host.Path)
		if err != nil {
			return fmt.Errorf("cannot read dir. Cause: %s", err.Error())
		}
		for _, f := range files {
			err := func() error {
				if f.IsDir() || !serviceRegexp.MatchString(f.Name()) {
					return nil
				}
				path := filepath.Join(host.Path, f.Name())
				fileName := f.Name()
				compressed := false
				if filepath.Ext(fileName) == ".gz" {
					compressed = true
					fileName = strings.TrimSuffix(fileName, ".gz")
				}
				fileName = strings.TrimSuffix(fileName, ".service")
				var fileReader io.ReadCloser
				var gzReader io.ReadCloser
				var reader io.ReadCloser
				if fileReader, err = os.Open(path); err != nil {
					return fmt.Errorf("cannot open file %s. Cause: %s", path, err.Error())
				}
				defer func() { _ = fileReader.Close() }()
				reader = fileReader
				if compressed {
					if gzReader, err = gzip.NewReader(fileReader); err != nil {
						return fmt.Errorf("cannot open file %s. Cause: %s", path, err.Error())
					}
					defer func() { _ = gzReader.Close() }()
					reader = gzReader
				}
				scanner := bufio.NewScanner(reader)
				var lineCount int
				for scanner.Scan() {
					lineCount++
				}
				err := csvWriter.Write([]string{string(host.IP), string(host.Type), fileName, strconv.Itoa(lineCount)})
				if err != nil {
					return fmt.Errorf("cannot write CSV. Cause: %s", err.Error())
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
