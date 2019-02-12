package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Config holds the configuration variables
type Config struct {
	PrivateKeyFile  string
	AwsRegion       string
	AwsClientID     string
	AwsClientSecret string
	AwsBucket       string
}

type Action struct {
	fileName string
	put      bool
	get      bool
}

var (
	privateKeyFile = "safestore.priv"
)

var config Config
var action Action

func readConfig(file string) bool {
	//config = Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Failed to load config: %s", file)
		return false
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Failed to unmarshal config: %s", err)
		return true
	}
	return false
}

func writeConfig(file string) {
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Could not converd config to json: %s", err)
		os.Exit(255)
	}

	ioutil.WriteFile(file, data, 0600)
}

func init() {
	configFile := flag.String("configfile", fmt.Sprintf("%s/.safestore", homedir()), "default location of the config file")
	privateKeyFile := flag.String("privatekey", fmt.Sprintf("%s/safestore.priv", homedir()), "path to the primary key (if it doesn't exist it will be generated)")
	awsBucket := flag.String("bucket", "", "name of the bucket to store the file")
	awsClientID := flag.String("clientid", "", "AWS client ID")
	awsClientSecret := flag.String("clientsecret", "", "AWS client Secret")
	awsRegion := flag.String("region", "eu-west-1", "AWS region")
	put := flag.Bool("put", false, "Put file in storage")
	get := flag.Bool("get", false, "Get file from storage")
	fileName := flag.String("file", "", "File to put/get")

	flag.Parse()

	configFound := readConfig(*configFile)
	if *awsBucket != "" {
		config.AwsBucket = *awsBucket
	}
	if *awsClientID != "" {
		config.AwsClientID = *awsClientID
	}
	if *awsClientSecret != "" {
		config.AwsClientSecret = *awsClientSecret
	}
	if *awsBucket != "" {
		config.AwsBucket = *awsBucket
	}
	if *awsRegion != "" {
		config.AwsRegion = *awsRegion
	}
	if *privateKeyFile != "" {
		config.PrivateKeyFile = *privateKeyFile
	}
	if *fileName != "" {
		action.fileName = *fileName
	}
	if *put != false {
		action.put = true
	}
	if *get != false {
		action.get = true
	}

	if config.AwsBucket == "" {
		log.Fatal("bucket must be supplied")
		os.Exit(1)
	}

	if config.AwsClientID == "" {
		log.Fatal("clientid must be supplied")
		os.Exit(1)
	}

	if config.AwsClientSecret == "" {
		log.Fatal("bucket must be supplied")
		os.Exit(1)
	}

	if action.fileName == "" {
		log.Fatal("specify the file to put or get")
		os.Exit(1)
	}

	if action.put == action.get {
		log.Fatal("specify one of -get or -put to receive or place a file in the bucket")
		os.Exit(1)
	}

	if !configFound {
		writeConfig(*configFile)
	}

}

func homedir() string {
	return os.Getenv("HOME")
}
