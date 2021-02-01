package spec

// ParsingError is an error from parsing interface methods
type ParsingError string

func (err ParsingError) Error() string {
	switch err {
	case UnknownOperationError:
		return "unknown operation"
	case UnsupportedNameError:
		return "method name is not supported"
	case UnsupportedReturnError:
		return "this type of return is not supported"
	case InvalidQueryError:
		return "invalid query"
	case InvalidParamError:
		return "parameters do not match the query"
	case InvalidUpdateFieldsError:
		return "update fields is invalid"
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
	UnsupportedReturnError    ParsingError = "ERROR_UNSUPPORTED_RETURN"
	InvalidQueryError         ParsingError = "ERROR_INVALID_QUERY"
	InvalidParamError         ParsingError = "ERROR_INVALID_PARAM"
	InvalidUpdateFieldsError  ParsingError = "ERROR_INVALID_UPDATE_FIELDS"
	ContextParamRequiredError ParsingError = "ERROR_CONTEXT_PARAM_REQUIRED"
	StructFieldNotFoundError  ParsingError = "ERROR_STRUCT_FIELD_NOT_FOUND"
)
