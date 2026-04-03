# Copilot Instructions for go-dmesg

## Project Scope
- This repository provides a small Go library and command for reading and decoding Linux kernel messages from /dev/kmsg.
- Keep implementations dependency-light. Prefer the Go standard library unless a dependency is clearly justified.

## Code Organization
- Library code lives under pkg/dmesg.
- CLI sample lives under cmd/go-dmesg.go.
- Keep public APIs in pkg/dmesg stable unless a change is explicitly requested.

## Scanner and Decoder Conventions
- Prefer direct /dev/kmsg handling through the existing scanner abstraction.
- Ensure resources are closed by the creator of the scanner.
- When modifying scanner behavior, preserve non-blocking semantics.

## Testing Guidance
- Add unit tests for logic and error paths using mocks/fakes where possible.
- Keep the root-gated scanner integration test pattern for /dev/kmsg access checks.
- Do not remove permission checks for tests that require elevated privileges.

## CI and Verification
- Before finalizing changes, run make ci.
- CI is expected to run formatting checks, vet, tests with coverage output, coverage report, and build.
- Prefer updating existing Makefile targets instead of adding ad-hoc shell commands in docs or workflows.

## Style
- Follow idiomatic Go naming and formatting.
- Keep comments concise and focused on intent.
- Avoid broad refactors unrelated to the requested change.
