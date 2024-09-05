package cmd

import (
	"github.com/geetchoubey/awsls/pkg/awsutil"
	"github.com/geetchoubey/awsls/pkg/config"
	"github.com/geetchoubey/awsls/pkg/types"
	"github.com/geetchoubey/awsls/resources"
)

type Scanner struct {
	Parameters    ScannerParameters
	Account       awsutil.Account
	Config        *config.Scanner
	ResourceTypes types.Collection

	items Queue
}

func NewScanner(account awsutil.Account) *Scanner {
	scanner := Scanner{
		Account: account,
	}

	return &scanner
}

func (s *Scanner) Run() error {
	if err := s.Scan(); err != nil {
		return err
	}
	return nil
}

func (s *Scanner) Scan() error {
	// accountConfig := s.Config.Accounts[s.Account.ID()]

	resourceTypes := resources.GetListerNames()

	queue := make(Queue, 0)

	for _, regionName := range s.Config.Regions {
		region := NewRegion(regionName, s.Account.ResourceTypeToServiceType, s.Account.NewSession)

		items := Scan(region, resourceTypes)

		for item := range items {
			queue = append(queue, item)
		}
	}
	s.items = queue
	return nil
}
