# repogen

<a href="https://github.com/sunboyy/repogen/actions?query=workflow%3Abuild">
    <img src="https://github.com/sunboyy/repogen/workflows/build/badge.svg" alt="build status badge">
</a>
<a href="https://codecov.io/gh/sunboyy/repogen">
    <img src="https://codecov.io/gh/sunboyy/repogen/branch/main/graph/badge.svg?token=9BD5Y8X7NO"/>
</a>

Repogen is a code generator for database repository in Golang inspired by Spring Data JPA. (WIP)

## Features

Repogen is a library that generates MongoDB repository implementation from repository interface by using method name pattern.

- CRUD functionality
- Method signature validation
- Supports single-entity and multiple-entity operations
- Supports many comparison operators

## Getting Started

This getting started tutorial shows a simple example on how to use repogen. You can also see the working code inside `examples` directory for more information.

### Step 1: Download and install repogen

Run this command in your terminal to download and install repogen

```
$ go get github.com/sunboyy/repogen
```

### Step 2: Write a repository specification

Write repository specification as an interface in the same file as the model struct. There are 4 types of operations that are currently supported and are determined by the first word of the method name. Single-entity and multiple-entity modes are determined be the first return value. More complex queries can also be written.

```go
// You write this interface specification (comment is optional)
type UserRepository interface {
	// InsertOne stores userModel into the database and returns inserted ID if insertion
	// succeeds and returns error if insertion fails.
	InsertOne(ctx context.Context, userModel *UserModel) (interface{}, error)

	// FindByUsername queries user by username. If a user with specified username exists,
	// the user will be returned. Otherwise, error will be returned.
	FindByUsername(ctx context.Context, username string) (*UserModel, error)

	// UpdateDisplayNameByID updates a user with the specified ID with a new display name.
	// If there is a user matches the query, it will return true. Error will be returned
	// only when error occurs while accessing the database.
	UpdateDisplayNameByID(ctx context.Context, displayName string, id primitive.ObjectID) (bool, error)

	// DeleteByCity deletes users that have `city` value match the parameter and returns
	// the match count. The error will be returned only when error occurs while accessing
	// the database. This is a MANY mode because the first return type is an integer.
	DeleteByCity(ctx context.Context, city string) (int, error)
}
```

### Step 3: Run the repogen

Run this command to generate a repository implementation from the specification.

```
$ repogen -src=<src_file> -dest=<dest_file> -model=<model_struct> -repo=<repo_interface>
```

- `<src_file>` is the file that contains struct model and repository interface code
- `<dest_file>` is the destination path for the repository implementation to be generated
- `<model_struct>` is the name of the model struct for generating the repository
- `<repo_interface>` is the name of the repository interface to generate implementation from

For example:

```
$ repogen -src=examples/getting-started/user.go -dest=examples/getting-started/user_repo.go -model=UserModel -repo=UserRepository
```

You can also write the above command in the `go:generate` format inside Go files in order to generate the implementation when `go generate` command is executed.
