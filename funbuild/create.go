package funbuild

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.Und)

// Create generates a service from templates using the given config.
func Create(cfg Config) error {
	if err := validatePkgName(cfg.PkgName); err != nil {
		return err
	}

	serviceType, err := normalizeServiceType(cfg.ServiceType)
	if err != nil {
		return err
	}
	cfg.ServiceType = serviceType

	readmeTitle := titleCaser.String(strings.ReplaceAll(cfg.PkgName, "-", " "))
	vars := templateVars{
		PkgName:          cfg.PkgName,
		ModulePath:       modulePrefix + "/" + cfg.PkgName,
		ReadmeTitle:      readmeTitle,
		SonarProjectName: strings.ReplaceAll(readmeTitle, " ", " - "),
		ServiceType:      serviceType,
	}

	if err := generateService(cfg, vars); err != nil {
		return err
	}

	outDir := cfg.Output
	if outDir == "" {
		outDir = "./" + cfg.PkgName
	}
	fmt.Printf("generated %s service at %s\n", serviceType, outDir)
	return nil
}
