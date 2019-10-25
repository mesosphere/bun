package checks

const searchChecksYAML = `
- name: unmount-volume
  description: Checks if Mesos agents had problems unmounting local persistent volumes. MESOS-8830
  fileTypeName: mesos-agent-log
  errorPattern: 'Failed to remove rootfs mount point'
  cure: Please, refer to the KB article https://support.d2iq.com/s/article/DC-OS-Impacted-by-a-Mesos-Agent-Garbage-Collection-Issue and MESOS-8830 

- name: exhibitor-disk-space
  description: Check disk space errors in Exhibitor logs
  fileTypeName: exhibitor-log
  errorPattern: 'No space left on device'
  cure: Please check that there is sufficient free space on the disk.

- name: migration-in-progress
  description: Detect marathon-upgrade-in-progress flag on failed cluster after upgrade
  fileTypeName: marathon
  errorPattern: 'Migration Failed: Migration is already in progress'
  cure: Please refer to the KB article https://support.d2iq.com/s/article/marathon-migration-in-progress-error

- name: networking-errors
  description: Identify errors in dcos-net logs
  fileTypeName: net-log
  errorPattern: '\[(?P<Level>error|emergency|critical|alert)\]'
  isErrorPatternRegexp: true
  cure: 'Please, collect the crash dumps with "sudo tar -czvf 172.29.108.26_master_dcos_net.tgz -C /opt/mesosphere/active/dcos-net/ ." and contact the networking team.'

- name: zookeeper-fsync
  description: Detects ZooKeeper problems with the write-ahead log
  fileTypeName: exhibitor-log
  errorPattern: 'fsync-ing the write ahead log in'
  max: 1
  cure: Zookeeper fsync threshold exceeded events detected. Zookeeper is swapping or disk IO is saturated. See more here https://jira.mesosphere.com/browse/COPS-4403'

- name: cockroach-time-sync
  descriptions: Detects CockroachDB time sync issues
  fileTypeName: cockroach-log
  errorPattern: 'fewer than half the known nodes are within the maximum offset'
  cure: CockroachDB logs indicate that there is or was an issue with time sync. Please ensure that time is in sync and CockroachDB is healthy on all Masters

- name: time-sync
  description: Checks if time is syncronized on the host machine.
  fileTypeName: net-log
  errorPattern: '(internal consistency is broken|Unable to determine clock sync|Time is not synchronized|Clock is less stable than allowed|Clock is out of sync)'
  isErrorPatternRegexp: true
  curePattern: 'Time is in sync'
  cure: Check NTP settings and NTP server availability.

- name: zookeeper-instances
  description: Checks if all ZooKeeper instances are up and running
  fileTypeName: net-log
  errorPattern: 'Exception: Expected.*servers'
  isErrorPatternRegexp: true
  curePattern: 'Zookeeper connection established, state: CONNECTED'
  cure: Make all ZooKeeper instances run and available for each other through the network.

- name: mesos-agent-invalid-cert
  description: Checks if there are errors for invalid certificate when fetching artifacts
  fileTypeName: mesos-agent-log
  errorPattern: 'Container.*Failed to perform ''curl''.*SSL certificate problem: self signed certificate' 
  isErrorPatternRegexp: true
  cure: 'Mesos agent is using certificates which does not allow to fetch an artifact from some repository. Please see https://jira.mesosphere.com/browse/COPS-2315 and https://jira.mesosphere.com/browse/COPS-2106 for more information.'

- name: overlay-network-recovery
  descriptions: Checks if the DC/OS overlay network master is in recovery state
  fileTypeName: mesos-master-log
  errorPattern: 'overlay-master in.*RECOVERING.*state'
  isErrorPatternRegexp: true
  curePattern: 'Moving overlay-master.* to .*RECOVERED.* state.'
  isCurePatternRegexp: true
  cure: 'Mesos master Overlay module cannot recover. Please see the KB articles https://support.d2iq.com/s/article/Known-Issue-Invalid-DNS-Resolvers-MSPH-2018-0012 and https://support.d2iq.com/s/article/Critical-Issue-with-Overlay-Networking for more information.'

- name: kmem-errors
  description: Detects kernel memory (kmem) errors in dmesg log
  fileTypeName: dmesg-log
  errorPattern: '(SLUB: Unable to allocate memory on node -1|task .+ blocked for more than .+ seconds)'
  isErrorPatternRegexp: true
  cure: 'Please see KB articles https://support.mesosphere.com/s/article/Critical-Issue-KMEM-MSPH-2018-0006 and https://support.mesosphere.com/s/article/Known-Issue-KMEM-with-Kubernetes-MSPH-2019-0002' 

- name: oom-kills
  description: Detects out of memory kills in dmesg log
  fileTypeName: dmesg-log
  errorPattern: 'invoked oom-killer'
  cure: 'The operating system is killing processes which exceed system or container memory limits. Please check which processes are getting killed. If it is a DC/OS container, increase its memory limit.'

- name: docker-running
  description: Checks if docker is running
  fileTypeName: ps
  errorPattern: 'dockerd'
  failIfNotFound: true
  cure: 'Docker daemon should be running on all DC/OS nodes.'
`
