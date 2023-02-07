package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/letsgodevops/ssm-sync/ssm"
	"github.com/letsgodevops/ssm-sync/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/sirupsen/logrus"
)

func checkDefaults() {
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			fmt.Printf("'-%s' not set\n", f.Name)
			os.Exit(1)
		}
	})
}

func usage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func parameterFound(err error) (bool, error) {
	if err == nil {
		return true, nil
	} else {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Printf("Error: '%s'\n", awsErr.Code())
			if awsErr.Code() == "ParameterNotFound" {
				logrus.Debug("Parameter not found in SSM")
				return false, nil
			}
		}
	}
	return false, err
}

func main() {
	logrus.Infof("Starting ssm-sync, version: %s", version)

	region := flag.String("region", "eu-west-1", "AWS Region")
	ssmKeyName := flag.String("ssm", "", "SSM Key name")
	filePath := flag.String("file", "", "File path on disk")
	env := flag.String("env", "", "Environment name")

	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Usage = usage
	flag.Parse()

	checkDefaults()

	logrus.SetLevel(logrus.WarnLevel)
	if *verbose {
		logrus.SetLevel(logrus.InfoLevel)
	}

	ssmClient, err := ssm.New(*region, aws.Config{})
	if err != nil {
		logrus.WithError(err).Errorln("Failed to obtain AWS:SSM client")
		os.Exit(1)
	}

	ssmValue, err := ssmClient.GetObject(&types.GetObjectInput{
		Key: *ssmKeyName,
	})
	ssmParamFound, err := parameterFound(err)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to read parameter '%s'", *ssmKeyName)
		os.Exit(1)
	}

	if ssmParamFound {
		// Overwrite file on disk
		err := os.WriteFile(*filePath, []byte(ssmValue.Value), 0400)
		if err != nil {
			logrus.WithError(err).Errorf("Error while writing to file %s", *filePath)
			os.Exit(1)
		}
		logrus.Info("Downloading of file content finished successfully")
	} else {
		// Upload file content
		data, err := os.ReadFile(*filePath)
		if err != nil {
			logrus.WithError(err).Errorf("Error while reading the file %s", *filePath)
			os.Exit(1)
		}

		err = ssmClient.PutObject(&types.PutObjectInput{
			Key:         *ssmKeyName,
			Value:       string(data),
			Application: "file-sync",
			Environment: *env,
		})
		if err != nil {
			logrus.WithError(err).Errorf("Error while uploading file to ssm %s", *ssmKeyName)
			os.Exit(1)
		}
		logrus.Info("File uploaded successfully")
	}
}
