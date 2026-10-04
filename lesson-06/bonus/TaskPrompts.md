# Homework 6 prompts

## Part 2 — blind validation-test prompt

The following prompt intentionally gives the test generator only the public
contract, not the implementation:

> Generate a comprehensive table-driven Go test suite named `TestValidateConfig`
> for the function `ValidateConfig(cfg Config) error` below. Include at least
> eight distinct cases. Cover both inclusive port boundaries, invalid port
> ranges including `0`, `80`, and `99999`, empty input, unsupported environment
> names, unusual characters, and every allowed environment. Use `t.Run` for
> each case and make the expected result explicit. Do not assume behavior that
> is not stated by the validation rules.
>
> ```go
> type Config struct {
>     ServerPort    int    `json:"server_port"`
>     Environment   string `json:"environment"`
>     DatabaseURL   string `json:"database_url,omitempty"`
>     DebugMode     bool   `json:"debug_mode"`
>     AdminPassword string `json:"-"`
> }
>
> func ValidateConfig(cfg Config) error
> ```
>
> Validation rules:
>
> - `ServerPort` must be between `1024` and `65535`, inclusive.
> - `Environment` must be exactly `development`, `staging`, or `production`.

### Part 2 implementation

`ValidateConfig` is implemented in `internal/config/config.go`. The test suite
in `config_test.go` exercises valid boundary values, all allowed environments,
invalid ports, empty input, unsupported names, unusual characters, and
whitespace. The test file also contains a short comparison note documenting
which edge cases the blind prompt emphasized.

## Part 3 — refactoring and quality-assurance prompt

> Analyze the configuration package in the current directory and refactor it
> according to these requirements:
>
> 1. Replace in-memory `json.Unmarshal` and `json.Marshal` calls with streaming
>    `json.NewDecoder` and `json.NewEncoder`. Preserve the existing
>    `SaveConfig(path string, cfg Config) error` and
>    `LoadConfig(path string) (Config, error)` APIs, including useful wrapped
>    errors and pretty-printed output.
> 2. Introduce an implicit `FileStorage` interface that abstracts file opening
>    and creation. Add injectable helpers if needed so the public APIs can keep
>    their existing signatures while tests can provide a storage implementation.
> 3. Add a handwritten `mockStorage` in the test package, without external
>    mocking frameworks. Use it to simulate physical failures such as permission
>    denied while opening or creating a file, disk full while writing, and a
>    close failure. Verify that errors remain discoverable with `errors.Is`.
> 4. Add targeted table-driven tests for all new success and failure paths. Run
>    `go test -coverprofile=cover.out ./...` and ensure the package has more than
>    80% statement coverage. Do not change dependencies.

### Part 3 implementation

The implementation uses `json.NewEncoder` with indentation for saves and
`json.NewDecoder` for loads. `FileStorage`, `SaveConfigWithStorage`, and
`LoadConfigWithStorage` provide dependency injection while the original APIs
use the operating-system-backed storage implementation. Tests use a handwritten
mock to cover permission, disk-full, and close errors and verify wrapped errors.
