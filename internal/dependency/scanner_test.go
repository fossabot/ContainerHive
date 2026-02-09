package dependency

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDockerfileForHiveRefs(t *testing.T) {
	t.Run("finds single __hive__ reference", func(t *testing.T) {
		dir := t.TempDir()
		df := filepath.Join(dir, "Dockerfile")
		os.WriteFile(df, []byte("FROM __hive__/ubuntu:22.04\nRUN echo hello"), 0644)

		refs, err := ScanDockerfileForHiveRefs(df)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(refs) != 1 {
			t.Fatalf("expected 1 ref, got %d", len(refs))
		}
		if refs[0].ImageName != "ubuntu" || refs[0].Tag != "22.04" {
			t.Errorf("expected ubuntu:22.04, got %s:%s", refs[0].ImageName, refs[0].Tag)
		}
	})

	t.Run("finds multiple __hive__ references", func(t *testing.T) {
		dir := t.TempDir()
		df := filepath.Join(dir, "Dockerfile")
		content := "FROM __hive__/ubuntu:22.04 AS base\nFROM __hive__/node:20\nRUN echo hello"
		os.WriteFile(df, []byte(content), 0644)

		refs, err := ScanDockerfileForHiveRefs(df)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(refs) != 2 {
			t.Fatalf("expected 2 refs, got %d", len(refs))
		}
	})

	t.Run("returns empty for no __hive__ references", func(t *testing.T) {
		dir := t.TempDir()
		df := filepath.Join(dir, "Dockerfile")
		os.WriteFile(df, []byte("FROM ubuntu:22.04\nRUN echo hello"), 0644)

		refs, err := ScanDockerfileForHiveRefs(df)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(refs) != 0 {
			t.Errorf("expected 0 refs, got %d", len(refs))
		}
	})

	t.Run("handles FROM with AS alias", func(t *testing.T) {
		dir := t.TempDir()
		df := filepath.Join(dir, "Dockerfile")
		os.WriteFile(df, []byte("FROM __hive__/ubuntu:22.04 AS builder"), 0644)

		refs, err := ScanDockerfileForHiveRefs(df)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(refs) != 1 {
			t.Fatalf("expected 1 ref, got %d", len(refs))
		}
		if refs[0].ImageName != "ubuntu" {
			t.Errorf("expected ubuntu, got %s", refs[0].ImageName)
		}
	})
}

func TestScanRenderedProject(t *testing.T) {
	t.Run("builds graph from rendered dist directory", func(t *testing.T) {
		dir := t.TempDir()

		os.MkdirAll(filepath.Join(dir, "ubuntu", "22.04"), 0755)
		os.WriteFile(filepath.Join(dir, "ubuntu", "22.04", "Dockerfile"), []byte("FROM ubuntu:22.04"), 0644)

		os.MkdirAll(filepath.Join(dir, "python", "3.13"), 0755)
		os.WriteFile(filepath.Join(dir, "python", "3.13", "Dockerfile"), []byte("FROM __hive__/ubuntu:22.04"), 0644)

		graph, err := ScanRenderedProject(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !graph.HasDependencies() {
			t.Error("expected dependencies")
		}

		order, err := graph.TopologicalSort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ubuntuIdx, pythonIdx := -1, -1
		for i, name := range order {
			if name == "ubuntu" {
				ubuntuIdx = i
			}
			if name == "python" {
				pythonIdx = i
			}
		}
		if ubuntuIdx > pythonIdx {
			t.Errorf("ubuntu (idx=%d) must come before python (idx=%d)", ubuntuIdx, pythonIdx)
		}
	})

	t.Run("handles project with no __hive__ references", func(t *testing.T) {
		dir := t.TempDir()

		os.MkdirAll(filepath.Join(dir, "nginx", "1.27"), 0755)
		os.WriteFile(filepath.Join(dir, "nginx", "1.27", "Dockerfile"), []byte("FROM nginx:alpine"), 0644)

		graph, err := ScanRenderedProject(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if graph.HasDependencies() {
			t.Error("expected no dependencies")
		}
	})
}
