# Bun

Command-line program which detects the most common problems in a DC/OS cluster
by analyzing its [diagnostics bundle](https://docs.mesosphere.com/1.11/cli/command-reference/dcos-node/dcos-node-diagnostics-create/).

```
$ bun
+-------------+----------------------------------------------------------+
| Check       | docker-not-running                                       |
+-------------+----------------------------------------------------------+
| Status      | [UNDEFINED]                                              |
+-------------+----------------------------------------------------------+
| Description | Checks if docker is running                              |
+-------------+----------------------------------------------------------+
| Summary     | Couldn't check any hosts because of the error(s). Launch |
|             | this check with the -v flag to see the details.          |
+-------------+----------------------------------------------------------+

+-------------+----------------------------------------------------------+
| Check       | firewalld-running                                        |
+-------------+----------------------------------------------------------+
| Status      | [UNDEFINED]                                              |
+-------------+----------------------------------------------------------+
| Description | Detects if firewalld is running on a DC/OS node          |
+-------------+----------------------------------------------------------+
| Summary     | Couldn't check any hosts because of the error(s). Launch |
|             | this check with the -v flag to see the details.          |
+-------------+----------------------------------------------------------+

+------------------------+-----------------------------------------------------------+
| Check                  | mesos-agent-invalid-cert                                  |
+------------------------+-----------------------------------------------------------+
| Status                 | [PROBLEM]                                                 |
+------------------------+-----------------------------------------------------------+
| Description            | Checks if there are errors for invalid certificate when   |
|                        | fetching artifacts                                        |
+------------------------+-----------------------------------------------------------+
| Cure                   | Mesos agent is using certificates which does not allow to |
|                        | fetch an artifact from some repository. Please see        |
|                        | https://jira.mesosphere.com/browse/COPS-2315 and          |
|                        | https://jira.mesosphere.com/browse/COPS-2106 for more     |
|                        | information.                                              |
+------------------------+-----------------------------------------------------------+
| Summary                | Error pattern "Container.*Failed to perform 'curl'.*SSL   |
|                        | certificate problem: self signed certificate" found.      |
+------------------------+-----------------------------------------------------------+
| [P] agent 10.10.10.105 | Error pattern occurred 10 time(s) in file                 |
|                        | dcos-mesos-slave.service.gz                               |
+------------------------+-----------------------------------------------------------+
| [P] agent 10.10.10.104 | Error pattern occurred 1 time(s) in file                  |
|                        | dcos-mesos-slave.service.gz                               |
+------------------------+-----------------------------------------------------------+

+------------------------+--------------------------------------------------------------+
| Check                  | networking-errors                                            |
+------------------------+--------------------------------------------------------------+
| Status                 | [PROBLEM]                                                    |
+------------------------+--------------------------------------------------------------+
| Description            | Identifies errors in dcos-net logs                           |
+------------------------+--------------------------------------------------------------+
| Cure                   | Please, collect the crash dumps with "sudo tar -czvf         |
|                        | 172.29.108.26_master_dcos_net.tgz -C                         |
|                        | /opt/mesosphere/active/dcos-net/ ." and contact the          |
|                        | networking team.                                             |
+------------------------+--------------------------------------------------------------+
| Summary                | Error pattern                                                |
|                        | "\[(?P<Level>error|emergency|critical|alert)\]" found.       |
+------------------------+--------------------------------------------------------------+
| [P] agent 10.10.10.105 | Error pattern occurred 1 time(s) in file dcos-net.service.gz |
+------------------------+--------------------------------------------------------------+

+-------------+-------------------------------------------------------------+
| Check       | nscd-running                                                |
+-------------+-------------------------------------------------------------+
| Status      | [UNDEFINED]                                                 |
+-------------+-------------------------------------------------------------+
| Description | Detects if Name Service Cache Daemon (nscd) is running on a |
|             | DC/OS node                                                  |
+-------------+-------------------------------------------------------------+
| Summary     | Couldn't check any hosts because of the error(s). Launch    |
|             | this check with the -v flag to see the details.             |
+-------------+-------------------------------------------------------------+

+------------------------+------------------------------------------------------------+
| Check                  | oom-kills                                                  |
+------------------------+------------------------------------------------------------+
| Status                 | [PROBLEM]                                                  |
+------------------------+------------------------------------------------------------+
| Description            | Detects out of memory kills in dmesg log                   |
+------------------------+------------------------------------------------------------+
| Cure                   | The operating system is killing processes which exceed     |
|                        | system or container memory limits. Please check which      |
|                        | processes are getting killed. If it is a DC/OS container,  |
|                        | increase its memory limit.                                 |
+------------------------+------------------------------------------------------------+
| Summary                | Error pattern "invoked oom-killer" found.                  |
+------------------------+------------------------------------------------------------+
| [P] agent 10.10.10.104 | Error pattern occurred 3 time(s) in file dmesg-0.output.gz |
+------------------------+------------------------------------------------------------+
| [P] agent 10.10.10.105 | Error pattern occurred 2 time(s) in file dmesg-0.output.gz |
+------------------------+------------------------------------------------------------+

+-----------+----+
|  SUMMARY  |    |
+-----------+----+
| Failed    |  3 |
| Undefined |  3 |
| Passed    | 16 |
+-----------+----+
|   TOTAL   | 22 |
+-----------+----+
```

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
## Update

Bun checks for its new versions and updates itself automatically with your permission.

## Usage

```bash
$ bun -p <path to bundle directory>
```

or if the working directory is the bundle directory simply:

```bash
$ bun
```

Please, launch the following command to learn more:

```
$ bun --help
```

## How to contribute

Please, report bugs and share your ideas for new features via the [issue page](https://github.com/mesosphere/bun/issues).

The project is written in Go; please, use [the latest version](https://golang.org/dl/) of the compiler.

To add a new feature or fix a bug, clone the repository:
`git clone https://github.com/mesosphere/bun.git` and use your favorite
editor or IDE.

To test your changes, simply build the CLI and launch it against some bundle:

```
$ go build
$ ./go -p <path to a bundle directory>
```

### Bundle files

Names of DC/OS diagnostics bundle files may vary from one version of DC/OS to another, moreover, they are not always
descriptive or handy. That is why in Bun we give each file a human-readable self-explanatory ID and use these IDs
to refer to the bundle files. File `files_type_yaml.go` contains description of bundle files.

### How to add new checks

The core abstraction of the Bun tool is `checks.Check`:

```go
package checks

type Check struct {
	Name           string 
	Description    string
	Cure           string
	OKSummary      string
	ProblemSummary string
	Run            CheckBundleFunc 
}

type CheckBundleFunc func(bundle.Bundle) Results

type Result struct {
	Status Status
	Value  interface{}
	Host   bundle.Host
}
```

To add a new check you need to create an instance of that struct, describe the check by specifying its string fields,
and provide a Run function, which does actual testing. 

To make adding checks easier, Bun source code provides some help
like ability to declare checks as a YAML object, or Run function builder. You can read about that in the next sections.

#### Search check

Search checks are looking for a specified strings or regular expressions in a bundle file to detect or rule out a 
specific problem. Also, search checks is very easy to add -- you don't even need to write a code. 

To create a new search check, simply add a new object to the YAML document in the
`checks/search_checks_yaml.go` file. For example: 

```yaml
- name: exhibitor-disk-space
  description: Checks for disk space errors in Exhibitor logs
  fileTypeName: exhibitor-log
  errorPattern: 'No space left on device'
  cure: Please check that there is sufficient free space on the disk.
```

To avoid false positives, you can specify a a string or regular expression, which manifests that the
problem is gone. For example, the following check will not fail if the string "Time is in sync" appears 
in the networking log after the last "Checks if time is synchronised on the host machine." message.

```yaml
- name: time-sync
  description: Checks if time is synchronised on the host machine.
  fileTypeName: net-log
  errorPattern: '(internal consistency is broken|Unable to determine clock sync|Time is not synchronized|Clock is less stable than allowed|Clock is out of sync)'
  isErrorPatternRegexp: true
  curePattern: 'Time is in sync'
  cure: Check NTP settings and NTP server availability.
```

#### Check if a condition on each node

If you need to check that a certain condition is satisfied on each DC/OS node of a given type (i.e.: master, agent, or public agent), you can 
use the `checks.CheckFuncBuilder`. With its help, you only need to create a function which checks for the condition on
one node. The builder will do the rest. For example, the following check detects a situation when Mesos mailboxes have
too many messages:

```yaml
...
	builder := checks.CheckFuncBuilder{
		CheckMasters:      collect,
		CheckAgents:       collect,
		CheckPublicAgents: collect,
	}
	check := checks.Check{
		Name: "mesos-actor-mailboxes",
		Description: "Checks if actor mailboxes in the Mesos process " +
			"have a reasonable amount of messages",
		Cure: "Check I/O on the correspondent hosts and if something is overloading Mesos agents or masters" +
			" with API calls.",
		OKSummary:      "All Mesos actors are fine.",
		ProblemSummary: "Some Mesos actors are backlogged.",
		Run:            builder.Build(),
	}
...

type MesosActor struct {
	ID     string `json:"id"`
	Events []struct{}
}

func collect(host bundle.Host) checks.Result {
	var actors []MesosActor
	if err := host.ReadJSON("mesos-processes", &actors); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}
	var mailboxes []string
	for _, a := range actors {
		if len(a.Events) > maxEvents {
			mailboxes = append(mailboxes, fmt.Sprintf("(Mesos) %v@%v: mailbox size = %v (> %v)",
				a.ID, host.IP, len(a.Events), maxEvents))
		}
	}
	if len(mailboxes) > 0 {
		return checks.Result{
			Host:   host,
			Status: checks.SProblem,
			Value:  mailboxes,
		}
	}
	return checks.Result{
		Status: checks.SOK,
	}
}
```

If your check needs to analyse the data collected on each node, you can implement an Aggregate function instead of
using the the default one; please see an example in the `dcos-version` (`checks/dcosversion/check.go`) check.

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
