## Architecture

Overview of Mender client architecture.

```

                                  +--------------------+
                                  |       SERVER       |
                                  |                    |
                                  +-------^---+--------+
                                          |   |
                                          |   |
                                          |   |
         +------------------+     +-------+---v--------+     +-------------------+
         |      daemon      |     |      client        |     |      mender       |
         |                  |     |                    |     |                   |
         +------------------+     +--------------------+     +-------------------+
         |                  |     | Updater            |     | Controller        |
         |                  |     |                    |     |                   |
         +------------------+     +--------------------+     +-------------------+
         | NewDaemon        |     | NewUpdater         |     | NewMender         |
         | Run              |     | NewHttpsClient     |     | GetState          |
         | StopDaemon       <-----+ NewHttpClient      |     | GetCurrentImageID |
         |                  |     | GetScheduledUpdate |     | LoadConfig        |
         |                  |     | FetchUpdate        |     | GetUpdaterConfig  |
         |                  |     | Bootstrap          |     | GetDaemonConfig   |
         |                  |     |                    |     |                   |
         +-------------^-^--+     +--------------------+     +-----+---^---------+
                       | |                                         |   |
                       | |-----------------------------------------+   |     MENDER SERVER INTERFACE
                       |                                               |
+--------------------------------------------------------------------------------------------------+
                       |                                               |
                       |                                               |          HARDWARE INTERFACE
                       |                                               |
                       |                                               |
                       |                                               |
+------------------+   |   +------------------------+        +---------+----------+
|    partitions    |   |   |         device         |        |      bootenv       |
|                  |   |   |                        |        |                    |
+------------------+   |   +------------------------+        +--------------------+
|                  |   |   | UInstaller             |        | BootEnvReadWritter |
|                  |   |   | UInstalCommitRebooter  |        |                    |
+------------------+   |   +------------------------+        +--------------------+
| GetInactive      |   |   | NewDevice              |        | NewEnvironment     |
| GetActive        |   +---+ Reboot                 |        | ReadEnv            |
|                  |       | InstallUpdate          <--------+ WriteEnv           |
|                  |       | EnableUpdatedPartition |        |                    |
|                  +-------> CommitUpdate           |        |                    |
|                  +-------> FetchUpdateFromFile    |        |                    |
|                  |       |                        |        |                    |
+------------------+       +------------------------+        +--------------------+

```

## Hacking on Mender client locally

Build mender with support for local development:
```
$ make LOCAL=1
```

Next step is to create required prefix tree that looks similar to what is present on the device. We want to reach a tree like this:

```
prefix (could be /usr, or other)
├── [drwxr-xr-x]  etc
│   └── [drwxr-xr-x]  mender
│       └── [-rw-r--r--]  mender.conf
├── [drwxr-xr-x]  share
│   └── [drwxr-xr-x]  mender
│       ├── [drwxr-xr-x]  identity
│       │   └── [-rwxr-xr-x]  mender-device-identity
│       └── [drwxr-xr-x]  inventory
│           ├── [-rwxr-xr-x]  mender-inventory-hostinfo
│           └── [-rwxr-xr-x]  mender-inventory-network
└── [drwxr-xr-x]  var
    └── [drwxr-xr-x]  lib
        └── [drwxr-xr-x]  mender
            └── [-rw-r--r--]  authtentoken
```

First, create necessary directories:

```
$ mkdir -p $PWD/prefix/etc/mender \
    $PWD/prefix/share/mender/{identity,inventory} \
    $PWD/prefix/var/lib/mender
```

Provide a dummy tenant token:

```
$ echo dummy > $PWD/prefix/var/lib/mender/authtentoken
```

Use example config file as template (edit as needed, ex. point the client to
http://docker.mender.io:9080):

```
$ cp mender.conf.example $PWD/prefix/etc/mender/mender.conf
```

Use example device identity script:

```
$ cp support/mender-device-identity $PWD/prefix/share/mender/identity/
```

Use example inventory scripts:

```
$ cp support/mender-inventory-* $PWD/prefix/share/mender/inventory/
```

Finally start mender, setting `MENDER_PREFIX` environment variable to point to
our local prefix tree:

```
$ MENDER_PREFIX=$PWD/prefix ./mender -daemon -debug
```

### Mender prefix tree breakdown

```
prefix
├── dev
│   ├── active -> mmcblk0p2       <-- 'active' partition
│   ├── inactive -> mmcblk0p1     <-- 'inactive' partition
│   ├── mmcblk0p1                 <-- fake ..
│   └── mmcblk0p2                 <--   partitions
├── etc
│   └── mender
│       ├── build_mender          <-- build manifest file
│       └── mender.conf           <-- client configuration file
├── share
│   └── mender
│       ├── identity
│       │   └── mender-device-identity         <-- identity script
│       └── inventory                          <-- inventory data scripts
│           ├── mender-inventory-hostinfo
│           └── mender-inventory-network
└── var
    └── lib
        └── mender                <-- client state directory
            ├── authseq           <-- authorization sequence
            ├── authtentoken      <-- tenan token
            ├── authtoken         <-- API authorization token
            ├── deployments.0001.31e6216f-1180-44c8-8e3a-2950e557bb46.log  <-- logs from ..
            ├── deployments.0002.adcc2dc4-6f98-43e4-92a9-dfbae40bc059.log  <--   failed deployments
            ├── deployments.0003.322a879a-f165-41d9-8e0c-4fc706e8bef8.log
            ├── deployments.0004.0f1fc042-6b93-497b-b4a3-7072bab31083.log
            ├── deployments.0005.e8994961-e61b-426d-8db6-0a10de817db4.log
            ├── fake-env          <-- fake bootloader environment
            ├── mender-agent.pem  <-- device key
            └── state             <-- update state
```
