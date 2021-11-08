// Golang
//go:generate protoc -I. --go_opt=paths=source_relative --go_out=plugins=grpc:./generated --go_opt=paths=source_relative ./s3.proto
//go:generate protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative,allow_delete_body=true:./generated --swagger_out=logtostderr=true,allow_delete_body=true,repeated_path_param_separator=ssv:./generated ./s3.proto
package proto
