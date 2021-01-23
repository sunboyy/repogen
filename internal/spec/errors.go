package spec

// ParsingError is an error from parsing interface methods
type ParsingError string

func (err ParsingError) Error() string {
	switch err {
	case UnknownOperationError:
		return "unknown operation"
	case UnsupportedNameError:
		return "method name is not supported"
	case InvalidQueryError:
		return "invalid query"
	case InvalidParamError:
		return "parameters do not match the query"
	case UnsupportedReturnError:
		return "this type of return is not supported"
	case ContextParamRequiredError:
		return "context parameter is required"
	case StructFieldNotFoundError:
		return "struct field not found"
	}
	return string(err)
}

// parsing error constants
const (
	UnknownOperationError     ParsingError = "ERROR_UNKNOWN_OPERATION"
	UnsupportedNameError      ParsingError = "ERROR_UNSUPPORTED"
	InvalidQueryError         ParsingError = "ERROR_INVALID_QUERY"
	InvalidParamError         ParsingError = "ERROR_INVALID_PARAM"
	UnsupportedReturnError    ParsingError = "ERROR_INVALID_RETURN"
	ContextParamRequiredError ParsingError = "ERROR_CONTEXT_PARAM_REQUIRED"
	StructFieldNotFoundError  ParsingError = "ERROR_STRUCT_FIELD_NOT_FOUND"
)