package main

import (
	"context"
	"testing"
	"time"

	pb "api-server-poc/proto/generated"

	"google.golang.org/grpc"
)

var (
	TestDataSet_Upload   []*pb.UploadRequest
	TestDataSet_Downlaod []*pb.DownloadRequest

	MaxDataSetSize = 100
)

func init() {
	TestDataSet_Upload = make([]*pb.UploadRequest, 0)
	TestDataSet_Downlaod = make([]*pb.DownloadRequest, 0)

	for i := 0; i < MaxDataSetSize; i++ {
		uploadData := &pb.UploadRequest{}
		downloadData := &pb.DownloadRequest{}

		TestDataSet_Upload = append(TestDataSet_Upload, uploadData)
		TestDataSet_Downlaod = append(TestDataSet_Downlaod, downloadData)
	}
}

func Test_GRPC_Upload(t *testing.T) {
	conn, err := grpc.Dial("localhost:12000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewS3ManagerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, data := range TestDataSet_Upload {
		r, _ := c.Upload(ctx, data)
		if err != nil {
			t.Errorf("could not request: %v", err)
			return
		}
		t.Logf("Upload Result: %v", r)
	}

}

func Test_GRPC_Download(t *testing.T) {
	conn, err := grpc.Dial("localhost:12000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewS3ManagerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, data := range TestDataSet_Downlaod {
		r, _ := c.Download(ctx, data)
		if err != nil {
			t.Errorf("could not request: %v", err)
			return
		}
		t.Logf("Download Result: %v", r)
	}

}
