# docker-file-log-driver

File log driver for Docker that sends all of the containers output to a File server. The code is heavily inspired by https://github.com/cpuguy83/docker-log-driver-test.

## Background

We use File as a reliable and very fast transport for logs, but not as a storage system. The excellent
Logstash File input plugin (https://www.elastic.co/guide/en/logstash/current/plugins-inputs-redis.html) is highly recommended
to pick up the logs and transport them to whatever output you like, for example Elasticsearch
(https://www.elastic.co/guide/en/logstash/current/plugins-outputs-elasticsearch.html).

## Features

* Send containers stdout/stderr to a File list (via `RPUSH`)
* Integrates seamlessly with orchestration platforms like Kubernetes, Mesos/Marathon or Docker Swarm
* Log lines come as JSON messages with all important container meta data
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
* Message payload may be arbitrarily complex (e.g. JSON encoded)
* Configure logging setup either globally through Docker `config.json` or per container (`--log-opt` style)
* Customizable File connection timeouts
* Automatic File reconnects

## Requirements

* Docker >= 17.05 (since that brings log driver plugin support).
* Docker for Windows isn't supported at the moment (see https://docs.docker.com/engine/extend/)

## Install

```
$ docker plugin install deep-compute/docker-file-log-driver:0.0.1 --alias file-log-driver
Plugin "deep-compute/docker-file-log-driver:0.0.1" is requesting the following privileges:
 - network: [host]
Do you grant the above permissions? [y/N] y
932e2beac35d: Download complete
Digest: sha256:5ae2850ddd18571da68827bac08cccbaa36d1b02022802e90028dfe931fc1e9f
Status: Downloaded newer image for deep-compute/docker-file-log-driver:0.0.1
Installed plugin deep-compute/docker-file-log-driver:0.0.1
```

## Usage

### Basic usage

Run a container using this plugin:

```
$ docker run --log-driver file-log-driver --log-opt file-address=some.file.server:6379 --log-opt file-password=secure --log-opt file-list=logs alpine date
```

Observe the File list named `logs` in database `0` of your File instance. You should see a JSON like the following:

```
{
  "message": "Thu Jun  8 14:36:01 UTC 2017",
  "container_id": "20cd880a54679f26c85edc53a7fbe7079d7e5b88a883f245a0c215d0afd3e600",
  "container_name": "hungry_cori",
  "container_created": "2017-06-08T14:36:01.08719318Z",
  "image_id": "sha256:baa5d63471ead618ff91ddfacf1e2c81bf0612bfeb1daf00eb0843a41fbfade3",
  "image_name": "alpine",
  "command": "date",
  "tag": "20cd880a5467",
  "extra": {},
  "host": "hostname",
  "timestamp": "2017-06-08T14:36:01.64922286Z"
}
```

### Advanced usage

This example shows the usage of

* Custom log tags (c.f. https://docs.docker.com/engine/admin/logging/log_tags/)
* Container label logging
* Container environment variable logging

```
$ docker run --label foo=abc --label bar=xyz -e SOME_ENV_VAR=foobar --log-driver file-log-driver --log-opt redis-address=some.redis.server:6379 --log-opt redis-password=secure --log-opt redis-list=logs --log-opt "tag={{.ImageName}}/{{.Name}}/{{.ID}}" --log-opt labels=foo,bar --log-opt env=SOME_ENV_VAR alpine date
```

Observe the File list named `logs` in database `0` of your File instance. You should see a JSON like the following:

```
{
  "message": "Thu Jun  8 14:45:22 UTC 2017",
  "container_id": "35fb6802c95dce77147b89f595ed7528ebf364c8a795ef1468f52464bfd8fb23",
  "container_name": "trusting_tesla",
  "container_created": "2017-06-08T14:45:22.253419923Z",
  "image_id": "sha256:baa5d63471ead618ff91ddfacf1e2c81bf0612bfeb1daf00eb0843a41fbfade3",
  "image_name": "alpine",
  "command": "date",
  "tag": "alpine/trusting_tesla/35fb6802c95d",
  "extra": {
    "SOME_ENV_VAR": "foobar",
    "bar": "xyz",
    "foo": "abc"
  },
  "host": "workstation",
  "timestamp": "2017-06-08T14:45:22.625446437Z"
}
```

### Options

All available options are documented here and can be set via `--log-opt KEY=VALUE`. Timeouts need to be specified in a format supported by https://golang.org/pkg/time/#ParseDuration.

|Key|Default|Description|
|---|---|---|
|redis-address||TCP address to connect to in the form `host:port`|
|redis-password||File password|
|redis-database|0|File database index|
|redis-list||File variable to append logs to|
|redis-connect-timeout|5s|Timeout when connecting to File|
|redis-read-timeout|5s|Timeout when reading from File|
|redis-write-timeout|5s|Timeout when writing to File|

## Uninstall

To uninstall, please make sure that no containers are still using this plugin. After that, disable and remove the plugin like this:

```
$ docker plugin disable file-log-driver
$ docker plugin rm file-log-driver
```

## Hack it

You're more than welcome to hack on this. PRs are also welcome :-)

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
