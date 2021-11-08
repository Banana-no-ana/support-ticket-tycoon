for X in $(find . -name "*.proto" | sed "s|^\./||"); do
	protoc -I$(pwd) --go_out=paths=source_relative:. --go-grpc_out=. --go-json_out=. --go-grpc_opt=paths=source_relative $X
done
