package detector

import (
	"errors"
	"os"
	"testing"
	"time"

	"lobster/internal/platform"
	"lobster/internal/products"
)

type fakeProduct struct {
	detectPlan products.DetectPlan
}

func (f fakeProduct) Key() string { return "fake" }

func (f fakeProduct) DisplayName() string { return "Fake" }

func (f fakeProduct) Summary() string { return "fake summary" }

func (f fakeProduct) InstallPlan(platform.Info) (products.InstallPlan, error) {
	return products.InstallPlan{}, nil
}

func (f fakeProduct) DetectPlan(platform.Info) products.DetectPlan {
	return f.detectPlan
}

func (f fakeProduct) LaunchPlan(platform.Info) products.LaunchPlan {
	return products.LaunchPlan{}
}

func TestCheckWhenCommandAvailable(t *testing.T) {
	originalLookPath := lookPath
	originalStatPath := statPath
	t.Cleanup(func() {
		lookPath = originalLookPath
		statPath = originalStatPath
	})

	lookPath = func(file string) (string, error) {
		if file == "codebuddy" {
			return "/tmp/codebuddy", nil
		}
		return "", errors.New("not found")
	}
	statPath = func(string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	status := Check(platform.Info{OS: platform.Darwin, Arch: "arm64"}, fakeProduct{
		detectPlan: products.DetectPlan{
			Commands: []string{"codebuddy", "workbuddy"},
			Paths:    []string{"/tmp/workbuddy"},
			Notes:    []string{"note"},
		},
	})

	if !status.Installed {
		t.Fatalf("命令存在时应视为已安装")
	}
	if !status.CommandAvailable {
		t.Fatalf("命令存在时应标记为可用")
	}
	if status.CommandPath != "/tmp/codebuddy" {
		t.Fatalf("命令路径不符合预期，实际：%s", status.CommandPath)
	}
	if status.MatchedCommand != "codebuddy" {
		t.Fatalf("命中的命令名不符合预期，实际：%s", status.MatchedCommand)
	}
}

func TestCheckWhenOnlyPathEvidenceExists(t *testing.T) {
	originalLookPath := lookPath
	originalStatPath := statPath
	t.Cleanup(func() {
		lookPath = originalLookPath
		statPath = originalStatPath
	})

	lookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}
	statPath = func(name string) (os.FileInfo, error) {
		if name == "/tmp/workbuddy" {
			return fakeFileInfo{name: "workbuddy"}, nil
		}
		return nil, os.ErrNotExist
	}

	status := Check(platform.Info{OS: platform.Darwin, Arch: "arm64"}, fakeProduct{
		detectPlan: products.DetectPlan{
			Commands: []string{"codebuddy"},
			Paths:    []string{"/tmp/workbuddy"},
			Notes:    []string{"note"},
		},
	})

	if !status.Installed {
		t.Fatalf("存在路径痕迹时应视为已检测到安装痕迹")
	}
	if status.CommandAvailable {
		t.Fatalf("只有路径痕迹时不应标记命令可用")
	}
	if !status.HasPathEvidence {
		t.Fatalf("存在路径痕迹时应标记 HasPathEvidence")
	}
	if len(status.FoundPaths) != 1 || status.FoundPaths[0] != "/tmp/workbuddy" {
		t.Fatalf("命中路径不符合预期，实际：%v", status.FoundPaths)
	}
	if len(status.Warnings) == 0 {
		t.Fatalf("只有路径痕迹时应给出 warning")
	}
}

func TestCheckWhenNothingDetected(t *testing.T) {
	originalLookPath := lookPath
	originalStatPath := statPath
	t.Cleanup(func() {
		lookPath = originalLookPath
		statPath = originalStatPath
	})

	lookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}
	statPath = func(string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	status := Check(platform.Info{OS: platform.Darwin, Arch: "arm64"}, fakeProduct{
		detectPlan: products.DetectPlan{
			Commands: []string{"codebuddy"},
			Paths:    []string{"/tmp/workbuddy"},
			Notes:    []string{"note"},
		},
	})

	if status.Installed {
		t.Fatalf("完全未命中时不应标记已安装")
	}
	if len(status.Warnings) == 0 {
		t.Fatalf("完全未命中时应给出 warning")
	}
	if status.Notes[0] != "note" {
		t.Fatalf("notes 应保留 detect plan 内容，实际：%v", status.Notes)
	}
}

type fakeFileInfo struct {
	name string
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() os.FileMode  { return 0 }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }
