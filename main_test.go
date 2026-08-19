package main

import "testing"

func TestParseNewArgs(t *testing.T) {
	cfg := parseNewArgs([]string{"my-service", "--output", "/tmp/out"})
	if cfg.PkgName != "my-service" {
		t.Fatalf("pkg = %q", cfg.PkgName)
	}
	if cfg.Output != "/tmp/out" {
		t.Fatalf("output = %q", cfg.Output)
	}
}

func TestParseNewArgsWithType(t *testing.T) {
	cfg := parseNewArgs([]string{"worker", "--type=cons", "--output", "/tmp/worker"})
	if cfg.PkgName != "worker" {
		t.Fatalf("pkg = %q", cfg.PkgName)
	}
	if cfg.ServiceType != "cons" {
		t.Fatalf("type = %q", cfg.ServiceType)
	}
	if cfg.Output != "/tmp/worker" {
		t.Fatalf("output = %q", cfg.Output)
	}
}
