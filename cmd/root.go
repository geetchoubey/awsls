package cmd

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/geetchoubey/awsls/pkg/awsutil"
	"github.com/geetchoubey/awsls/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var (
		params        ScannerParameters
		defaultRegion string
		verbose       bool
		creds         awsutil.Credentials
	)

	command := &cobra.Command{
		Use:   "awsls",
		Short: "awsls lists all resources from AWS",
		Long:  "A tool that lists every resources from an AWS Account.",
	}

	command.PreRun = func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		log.SetFormatter(&log.TextFormatter{
			EnvironmentOverrideColors: true,
		})

	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		var err error

		err = params.Validate()
		if err != nil {
			return err
		}

		if !creds.HasKeys() && !creds.HasProfile() && defaultRegion != "" {
			creds.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
			creds.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		}
		err = creds.Validate()
		if err != nil {
			return err
		}

		config, err := config.Load(params.ConfigPath)
		if err != nil {
			log.Errorf("Failed to parse config file %s", params.ConfigPath)
			return err
		}

		if defaultRegion != "" {
			awsutil.DefaultRegionID = defaultRegion
			switch defaultRegion {
			case endpoints.UsEast1RegionID, endpoints.UsEast2RegionID, endpoints.UsWest1RegionID, endpoints.UsWest2RegionID:
				awsutil.DefaultAWSPartitionID = endpoints.AwsPartitionID
			case endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
				awsutil.DefaultAWSPartitionID = endpoints.AwsUsGovPartitionID
			case endpoints.CnNorth1RegionID, endpoints.CnNorthwest1RegionID:
				awsutil.DefaultAWSPartitionID = endpoints.AwsCnPartitionID
			default:

			}
		}

		account, err := awsutil.NewAccount(creds)
		if err != nil {
			return err
		}
		s := NewScanner(*account)

		s.Config = config

		return s.Run()
	}

	command.PersistentFlags().StringVar(&defaultRegion, "default-region", "", "Custom default region name")
	command.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enables debug output")
	command.PersistentFlags().StringVarP(
		&params.ConfigPath, "config", "c", "",
		"(required) Path to the nuke config file.")
	return command
}
