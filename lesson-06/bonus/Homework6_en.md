# Homework 6 — File I/O, JSON Serialization, and Table-Driven Testing in Go

**Estimated time:** 2 hours (strict time limit)

**Objective:** Learn to combine low-level disk operations, structured data
parsing, interface design for isolating errors with mocks, and modern
AI-driven TDD practices.

---

## Part 1: Manual Execution (Coding from Scratch) — ~45 minutes

> Complete this part entirely by hand, without AI autocomplete or chat tools,
> to reinforce syntax memory and architectural understanding.

### Task: Build a Server Configuration Manager

Create a package for managing a server configuration file in JSON format.

1. **Define the `Config` struct** in `config.go`:

   - Add the following fields: `ServerPort` (`int`), `Environment` (`string`),
     `DatabaseURL` (`string`), and `DebugMode` (`bool`).
   - Add struct tags so that all JSON keys use `snake_case`. If `DatabaseURL`
     is empty, omit it from the output with `omitempty`.
   - Add an `AdminPassword` (`string`) field. It must never be serialized into
     JSON, so use the `json:"-"` tag.

2. **Implement two functions**:

   - `SaveConfig(path string, cfg Config) error`: Create a file at the specified
     path and write the pretty-printed JSON configuration using
     `json.MarshalIndent`. Safely close the file descriptor with a deferred
     call. Wrap all file creation and write errors with useful context.
   - `LoadConfig(path string) (Config, error)`: Read the configuration from a
     file. If the file does not exist, return a clear wrapped error using
     `fmt.Errorf` with the `%w` verb, so callers can check it with
     `errors.Is(err, os.ErrNotExist)`.

3. **Write a table-driven test** in `config_test.go`:

   - Write `TestLoadConfig` manually.
   - Create a slice of test cases with at least three cases: successful read,
     file not found, and corrupted JSON.
   - Use `t.Run` to isolate each subtest. Use `t.Fatalf` for critical assertions,
     such as failing to set up a temporary file, and `t.Errorf` for non-fatal
     field mismatches.

---

## Part 2: Prompt Engineering (Co-piloting with an LLM) — ~30 minutes

> This part teaches you how to collaborate with an LLM as a pairing partner,
> focusing on generating comprehensive edge cases.

### Task: Blind Testing and Input Validation

Expand your codebase by adding a configuration validation function:

```go
func ValidateConfig(cfg Config) error
```

Validation rules:

- `ServerPort` must be between `1024` and `65535`.
- `Environment` must be one of `development`, `staging`, or `production`.

1. **Draft a prompt** for an LLM, such as ChatGPT or Claude, following the
   blind AI test-generation pattern:

   - Provide only the function signature, `ValidateConfig(cfg Config) error`,
     and the `Config` struct definition. Do not share your validation
     implementation.
   - Instruct the AI to generate a comprehensive table-driven test suite,
     `TestValidateConfig`, with at least eight distinct edge cases covering
     boundary values, empty inputs, unusual characters, and invalid port ranges
     such as `0`, `80`, and `99999`.

2. **Execute the generated tests** against your manual validation code.

3. **Compare the results:** Which boundary conditions did the AI think of that
   you missed? Did the AI hallucinate any incorrect test assertions? Write a
   brief comparison note of three to four sentences directly in comments in
   your test file.

---

## Part 3: Agentic AI Usage (Automation and Refactoring) — ~45 minutes

> This part trains you to delegate high-level refactoring tasks to autonomous
> AI agents, such as Cursor Composer, Copilot Workspace, or Claude Engineer,
> while supervising quality and coverage.

### Task: Performance Tuning, Mocking Disk Failures, and Coverage Verification

Assign your AI agent the following comprehensive refactoring and quality
assurance prompt:

```text
Analyze the config package in the current directory and refactor it based on
the following requirements:

1. Replace the in-memory json.Unmarshal/json.Marshal calls with streaming
   json.NewDecoder and json.NewEncoder to optimize RAM usage for large
   configuration payloads.
2. Introduce an implicit interface named FileStorage to abstract all disk
   read/write operations.
3. Write a handwritten mock struct named mockStorage, without pulling in
   external mocking frameworks, that allows simulating filesystem failures
   such as "Permission Denied" or "Disk Full". Add a test verifying that the
   package gracefully handles and wraps these physical disk errors.
4. Run the test coverage tool and ensure total coverage is above 80%. If
   coverage is lower, generate targeted table-driven test cases to cover the
   missing logic paths.
```

### Verification pass

- Review the changes made by the AI agent to ensure that the business logic was
  not compromised.
- Manually execute the coverage analysis commands to verify the agent's work:

  ```bash
  go test -coverprofile=cover.out ./...
  go tool cover -html=cover.out
  ```

- Open the coverage report in your browser, look for red lines indicating
  untested code branches, and confirm that overall coverage exceeds 80%.
