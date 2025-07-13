# Go Modules Sample

## Steps

1. Install the cmd `go install github.com/asmsh/flagged/cmd/genflagged`
2. `go generate` This should create `permissions_flagged.go`.
3. `go mod tidy` to add the `github.com/asmsh/flagged` dependency.
4. `go run .` to see it in action.
