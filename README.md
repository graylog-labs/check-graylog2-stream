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
  -url="http://localhost:12900": URL to Graylog2 API (optional)
  -user="<username>": API username (mandatory)
```

To check a stream on a remote server configure this check in your monitoring config

```shell
check-graylog2-stream -stream=545b8c15e4b07ae85aee40d1 -user=admin -password=secret -url='http://172.16.0.1:12900'
```

You can also check a single condition of a stream, through HTTP

```shell
check-graylog2-stream -stream=545b8c15e4b07ae85aee40d1 -user=admin -password=secret -url='http://172.16.0.1:12900' -condition=eeae1109-7cba-4fa0-a35a-8aa7d162ed54
```

or HTTPS (see rest_listen_uri and rest_enable_tls in graylog-server configuration)

```shell
check-graylog2-stream -stream=545b8c15e4b07ae85aee40d1 -user=admin -password=secret -url='https://172.16.0.1:12900' -condition=eeae1109-7cba-4fa0-a35a-8aa7d162ed54
```

To figure out which stream or condition IDs to use, query the Graylog2 API
```shell
curl -i --user admin:secret -H 'Accept: application/json' 'http://172.16.0.1:12900/streams?pretty=true'
```

### Build
If you need to build the check yourself, you can do it like this:

```shell
go get github.com/fractalcat/nagiosplugin
go build -o bin/check-graylog2-stream src/check-graylog2-stream/check-graylog2-stream.go
```

### Icinga 2 Integration

Make sure to download and install the plugin into `PluginDir`.
This constant is defined in the `constants.conf` file (default
in `/etc/icinga2/constants.conf`). More details in the
[Icinga 2 documentation](http://docs.icinga.org).

#### New CheckCommand Definition

Add that to `conf.d/commands.conf`  or a similar
included file.

```shell
object CheckCommand "graylog2-stream" {
  import "plugin-check-command"

  command = [ PluginDir + "/check-graylog2-stream" ]

  arguments = {
    "-stream" = "$graylog2_stream_id$"
    "-user" = "$graylog2_api_username$"
    "-password" = "$graylog2_api_password$"
    "-url" = "$graylog2_api_url$"
  }

  // default values
  vars.graylog2_api_url = "http://localhost:12900"
  vars.graylog2_api_username = "admin"
  vars.graylog2_api_password = "yourpassword"
}
```

#### Host and Service Definition

Depending your monitoring configuration strategy, you can define
it like the following example for Icinga2 2.2+:

```shell
object Host "graylog2-host" {
  address = "127.0.0.1"
  check_command = "hostalive"

  /* `icinga2` is the stream name. Used in service apply for rule below */
  vars.streams["icinga2"] = {
    graylog2_stream_id = "54610d26e4b059482bbfab0f"
  }
}

template Service "graylog2-service" {
  check_interval = 30s
  retry_interval = 30s
  max_check_attempts = 3
  enable_flapping = false
  enable_notifications = true
}

const GraylogStreamApiUrl = "http://127.0.0.1:9000/streams/"

apply Service "alert-" for (stream => config in host.vars.streams) {
  import "graylog2-service"

  check_command = "graylog2-stream"

  vars += config

  notes = "My " + stream + " graylog2 alert stream checker."
  notes_url = GraylogStreamApiUrl + stream + "/alerts"

  assign where host.vars.streams
}
```




