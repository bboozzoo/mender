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

Use example config file as template (edit as needed):

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
