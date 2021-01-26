package mongo

// GenerationError is an error from generating MongoDB repository
type GenerationError string

func (err GenerationError) Error() string {
	switch err {
	case OperationNotSupportedError:
		return "operation not supported"
	case BsonTagNotFoundError:
		return "bson tag not found"
	}
	return string(err)
}

// generation error constants
const (
	OperationNotSupportedError GenerationError = "ERROR_OPERATION_NOT_SUPPORTED"
	BsonTagNotFoundError       GenerationError = "ERROR_BSON_TAG_NOT_FOUND"
)
