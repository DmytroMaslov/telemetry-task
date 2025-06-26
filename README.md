# telemetry task

This repo contains implementation of following components:

- sink
- sensor
- decrypt (for testing purpose)

## Requirements

- GO 1.23
- Docker(optional) for proto rebuild
- OpenSSL

## Installation

### Clone the repo

```sh
git clone https://github.com/DmytroMaslov/telemetry-task.git
```

### Generate certificates

```sh
make gen-cert
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
cert: ./certs/server-cert.pem
key: ./certs/server-key.pem
```

here:

- `file_path` - file where metric data will be collected
- `buffer_size` - size of sink internal (in memory) buffer; in bytes
- `flush_interval` - how often data from buffer is flushed to disk; in milliseconds
- `rate_limit` - max input flow rate in bytes/sec
- `cert` - path to cert
- `key` - path to key
NOTE: if `cert` or `key` are empty or do not exist sink will run without TLS

### Run the sensor

```sh
make run-sensor
```

or

```sh
./tmp/bin/sensor --addr localhost:8080 --name sensorName --rate 100 --cert=./certs/ca-cert.pem
```

here:

- `addr` - address of sink server
- `name` - name of the sensor
- `rate` - metrics per second
- `cert` - path to cert file
NOTE: if `cert` is empty or does not exist, sensor will be run in insecure mode

### Decrypt result

```sh
make decrypt
```

or

```sh
./tmp/bin/decrypt --input=./tmp/metrics.txt --output=./tmp/metrics_decrypted.txt
```

here:

- `input` - path to file with encrypted data
- `output` - path to file where decrypted data will be saved