### Assure unit test coverage meets minimums
- `go test ./... -coverprofile cover.out | tee <path_to_coverage_output>`
- `go run ./ci/unit_test_coverage/main.go --coverage <path_to_coverage_output> --limits <path_to_limits_file>`

### Bypassing in an emergency
- The `unit_test_coverage` application has a `--bypass` flag that can be set. To push code that lowers test coverage, either enable the `bypass` flag in the CI script, or manually alter the `limits.json` file.

### Updating coverage
- If you've increased coverage and want to set a new high threshold, generate a coverage file and use `unit_test_coverage` to update the `limits.json` file:
  - `go test ./... -coverprofile cover.out | tee <path_to_coverage_output>`
  - `go run ./ci/unit_test_coverage/main.go --update --coverage <path_to_coverage_output> --limits <path_to_limits_file>`