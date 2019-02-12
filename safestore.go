package main

import (
	"crypto/sha256"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	log.Printf("config: %+v", config.AwsBucket)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
		//Credentials: credentials.NewSharedCredentials("", "test-account"),

		Credentials: credentials.NewStaticCredentials(config.AwsClientID, config.AwsClientSecret, ""),
	})
	if err != nil {
		log.Fatalf("failed to setup aws credentials: %s", err)
		os.Exit(1)
	}

	if action.put {
		err := put(sess)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
	}

	if action.get {
		err := get(sess)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
	}
}

func put(sess *session.Session) error {
	file, err := os.Open(action.fileName)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	/*
		iv := make([]byte, IV_SIZE)
		_, err = rand.Read(iv)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}*/

	/*keyAes, _ := hex.DecodeString("6368616e676520746869732070617373")
	keyHmac := keyAes // don't do this

	aes, err := aes.NewCipher(keyAes)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}*/

	password := "test"
	salt := []byte(action.fileName)
	dk := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
	log.Printf("key: %x", dk)

	//ctr := cipher.NewCTR(aes, iv)
	//hmac := hmac.New(sha256.New, keyHmac)

	reader := &CustomReader{
		fp:   file,
		size: fileInfo.Size(),
		/*iv:   iv,
		ctr:  ctr,
		hmac: hmac,*/
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		//u.PartSize = 4096
		u.LeavePartsOnError = true
	})

	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(config.AwsBucket),
		Key:                  aws.String(action.fileName),
		Body:                 reader,
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String(string(dk)),
		//SSECustomerKeyMD5:    aws.String(md5.Sum(dk)),
	})

	if err != nil {
		return err
	}

	log.Printf("Output: %+v", output)
	return nil

}

func get(sess *session.Session) error {

	file, err := os.Create(action.fileName)
	if err != nil {
		return err
	}
	password := "test"
	salt := []byte("a-0einfa09unpc09c34uajh;obhc84bfob;abfhlh89b4fnab03chr8bq23orhQbchrr38hbaBY*CHRW#BOHABw3RPYbhaCOUWGBR32rbc89bb9")
	dk := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
	log.Printf("key: %x", dk)

	downloader := s3manager.NewDownloader(sess, func(u *s3manager.Downloader) {
		u.PartSize = 5 * 1024 * 1024
	})

	output, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket:               aws.String(config.AwsBucket),
		Key:                  aws.String(action.fileName),
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String(string(dk)),
		//SSECustomerKeyMD5:    aws.String(md5.Sum(dk)),
	})

	if err != nil {
		log.Printf("output: %+v err:%+v", output, err)
		return err
	}

	log.Printf("Output: %+v", output)
	return nil

}
