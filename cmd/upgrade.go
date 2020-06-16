package cmd

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	version "github.com/hashicorp/go-version"
)

var httpClient = http.Client{Timeout: 5 * time.Second}

func CheckNewRelease() {
	thisVersion, err := version.NewSemver("v" + Version)
	if err != nil {
		return
	}
	gh := gitHub{"mesosphere", "bun"}
	latestVersion, err := gh.latestVersion()
	if err != nil {
		return
	}
	if thisVersion.GreaterThanOrEqual(latestVersion) {
		return
	}
	fmt.Fprintf(os.Stderr, "Good news, everyone! New Bun version (v%s) is out, "+
		"you can read the release notes here: https://github.com/mesosphere/bun/releases/latest\n",
		latestVersion.String())
	if outputRedirectedToFile() {
		return
	}
	fmt.Print("Would you like to upgrade to the latest version? [Y/n] ")
	ok, err := askForConfirmation()
	if err != nil {
		fmt.Println("Couldn't read user input:", err.Error())
	}
	if !ok {
		return
	}
	fmt.Println("Upgrading...")
	if err := gh.upgradeExecutable("bun"); err != nil {
		fmt.Println("Couldn't upgrade to the newer version:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully upgraded to the newer version.")
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Println("Couldn't launch to the newer version:", err.Error())
			os.Exit(1)
		}
	}
	os.Exit(0)
}

type gitHub struct {
	Account string
	Repo    string
}

type tag struct {
	Name string
}

func (g gitHub) latestVersion() (*version.Version, error) {
	resp, err := httpClient.Get("https://api.github.com/repos/" + g.Account + "/" + g.Repo + "/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	tags := make([]tag, 0)
	if err := decoder.Decode(&tags); err != nil {
		return nil, err
	}
	latest := version.Must(version.NewSemver("v0.0.0"))
	for _, tag := range tags {
		v, err := version.NewSemver(tag.Name)
		if err != nil {
			continue
		}
		if latest.LessThan(v) {
			latest = v
		}
	}
	return latest, nil
}

func (g gitHub) upgradeExecutable(name string) error {
	newReleaseURL := "https://github.com/" + g.Account + "/" + g.Repo + "/releases/latest/download/" + name + "_" +
		runtime.GOOS + "_" + runtime.GOARCH + ".tar.gz"
	tempFile, err := ioutil.TempFile("", name)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	httpClient.Timeout = 5 * time.Minute
	resp, err := httpClient.Get(newReleaseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzReader)
	_, err = tarReader.Next()
	if err != nil {
		return err
	}
	if _, err = io.Copy(tempFile, tarReader); err != nil {
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return err
	}
	if err := os.Rename(tempFile.Name(), executable); err != nil {
		return err
	}
	if err := os.Chmod(executable, 755); err != nil {
		return err
	}
	return nil
}

func askForConfirmation() (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := scanner.Text()
	yesResponses := []string{"y", "Y", "yes", "Yes", "YES", ""}
	noResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(yesResponses, response) {
		return true, nil
	} else if containsString(noResponses, response) {
		return false, nil
	} else {
		fmt.Print("Please type \"yes\" or \"no\" and then press Enter: ")
		return askForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}
