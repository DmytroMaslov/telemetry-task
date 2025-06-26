# telemetry task

This repo contains implementation of two components:

- sink
- sensor

## Requirements

- GO 1.23
- Docker(optional) for proto rebuild

## Installation

### Clone the repo

```sh
git clone https://github.com/DmytroMaslov/telemetry-task.git
```

### Build binary

```sh
make build
```

### Run the sink

```sh
make run-sink
```

or

```sh
./tmp/bin/sink --config  ./artifacts/sink/config.yaml
```

here:

`config` - path to configuration file in yaml format

example of configuration:

```yaml
bind_address: "localhost:8080"
file_path: "./tmp/metrics.txt"
buffer_size: 1024
flush_interval: 100
rate_limit: 1048576
```

here:

- `file_path` - file where metric data will be collected
- `buffer_size` - size of sink internal (in memory) buffer; in bytes
- `flush_interval` - how often data from buffer is flushed to disk; in milliseconds
- `rate_limit` - max input flow rate in bytes/sec

### Run the sensor

```sh
make run-sensor
```

or

```sh
./tmp/bin/sensor --addr localhost:8080 --name sensorName
```

here:

- `addr` - address of sink server
- `name` - name of the sensor
