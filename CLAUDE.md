# CLAUDE.md - AI Assistant Guide for aws-services

**Last Updated**: 2025-11-17
**Repository**: https://github.com/rollwagen/aws-services
**Purpose**: CLI tool to list AWS service availability across regions

## Table of Contents

- [Project Overview](#project-overview)
- [Codebase Structure](#codebase-structure)
- [Technology Stack](#technology-stack)
- [Development Workflows](#development-workflows)
- [Code Architecture](#code-architecture)
- [Key Conventions](#key-conventions)
- [Common Development Tasks](#common-development-tasks)
- [Testing & Quality Assurance](#testing--quality-assurance)
- [Release Process](#release-process)
- [Important Notes for AI Assistants](#important-notes-for-ai-assistants)

---

## Project Overview

**aws-services** is a Go-based CLI tool that helps users check the availability of AWS services across different regions. It provides both interactive and command-line modes for querying AWS infrastructure information.

### Key Features
- Interactive service selection with visual feedback
- Concurrent availability checking across all AWS regions
- Lists all AWS services and regions
- Real-time progress indicators
- Colored output for easy readability

### Distribution
- **Homebrew**: `brew install rollwagen/tap/aws-services`
- **Go Install**: `go run github.com/rollwagen/aws-services@latest`
- **GitHub Releases**: Pre-built binaries for Linux and macOS (amd64, arm64)

---

## Codebase Structure

```
aws-services/
├── cmd/                          # Cobra command definitions
│   ├── root.go                   # Main command + interactive mode (94 lines)
│   ├── list.go                   # Parent command for list subcommands (15 lines)
│   ├── services.go               # Lists all AWS services (41 lines)
│   └── regions.go                # Lists all AWS regions (25 lines)
│
├── pkg/                          # Reusable packages
│   ├── service/
│   │   └── infrastructure.go    # Core AWS querying logic (197 lines)
│   └── prompter/
│       └── prompter.go          # Interactive prompt wrapper (38 lines)
│
├── assets/                       # Build artifacts
│   └── cosign/                   # Code signing materials
│
├── .github/
│   ├── workflows/                # CI/CD pipelines
│   │   ├── lint.yml             # golangci-lint
│   │   ├── release.yml          # GoReleaser automation
│   │   ├── codeql.yml           # Security scanning
│   │   └── semgrep.yml          # Static analysis
│   └── dependabot.yml           # Dependency updates
│
├── main.go                       # Entry point (10 lines)
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
│
├── .golangci.yml                 # Linter configuration
├── .goreleaser.yaml              # Release configuration
├── .pre-commit-config.yaml       # Pre-commit hooks
├── .checkov.yaml                 # Infrastructure security
├── .yamllint.yaml                # YAML linting
│
├── release.sh                    # Manual release script
└── README.md                     # User documentation
```

### File Statistics
- **Total Go files**: 7
- **Total code lines**: ~420 (excluding blanks/comments)
- **Largest file**: `pkg/service/infrastructure.go` (197 lines)
- **Test files**: 0 (no tests currently)

---

## Technology Stack

### Go Version
- **Minimum**: Go 1.23.0
- **Toolchain**: go1.24.0

### Core Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `spf13/cobra` | v1.9.1 | CLI framework and command structure |
| `aws/aws-sdk-go-v2` | v1.36.3 | AWS SDK for service queries |
| `aws-sdk-go-v2/service/ssm` | - | SSM Parameter Store (data source) |
| `AlecAivazis/survey/v2` | - | Interactive prompts |
| `briandowns/spinner` | v1.23.2 | Terminal spinner for progress |
| `fatih/color` | v1.18.0 | Colored terminal output |
| `samber/lo` | v1.49.1 | Functional programming utilities |
| `sourcegraph/conc` | v0.3.0 | Concurrency primitives |

### Supporting Libraries
- AWS config, credentials, STS, SSO support
- Terminal utilities (color detection, isatty)
- Standard library: fmt, os, strings, sync, time

### Build Tools
- **GoReleaser**: Multi-platform builds and distribution
- **Cosign**: Binary signing for security
- **golangci-lint**: Code quality enforcement
- **Pre-commit**: Git hooks for quality checks

---

## Development Workflows

### GitHub Actions CI/CD

#### 1. Lint Workflow (`lint.yml`)
- **Triggers**: Push to main, tags (v*), pull requests
- **Actions**: golangci-lint with Go stable version
- **Purpose**: Enforce code quality standards

#### 2. Release Workflow (`release.yml`)
- **Triggers**: Tag pushes (v*)
- **Steps**:
  1. Checkout code
  2. Setup Go (>=1.20.0)
  3. Install Cosign
  4. Run GoReleaser with signing
- **Secrets Required**: `GORELEASER_TOKEN`, `COSIGN_PWD`
- **Timeout**: 60 minutes

#### 3. CodeQL Workflow (`codeql.yml`)
- **Triggers**: Push to main, PRs, weekly schedule (Mon 19:40 UTC)
- **Purpose**: Security vulnerability scanning
- **Language**: Go

#### 4. Semgrep Workflow (`semgrep.yml`)
- **Triggers**: Manual dispatch, PRs, push to main, daily (08:58 UTC)
- **Purpose**: Static analysis for security and code quality
- **Secret**: `SEMGREP_APP_TOKEN`

#### 5. Dependabot (`dependabot.yml`)
- **Go modules**: Weekly updates (max 10 PRs)
- **GitHub Actions**: Monthly updates

### Pre-commit Hooks

Required hooks (run on every commit):
```yaml
# Go-specific
- go-imports          # Auto-import management
- go-mod-tidy         # Clean go.mod
- golangci-lint       # Full lint suite
- gofumpt             # Stricter formatting

# General
- trailing-whitespace
- end-of-file-fixer
- check-yaml
- check-added-large-files
- detect-private-key
- detect-aws-credentials (with --allow-missing-credentials)
- commitlint          # Conventional commit messages
```

### Linter Configuration (`.golangci.yml`)

**Enabled Linters** (16):
- `bodyclose` - HTTP response body closure
- `errcheck` - Unchecked errors
- `goconst` - Repeated strings that could be constants
- `gocritic` - Opinionated checks
- `gosec` - Security issues
- `govet` - Standard Go vet
- `ineffassign` - Ineffectual assignments
- `revive` - Enhanced golint replacement
- `staticcheck` - Advanced static analysis
- `unconvert` - Unnecessary type conversions
- `unparam` - Unused function parameters
- `unused` - Unused code

**Excluded Paths**:
- `scripts/`, `third_party/`, `builtin/`, `examples/`

---

## Code Architecture

### Entry Point Flow

```
main.go
  └─> cmd.Execute()
      └─> rootCmd (interactive mode) OR list subcommands
```

### Command Structure (Cobra)

```
aws-services (root command - cmd/root.go)
├─> Default: Interactive mode
│   ├─> Fetch services from SSM Parameter Store
│   ├─> Show interactive selection prompt
│   ├─> Query availability across all regions concurrently
│   └─> Display sorted results with visual indicators
│
└─> list (parent command - cmd/list.go)
    ├─> services (cmd/services.go) - List all AWS services
    └─> regions (cmd/regions.go) - List all AWS regions
```

### Data Source: AWS SSM Parameter Store

**Critical Implementation Detail**: This tool uses AWS Systems Manager Parameter Store as its data source for infrastructure metadata.

**SSM Paths**:
```go
// Regions list
/aws/service/global-infrastructure/regions/

// Services per region
/aws/service/global-infrastructure/regions/{region}/services/

// Reference region (most complete)
us-east-1
```

### Core Functions (`pkg/service/infrastructure.go`)

#### 1. `Services() ([]string, error)`
- **Purpose**: Returns list of all available AWS services
- **Method**: Queries SSM Parameter Store for us-east-1 services
- **Pagination**: Max 10 results per request
- **Returns**: Sorted slice of service names
- **Error Handling**: Returns wrapped error if SSM query fails

#### 2. `Regions() ([]string, error)`
- **Purpose**: Returns list of all AWS regions
- **Method**: Queries SSM Parameter Store for region list
- **Processing**: Extracts region names from parameter paths
- **Returns**: Sorted slice of region codes (e.g., "us-east-1")

#### 3. `AvailabilityPerRegion(service string, progress chan string) map[string]bool`
- **Purpose**: Checks service availability across all regions concurrently
- **Concurrency**: Pool of 3 goroutines (configurable via `concurrency` constant)
- **Thread Safety**: Uses `sync.Map` for concurrent writes
- **Progress**: Sends region names to channel as they're checked
- **Returns**: Map of region -> availability boolean

#### 4. `isAvailable(region, service string, serviceAvailability *sync.Map)`
- **Purpose**: Checks if service exists in specific region
- **Method**: Queries region-specific SSM path with pagination
- **Optimization**: Early return on first match
- **Side Effects**: Writes result to sync.Map

#### 5. `Names() ([]string, error)` (Legacy/Alternative)
- **Purpose**: Alternative service fetching from AWS API endpoint
- **Endpoint**: `api.regional-table.region-services.aws.a2z.com`
- **Status**: Not actively used in current implementation
- **Processing**: Parses JSON and deduplicates service names

### Concurrency Pattern

**Worker Pool Implementation**:
```go
const concurrency = 3  // Number of concurrent goroutines

concurrencyPool := pool.New().WithMaxGoroutines(concurrency)

for _, r := range regions {
    region := r  // Capture loop variable (Go best practice)
    concurrencyPool.Go(func() {
        isAvailable(region, service, &serviceAvailability)
    })
}

concurrencyPool.Wait()  // Block until all goroutines complete
```

**Why 3 goroutines?**
- Balance between AWS API rate limits and performance
- Prevents overwhelming SSM Parameter Store API
- Provides reasonable parallelism for ~20-30 regions

### User Experience Flow

**Interactive Mode** (default execution):
1. Display spinner: "Retrieving list of services..."
2. Fetch services via `service.Services()`
3. Show interactive selection prompt (15 items per page)
4. On selection:
   - Start progress spinner with elapsed time
   - Check availability concurrently
   - Update spinner with current region being checked
5. Display results:
   - Sorted by region name
   - Green ✔ for available, Red ✖ for unavailable
   - Clean, aligned output

### Prompter Abstraction (`pkg/prompter/prompter.go`)

**Design Pattern**: Interface-based for testability

```go
type Prompter interface {
    Select(message string, options []string) (int, error)
}

type Select struct{}  // Concrete implementation

func (s Select) Select(message string, options []string) (int, error) {
    // Wraps AlecAivazis/survey library
    // 15 items per page
    // Returns selected index or error
}
```

**Purpose**: Decouples UI library from command logic, enabling mocking in tests (when added).

---

## Key Conventions

### Package Organization

**Principle**: Separation of concerns

- **`cmd/`**: CLI layer - Cobra commands, user interaction, output formatting
- **`pkg/`**: Business logic - AWS API interactions, data processing
- **`main.go`**: Minimal entry point - delegates to cmd package

### Command Registration Pattern

**Convention**: Use `init()` functions for command tree construction

```go
// In cmd/list.go
func init() {
    rootCmd.AddCommand(listCmd)
}

// In cmd/services.go
func init() {
    listCmd.AddCommand(servicesCmd)
}

// In cmd/regions.go
func init() {
    listCmd.AddCommand(regionsCmd)
}
```

**Rationale**: Automatic registration, clear hierarchy, follows Cobra best practices.

### Error Handling

**Mixed Approaches** (be aware when modifying):

1. **Panic on Critical Errors** (infrastructure.go):
   ```go
   // Used for unrecoverable AWS client creation failures
   if err != nil {
       panic(err)  // Lines 102, 116
   }
   ```
   **When to use**: Configuration errors, impossible states

2. **Return Errors** (most functions):
   ```go
   func Services() ([]string, error) {
       if err != nil {
           return nil, fmt.Errorf("could not retrieve params: %w", err)
       }
   }
   ```
   **When to use**: Expected failures, recoverable errors

3. **Print and Exit** (cmd layer):
   ```go
   if err != nil {
       fmt.Fprintln(os.Stderr, err)
       os.Exit(1)
   }
   ```
   **When to use**: Top-level command failures

4. **Ignore Errors** (rare):
   ```go
   regions, _ := service.Regions()  // cmd/regions.go:16
   ```
   **Warning**: Only used in regions list command - acceptable since it prints whatever it gets

**Recommendation for New Code**: Return errors from pkg/, handle at cmd/ level.

### Concurrency Conventions

**Pattern**: Worker pool with sync.Map

```go
// Always capture loop variables
for _, item := range items {
    i := item  // Capture for goroutine
    pool.Go(func() {
        process(i)
    })
}
```

**Thread-Safe Data**:
- Use `sync.Map` for concurrent writes
- Convert to regular map after synchronization point
- Channel for progress communication (non-blocking sends)

### Code Style

**Enforced via tooling**:
- **Formatting**: gofumpt (stricter than gofmt)
- **Imports**: goimports (automatic grouping and sorting)
- **Linting**: golangci-lint with 16 enabled linters

**Naming Conventions**:
- Exported: `PascalCase` (functions, types, constants)
- Private: `camelCase` (internal functions, variables)
- Acronyms: `SSM`, `AWS`, `API` (all caps when at start)
- Descriptive: Avoid single-letter vars except loop indices

**Comments**:
- Exported functions should have doc comments
- Format: `// FunctionName does something`
- Keep concise but informative

### Configuration Management

**AWS Credentials**:
- Uses standard AWS SDK credential chain
- Supports: environment variables, ~/.aws/credentials, IAM roles, SSO
- No hardcoded credentials (enforced by pre-commit hooks)

**No Config Files**:
- Tool has no configuration files
- All behavior controlled via AWS credentials and CLI flags
- Philosophy: Simple, portable, minimal setup

---

## Common Development Tasks

### Setup Development Environment

```bash
# Clone repository
git clone https://github.com/rollwagen/aws-services.git
cd aws-services

# Install dependencies
go mod download

# Install pre-commit hooks (recommended)
pip install pre-commit
pre-commit install

# Configure AWS credentials (required for testing)
aws configure
# OR use environment variables
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_REGION=us-east-1
```

### Build and Run Locally

```bash
# Build binary
go build -o aws-services

# Run directly with Go
go run main.go

# Run with arguments
go run main.go list services
go run main.go list regions

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o aws-services-linux
GOOS=darwin GOARCH=arm64 go build -o aws-services-macos-arm64
```

### Linting and Code Quality

```bash
# Run golangci-lint
golangci-lint run

# Auto-fix issues
golangci-lint run --fix

# Run specific linters
golangci-lint run --enable-only=gosec

# Format code
gofumpt -w .
goimports -w .

# Run pre-commit hooks manually
pre-commit run --all-files
```

### Dependency Management

```bash
# Add new dependency
go get github.com/some/package@version

# Update dependency
go get -u github.com/some/package

# Update all dependencies
go get -u ./...

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# View dependency graph
go mod graph
```

### Testing AWS Functionality

**Note**: Currently no automated tests exist.

```bash
# Manual testing commands
./aws-services                    # Test interactive mode
./aws-services list services      # Test services listing
./aws-services list regions       # Test regions listing

# Test with specific AWS profile
AWS_PROFILE=myprofile ./aws-services

# Test with SSO
aws sso login --profile myprofile
AWS_PROFILE=myprofile ./aws-services
```

### Debugging

```bash
# Enable verbose AWS SDK logging
export AWS_SDK_LOAD_CONFIG=1
export AWS_LOG_LEVEL=debug

# Run with delve debugger
dlv debug main.go

# Build with debug symbols
go build -gcflags="all=-N -l" -o aws-services-debug
```

---

## Testing & Quality Assurance

### Current State
**⚠️ IMPORTANT**: This project currently has **no automated tests**.

**Missing Test Coverage**:
- No unit tests (`*_test.go` files)
- No integration tests
- No test coverage metrics
- No mocking infrastructure

**Pre-commit config shows**: `# - id: go-unit-tests` (commented out)

### Testing Strategy (for future implementation)

**Recommended Test Structure**:
```
pkg/
├── service/
│   ├── infrastructure.go
│   └── infrastructure_test.go      # Unit tests
│
└── prompter/
    ├── prompter.go
    └── prompter_test.go            # Interface mocking

cmd/
├── root_test.go                     # Command integration tests
├── services_test.go
└── regions_test.go
```

**Test Considerations**:

1. **AWS API Mocking**:
   - Use `aws-sdk-go-v2/aws/fakes` for mocking SSM client
   - Mock SSM responses for deterministic tests
   - Test error conditions (rate limits, permissions)

2. **Prompter Interface**:
   - Already designed for testability
   - Create mock implementation for command tests
   - Test user selection flows

3. **Concurrency Tests**:
   - Test worker pool behavior
   - Verify thread-safe operations
   - Test progress channel communication

4. **CLI Tests**:
   - Use cobra's testing utilities
   - Capture stdout/stderr
   - Test command flags and arguments

**Example Test Structure**:
```go
// pkg/service/infrastructure_test.go
func TestServices(t *testing.T) {
    // Mock SSM client
    // Call Services()
    // Assert expected service list
}

func TestAvailabilityPerRegion(t *testing.T) {
    // Mock SSM responses
    // Test concurrent execution
    // Verify progress channel updates
}
```

### Quality Assurance Tools

**Automated Checks** (enforced):
- golangci-lint (16 linters)
- CodeQL security scanning
- Semgrep static analysis
- Pre-commit hooks (formatting, secrets detection)
- Dependabot (dependency updates)

**Manual Checks** (recommended):
- Test interactive mode manually before releases
- Verify AWS API responses for new regions/services
- Check output formatting across different terminals
- Test with different AWS credential methods

---

## Release Process

### Automated Release (Recommended)

**Trigger**: Push a semantic version tag

```bash
# Using the provided release script
./release.sh 1.2.3

# Manual tag creation
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

**GitHub Actions Workflow** (`release.yml`):
1. Checkout code
2. Setup Go (>=1.20.0)
3. Install Cosign for signing
4. Run GoReleaser:
   - Build for linux/darwin on amd64/arm64
   - Generate shell completions (bash, zsh, fish)
   - Create archives with README, LICENSE, completions
   - Sign checksums with Cosign
   - Publish to GitHub Releases
   - Update Homebrew tap (rollwagen/homebrew-tap)
   - Generate changelog (excludes `docs:`, `test:` commits)

### GoReleaser Configuration (`.goreleaser.yaml`)

**Build Settings**:
```yaml
env:
  - CGO_ENABLED=0          # Static binaries

ldflags:
  - -s -w                  # Strip symbols
  - -X main.version={{.Version}}
  - -X main.commit={{.Commit}}
  - -X main.date={{.Date}}
```

**Distributions**:
- GitHub Releases (primary)
- Homebrew tap (automatic)
- Checksums signed with Cosign

**Pre-release Hooks**:
```bash
go mod tidy
go run main.go completion bash > completions/aws-services.bash
go run main.go completion zsh > completions/aws-services.zsh
go run main.go completion fish > completions/aws-services.fish
```

### Version Numbering

**Format**: Semantic Versioning (MAJOR.MINOR.PATCH)

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

**Tag Format**: `v1.2.3` (note the 'v' prefix)

### Release Checklist

1. **Pre-release**:
   - [ ] Update README if needed
   - [ ] Ensure CHANGELOG reflects changes
   - [ ] Run `go mod tidy`
   - [ ] Run all linters locally
   - [ ] Manual testing of interactive mode
   - [ ] Test with different AWS profiles/regions

2. **Release**:
   - [ ] Determine version number (SemVer)
   - [ ] Run `./release.sh X.Y.Z` OR manually create tag
   - [ ] Verify GitHub Actions workflow succeeds
   - [ ] Check GitHub Releases page for artifacts

3. **Post-release**:
   - [ ] Test Homebrew installation: `brew upgrade aws-services`
   - [ ] Verify binary signatures
   - [ ] Update documentation if needed
   - [ ] Announce release (if major version)

### Rollback Strategy

**If release fails or has critical bugs**:
```bash
# Delete remote tag
git push --delete origin v1.2.3

# Delete local tag
git tag -d v1.2.3

# Delete GitHub Release (via UI or gh CLI)
gh release delete v1.2.3

# Fix issues, then re-tag
./release.sh 1.2.4
```

---

## Important Notes for AI Assistants

### Critical Context

1. **No Tests Exist**: When adding features, consider whether tests should be added. The infrastructure for testability exists (prompter interface, modular design) but no tests are implemented.

2. **AWS SSM is the Data Source**: All service/region data comes from SSM Parameter Store, not AWS SDK service enumeration. This is a critical architectural choice that should not be changed without understanding the implications.

3. **Concurrency is Limited to 3**: The worker pool size is intentionally set to 3 to respect AWS API rate limits. Do not increase without testing against actual AWS accounts.

4. **Error Handling is Mixed**: The codebase uses panic(), return errors, and os.Exit() in different contexts. Follow the existing patterns in each layer.

5. **No Configuration Files**: This is a design choice for simplicity. Don't add config files unless absolutely necessary and discussed with maintainers.

### When Making Changes

**DO**:
- Run pre-commit hooks before committing (`pre-commit run --all-files`)
- Follow the existing error handling patterns for each layer
- Use the established concurrency patterns (worker pools, sync.Map)
- Test manually with real AWS credentials
- Update comments for exported functions
- Keep code simple and readable
- Consider AWS API rate limits in design

**DON'T**:
- Add tests without discussing the testing strategy (none exist currently)
- Change the SSM data source without strong justification
- Increase concurrency without testing rate limits
- Add configuration files without discussion
- Introduce breaking changes without major version bump
- Ignore linter warnings (fix or explicitly exclude)
- Commit secrets or credentials (hooks will catch this)

### Code Modification Guidelines

**Adding a New Command**:
1. Create file in `cmd/` (e.g., `cmd/newcommand.go`)
2. Define cobra.Command struct
3. Add `init()` function to register with parent command
4. Implement `Run` or `RunE` function
5. Handle errors appropriately (print to stderr, exit 1)
6. Update shell completions in GoReleaser config

**Adding New AWS Functionality**:
1. Add to `pkg/service/infrastructure.go`
2. Use existing SSM client creation pattern
3. Follow pagination pattern for SSM queries
4. Return errors, don't panic (except for client creation)
5. Consider concurrency if querying multiple regions
6. Update `cmd/` layer to use new functionality

**Modifying Output Formatting**:
1. Changes likely in `cmd/root.go`
2. Use `fatih/color` for colored output
3. Maintain alignment for table-like output
4. Test in multiple terminal types
5. Consider colorblind-friendly indicators

### Security Considerations

**Enforced by Pre-commit Hooks**:
- No AWS credentials in code
- No private keys
- No large files (accidental binary commits)

**CodeQL and Semgrep**:
- Scan for security vulnerabilities
- Check for common Go security issues
- Validate AWS SDK usage patterns

**When Adding Dependencies**:
- Verify package is maintained and trusted
- Check for known vulnerabilities (dependabot will help)
- Minimize dependencies to reduce attack surface
- Run `go mod tidy` after adding

### Performance Considerations

**Current Performance Characteristics**:
- Service list: ~1-2 seconds (SSM query for us-east-1)
- Region list: ~1-2 seconds (SSM query for all regions)
- Availability check: ~10-15 seconds (3 concurrent queries across ~25 regions)

**Optimization Opportunities** (if needed):
- Increase concurrency (test rate limits first)
- Cache SSM responses (add cache invalidation strategy)
- Parallel service and region fetching
- Progress bar instead of spinner (more informative)

**Don't Optimize Prematurely**:
- Current performance is acceptable for CLI tool
- Complexity tradeoffs may not be worth it
- User experience (UX) is already good

### AWS SDK Best Practices

**Credential Handling**:
- Never hardcode credentials
- Use default credential chain
- Support all standard AWS credential methods
- Let SDK handle credential refresh

**Error Handling**:
- Wrap AWS errors with context: `fmt.Errorf("operation failed: %w", err)`
- Check for specific AWS error types if needed
- Provide helpful error messages to users
- Don't expose internal AWS details in user-facing errors

**API Usage**:
- Use pagination for all list operations
- Respect rate limits (hence the concurrency=3)
- Consider regional endpoints (already doing this)
- Handle throttling errors gracefully (exponential backoff)

### Maintenance Notes

**Dependencies to Watch**:
- `aws-sdk-go-v2`: Major updates may require code changes
- `cobra`: CLI framework changes rare but impactful
- `survey`: UI library updates may affect interactive mode

**Go Version Updates**:
- Currently requires 1.23.0, toolchain 1.24.0
- Update conservatively (wait for minor releases to stabilize)
- Test thoroughly after Go version bumps
- Update in go.mod and GitHub Actions workflows

**GitHub Actions**:
- Dependabot keeps actions up to date
- Review action updates for breaking changes
- Secrets (`GORELEASER_TOKEN`, `COSIGN_PWD`) must remain valid

### Common Pitfalls to Avoid

1. **Loop Variable Capture**: Always capture loop variables in goroutines
   ```go
   for _, item := range items {
       i := item  // Capture before goroutine
       go func() { use(i) }()
   }
   ```

2. **SSM Pagination**: Always paginate SSM queries, don't assume <10 results
   ```go
   for paginator.HasMorePages() {
       page, err := paginator.NextPage(ctx)
       // Process page
   }
   ```

3. **Error Wrapping**: Use `%w` for error wrapping, not `%v`
   ```go
   return fmt.Errorf("failed to fetch: %w", err)  // Correct
   return fmt.Errorf("failed to fetch: %v", err)  // Wrong
   ```

4. **Spinner Cleanup**: Always stop spinners, use defer
   ```go
   spinner.Start()
   defer spinner.Stop()
   ```

5. **Channel Communication**: Use non-blocking sends for progress
   ```go
   select {
   case progressChan <- region:
   default: // Don't block if nobody is listening
   }
   ```

### Useful Commands for Development

```bash
# Check what GoReleaser would do (dry run)
goreleaser release --snapshot --clean

# Build local snapshot release
goreleaser build --snapshot --clean

# Test shell completions
source <(./aws-services completion bash)

# View effective golangci-lint config
golangci-lint config path

# List all linters
golangci-lint linters

# Profile memory/CPU (requires pprof additions)
go build -o aws-services
./aws-services # with pprof enabled
go tool pprof [binary] [profile]

# View module dependencies
go list -m all
go mod graph | grep specific-package
```

### Documentation References

- **Cobra CLI**: https://cobra.dev/
- **AWS SDK Go v2**: https://aws.github.io/aws-sdk-go-v2/docs/
- **GoReleaser**: https://goreleaser.com/
- **golangci-lint**: https://golangci-lint.run/
- **Pre-commit**: https://pre-commit.com/

---

## Quick Reference

### Project Commands
```bash
go run main.go                    # Interactive mode
go run main.go list services      # List all services
go run main.go list regions       # List all regions
go build                          # Build binary
golangci-lint run                 # Run linter
./release.sh 1.2.3                # Create release
```

### Key Files to Understand
1. `pkg/service/infrastructure.go` - Core AWS logic
2. `cmd/root.go` - Interactive mode
3. `.golangci.yml` - Linter rules
4. `.goreleaser.yaml` - Release process
5. `go.mod` - Dependencies

### Architecture Summary
- **CLI Framework**: Cobra
- **Data Source**: AWS SSM Parameter Store
- **Concurrency**: 3 goroutines with sync.Map
- **User Input**: AlecAivazis/survey
- **Progress**: briandowns/spinner
- **Output**: fatih/color

### Workflow Summary
1. User runs tool → Interactive or list mode
2. Fetch services/regions from SSM
3. (Interactive) User selects service
4. Query all regions concurrently (3 workers)
5. Display sorted, colored results

---

**End of CLAUDE.md**

*This document should be updated whenever significant architectural changes are made to the codebase.*
