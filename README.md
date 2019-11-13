# Bun

Command-line program which detects the most common problems in a DC/OS cluster
by analyzing its [diagnostics bundle](https://docs.mesosphere.com/1.11/cli/command-reference/dcos-node/dcos-node-diagnostics-create/).

## Installation

### macOS

1. Download and unpack the binary:

```
$ curl -O -L https://github.com/mesosphere/bun/releases/latest/download/bun_darwin_amd64.tar.gz && tar -zxvf bun_darwin_amd64.tar.gz
```

2. Move the `bun` binary to one of the directories in the `PATH`.

### Linux

1. Download and unpack the binary:

```
$ curl -O -L https://github.com/mesosphere/bun/releases/latest/download/bun_linux_amd64.tar.gz && tar -zxvf bun_linux_amd64.tar.gz
```

2. Move the `bun` binary to one of the directories in the `PATH`.

### Windows

1. Download [the command](https://github.com/mesosphere/bun/releases/latest/download/bun_windows_amd64.tar.gz)
2. Unzip it and move the `bun` binary to one of the folders in the `PATH`.

### From sources

1. Install [Go compiler](https://golang.org/dl/).
2. Run the following command in your terminal:

```bash
$ go get github.com/mesosphere/bun
```

## Usage

```bash
$ bun -p <path to bundle directory>
```

or if you working directory is the bundle directory simply:

```bash
$ bun
```


The output should look like that:

```
[PROBLEM] "dcos-version" - Versions are different.
---------------
Problem details
---------------
master 172.20.0.23 has DC/OS version 1.11.0
master 172.20.0.24 has DC/OS version 1.11.0
agent 172.20.0.27 has DC/OS version 1.11.0
agent 172.20.0.28 has DC/OS version 1.11.0
agent 172.20.0.29 has DC/OS version 1.11.0
agent 172.20.0.21 has DC/OS version 1.10.1
agent 172.20.0.25 has DC/OS version 1.11.0
public agent 172.20.0.26 has DC/OS version 1.11.0

[PROBLEM] "health" - Problems were found.
---------------
Problem details
---------------
agent 172.20.0.21: The following components are not healthy:
dcos-docker-gc.service: health = 1

[OK] "mesos-actor-mailboxes" - All Mesos actors are fine.
[OK] "node-count" - Masters: 3, Agents: 5, Public Agents: 1, Total: 9
```

Please, launch the following command to learn more:

```
$ bun --help
```

## How to contribute

Please, report bugs and share your ideas for new features via the [issue page](https://github.com/mesosphere/bun/issues).

The project is written in Go; please, use [the latest version](https://golang.org/dl/) of the compiler.

To add a new feature or fix a bug, simply
`git clone https://github.com/mesosphere/bun.git` and use your favorite
editor or IDE.

To test your changes, simply build the CLI and launch it against some bundle:

```
$ go build
$ ./go -p <path to a bundle directory>
```

### How to add new checks

Each check uses one or more bundle files. Please, refer to the `filetypes/files_type_yaml.go`
file to find out a name of the file type by the file name and vice versa.

#### Simple search check

Simple search check fails when a specified string is found in a
specified file of a specified type. Use it when you'd like to check if the log file
contains a record which appears only when the error occurs in the correspondent
DC/OS component.

To create a new search check, simply add a YAML definition to the YAML document in the
`checks/search_checks_yaml.go` file:

```yaml
...
  - name: disk-space-exhibitor
    description: Check disk space errors in Exhibitor logs
    fileTypeName: exhibitor-log
    errorPattern: No space left on device
```

#### Check if a certain condition is fulfilled on each node

If you would like to check for a certain condition on each node of a certain role
(i.e.: master, agent or public agent), please use the `bun.CheckBuilder` with a default
aggregate function:

```go
package health

import (
	"fmt"
	"github.com/adyatlov/bun/filetypes"
	"strings"

	"github.com/adyatlov/bun"
)

func init() {
	builder := bun.CheckBuilder{
		Name:                    "diagnostics-health",
		Description:             "Check if all DC/OS components are healthy",
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
		Aggregate:               bun.DefaultAggregate,
	}
	check := builder.Build()
	bun.RegisterCheck(check)
}

func collect(host bun.Host) (ok bool, details interface{}, err error) {
	h := filetypes.Host{}
	if err = host.ReadJSON("diagnostics-health", &h); err != nil {
		return
	}
	unhealthy := []string{}
	for _, u := range h.Units {
		if u.Health != 0 {
			unhealthy = append(unhealthy,
				fmt.Sprintf("%v: health = %v", u.ID, u.Health))
		}
	}
	if len(unhealthy) > 0 {
		details = fmt.Sprintf("The following components are not healthy:\n%v",
			strings.Join(unhealthy, "\n"))
		ok = false
	} else {
		ok = true
	}
	return
}
```

#### More complex checks

If you need a check which requires analysis of a collected data, you can use a custom
aggregate function:

```go
package dcosversion

import (
	"fmt"
	"github.com/adyatlov/bun"
	"github.com/adyatlov/bun/filetypes"
)

func init() {
	builder := bun.CheckBuilder{
		Name: "dcos-version",
		Description: "Verify that all hosts in the cluster have the " +
			"same DC/OS version installed",
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
		Aggregate:               aggregate,
	}
	check := builder.Build()
	bun.RegisterCheck(check)
}

func collect(host bun.Host) (ok bool, details interface{}, err error) {
	v := filetypes.Version{}
	if err = host.ReadJSON("dcos-version", &v); err != nil {
		return
	}
	details = v.Version
	ok = true
	return
}

func aggregate(c *bun.Check, b bun.CheckBuilder) {
	version := ""
	// Compare versions
	details := []string{}
	ok := true
	for _, r := range b.OKs {
		v := r.Details.(string)
		if version == "" {
			version = v
		}
		if v != version {
			ok = false
		}
		details = append(details, fmt.Sprintf("%v %v has DC/OS version %v",
			r.Host.Type, r.Host.IP, v))
	}
	// No need to interpret problems, as we didn't create it in the host check.
	if ok {
		c.OKs = details
		c.Summary = fmt.Sprintf("All versions are the same: %v.", version)
	} else {
		c.Problems = details
		c.Summary = "Versions are different."
	}
}
```

### How to release

1. Install [GoReleaser](https://goreleaser.com/install/).
2. Create [Github personal access token](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)
    with the `repo` scope and export it as an environment variable called `GITHUB_TOKEN`:

  	```bash
  	$ export GITHUB_TOKEN=<your personal GitHub access token>
  	```

    Please find more information about this step [here](https://goreleaser.com/environment/).
3. Create a Git tag which adheres to [semantic versioning](https://semver.org/) and
    push it to GitHub:

    ```bash
    $ git tag -a v1.9.8 -m "Release v1.9.8"
    $ git push origin v1.9.8
    ```

    If you made a mistake on this step, you can delete the tag remotely and locally:

    ```bash
    $ git push origin :refs/tags/v1.9.8
    $ git tag --delete v1.9.8
    ```

4. Test that the build works with the following command:

    ```bash
    $ goreleaser release --skip-publish --rm-dist
    ```

5. If everything is fine publish the build with the following command:

    ```bash
	$ goreleaser release --rm-dist
    ```
