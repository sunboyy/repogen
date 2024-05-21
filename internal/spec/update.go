package spec

import "go/types"

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
func (o UpdateOperator) ArgumentType(fieldType types.Type) types.Type {
	switch o {
	case UpdateOperatorPush:
		sliceType := fieldType.(*types.Slice)
		return sliceType.Elem()
	default:
		return fieldType
	}
}

func (p interfaceMethodParser) parseUpdateOperation(tokens []string) (Operation, error) {
	mode, err := p.extractIntOrBoolReturns(p.Signature.Results())
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

	if err := p.validateQueryFromParams(p.Signature.Params(), 1+update.NumberOfArguments(), querySpec); err != nil {
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
		expectedType := types.NewPointer(p.NamedStruct)
		if p.Signature.Params().Len() <= 1 || !types.Identical(p.Signature.Params().At(1).Type(), expectedType) {
			return nil, ErrInvalidUpdateFields
		}
		return UpdateModel{}, nil
	}

	updateFields, err := p.parseUpdateFieldsFromTokens(tokens)
	if err != nil {
		return nil, err
	}

	if err := p.validateUpdateFieldsWithParams(updateFields); err != nil {
		return nil, err
	}

	return updateFields, nil
}

func (p interfaceMethodParser) parseUpdateFieldsFromTokens(tokens []string) (UpdateFields, error) {
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

	fieldReference, ok := resolveStructField(p.UnderlyingStruct, t)
	if !ok {
		return UpdateField{}, NewStructFieldNotFoundError(t)
	}

	if !p.validateUpdateOperator(fieldReference.ReferencedField().Var.Type(), operator) {
		return UpdateField{}, NewIncompatibleUpdateOperatorError(operator, fieldReference)
	}

	return UpdateField{
		FieldReference: fieldReference,
		ParamIndex:     paramIndex,
		Operator:       operator,
	}, nil
}

func (p interfaceMethodParser) validateUpdateOperator(referencedType types.Type, operator UpdateOperator) bool {
	switch operator {
	case UpdateOperatorPush:
		_, ok := referencedType.(*types.Slice)
		return ok

	case UpdateOperatorInc:
		switch t := referencedType.(type) {
		case *types.Basic:
			return t.Info()&types.IsNumeric != 0

		case *types.Pointer:
			return p.validateUpdateOperator(t.Elem(), operator)

		case *types.Named:
			return p.validateUpdateOperator(t.Underlying(), operator)

		default:
			return false
		}
	}
	return true
}

func (p interfaceMethodParser) validateUpdateFieldsWithParams(updateFields UpdateFields) error {
	for _, field := range updateFields {
		if p.Signature.Params().Len() < field.ParamIndex+field.Operator.NumberOfArguments() {
			return ErrInvalidUpdateFields
		}

		expectedType := field.Operator.ArgumentType(field.FieldReference.ReferencedField().Var.Type())

		for i := 0; i < field.Operator.NumberOfArguments(); i++ {
			if !types.Identical(p.Signature.Params().At(field.ParamIndex+i).Type(), expectedType) {
				return NewArgumentTypeNotMatchedError(field.FieldReference.ReferencingCode(), expectedType,
					p.Signature.Params().At(field.ParamIndex+i).Type())
			}
		}
	}

	return nil
}
