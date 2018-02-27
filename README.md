# docker-file-log-driver

File log driver for Docker that sends all of the containers output to a specified File. The code is inspired by https://github.com/pressrelations/docker-redis-log-driver.

## Background

We use File as a reliable and simple storage for logs, Running docker containers stores logs in files with container id file name which is hard to locate because are stored in: 
```
/var/lib/docker/containers/<container id>/<container id>-json.log
```

The excellent [Logagg](https://github.com/deep-compute/logagg) is highly recommended to pick up the logs and transport them to whatever datbase you like.

## Features

* Send containers stdout/stderr to a File.
* Integrates seamlessly with orchestration platforms like Kubernetes, Mesos/Marathon or Docker Swarm

* Output format of logs (dictionary)
  * Level of log at docker container (`level`)
  * JSON messages with all important container meta data (`msg`)
    * Container ID (`container_id`)
    * Container name (`container_name`)
    * Container creation date (`container_created`)
    * Image ID (`image_id`)
    * Image name (`image_name`)
    * Command including `ENTRYPOINT` and arguments (`command`)
    * Log tag as provided via `--log-tag` option (`tag`)
    * Extra information as defined via `--log-opt labels=` or `--log-opt env=` (`extra`)
    * Host that container runs on (`host`)
    * Timestamp when log was generated (`timestamp`)
    * Raw log line produced by child process (`message`)
  * Time of log in docker contrainer (`time`)


* `message` payload may be arbitrarily complex (e.g. JSON encoded)
* You can also configure logging setup either globally through Docker `config.json` or per container (`--log-opt` style)

## Requirements

* Docker >= 17.05 (since that brings log driver plugin support).
* [How to install](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-docker-ce-1)
* Check docker version
    ```
    $ sudo docker version
    Client:
    Version:	17.12.0-ce
    API version:	1.35
    Go version:	go1.9.2
    Git commit:	c97c6d6
    Built:	Wed Dec 27 20:11:19 2017
    OS/Arch:	linux/amd64

    Server:
    Engine:
    Version:	17.12.0-ce
    API version:	1.35 (minimum version 1.12)
    Go version:	go1.9.2
    Git commit:	c97c6d6
    Built:	Wed Dec 27 20:09:53 2017
    OS/Arch:	linux/amd64
    Experimental:	false
    ```
* Docker for Windows isn't supported at the moment (see https://docs.docker.com/engine/extend/)

## Install

```
$ docker plugin install deepcompute/docker-file-log-driver:1.0 --alias file-log-driver
Plugin "deepcompute/docker-file-log-driver:1.0" is requesting the following privileges:
 - network: [host]
 - mount: [/var/log]
Do you grant the above permissions? [y/N] y
1.0: Pulling from deepcompute/docker-file-log-driver
a019fc3de34c: Download complete 
Digest: sha256:5b785ded313acd0881c589c5f588f19b3ec3b5300230684a5a7ab1ed1c65e400
Status: Downloaded newer image for deepcompute/docker-file-log-driver:1.0
Installed plugin deepcompute/docker-file-log-driver:1.0
```
## Check
```
$ docker plugin ls
ID                  NAME                     DESCRIPTION         ENABLED
8c7587db6fdc        file-log-driver:latest   File log driver     true

```

## Usage

### Basic usage

Run a container using this plugin:

```
$ docker run --log-driver file-log-driver --log-opt fpath=/testing/test.log alpine date
Tue Feb 27 06:13:36 UTC 2018
```
**Note:** log file `--log-opt fpath` is stored inside **/var/log/fpath**
i.e. fpath=/testing/test.log originally is stored in **/var/log/testing/test.log**

Observe the logs inside path **/var/log/fpath**

```
$ sudo cat /var/log/testing/test.log |jq -r '.msg'| jq -r '.'
{
  "message": "Tue Feb 27 06:13:36 UTC 2018",
  "container_id": "dd3662cd039c3ffb39cdddfe16bf8ca4ee7eeae25080df20f51b1bdf7d6b2f1a",
  "container_name": "blissful_varahamihira",
  "container_created": "2018-02-27T06:13:35.445157368Z",
  "image_id": "sha256:3fd9065eaf02feaf94d68376da52541925650b81698c53c6824d92ff63f98353",
  "image_name": "alpine",
  "command": "date",
  "tag": "dd3662cd039c",
  "extra": {},
  "host": "deepcompute-ThinkPad-E470",
  "timestamp": "2018-02-27T06:13:36.58459717Z"
}
```

### Advanced usage

This example shows the usage of

* Custom log tags (c.f. https://docs.docker.com/engine/admin/logging/log_tags/)
* Container label logging
* Container environment variable logging

```
$ docker run --label foo=abc --label bar=xyz -e SOME_ENV_VAR=foobar --log-driver file-log-driver --log-opt fpath=/testing/test2.log --log-opt "tag={{.ImageName}}/{{.Name}}/{{.ID}}" --log-opt labels=foo,bar --log-opt env=SOME_ENV_VAR alpine date
Tue Feb 27 07:20:32 UTC 2018
```

Observe the logs inside path **/var/log/fpath**

```
$ sudo cat /var/log/testing/test2.log |jq -r '.msg'| jq -r '.'
{
  "message": "Tue Feb 27 07:20:32 UTC 2018",
  "container_id": "3332274db729219ed738458eb120ddc64436513199d182fa7a9fe635363983ce",
  "container_name": "serene_leavitt",
  "container_created": "2018-02-27T07:20:31.729281471Z",
  "image_id": "sha256:3fd9065eaf02feaf94d68376da52541925650b81698c53c6824d92ff63f98353",
  "image_name": "alpine",
  "command": "date",
  "tag": "alpine/serene_leavitt/3332274db729",
  "extra": {
    "SOME_ENV_VAR": "foobar",
    "bar": "xyz",
    "foo": "abc"
  },
  "host": "deepcompute-ThinkPad-E470",
  "timestamp": "2018-02-27T07:20:32.948198669Z"
}
```

### Options

All available options are documented here and can be set via `--log-opt KEY=VALUE`. Timeouts need to be specified in a format supported by https://golang.org/pkg/time/#ParseDuration.

|Key|Default|Description|
|---|---|---|
|`fpath`|/var/log/docker/docker_file_log_driver_default.log|File path of the log file inside /var/log|
|`max-size`|10|size in mb of each log file|
|`max-backups`|10|number of log file backups after `max-size` is reched|
|`max-age`|100|number of days log files are kept in the file system before dleting|

## Uninstall

To uninstall, please make sure that no containers are still using this plugin. After that, disable and remove the plugin like this:

```
$ docker plugin disable file-log-driver
$ docker plugin rm file-log-driver
```

## Hack it

You're more than welcome to hack on this.:-)

```
$ git clone https://github.com/deep-compute/docker-file-log-driver
$ cd docker-file-log-driver
$ docker build -t docker-file-log-driver .
$ ID=$(docker create docker-file-log-driver true)
$ mkdir rootfs
$ docker export $ID | tar -x -C rootfs/
$ docker plugin create file-log-driver .
$ docker plugin enable file-log-driver
```

