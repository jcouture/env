# Run Aqua Security’s Trivy to catch possible vulnerabilities in the codebase
@audit:
  docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -v {{justfile_directory()}}:/path aquasec/trivy fs --scanners config,secret,vuln /path

# Update dependencies
@go-mod-update:
  go get -u
  go mod tidy

# Run unit tests
@test:
  go test -v -cover  ./... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
