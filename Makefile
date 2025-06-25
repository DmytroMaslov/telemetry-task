SENSOR_NAME= demo
SERVER_CONFIG_FILE= ./artifacts/sink/config.yaml

build:
	GOBIN=${PWD}/tmp/bin go install -mod=vendor ./cmd/...

run-sensor:
	go run \
		-mod=vendor \
		./cmd/sensor \
		--name=${SENSOR_NAME} \
		--addr="localhost:8080"

run-sink:
	go run \
		-mod=vendor \
		./cmd/sink \
		--config=${SERVER_CONFIG_FILE}

clean:
	rm -f -r ./tmp


generate-proto:
	docker run --platform linux/amd64 -v ${PWD}:/defs namely/protoc-all -f protocol/telemetry.proto -l go -o protocol 