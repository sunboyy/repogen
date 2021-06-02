# repogen

<a href="https://github.com/sunboyy/repogen/actions?query=workflow%3Abuild" target="_blank">
    <img src="https://github.com/sunboyy/repogen/workflows/build/badge.svg" alt="build status badge">
</a>
<a href="https://codecov.io/gh/sunboyy/repogen" target="_blank">
    <img src="https://codecov.io/gh/sunboyy/repogen/branch/main/graph/badge.svg?token=9BD5Y8X7NO"/>
</a>
<a href="https://codeclimate.com/github/sunboyy/repogen/maintainability" target="_blank">
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

Run the repogen to generate a repository implementation from the interface. The following command is an example to generate `UserRepository` interface implementation defined in `examples/getting-started/user.go` to the destination file `examples/getting-started/user_repo.go`. See [Usage](#Usage) section below for more detailed information.

```
$ repogen -src=examples/getting-started/user.go -dest=examples/getting-started/user_repo.go \
        -model=UserModel -repo=UserRepository
```

You can also write the above command in the `go:generate` format inside Go files in order to generate the implementation when `go generate` command is executed.

## Usage

### Running Options

The `repogen` command is used to generate source code for a given Go file containing repository interface to be implemented. Run `repogen -h` to see all available options while the necessary options for code generation are described as follows:

- `-src`: A Go file containing struct model and repository interface to be implemented
- `-dest`: A file to which to write the resulting source code. If not specified, the source code will be printed to the standard output.
- `-model`: The name of the base struct model that represents the data stored in MongoDB for a specific collection.
- `-repo`: The name of the repository interface that you want to be implemented.

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

As the `Insert` operation has a limited use case, we do not want to limit you on how you name your method. Any method that has the name starting with the word `Insert` is always valid. For example, you can name your method `InsertAWholeBunchOfDocuments` and it will work as long as you specify method parameters and return types correctly.

#### Find operation

A `Find` operation also has two modes like `Insert` operation: single-entity and multiple-entity. However `Find` operation can be very simple or complex depending on how complex the query is. In this section, we will show you how to write single-modes and multiple-entity modes of find method with a simple query. For more information about more complex queries, please consult the [query specification](#query-specification) section in this document.

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

Find operation also supports sorting results for both single-entity and multiple-entity operations. To sort, append the existing method name with `OrderBy`, followed by the field names to sort. The order will be default to ascending. In case that you want descending order, write `Desc` after the field name. For example:

```go
// This will sort results by age ascendingly
FindByCityOrderByAge(ctx context.Context, city string) ([]*Model, error)

// This will also sort results by age ascendingly
FindByCityOrderByAgeAsc(ctx context.Context, city string) ([]*Model, error)

// This will sort results by age descendingly
FindByCityOrderByAgeDesc(ctx context.Context, city string) ([]*Model, error)
```

#### Update operation

An `Update` operation also has single-entity and multiple-entity operations. An `Update` operation also supports querying like `Find` operation. Specifying the query is the same as in `Find` method. However, an `Update` operation requires more parameters than `Find` method depending on update type. There are two update types provided.

1. Model-type update

This type of update is for changing the whole model, replacing all the fields except the field with bson `omitempty` tag when the value is not provided. To write this type of update, write `Update` followed by query like in find method.

```go
// UpdateByID updates a single document by ID
UpdateByID(ctx context.Context, model *Model, id primitive.ObjectID) (bool, error)
```

2. Fields-type update

This type of update is for changing only some fields in the model. To write this type of update, specify the fields to update explicitly after `Update`. Updating multiple fields are allowed by concatinating each field name with `And` as follows:

```go
// UpdateDisplayNameAndCityByID updates a single document with a new display name and
// city by ID
UpdateDisplayNameAndCityByID(ctx context.Context, displayName string, city string,
	id primitive.ObjectID) (bool, error)

// UpdateGenderByCity updates Gender field of documents with matching city parameter
UpdateGenderByCity(ctx context.Context, gender Gender, city string) (int, error)
```

The update operator will be default to `$set` operator. In case that you want to use other operators, you can append the update field by the keyword that specifies the update operator. The current supported ones other than `$set` are `$push` and `$inc`. Write `Push` or `Inc` after the field name to apply those operators. Keep in mind that different operators requires different type of arguments such as an array-type for `Push` and a number-type for `Inc`.

```go
// UpdateConsentHistoryPushByID appends consentHistory to the ConsentHistory field
UpdateConsentHistoryPushByID(ctx context.Context, consentHistory ConsentHistory,
	id primitive.ObjectID) (bool, error)

// UpdateAgeIncByID increments age value by `incAge`
UpdateAgeIncByID(ctx context.Context, incAge int, id primitive.ObjectID) (bool, error)
```

For all types of updates, repogen determines a single-entity operation or a multiple-entity by checking the first return value. If it is of type `bool`, the method will be single-entity operation. If it is of type `int`, the method will be multiple-entity operation. For single-entity operation, the method returns true if there is a matching document. For multiple-entity operation, the integer return shows the number of matched documents.

The requirement of the `Update` operation method is that there must be only two return values, the second return value must be of type `error` and the first method parameter must be of type `context.Context`. The requirement of number of method parameters depends on the update operation and the query.

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

A `Count` operation is also similar to `Find` operation except it has only multiple-entity mode and does not support sorting. This means that the method returns are always the same for any count operations. The method name pattern and the parameters are the same as `Find` operation without the sort parameters.

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

When you specify the query like `ByAge`, it finds documents that contains age value **equal to** the provided parameter value. However, there are other types of comparators supported in the following table:

| Keyword            | Meaning         | Sample                               |
|--------------------|-----------------|--------------------------------------|
| -                  | == $1           | `FindByUsername(ctx, $1)`            |
| `LessThan`         | < $1            | `FindByAgeLessThan(ctx, $1)`         |
| `LessThanEqual`    | <= $1           | `FindByAgeLessThanEqual(ctx, $1)`    |
| `GreaterThan`      | > $1            | `FindByAgeGreaterThan(ctx, $1)`      |
| `GreaterThanEqual` | >= $1           | `FindByAgeGreaterThanEqual(ctx, $1)` |
| `Between`          | >= $1 and <= $2 | `FindByAgeBetween(ctx, $1, $2)`      |
| `In`               | in slice $1     | `FindByCityIn(ctx, $1)`              |
| `NotIn`            | not in slice $1 | `FindByCityNotIn(ctx, $1)`           |
| `True`             | == `true`       | `FindByEnabledTrue(ctx)`             |
| `False`            | == `false`      | `FindByEnabledFalse(ctx)`            |

To apply these comparators to the query, place the keyword after the field name such as `ByAgeGreaterThan`. You can also use comparators along with `And` and `Or` operators. For example, `ByGenderNotOrAgeLessThan` will apply `Not` comparator to the `Gender` field and `LessThan` comparator to the `Age` field.

`Between`, `In`, `NotIn`, `True` and `False` comparators are special in terms of parameter requirements. `Between` needs two parameters to perform the query, `In` and `NotIn` needs a slice instead of its raw type and `True` and `False` doesn't need any parameter. The example is provided below:

```go
FindByAgeBetween(ctx context.Context, fromAge int, toAge int) ([]*UserModel, error)

FindByCityIn(ctx context.Context, cities []string) ([]*UserModel, error)
FindByCityNotIn(ctx context.Context, cities []string) ([]*UserModel, error)

FindByEnabledTrue(ctx context.Context) ([]*UserModel, error)
FindByEnabledFalse(ctx context.Context) ([]*UserModel, error)
```

Assuming that the `Age` field in the `UserModel` struct is of type `int`, it requires that there must be two `int` parameters provided for `Age` field in the method. And assuming that the `City` field in the `UserModel` struct is of type `string`, it requires that the parameter that is provided to the query must be of slice type.

### Field Referencing

To query, update or sort, you have to specify struct fields that you want to use. Repogen determines struct field by the field name. For example, the method name `FindByPhoneNumber` refer to the field named `PhoneNumber`. Repogen tries to find the properties of the struct field named `PhoneNumber` for further processing.

In the real world applications, you might want to refer to a struct field of another struct that is used in the base model, for example:

```go
type ContactModel struct {
	Phone string `bson:"phone"`
	Email string `bson:"email"`
}

type UserModel struct {
	Contact ContactModel `bson:"contact"`
}
```

If you want to deeply refer to a struct field, you can do it by concatenating the field names that is the reference to the field you want. For example:

- To find by phone number, write `FindByContactPhone`
- To update phone number by ID, write `UpdateContactPhoneByID`
- To find all and sort results by phone number, write `FindAllOrderByContactPhone`

Deep referencing is supported for query fields, sort fields and update fields. However, the `inline` option for bson struct tag is not currently supported.

## License

Licensed under [MIT](https://github.com/sunboyy/repogen/blob/main/LICENSE)
