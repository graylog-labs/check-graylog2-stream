Check Graylog2 Stream
=====================

A simple Icinga/Icinga2/Nagios check to monitor Graylog2 stream alerts.

### Install
See [releases](http://github.com/graylog2/check-graylog2-stream/releases) for pre-compiled
binaries. No dependencies are needed.

### Usage
Following options can be set

```shell
$ ./bin/check-graylog2-stream
usage:
  -condition="<ID>": Condition ID, set only to check a single alert (optional)
  -password="<password>": API password (mandatory)
  -stream="<ID>": Stream ID (mandatory)
  -url="http://localhost:12900": URL to Graylog2 api (optional)
  -user="<username>": API username (mandatory)
```

To check a stream on a remote server configure this check in your monitoring config

```shell
check-graylog2-stream -stream=545b8c15e4b07ae85aee40d1 -user=admin -password=secret -url='http://172.16.0.1:12900'
```

You can also check a single condition of a stream

```shell
check-graylog2-stream -stream=545b8c15e4b07ae85aee40d1 -user=admin -password=secret -url='http://172.16.0.1:12900' -condition=eeae1109-7cba-4fa0-a35a-8aa7d162ed54
```

To figure out which stream or codition IDs to use, query the Graylog2 API
```shell
curl -i --user admin:secret -H 'Accept: application/json' 'http://172.16.0.1:12900/streams?pretty=true'
```

### Build
If you need to build the check yourself, you can do it like this:

```shell
go get github.com/fractalcat/nagiosplugin
go build -o bin/check-graylog2-stream src/check-graylog2-stream/check-graylog2-stream.go
```
