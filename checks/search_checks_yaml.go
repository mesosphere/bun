package checks

const searchChecksYAML = `
- name: unmount-volume
  description: Checks if Mesos agents had problems unmounting local persistent volumes. MESOS-8830
  fileTypeName: mesos-agent-log
  errorPattern: 'Failed to remove rootfs mount point'
- name: exhibitor-disk-space
  description: Check disk space errors in Exhibitor logs
  fileTypeName: exhibitor-log
  errorPattern: 'No space left on device'
- name: migration-in-progress
  description: Detect marathon-upgrade-in-progress flag on failed cluster after upgrade
  fileTypeName: marathon
  errorPattern: 'Migration Failed: Migration is already in progress'
- name: networking-errors
  description: Identify errors in dcos-net logs
  fileTypeName: net-log
  errorPattern: '\[(?P<Level>error|emergency|critical|alert)\]'
  isErrorPatternRegexp: true
- name: zookeeper-fsync
  description: Detects ZooKeeper problems with the write-ahead log
  cure: Zookeeper fsync threshold exceeded events detected. Zookeeper is swapping or disk IO is saturated. See more here https://jira.mesosphere.com/browse/COPS-4403'
  fileTypeName: exhibitor-log
  errorPattern: 'fsync-ing the write ahead log in'
  max: 1
- name: cockroach-time-sync
  descriptions: Detects CockroachDB time sync issues
  cure: CockroachDB logs indicate that there is or was an issue with time sync. Please ensure that time is in sync and CockroachDB is healthy on all Masters
  fileTypeName: cockroach-log
  errorPattern: 'fewer than half the known nodes are within the maximum offset'
- name: time-sync
  description: Checks if time is syncronized on the host machine.
  cure: Check NTP settings and NTP server availability.
  fileTypeName: net-log
  errorPattern: '(internal consistency is broken|Unable to determine clock sync|Time is not synchronized|Clock is less stable than allowed|Clock is out of sync)'
  isErrorPatternRegexp: true
  curePattern: 'Time is in sync'
- name: zookeeper-instances
  description: Checks if all ZooKeeper instances are up and running
  cure: Make all ZooKeeper instances run and available for each other through the network.
  fileTypeName: net-log
  errorPattern: 'Exception: Expected.*servers'
  isErrorPatternRegexp: true
  curePattern: 'Zookeeper connection established, state: CONNECTED'
`
