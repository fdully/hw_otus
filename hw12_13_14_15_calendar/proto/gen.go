package proto

// Generating our protofiles and our swagger json docs
//go:generate sh -c "protoc --proto_path=. *.proto --go_out=plugins=grpc:../internal/pb --grpc-gateway_out=../internal/pb --openapiv2_out=:../swaggerui/calendar"

// Generating byte data to for binding to app
//go:generate sh -c "statik -src ../swaggerui -dest ../internal/api/swaggerui"
