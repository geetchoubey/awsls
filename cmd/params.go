package cmd

import (
	"fmt"
	"strings"
)

type ScannerParameters struct {
	ConfigPath string
}

func (p *ScannerParameters) Validate() error {
	if strings.TrimSpace(p.ConfigPath) == "" {
		return fmt.Errorf("You have to specify the --config flag.")
	}

	return nil
}
