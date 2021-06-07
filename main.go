package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
)

var InternalServerError = events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}

type ResizableImage struct{
	Source string
	Length int64
	Height int64
}

func extractImageFromSource(sourceFile string) (image.Image, error) {
	//open file
	file, err := os.Open(sourceFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//get file size
	fileConfig, err := file.Stat()
	if err != nil {
		return nil, err
	}

	//get byte slice for file
	data := make([]byte, fileConfig.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	//decode imageJPEG from byte slice
	image, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return image, nil
}

func resizeImage(sourceFile string, dx uint, dy uint) (image.Image, error) {
	imageJPEG, err := extractImageFromSource(sourceFile)
	if err != nil {
		return nil, err
	}
	//change imageJPEG size
	resizedImage := resize.Resize(dx, dy, imageJPEG, resize.Lanczos3)

	return resizedImage, nil
}

func resizeHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody ResizableImage
	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		return InternalServerError, err
	}

	image, err := resizeImage(reqBody.Source, uint(reqBody.Length), uint(reqBody.Height))
	if err != nil {
		return InternalServerError, err
	}

	//TODO: write to S3 bucket
}


func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "POST" {
		switch request.Resource {
		case "resize":
			//return resizeHandler(request)
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
	req := events.APIGatewayProxyRequest{Body:`{ "source":"/home/karadzhovai/Downloads/newnew.jpg" , "length":500 , "height":500 }`}
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
	fmt.Println( imageJPEG.Bounds())

	//change imageJPEG size
	smallerImage := resize.Resize(uint(imageJPEG.Bounds().Dx()/4), uint(imageJPEG.Bounds().Dy()/4), imageJPEG, resize.Lanczos3)

	//change file type (if needed)
	pngFile, _ := os.Create("/home/karadzhovai/Downloads/newnew.jpg")
	defer pngFile.Close()


	jpeg.Encode(pngFile, smallerImage, nil)


}



