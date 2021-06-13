package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"os"
)

func resizeHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody ResizableImage
	//get request body
	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		return InternalServerErrorResponse, err
	}

	image, err := extractImageFromSource(reqBody.Source)
	if err != nil {
		return InternalServerErrorResponse, err
	}

	// image resizing
	resizedImage := resizeImage(image, uint(reqBody.Length), uint(reqBody.Height))

	// encode into right format
	tempFile := "temp_file" + "." + reqBody.ImgType
	encodeImage(resizedImage, tempFile, reqBody.ImgType)

	err = os.Remove(tempFile)
	if err != nil {
		return InternalServerErrorResponse, err
	}

	return events.APIGatewayProxyResponse{Body: "Successful resize", StatusCode: 200}, nil

}
