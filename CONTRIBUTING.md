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
