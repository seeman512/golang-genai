package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		want         Config
		wantErr      bool
		wantNotExist bool
	}{
		{
			name: "success",
			content: `{
  "server_port": 8080,
  "environment": "production",
  "database_url": "postgres://localhost/app",
  "debug_mode": true,
  "admin_password": "must not be loaded"
}`,
			want: Config{
				ServerPort:  8080,
				Environment: "production",
				DatabaseURL: "postgres://localhost/app",
				DebugMode:   true,
			},
		},
		{
			name:         "file not found",
			wantErr:      true,
			wantNotExist: true,
		},
		{
			name:    "corrupted json",
			content: `{"server_port": 8080`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if tc.content != "" {
				if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
					t.Fatalf("failed to create test configuration: %v", err)
				}
			}

			got, err := LoadConfig(path)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("LoadConfig(%q) returned nil error, want error", path)
				}
				if tc.wantNotExist && !errors.Is(err, os.ErrNotExist) {
					t.Errorf("LoadConfig(%q) error = %v, want an error matching os.ErrNotExist", path, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("LoadConfig(%q) returned unexpected error: %v", path, err)
			}
			if got != tc.want {
				t.Errorf("LoadConfig(%q) = %+v, want %+v", path, got, tc.want)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{
			name: "all fields except admin password",
			cfg: Config{
				ServerPort:    8080,
				Environment:   "production",
				DatabaseURL:   "postgres://localhost/app",
				DebugMode:     true,
				AdminPassword: "secret",
			},
			want: "{\n  \"server_port\": 8080,\n  \"environment\": \"production\",\n  \"database_url\": \"postgres://localhost/app\",\n  \"debug_mode\": true\n}",
		},
		{
			name: "empty database url is omitted",
			cfg: Config{
				ServerPort:    3000,
				Environment:   "development",
				DebugMode:     false,
				AdminPassword: "secret",
			},
			want: "{\n  \"server_port\": 3000,\n  \"environment\": \"development\",\n  \"debug_mode\": false\n}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")

			if err := SaveConfig(path, tc.cfg); err != nil {
				t.Fatalf("SaveConfig(%q) returned unexpected error: %v", path, err)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read saved configuration: %v", err)
			}
			if got := string(data); got != tc.want {
				t.Errorf("saved configuration = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSaveConfigWriteError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "config.json")

	err := SaveConfig(path, Config{})
	if err == nil {
		t.Fatal("SaveConfig returned nil error for a path with a missing parent directory")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("SaveConfig error = %v, want an error matching os.ErrNotExist", err)
	}
}
