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
	privateKeyFile  *string
	awsRegion       *string
	awsClientID     *string
	awsClientSecret *string
	file            *string
	bucket          *string
}

var (
	privateKeyFile = "safestore.priv"
)

var config Config

func readConfig() bool {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/.safestore", homedir()))
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return true
	}
	return false
}

func writeConfig() {
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Could not converd config to json: %s", err)
		os.Exit(255)
	}

	ioutil.WriteFile(*config.file, data, 0600)
}

func init() {
	configFound := readConfig()

	config.privateKeyFile = flag.String("privatekey", fmt.Sprintf("%s/safestore.priv", homedir()), "path to the primary key (if it doesn't exist it will be generated)")
	config.bucket = flag.String("bucket", "", "name of the bucket to store the file")
	config.file = flag.String("file", fmt.Sprintf("%s/.safestore", homedir()), "default location of the config file")

	if *config.bucket == "" {
		log.Fatal("Bucket must be supplied")
		os.Exit(1)
	}

	flag.Parse()

	if !configFound {
		writeConfig()
	}

	/*
			sess, err := session.NewSession(&aws.Config{
		    Region:      aws.String("us-west-2"),
		    Credentials: credentials.NewSharedCredentials("", "test-account"),
		})
	*/

}

func homedir() string {
	return os.Getenv("HOME")
}
