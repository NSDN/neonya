package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifySpecificationAcceptsRenderedDocument(t *testing.T) {
	document, err := renderSpecification()
	if err != nil {
		t.Fatalf("render specification: %v", err)
	}

	path := filepath.Join(t.TempDir(), "openapi.yaml")
	if err := os.WriteFile(path, document, 0o644); err != nil {
		t.Fatalf("write specification: %v", err)
	}

	if err := verifySpecification(path, document); err != nil {
		t.Fatalf("identical document rejected: %v", err)
	}
}

func TestVerifySpecificationIgnoresCarriageReturns(t *testing.T) {
	document, err := renderSpecification()
	if err != nil {
		t.Fatalf("render specification: %v", err)
	}

	path := filepath.Join(t.TempDir(), "openapi.yaml")
	windowsDocument := []byte(strings.ReplaceAll(string(document), "\n", "\r\n"))
	if err := os.WriteFile(path, windowsDocument, 0o644); err != nil {
		t.Fatalf("write specification: %v", err)
	}

	if err := verifySpecification(path, document); err != nil {
		t.Fatalf("CRLF document rejected: %v", err)
	}
}

func TestVerifySpecificationRejectsDriftedDocument(t *testing.T) {
	document, err := renderSpecification()
	if err != nil {
		t.Fatalf("render specification: %v", err)
	}

	drifted := strings.Replace(string(document), "version: 1.0.0", "version: 9.9.9", 1)
	if drifted == string(document) {
		t.Fatal("expected the rendered document to contain the contract version")
	}

	path := filepath.Join(t.TempDir(), "openapi.yaml")
	if err := os.WriteFile(path, []byte(drifted), 0o644); err != nil {
		t.Fatalf("write specification: %v", err)
	}

	err = verifySpecification(path, document)
	if err == nil {
		t.Fatal("drifted document accepted")
	}
	if !strings.Contains(err.Error(), "task gen") {
		t.Fatalf("error %q should point at task gen", err)
	}
}

func TestVerifySpecificationRejectsMissingDocument(t *testing.T) {
	document, err := renderSpecification()
	if err != nil {
		t.Fatalf("render specification: %v", err)
	}

	path := filepath.Join(t.TempDir(), "absent.yaml")
	if err := verifySpecification(path, document); err == nil {
		t.Fatal("missing document accepted")
	}
}
