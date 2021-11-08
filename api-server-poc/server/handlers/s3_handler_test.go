package handlers

import (
	pb "api-server-poc/proto/generated"
	"context"
	"reflect"
	"testing"
)

func TestS3ManagerServiceServer_Upload(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *pb.UploadRequest
	}
	tests := []struct {
		name    string
		server  *S3ManagerServiceServer
		args    args
		want    *pb.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &S3ManagerServiceServer{}
			got, err := server.Upload(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3ManagerServiceServer.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S3ManagerServiceServer.Upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3ManagerServiceServer_Download(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *pb.DownloadRequest
	}
	tests := []struct {
		name    string
		server  *S3ManagerServiceServer
		args    args
		want    *pb.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &S3ManagerServiceServer{}
			got, err := server.Download(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3ManagerServiceServer.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S3ManagerServiceServer.Download() = %v, want %v", got, tt.want)
			}
		})
	}
}
