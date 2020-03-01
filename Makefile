proto:
	protoc -I./proto --go_out=plugins=grpc:./proto proto/*.proto

proto-docs:
	@mkdir -p ./proto/docs
	@protoc -I./proto --doc_out=./proto/docs --doc_opt=html,index.html proto/*.proto
