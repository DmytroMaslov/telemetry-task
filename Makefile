build-sensor:
	go build -o ./tmp/bin/sensor ./cmd/sensor/*
	chmod +x ./tmp/bin/sensor

build-sink:
	go build -o ./tmp/bin/sink ./cmd/sink/*
	chmod +x ./tmp/bin/sink

clean:
	rm -f -r ./tmp

generate-proto:
#	protoc --proto_path=. --go_out=. protocol/telemetry.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	protocol/telemetry.proto