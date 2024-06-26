# Contributing

## Development guidelines

### Formatter

We use `golangci-lint` to validate the code formatting. Make sure to have your code formatted before opening pull requests.

### Line length limit

Every code line should not exceed 80 characters when possible, but the hard limit is at 120 characters. Lines with documentation have an exception with the hard limit at 80 characters. A tab is treated as 4 spaces. It is recommended to set ruler guide at 80 and 120 characters on your editor.

### Wrapping long struct initialization

If the struct initialization is too long for the 80 character limit, wrap the initialization by assigning each field on its own line and end the initialization on its own line.

**WRONG**

```go
a := User{
    Firstname: "John", Lastname: "Doe",
}

a := User{
    Firstname: "John",
    Lastname:  "Doe"}

a := User{Firstname: "John",
    Lastname: "Doe",
}
```

**RIGHT**

```go
a := User{
    Firstname: "John",
    Lastname:  "Doe",
}
```

### Wrapping long function definitions

If the function definition is too long for the 80 character limit, wrap the function definition by writing with as few lines as possible while avoiding trailing open parenthesis and traling commas after the last parameter and the last return type.

**WRONG**

```go
func doSomething(
    a,
    b,
    c,
) (d, error) {

func doSomething(a, b, c,
) (d, error) {

func doSomething(a, b, c) (
    d, error) {
```

**RIGHT**

```go
func doSomething(a, b,
    c) (d, error) {

func doSomething(a, b, c) (d,
    error) {
```

If the function definition spans multiple lines, the function body should start with an empty line to help distinguishing two elements.

## Releasing Steps

1. Make sure the dependencies are up to date
2. Make sure the examples are generated by the latest version
3. Make sure the `README.md` file reflects the current release
4. Update the unreleased version specified in `CHANGELOG.md`
5. Remove `-next` from version variable in `main.go`
6. Create a git tag `vX.X.X` and push to GitHub
7. Run `goreleaser release` with GitHub token provided as an environment variable `GITHUB_TOKEN`
8. Bump version with `-next` suffix in version variable in `main.go`
