package spec

import "github.com/sunboyy/repogen/internal/code"

// UpdateOperation is a method specification for update operations
type UpdateOperation struct {
	Update Update
	Mode   QueryMode
	Query  QuerySpec
}

// Name returns "Update" operation name
func (o UpdateOperation) Name() string {
	return "Update"
}

// Update is an interface of update operation type
type Update interface {
	Name() string
	NumberOfArguments() int
}

// UpdateModel is a type of update operation that update the whole model
type UpdateModel struct {
}

// Name returns UpdateModel name 'Model'
func (u UpdateModel) Name() string {
	return "Model"
}

// NumberOfArguments returns 1
func (u UpdateModel) NumberOfArguments() int {
	return 1
}

// UpdateFields is a type of update operation that update specific fields
type UpdateFields []UpdateField

// Name returns UpdateFields name 'Fields'
func (u UpdateFields) Name() string {
	return "Fields"
}

// NumberOfArguments returns number of update fields
func (u UpdateFields) NumberOfArguments() int {
	return len(u)
}

// UpdateField stores mapping between field name in the model and the parameter
// index.
type UpdateField struct {
	FieldReference FieldReference
	ParamIndex     int
	Operator       UpdateOperator
}

// UpdateOperator is a custom type that declares update operator to be used in
// an update operation
type UpdateOperator string

// UpdateOperator constants
const (
	UpdateOperatorSet  UpdateOperator = "SET"
	UpdateOperatorPush UpdateOperator = "PUSH"
	UpdateOperatorInc  UpdateOperator = "INC"
)

// NumberOfArguments returns number of arguments required to perform an update operation
func (o UpdateOperator) NumberOfArguments() int {
	return 1
}

// ArgumentType returns type that is required for function parameter
func (o UpdateOperator) ArgumentType(fieldType code.Type) code.Type {
	switch o {
	case UpdateOperatorPush:
		arrayType := fieldType.(code.ArrayType)
		return arrayType.ContainedType
	default:
		return fieldType
	}
}

func (p interfaceMethodParser) parseUpdateOperation(tokens []string) (Operation, error) {
	mode, err := p.extractIntOrBoolReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	updateTokens, queryTokens := p.splitUpdateAndQueryTokens(tokens)

	update, err := p.parseUpdate(updateTokens)
	if err != nil {
		return nil, err
	}

	querySpec, err := p.parseQuery(queryTokens, 1+update.NumberOfArguments())
	if err != nil {
		return nil, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[update.NumberOfArguments()+1:], querySpec); err != nil {
		return nil, err
	}

	return UpdateOperation{
		Update: update,
		Mode:   mode,
		Query:  querySpec,
	}, nil
}

func (p interfaceMethodParser) parseUpdate(tokens []string) (Update, error) {
	if len(tokens) == 0 {
		requiredType := code.PointerType{ContainedType: p.StructModel.ReferencedType()}
		if len(p.Method.Params) <= 1 || p.Method.Params[1].Type != requiredType {
			return nil, ErrInvalidUpdateFields
		}
		return UpdateModel{}, nil
	}

	updateFieldTokens, ok := splitByAnd(tokens)
	if !ok {
		return nil, ErrInvalidUpdateFields
	}

	var updateFields UpdateFields

	paramIndex := 1
	for _, updateFieldToken := range updateFieldTokens {
		updateField, err := p.parseUpdateField(updateFieldToken, paramIndex)
		if err != nil {
			return nil, err
		}

		updateFields = append(updateFields, updateField)
		paramIndex += updateField.Operator.NumberOfArguments()
	}

	for _, field := range updateFields {
		if len(p.Method.Params) < field.ParamIndex+field.Operator.NumberOfArguments() {
			return nil, ErrInvalidUpdateFields
		}

		requiredType := field.Operator.ArgumentType(field.FieldReference.ReferencedField().Type)

		for i := 0; i < field.Operator.NumberOfArguments(); i++ {
			if requiredType != p.Method.Params[field.ParamIndex+i].Type {
				return nil, NewArgumentTypeNotMatchedError(field.FieldReference.ReferencingCode(), requiredType,
					p.Method.Params[field.ParamIndex+i].Type)
			}
		}
	}

	return updateFields, nil
}

func (p interfaceMethodParser) parseUpdateField(t []string,
	paramIndex int) (UpdateField, error) {

	if len(t) > 1 && t[len(t)-1] == "Push" {
		return p.createUpdateField(t[:len(t)-1], UpdateOperatorPush, paramIndex)
	}
	if len(t) > 1 && t[len(t)-1] == "Inc" {
		return p.createUpdateField(t[:len(t)-1], UpdateOperatorInc, paramIndex)
	}
	return p.createUpdateField(t, UpdateOperatorSet, paramIndex)
}

func (p interfaceMethodParser) createUpdateField(t []string,
	operator UpdateOperator, paramIndex int) (UpdateField, error) {

	fieldReference, ok := p.fieldResolver.ResolveStructField(p.StructModel, t)
	if !ok {
		return UpdateField{}, NewStructFieldNotFoundError(t)
	}

	if !p.validateUpdateOperator(fieldReference.ReferencedField().Type, operator) {
		return UpdateField{}, NewIncompatibleUpdateOperatorError(operator, fieldReference)
	}

	return UpdateField{
		FieldReference: fieldReference,
		ParamIndex:     paramIndex,
		Operator:       operator,
	}, nil
}

func (p interfaceMethodParser) validateUpdateOperator(referencedType code.Type, operator UpdateOperator) bool {
	switch operator {
	case UpdateOperatorPush:
		_, ok := referencedType.(code.ArrayType)
		return ok
	case UpdateOperatorInc:
		return referencedType.IsNumber()
	}
	return true
}
