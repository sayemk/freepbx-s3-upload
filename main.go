package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bigkevmcd/go-configparser"
	"io"
	"os"
	"strings"
)

func main() {
	//take input
	fileName := flag.String("file", "", "File to upload")

	flag.Parse()

	if *fileName == "" {
		fmt.Println("Missing file name")

	}

	configPath := os.Getenv("S3_CONFIG_PATH")

	if configPath == "" {
		configPath = "/etc/asterisk/s3_go.conf"
		//configPath = "s3_go.conf"
	}
	configParser, err := configparser.NewConfigParserFromFile(configPath)
	if err != nil {
		fmt.Println("Error parsing config file", err)
	}

	callFileSlice := strings.Split(*fileName, "-")
	fmt.Println("CallDate: ", callFileSlice[3])
	year := callFileSlice[3][:4]
	month := callFileSlice[3][4:6]
	day := callFileSlice[3][6:8]
	//fmt.Println("Year: ", year)
	//fmt.Println("Month: ", month)
	//fmt.Println("Day: ", day)
	s3Key := year + "/" + month + "/" + day + "/" + *fileName
	filePath := "/var/spool/asterisk/monitor/" + s3Key
	//filePath := s3Key

	accessKeyID, _ := configParser.Get("aws", "access_key_id")
	secretAccessKey, _ := configParser.Get("aws", "secret_access_key")
	//sessionToken, _ := config.Get("aws", "session_token")
	s3BucketName, _ := configParser.Get("aws", "s3_bucket_name")
	awsRegion, _ := configParser.Get("aws", "aws_region")

	fmt.Println(accessKeyID, secretAccessKey, s3BucketName, awsRegion)
	staticProvider := credentials.NewStaticCredentials(
		accessKeyID,
		secretAccessKey,
		"",
	)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: staticProvider},
	)
	if err != nil {
		fmt.Errorf("AWS Session : %v", err.Error())
		panic("AWS Session Error")
	}
	fmt.Println("AWS Session Success")

	svc := s3.New(sess)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		panic("Error opening file: " + err.Error())
	}
	defer file.Close()

	// Read the contents of the file into a buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		panic("Error reading file: " + err.Error())
	}

	// This uploads the contents of the buffer to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		panic("Error uploading file: " + err.Error())
	}
	fmt.Printf("UploadFile - filename: %v , Succesful. \n", fileName)

}
