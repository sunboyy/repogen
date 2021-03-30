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

// UpdateField stores mapping between field name in the model and the parameter index
type UpdateField struct {
	FieldReference FieldReference
	ParamIndex     int
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
			return nil, InvalidUpdateFieldsError
		}
		return UpdateModel{}, nil
	}

	updateFieldTokens, ok := splitByAnd(tokens)
	if !ok {
		return nil, InvalidUpdateFieldsError
	}

	var updateFields UpdateFields

	paramIndex := 1
	for _, updateFieldToken := range updateFieldTokens {
		updateFieldReference, ok := p.fieldResolver.ResolveStructField(p.StructModel, updateFieldToken)
		if !ok {
			return nil, NewStructFieldNotFoundError(updateFieldToken)
		}

		updateFields = append(updateFields, UpdateField{
			FieldReference: updateFieldReference,
			ParamIndex:     paramIndex,
		})
		paramIndex++
	}

	for _, field := range updateFields {
		if len(p.Method.Params) <= field.ParamIndex ||
			field.FieldReference.ReferencedField().Type != p.Method.Params[field.ParamIndex].Type {
			return nil, InvalidUpdateFieldsError
		}
	}

	return updateFields, nil
}
