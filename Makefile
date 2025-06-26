SENSOR_NAME= demo
SERVER_CONFIG_FILE= ./artifacts/sink/config.yaml
BIN_DIR= ${PWD}/tmp/bin
export SECRET_KEY=1234567890123456

build:
	GOBIN=${BIN_DIR} go install -mod=vendor ./cmd/...

run-sensor:
	${BIN_DIR}/sensor \
		--addr="localhost:8080" \
		--name=${SENSOR_NAME} \
		--rate=100			\
		--cert=./certs/ca-cert.pem

run-sink:
	${BIN_DIR}/sink \
		--config=${SERVER_CONFIG_FILE}

clean:
	rm -f -r ./tmp


generate-proto:
	docker run --platform linux/amd64 -v ${PWD}:/defs namely/protoc-all -f protocol/telemetry.proto -l go -o protocol 

run-sensor-test:
	go run \
		-race \
		-mod=vendor \
		./cmd/sensor \
		--name=${SENSOR_NAME} \
		--addr="localhost:8080" \
		--rate=100

run-sink-test:
	go run \
		-race \
		-mod=vendor \
		./cmd/sink \
		--config=${SERVER_CONFIG_FILE}

gen-cert:
	cd scripts; ./gen-certs.sh; cd ..

decrypt:
	${BIN_DIR}/decrypt \
		--input=./tmp/metrics.txt \
		--output=./tmp/metrics_decrypted.txt \