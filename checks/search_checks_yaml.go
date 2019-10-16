package checks

const searchChecksYAML = `
- name: unmount-volume
  description: Checks if Mesos agents had problems unmounting local persistent volumes. MESOS-8830
  fileTypeName: mesos-agent-log
  searchString: Failed to remove rootfs mount point
- name: exhibitor-disk-space
  description: Check disk space errors in Exhibitor logs
  fileTypeName: exhibitor-log
  searchString: No space left on device
- name: migration-in-progress
  description: Detect marathon-upgrade-in-progress flag on failed cluster after upgrade
  fileTypeName: marathon
  searchString: 'Migration Failed: Migration is already in progress'
- name: networking-errors
  description: Identify errors in dcos-net logs
  fileTypeName: net
  isRegexp: true
  searchString: \[(?P<Level>error|emergency|critical|alert)\]
- name: zookeeper-fsync
  description: Detects ZooKeeper problems with the write-ahead log
  cure: Zookeeper fsync threshold exceeded events detected. Zookeeper is swapping or disk IO is saturated. See more here https://jira.mesosphere.com/browse/COPS-4403"
  fileTypeName: exhibitor-log
  searchString: fsync-ing the write ahead log in
  max: 1
- name: cockroach-time-sync
  descriptions: Detects CockroachDB time sync issues
  cure: CockroachDB logs indicate that there is or was an issue with time sync. Please ensure that time is in sync and CockroachDB is healthy on all Masters
  fileTypeName: cockroach-log
  searchString: fewer than half the known nodes are within the maximum offset
`
