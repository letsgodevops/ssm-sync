package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/letsgodevops/ssm-sync/ssm"
	"github.com/letsgodevops/ssm-sync/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"golang.org/x/exp/slog"
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
				slog.Debug("Parameter not found in SSM")
				return false, nil
			}
		}
	}
	return false, err
}

func main() {
	slog.Info(fmt.Sprintf("Starting ssm-sync, version: %s", version))

	region := flag.String("region", "eu-west-1", "AWS Region")
	ssmKeyName := flag.String("ssm", "", "SSM Key name")
	filePath := flag.String("file", "", "File path on disk")
	kmsKeyAlias := flag.String("kmsKeyAlias", "", "KMS key alias name (would use alias/<name> in KMS)")

	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Usage = usage
	flag.Parse()

	checkDefaults()

	ssmClient, err := ssm.New(*region, aws.Config{})
	if err != nil {
		slog.Error("Failed to obtain AWS:SSM client", err)
		os.Exit(1)
	}

	ssmValue, err := ssmClient.GetObject(&types.GetObjectInput{
		Key: *ssmKeyName,
	})
	ssmParamFound, err := parameterFound(err)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to read parameter '%s'", *ssmKeyName), err)
		os.Exit(1)
	}

	if ssmParamFound {
		// Overwrite file on disk
		err := os.WriteFile(*filePath, []byte(ssmValue.Value), 0400)
		if err != nil {
			slog.Error(fmt.Sprintf("Error while writing to file %s", *filePath), err)
			os.Exit(1)
		}
		if *verbose {
			slog.Info("Downloading of file content finished successfully")
		}
	} else {
		// Upload file content
		data, err := os.ReadFile(*filePath)
		if err != nil {
			slog.Error(fmt.Sprintf("Error while reading the file %s", *filePath), err)
			os.Exit(1)
		}

		err = ssmClient.PutObject(&types.PutObjectInput{
			Key:         *ssmKeyName,
			Value:       string(data),
			Application: "file-sync",
			KmsKeyAlias: *kmsKeyAlias,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("Error while uploading file to ssm %s", *ssmKeyName), err)
			os.Exit(1)
		}
		if *verbose {
			slog.Info("File uploaded successfully")
		}
	}
}
