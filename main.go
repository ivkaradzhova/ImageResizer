package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

var (
	InternalServerErrorResponse        = events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}
	resizedImagesBucket         string = "ResizedImages"
	originalImagesBucket        string = "OriginalImages"
)

type ResizableImage struct {
	Source  string
	Length  int64
	Height  int64
	ImgType string
}

func uploadToS3(filename string, bucket string) error {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(myString),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "POST" {
		switch request.Resource {
		case "resize":
			return resizeHandler(request)
		default:
			return events.APIGatewayProxyResponse{Body: "Unsupported request", StatusCode: 400}, errors.New("unsupported request")
		}

	} else {
		return events.APIGatewayProxyResponse{Body: "Unsupported HTTP method", StatusCode: 400}, errors.New("unsupported HTTP method")
	}
	return events.APIGatewayProxyResponse{Body: "Unsupported HTTP method", StatusCode: 405}, errors.New("unsupported HTTP method")

}

const file = "/home/karadzhovai/Downloads/trail-l4MwmH8QIxk-unsplash.jpg"

func main() {
	req := events.APIGatewayProxyRequest{Body: `{ "source":"/home/karadzhovai/Downloads/newnew.jpg" , "length":500 , "height":500 }`}
	reqData := ResizableImage{}
	err := json.Unmarshal([]byte(req.Body), &reqData)
	fmt.Println("hehjsfhkse", reqData.Source, err)
	//open file
	file, _ := os.Open(file)
	defer file.Close()

	//get file size
	f, _ := file.Stat()
	fmt.Println(f.Size())

	//get byte slice for file
	data := make([]byte, f.Size())
	file.Read(data)
	//fmt.Println(count, data)

	//decode imageJPEG from byte slice
	imageJPEG, _ := jpeg.Decode(bytes.NewReader(data))
	fmt.Println(imageJPEG.Bounds())

	//change imageJPEG size
	smallerImage := resize.Resize(uint(imageJPEG.Bounds().Dx()/4), uint(imageJPEG.Bounds().Dy()/4), imageJPEG, resize.Lanczos3)

	//change file type (if needed)
	pngFile, _ := os.Create("/home/karadzhovai/Downloads/newnew.jpg")
	defer pngFile.Close()

	jpeg.Encode(pngFile, smallerImage, nil)

}
