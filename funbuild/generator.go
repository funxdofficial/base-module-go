package funbuild

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const modulePrefix = "github.com/funxdofficial"

var pkgNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)

type templateVars struct {
	PkgName          string
	ModulePath       string
	ReadmeTitle      string
	SonarProjectName string
	ServiceType      string
}

func validatePkgName(name string) error {
	if len(name) < 3 {
		return fmt.Errorf("pkg name must be at least 3 characters")
	}
	if !pkgNamePattern.MatchString(name) {
		return fmt.Errorf("pkg name must be lowercase alphanumeric with hyphens (e.g. my-service)")
	}
	return nil
}

func renderTemplate(content string, vars templateVars) string {
	replacer := strings.NewReplacer(
		"{{PKG_NAME}}", vars.PkgName,
		"{{MODULE_PATH}}", vars.ModulePath,
		"{{README_TITLE}}", vars.ReadmeTitle,
		"{{SONAR_PROJECT_NAME}}", vars.SonarProjectName,
		"{{SERVICE_TYPE}}", vars.ServiceType,
	)
	return replacer.Replace(content)
}

func templateRoot(serviceType string) string {
	return filepath.Join("tpls", serviceType)
}

func outputPath(root, embedPath, serviceType string) string {
	prefix := templateRoot(serviceType) + "/"
	rel := strings.TrimPrefix(embedPath, prefix)
	if rel == embedPath {
		rel = strings.TrimPrefix(embedPath, "tpls/")
	}
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" {
		return root
	}
	if strings.HasSuffix(rel, ".txt") {
		rel = strings.TrimSuffix(rel, ".txt")
	}
	return filepath.Join(root, rel)
}

func generateService(cfg Config, vars templateVars) error {
	outDir := cfg.Output
	if outDir == "" {
		outDir = filepath.Join(".", cfg.PkgName)
	}

	if info, err := os.Stat(outDir); err == nil && info.IsDir() {
		if err := os.RemoveAll(outDir); err != nil {
			return fmt.Errorf("remove existing output dir: %w", err)
		}
	}

	root := templateRoot(cfg.ServiceType)
	return fs.WalkDir(tpls, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".txt") {
			return nil
		}

		data, err := tpls.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read template %s: %w", path, err)
		}

		target := outputPath(outDir, path, cfg.ServiceType)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("mkdir for %s: %w", target, err)
		}

		rendered := renderTemplate(string(data), vars)
		if err := os.WriteFile(target, []byte(rendered), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", target, err)
		}
		return nil
	})
}
