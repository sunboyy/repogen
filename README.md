# repogen

<a href="https://github.com/sunboyy/repogen/actions?query=workflow%3Abuild">
    <img src="https://github.com/sunboyy/repogen/workflows/build/badge.svg" alt="build status badge">
</a>
<a href="https://codecov.io/gh/sunboyy/repogen">
    <img src="https://codecov.io/gh/sunboyy/repogen/branch/main/graph/badge.svg?token=9BD5Y8X7NO"/>
</a>
<a href="https://codeclimate.com/github/sunboyy/repogen/maintainability">
	<img src="https://api.codeclimate.com/v1/badges/d0270245c28814200c5f/maintainability" />
</a>

Repogen is a code generator for database repository in Golang inspired by Spring Data JPA.

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

Write repository specification as an interface in the same file as the model struct. There are 5 types of operations that are currently supported and are determined by the first word of the method name. Single-entity and multiple-entity modes are determined be the first return value. More complex queries can also be written.

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

	// CountByCity returns the number of rows that match the given city parameter. If an
	// error occurs while accessing the database, error value will be returned.
	CountByCity(ctx context.Context, city string) (int, error)
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
$ repogen -src=examples/getting-started/user.go -dest=examples/getting-started/user_repo.go \
        -model=UserModel -repo=UserRepository
```

You can also write the above command in the `go:generate` format inside Go files in order to generate the implementation when `go generate` command is executed.

## Usage

### Method Definition

To begin, your method name must be in pascal-case (camel-case with beginning uppercase letter). Repogen determines an operation for a method by getting the **first word** of the method name. There are 5 supported words which refer to 5 supported operations.

1. `Insert` - Stores new data to the database
2. `Find` - Retrives data from the database
3. `Update` - Changes some fields of the data in the database
4. `Delete` - Removes data from the database
5. `Count` - Retrieves number of matched documents in the database

Each of the operations has their own requirements for the method name, parameters and return values. Please consult the documentation for each operation for its requirements.

#### Insert operation

An `Insert` operation has a very limited use case, i.e. inserting a single document or multiple documents. So, it is quite limited in method parameters and method returns. An insert method can only have one of these signatures.

```go
// InsertOne inserts a single document
InsertOne(ctx context.Context, model *Model) (interface{}, error)

// InsertMany inserts multiple documents
InsertMany(ctx context.Context, models []*Model) ([]interface{}, error)
```

Repogen determines a single-entity operation or a multiple-entity by checking the second parameter and the first return value. However, the operation requires the first parameter to be of type `context.Context` and the second return value to be of type `error`.

As the `Insert` operation has a limited use case, we do not want to limit you on how you name your method. Any method that has the name starting with the word `Insert` is always valid. For example, you can name your method `InsertAWholeBunchOfDocuments` and it will work as long as you specify method parameters and returns correctly.

#### Find operation

A `Find` operation also has two modes like `Insert` operation: single-entity and multiple-entity. However `Find` operation can be very simple or complex depending on how complex the query is. In this section, we will show you how to write single-modes and multiple-entity modes of find method with a simple query. For more information about more complex queries, please consult the query specification section in this document.

```go
// FindByID gets a single document by ID
FindByID(ctx context.Context, id primitive.ObjectID) (*Model, error)

// FindByCity gets all documents that match city parameter
FindByCity(ctx context.Context, city string) ([]*Model, error)

// FindAll gets all documents
FindAll(ctx context.Context) ([]*Model, error)
```

Repogen determines a single-entity or a multiple-entity operation by checking the first return value. If it is a pointer of a model, the method will be single-entity operation. If it is a slice of pointers of a model, the method will be multiple-entity operation.

The requirement of the `Find` operation method is that there must be only two return values, the second return value must be of type `error` and the first method parameter must be of type `context.Context`. The requirement of number of method parameters depends on the query which will be described in the query specification section.

#### Update operation

An `Update` operation also has two modes like `Insert` and `Find` operations: single-entity and multiple-entity. An `Update` operation also supports querying like `Find` operation. However, an `Update` operation requires more parameters than `Find` method, i.e. new values of updating fields. Specifying the query is the same as in `Find` method but specifying the updating fields are a little different.

```go
// UpdateDisplayNameAndCityByID updates a single document with a new display name and
// city by ID
UpdateDisplayNameAndCityByID(ctx context.Context, displayName string, city string,
	id primitive.ObjectID) (bool, error)

// UpdateGenderByCity updates Gender field of documents with matching city parameter
UpdateGenderByCity(ctx context.Context, gender Gender, city string) (int, error)
```

Repogen determines a single-entity operation or a multiple-entity by checking the first return value. If it is of type `bool`, the method will be single-entity operation. If it is of type `int`, the method will be multiple-entity operation. For single-entity operation, the method returns true if there is a matching document. For multiple-entity operation, the integer return shows the number of matched documents.

The requirement of the `Update` operation method is that there must be only two return values, the second return value must be of type `error` and the first method parameter must be of type `context.Context`. The requirement of number of method parameters depends on the number of updating fields and the query. Updating fields must be directly after context parameter and query fields must be directly after updating fields.

#### Delete operation

A `Delete` operation is the very similar to `Find` operation. It has two modes. The method name pattern is the same. The method parameters and returns are also almost the same except that `Delete` operation has different first return value of the method. For single-entity operation, the method returns true if there is a matching document. For multiple-entity operation, the integer return shows the number of matched documents.

```go
// DeleteByID deletes a single document by ID
DeleteByID(ctx context.Context, id primitive.ObjectID) (bool, error)

// DeleteByCity deletes all documents that match city parameter
DeleteByCity(ctx context.Context, city string) (int, error)

// DeleteAll deletes all documents
DeleteAll(ctx context.Context) (int, error)
```

#### Count operation

A `Count` operation is also similar to `Find` operation except it has only multiple-entity mode. This means that the method returns are always the same for any count operations. The method name pattern and the parameters are the same as `Find` operation.

```go
// CountByGender returns number of documents that match gender parameter
CountByGender(ctx context.Context, gender Gender) (int, error)
```

### Query Specification

A query can be applied on `Find`, `Update`, `Delete` and `Count` operations. The query specification starts with `By` or `All` word in the method name.

- `All` is used for querying all documents of the given type in the database. It is simple because only one word `All` is enough for repogen to understand. For example, `FindAll`, `UpdateCityAll` and `DeleteAll`.
- `By` is used for querying by a set of fields with specific operators. It is more complicated than `All` query but not be too difficult to understand. For example, `FindByGenderAndCity` and `DeleteByAgeGreaterThan`.

#### Specifying fields to query

You can write a query by specifying field name after `By` such as `ByID`. In case that you have multiple fields to query, you can connect field names with `And` or `Or` word such as `ByCityAndGender` and `ByCityOrGender`. `And` and `Or` operators are different in their meaning so the query result will be also different.

After specifying query to the method name, you also need to provide method parameters that match the given query fields. The example is given below:

```go
FindByCityAndGender(ctx context.Context, city string, gender Gender) ([]*UserModel, error)
```

Assuming that the `City` field in the `UserModel` struct is of type `string` and the `Gender` field in the `UserModel` struct is of custom type `Gender`, you have to provide `string` and `Gender` type parameters in the method.

#### Comparators to each field

When you specify the query like `ByAge`, it finds documents that contains age value **equal to** the provided parameter value. However, there are other types of comparators provided for you to use as follows.

| Keyword            | Meaning         | Sample                               |
|--------------------|-----------------|--------------------------------------|
| -                  | == $1           | `FindByUsername(ctx, $1)`            |
| `LessThan`         | < $1            | `FindByAgeLessThan(ctx, $1)`         |
| `LessThanEqual`    | <= $1           | `FindByAgeLessThanEqual(ctx, $1)`    |
| `GreaterThan`      | > $1            | `FindByAgeGreaterThan(ctx, $1)`      |
| `GreaterThanEqual` | >= $1           | `FindByAgeGreaterThanEqual(ctx, $1)` |
| `Between`          | >= $1 and <= $2 | `FindByAgeBetween(ctx, $1, $2)`      |
| `In`               | in slice $1     | `FindByCityIn(ctx, $1)`              |

To apply these comparators to the query, place these words after the field name such as `ByAgeGreaterThan`. You can also use comparators along with `And` and `Or` operators. For example, `ByGenderNotOrAgeLessThan` will apply `Not` comparator to the `Gender` field and `LessThan` comparator to the `Age` field.

`Between` and `In` comparators are special in terms of parameter requirements. `Between` needs two parameters to perform the query and `In` needs a slice instead of its raw type. The example is provided below:

```go
FindByAgeBetween(ctx context.Context, fromAge int, toAge int) ([]*UserModel, error)

FindByCityIn(ctx context.Context, cities []string) ([]*UserModel, error)
```

Assuming that the `Age` field in the `UserModel` struct is of type `int`, it requires that there must be two `int` parameters provided for `Age` field in the method. And assuming that the `City` field in the `UserModel` struct is of type `string`, it requires that the parameter that is provided to the query must be of slice type.

## License

Licensed under [MIT](https://github.com/sunboyy/repogen/blob/main/LICENSE)
