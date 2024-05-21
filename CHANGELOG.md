# Changelog

## Unreleased

### Breaking Changes

- With the integration to the go/packages library, `-pkg` option now requires a **Go-style package path** instead of a path to the directory. If you want to use a relative path, you need to prefix it with `./` or `../`.

### Added

- Now supports interface embedding
- `-model-pkg` option to specify the package containing the model struct
- `-dest-pkg` option to specify the package path to write the generated code to

### Changed

- The constructor of the generated repository implementation now returns a pointer to the generated repository implementation instead of the repository interface.
- `-pkg` option now requires a **Go-style package path** instead of a path to the directory. Without specifying a dot (`./` or `../`) at the beginning of the path, Go will assume the absolute package path.

## v0.3.1 - 2024-01-14

### Changed

- Make find many operations return empty slice instead of nil to provide consistency when further encoded to JSON.

## v0.3.0 - 2023-10-31

### Breaking Changes

- `-src` option is now removed, replacing with `-pkg` option. Repogen now supports when model struct and repository interface is placed in separate files in the same package. Instead of specifying a single Go file to scan the model struct and repository interface, you can now specify a path to the package containing  with `-pkg`. If the `-pkg` option is not specified, it will default to the current directory.
- The keyword `One` after `Find` is not allowed anymore. If you have `FindOne` declared anywhere, replace them with `Find`.

### Added

- Package scanning: Instead of specifying a single Go file to scan the model struct and repository interface using `-src` option, you can now specify a path to the package containing the model struct and repository interface using `-pkg` option.
- Find many operations with limits: Write `TopN` keyword after `Find` to allow finding with limit N: e.g. `FindTop5AllOrderByAge`.
- Query comparators: `Exists` and `NotExists`

### Removed

- `-src` option is removed, in favor of `-pkg` option.

### Fixed

- Fixed decoding error when the struct tag has space characters

## v0.2.1 - 2021-11-27

### Fixed

- Fixed error when any method parameter is assigned with map type

## v0.2.0 - 2021-06-05

### Added

- Query operators: `NotIn`, `True` and `False`
- Update the whole model by not specifying queries: e.g. `UpdateByID`
- Sorting in find operations (both find one and find many)
- Deeply referencing to struct fields: e.g. ContactEmail = contact.email
- Update operators: `Push` and `Inc`

### Improved

- Improved error message from generation to help investigate repository interface errors

## v0.1.0 - 2021-02-16

Initial release

- Supports create, read, update, delete and count functionality
- Supports single-entity and many-entity operations
- Supports AND and OR query
- Supports many comparators (equal, less than, less than or equal, greater than, greater than or equal, not equal, between, in)
- Validate method signature depending on method operations
