package config

import (
	"bytes"
	"os"

	"github.com/geetchoubey/awsls/pkg/types"
	"gopkg.in/yaml.v3"
)

type ResourceTypes struct {
	Targets      types.Collection `yaml:"targets"`
	Excludes     types.Collection `yaml:"excludes"`
	CloudControl types.Collection `yaml:"cloud-control"`
}

type Account struct {
	// Filters       Filters       `yaml:"filters"`
	ResourceTypes ResourceTypes `yaml:"resource-types"`
	Presets       []string      `yaml:"presets"`
}

type Scanner struct {
	Regions         []string           `yaml:"regions"`
	Accounts        map[string]Account `yaml:"accounts"`
	ResourceTypes   ResourceTypes      `yaml:"resource-types"`
	CustomEndpoints CustomEndpoints    `yaml:"endpoints"`
}

type CustomService struct {
	Service               string `yaml:"service"`
	URL                   string `yaml:"url"`
	TLSInsecureSkipVerify bool   `yaml:"tls_insecure_skip_verify"`
}

type CustomServices []*CustomService

type CustomRegion struct {
	Region                string         `yaml:"region"`
	Services              CustomServices `yaml:"services"`
	TLSInsecureSkipVerify bool           `yaml:"tls_insecure_skip_verify"`
}

type CustomEndpoints []*CustomRegion

// GetRegion returns the custom region or nil when no such custom endpoints are defined for this region
func (endpoints CustomEndpoints) GetRegion(region string) *CustomRegion {
	for _, r := range endpoints {
		if r.Region == region {
			if r.TLSInsecureSkipVerify {
				for _, s := range r.Services {
					s.TLSInsecureSkipVerify = r.TLSInsecureSkipVerify
				}
			}
			return r
		}
	}
	return nil
}

// GetService returns the custom region or nil when no such custom endpoints are defined for this region
func (services CustomServices) GetService(serviceType string) *CustomService {
	for _, s := range services {
		if serviceType == s.Service {
			return s
		}
	}
	return nil
}

func (endpoints CustomEndpoints) GetURL(region, serviceType string) string {
	r := endpoints.GetRegion(region)
	if r == nil {
		return ""
	}
	s := r.Services.GetService(serviceType)
	if s == nil {
		return ""
	}
	return s.URL
}

func Load(path string) (*Scanner, error) {
	var err error

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(Scanner)
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
