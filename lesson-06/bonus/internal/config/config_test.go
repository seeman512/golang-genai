package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		{
			name:    "multiple json values",
			content: `{"server_port": 8080} {"server_port": 8081}`,
			wantErr: true,
		},
		{
			name:    "trailing invalid data",
			content: `{"server_port": 8080} trailing data`,
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
			want: "{\n  \"server_port\": 8080,\n  \"environment\": \"production\",\n  \"database_url\": \"postgres://localhost/app\",\n  \"debug_mode\": true\n}\n",
		},
		{
			name: "empty database url is omitted",
			cfg: Config{
				ServerPort:    3000,
				Environment:   "development",
				DebugMode:     false,
				AdminPassword: "secret",
			},
			want: "{\n  \"server_port\": 3000,\n  \"environment\": \"development\",\n  \"debug_mode\": false\n}\n",
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

// The blind-test prompt highlighted inclusive port boundaries that are easy to
// miss when implementing validation. It also added empty and unusual environment
// values, along with clearly invalid low and high ports. The generated cases
// were reviewed against the stated requirements before being added here.
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "minimum valid port",
			cfg:  Config{ServerPort: 1024, Environment: "development"},
		},
		{
			name: "maximum valid port",
			cfg:  Config{ServerPort: 65535, Environment: "production"},
		},
		{
			name: "valid staging configuration",
			cfg:  Config{ServerPort: 8080, Environment: "staging"},
		},
		{
			name:    "zero port",
			cfg:     Config{ServerPort: 0, Environment: "development"},
			wantErr: true,
		},
		{
			name:    "privileged port",
			cfg:     Config{ServerPort: 80, Environment: "development"},
			wantErr: true,
		},
		{
			name:    "port below minimum",
			cfg:     Config{ServerPort: 1023, Environment: "production"},
			wantErr: true,
		},
		{
			name:    "port above maximum",
			cfg:     Config{ServerPort: 65536, Environment: "production"},
			wantErr: true,
		},
		{
			name:    "large invalid port",
			cfg:     Config{ServerPort: 99999, Environment: "staging"},
			wantErr: true,
		},
		{
			name:    "empty environment",
			cfg:     Config{ServerPort: 8080},
			wantErr: true,
		},
		{
			name:    "unsupported environment",
			cfg:     Config{ServerPort: 8080, Environment: "testing"},
			wantErr: true,
		},
		{
			name:    "environment with unusual characters",
			cfg:     Config{ServerPort: 8080, Environment: "prod!@#$"},
			wantErr: true,
		},
		{
			name:    "environment with surrounding whitespace",
			cfg:     Config{ServerPort: 8080, Environment: " production "},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConfig(tc.cfg)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateConfig(%+v) error = %v, wantErr %v", tc.cfg, err, tc.wantErr)
			}
		})
	}
}

var errDiskFull = errors.New("disk full")

// mockStorage is a handwritten FileStorage implementation used to exercise
// filesystem failures without depending on host permissions or disk state.
type mockStorage struct {
	createErr error
	openErr   error
	writer    io.WriteCloser
	reader    io.ReadCloser
}

func (m mockStorage) Create(string) (io.WriteCloser, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.writer, nil
}

func (m mockStorage) Open(string) (io.ReadCloser, error) {
	if m.openErr != nil {
		return nil, m.openErr
	}
	return m.reader, nil
}

type mockWriteCloser struct {
	io.Writer
	closeErr error
}

func (m *mockWriteCloser) Close() error {
	return m.closeErr
}

type mockReadCloser struct {
	io.Reader
	closeErr error
}

func (m *mockReadCloser) Close() error {
	return m.closeErr
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestSaveConfigWithStorageErrors(t *testing.T) {
	tests := []struct {
		name    string
		storage FileStorage
		wantErr error
	}{
		{
			name:    "permission denied while creating",
			storage: mockStorage{createErr: os.ErrPermission},
			wantErr: os.ErrPermission,
		},
		{
			name: "disk full while writing",
			storage: mockStorage{
				writer: &mockWriteCloser{Writer: failingWriter{err: errDiskFull}},
			},
			wantErr: errDiskFull,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SaveConfigWithStorage(tc.storage, "config.json", Config{ServerPort: 8080})
			if err == nil {
				t.Fatal("SaveConfigWithStorage returned nil error, want error")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("SaveConfigWithStorage error = %v, want an error matching %v", err, tc.wantErr)
			}
		})
	}
}

func TestSaveConfigCloseError(t *testing.T) {
	closeErr := errors.New("close failed")
	storage := mockStorage{
		writer: &mockWriteCloser{Writer: io.Discard, closeErr: closeErr},
	}

	err := SaveConfigWithStorage(storage, "config.json", Config{ServerPort: 8080})
	if err == nil {
		t.Fatal("SaveConfigWithStorage returned nil error, want close error")
	}
	if !errors.Is(err, closeErr) {
		t.Errorf("SaveConfigWithStorage error = %v, want an error matching %v", err, closeErr)
	}
}

func TestLoadConfigWithStorageErrors(t *testing.T) {
	tests := []struct {
		name    string
		storage FileStorage
		wantErr error
	}{
		{
			name:    "permission denied while opening",
			storage: mockStorage{openErr: os.ErrPermission},
			wantErr: os.ErrPermission,
		},
		{
			name: "error while closing",
			storage: mockStorage{
				reader: &mockReadCloser{
					Reader:   strings.NewReader(`{"server_port": 8080, "environment": "production"}`),
					closeErr: errDiskFull,
				},
			},
			wantErr: errDiskFull,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadConfigWithStorage(tc.storage, "config.json")
			if err == nil {
				t.Fatal("LoadConfigWithStorage returned nil error, want error")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("LoadConfigWithStorage error = %v, want an error matching %v", err, tc.wantErr)
			}
		})
	}
}

func TestConfigWithNilStorage(t *testing.T) {
	if err := SaveConfigWithStorage(nil, "config.json", Config{}); err == nil || !strings.Contains(err.Error(), "file storage is nil") {
		t.Errorf("SaveConfigWithStorage(nil) error = %v, want a nil-storage error", err)
	}

	if _, err := LoadConfigWithStorage(nil, "config.json"); err == nil || !strings.Contains(err.Error(), "file storage is nil") {
		t.Errorf("LoadConfigWithStorage(nil) error = %v, want a nil-storage error", err)
	}
}
