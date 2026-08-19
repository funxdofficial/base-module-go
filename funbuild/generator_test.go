package funbuild

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func execCommand(t *testing.T, dir, name string, args ...string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	return cmd
}

func TestValidatePkgName(t *testing.T) {
	tests := []struct {
		name    string
		pkg     string
		wantErr bool
	}{
		{"valid", "my-service", false},
		{"valid short", "abc", false},
		{"too short", "ab", true},
		{"uppercase", "My-Service", true},
		{"leading hyphen", "-service", true},
		{"trailing hyphen", "service-", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePkgName(tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validatePkgName(%q) error = %v, wantErr %v", tt.pkg, err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeServiceType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"default empty", "", ServiceTypeREST, false},
		{"rest", "rest", ServiceTypeREST, false},
		{"consumer", "consumer", ServiceTypeConsumer, false},
		{"cons alias", "cons", ServiceTypeConsumer, false},
		{"invalid", "grpc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeServiceType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("normalizeServiceType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("normalizeServiceType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	vars := templateVars{
		PkgName:          "demo-svc",
		ModulePath:       "github.com/funxdofficial/demo-svc",
		ReadmeTitle:      "Demo Svc",
		SonarProjectName: "Demo - Svc",
		ServiceType:      ServiceTypeREST,
	}

	out := renderTemplate("module {{MODULE_PATH}} // {{PKG_NAME}} {{README_TITLE}} {{SERVICE_TYPE}}", vars)
	if !strings.Contains(out, "github.com/funxdofficial/demo-svc") {
		t.Fatalf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "Demo Svc") {
		t.Fatalf("missing readme title: %s", out)
	}
	if !strings.Contains(out, "rest") {
		t.Fatalf("missing service type: %s", out)
	}
}

func TestGenerateRESTService(t *testing.T) {
	outDir := t.TempDir()
	cfg := Config{
		PkgName:     "test-service",
		Output:      outDir,
		ServiceType: ServiceTypeREST,
	}
	vars := templateVars{
		PkgName:          cfg.PkgName,
		ModulePath:       modulePrefix + "/" + cfg.PkgName,
		ReadmeTitle:      "Test Service",
		SonarProjectName: "Test - Service",
		ServiceType:      ServiceTypeREST,
	}

	if err := generateService(cfg, vars); err != nil {
		t.Fatalf("generateService: %v", err)
	}

	checks := []string{
		"main.go",
		"go.mod",
		"Makefile",
		".env.example",
		"config/config.go",
		"container/container.go",
		"module/module.go",
		"model/item.go",
		"model/request/item.go",
		"model/response/item.go",
		"service/service.go",
		"service/middleware.go",
		"route/route.go",
		"route/item.go",
		"service/controller/controller.go",
		"service/controller/controller_item.go",
		"service/controller/controller_item_test.go",
		"service/usecase/usecase.go",
		"service/usecase/usecase_item.go",
		"service/repository/repository.go",
		"service/repository/repository_item.go",
		"pkg/helpers/datetime.go",
	}

	absent := []string{
		"service/scheduler.go",
	}

	for _, rel := range checks {
		path := filepath.Join(outDir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected file %s: %v", rel, err)
		}
	}

	for _, rel := range absent {
		path := filepath.Join(outDir, rel)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("unexpected file %s", rel)
		}
	}

	mainGo, err := os.ReadFile(filepath.Join(outDir, "main.go"))
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	if !strings.Contains(string(mainGo), "github.com/funxdofficial/test-service/config") {
		t.Fatalf("main.go module path not substituted: %s", mainGo)
	}

	goMod, err := os.ReadFile(filepath.Join(outDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if strings.Contains(string(goMod), "golang-module-scheduler") {
		t.Fatalf("rest go.mod should not include scheduler: %s", goMod)
	}
}

func TestGenerateConsumerService(t *testing.T) {
	outDir := t.TempDir()
	cfg := Config{
		PkgName:     "worker-service",
		Output:      outDir,
		ServiceType: ServiceTypeConsumer,
	}
	vars := templateVars{
		PkgName:          cfg.PkgName,
		ModulePath:       modulePrefix + "/" + cfg.PkgName,
		ReadmeTitle:      "Worker Service",
		SonarProjectName: "Worker - Service",
		ServiceType:      ServiceTypeConsumer,
	}

	if err := generateService(cfg, vars); err != nil {
		t.Fatalf("generateService: %v", err)
	}

	checks := []string{
		"main.go",
		"go.mod",
		"config/config.go",
		"service/service.go",
		"service/scheduler.go",
		"service/usecase/usecase.go",
		"service/usecase/usecase_item.go",
		"service/repository/repository.go",
		"service/repository/repository_item.go",
	}

	absent := []string{
		"route/route.go",
		"service/middleware.go",
	}

	for _, rel := range checks {
		path := filepath.Join(outDir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected file %s: %v", rel, err)
		}
	}

	for _, rel := range absent {
		path := filepath.Join(outDir, rel)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("unexpected file %s", rel)
		}
	}

	goMod, err := os.ReadFile(filepath.Join(outDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "golang-module-scheduler") {
		t.Fatalf("consumer go.mod should include scheduler: %s", goMod)
	}
	if strings.Contains(string(goMod), "labstack/echo") {
		t.Fatalf("consumer go.mod should not include echo: %s", goMod)
	}
}

func TestGenerateServiceBuild(t *testing.T) {
	if os.Getenv("BASE_MODULE_GO_INTEGRATION") != "1" {
		t.Skip("set BASE_MODULE_GO_INTEGRATION=1 to run build integration test")
	}

	for _, serviceType := range []string{ServiceTypeREST, ServiceTypeConsumer} {
		t.Run(serviceType, func(t *testing.T) {
			outDir := t.TempDir()
			pkgName := "integration-" + serviceType
			cfg := Config{PkgName: pkgName, Output: outDir, ServiceType: serviceType}
			if err := Create(cfg); err != nil {
				t.Fatalf("Create: %v", err)
			}

			schedulerPath, _ := filepath.Abs("../../golang-module-scheduler")
			syslogPath, _ := filepath.Abs("../../golang-module-syslog")
			replaceBlock := "\nreplace (\n\tgithub.com/funxdofficial/golang-module-scheduler => " + schedulerPath + "\n\tgithub.com/funxdofficial/golang-module-syslog => " + syslogPath + "\n)\n"

			goModPath := filepath.Join(outDir, "go.mod")
			data, err := os.ReadFile(goModPath)
			if err != nil {
				t.Fatalf("read go.mod: %v", err)
			}
			if err := os.WriteFile(goModPath, append(data, []byte(replaceBlock)...), 0o644); err != nil {
				t.Fatalf("write go.mod replace: %v", err)
			}

			run := func(name string, args ...string) {
				t.Helper()
				cmd := execCommand(t, outDir, name, args...)
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
				}
			}

			run("go", "mod", "tidy")
			run("go", "test", "./...")
			run("go", "build", "-o", filepath.Join(outDir, "bin", "app"), ".")
		})
	}
}
