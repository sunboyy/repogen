package spec

import (
	"github.com/fatih/camelcase"
	"github.com/sunboyy/repogen/internal/code"
)

// ParseInterfaceMethod returns repository method spec from declared interface method
func ParseInterfaceMethod(structModel code.Struct, method code.Method) (MethodSpec, error) {
	parser := interfaceMethodParser{
		StructModel: structModel,
		Method:      method,
	}

	return parser.Parse()
}

type interfaceMethodParser struct {
	StructModel code.Struct
	Method      code.Method
}

func (p interfaceMethodParser) Parse() (MethodSpec, error) {
	operation, err := p.parseMethod()
	if err != nil {
		return MethodSpec{}, err
	}

	return MethodSpec{
		Name:      p.Method.Name,
		Params:    p.Method.Params,
		Returns:   p.Method.Returns,
		Operation: operation,
	}, nil
}

func (p interfaceMethodParser) parseMethod() (Operation, error) {
	methodNameTokens := camelcase.Split(p.Method.Name)
	switch methodNameTokens[0] {
	case "Insert":
		return p.parseInsertOperation(methodNameTokens[1:])
	case "Find":
		return p.parseFindOperation(methodNameTokens[1:])
	case "Update":
		return p.parseUpdateOperation(methodNameTokens[1:])
	case "Delete":
		return p.parseDeleteOperation(methodNameTokens[1:])
	case "Count":
		return p.parseCountOperation(methodNameTokens[1:])
	}
	return nil, NewUnknownOperationError(methodNameTokens[0])
}

func (p interfaceMethodParser) parseInsertOperation(tokens []string) (Operation, error) {
	mode, err := p.extractInsertReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	pointerType := code.PointerType{ContainedType: p.StructModel.ReferencedType()}
	if mode == QueryModeOne && p.Method.Params[1].Type != pointerType {
		return nil, InvalidParamError
	}

	arrayType := code.ArrayType{ContainedType: pointerType}
	if mode == QueryModeMany && p.Method.Params[1].Type != arrayType {
		return nil, InvalidParamError
	}

	return InsertOperation{
		Mode: mode,
	}, nil
}

func (p interfaceMethodParser) extractInsertReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", UnsupportedReturnError
	}

	if returns[1] != code.SimpleType("error") {
		return "", UnsupportedReturnError
	}

	interfaceType, ok := returns[0].(code.InterfaceType)
	if ok {
		if len(interfaceType.Methods) != 0 {
			return "", UnsupportedReturnError
		}
		return QueryModeOne, nil
	}

	arrayType, ok := returns[0].(code.ArrayType)
	if ok {
		interfaceType, ok := arrayType.ContainedType.(code.InterfaceType)
		if !ok || len(interfaceType.Methods) != 0 {
			return "", UnsupportedReturnError
		}
		return QueryModeMany, nil
	}

	return "", UnsupportedReturnError
}

func (p interfaceMethodParser) parseFindOperation(tokens []string) (Operation, error) {
	mode, err := p.extractModelOrSliceReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	querySpec, err := parseQuery(tokens, 1)
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[1:], querySpec); err != nil {
		return nil, err
	}

	return FindOperation{
		Mode:  mode,
		Query: querySpec,
	}, nil
}

func (p interfaceMethodParser) extractModelOrSliceReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", UnsupportedReturnError
	}

	if returns[1] != code.SimpleType("error") {
		return "", UnsupportedReturnError
	}

	pointerType, ok := returns[0].(code.PointerType)
	if ok {
		simpleType := pointerType.ContainedType
		if simpleType == code.SimpleType(p.StructModel.Name) {
			return QueryModeOne, nil
		}
		return "", UnsupportedReturnError
	}

	arrayType, ok := returns[0].(code.ArrayType)
	if ok {
		pointerType, ok := arrayType.ContainedType.(code.PointerType)
		if ok {
			simpleType := pointerType.ContainedType
			if simpleType == code.SimpleType(p.StructModel.Name) {
				return QueryModeMany, nil
			}
			return "", UnsupportedReturnError
		}
	}

	return "", UnsupportedReturnError
}

func (p interfaceMethodParser) parseUpdateOperation(tokens []string) (Operation, error) {
	mode, err := p.extractIntOrBoolReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	updateFieldTokens, queryTokens := p.splitUpdateFieldAndQueryTokens(tokens)

	paramIndex := 1
	var fields []UpdateField
	var aggregatedToken string
	for _, token := range updateFieldTokens {
		if token != "And" {
			aggregatedToken += token
		} else if len(aggregatedToken) == 0 {
			return nil, InvalidUpdateFieldsError
		} else {
			fields = append(fields, UpdateField{Name: aggregatedToken, ParamIndex: paramIndex})
			paramIndex++
			aggregatedToken = ""
		}
	}
	if len(aggregatedToken) == 0 {
		return nil, InvalidUpdateFieldsError
	}
	fields = append(fields, UpdateField{Name: aggregatedToken, ParamIndex: paramIndex})

	querySpec, err := parseQuery(queryTokens, 1+len(fields))
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	for _, field := range fields {
		structField, ok := p.StructModel.Fields.ByName(field.Name)
		if !ok {
			return nil, NewStructFieldNotFoundError(field.Name)
		}

		if structField.Type != p.Method.Params[field.ParamIndex].Type {
			return nil, InvalidParamError
		}
	}

	if err := p.validateQueryFromParams(p.Method.Params[len(fields)+1:], querySpec); err != nil {
		return nil, err
	}

	return UpdateOperation{
		Fields: fields,
		Mode:   mode,
		Query:  querySpec,
	}, nil
}

func (p interfaceMethodParser) splitUpdateFieldAndQueryTokens(tokens []string) ([]string, []string) {
	var updateFieldTokens []string
	var queryTokens []string

	for i, token := range tokens {
		if token == "By" || token == "All" {
			queryTokens = tokens[i:]
			break
		} else {
			updateFieldTokens = append(updateFieldTokens, token)
		}
	}

	return updateFieldTokens, queryTokens
}

func (p interfaceMethodParser) parseDeleteOperation(tokens []string) (Operation, error) {
	mode, err := p.extractIntOrBoolReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	querySpec, err := parseQuery(tokens, 1)
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[1:], querySpec); err != nil {
		return nil, err
	}

	return DeleteOperation{
		Mode:  mode,
		Query: querySpec,
	}, nil
}

func (p interfaceMethodParser) parseCountOperation(tokens []string) (Operation, error) {
	if err := p.validateCountReturns(p.Method.Returns); err != nil {
		return nil, err
	}

	querySpec, err := parseQuery(tokens, 1)
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[1:], querySpec); err != nil {
		return nil, err
	}

	return CountOperation{
		Query: querySpec,
	}, nil
}

func (p interfaceMethodParser) validateCountReturns(returns []code.Type) error {
	if len(returns) != 2 {
		return UnsupportedReturnError
	}

	if returns[0] != code.SimpleType("int") {
		return UnsupportedReturnError
	}

	if returns[1] != code.SimpleType("error") {
		return UnsupportedReturnError
	}

	return nil
}

func (p interfaceMethodParser) extractIntOrBoolReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", UnsupportedReturnError
	}

	if returns[1] != code.SimpleType("error") {
		return "", UnsupportedReturnError
	}

	simpleType, ok := returns[0].(code.SimpleType)
	if ok {
		if simpleType == code.SimpleType("bool") {
			return QueryModeOne, nil
		}
		if simpleType == code.SimpleType("int") {
			return QueryModeMany, nil
		}
	}

	return "", UnsupportedReturnError
}

func (p interfaceMethodParser) validateContextParam() error {
	contextType := code.ExternalType{PackageAlias: "context", Name: "Context"}
	if len(p.Method.Params) == 0 || p.Method.Params[0].Type != contextType {
		return ContextParamRequiredError
	}
	return nil
}

func (p interfaceMethodParser) validateQueryFromParams(params []code.Param, querySpec QuerySpec) error {
	if querySpec.NumberOfArguments() != len(params) {
		return InvalidParamError
	}

	var currentParamIndex int
	for _, predicate := range querySpec.Predicates {
		structField, ok := p.StructModel.Fields.ByName(predicate.Field)
		if !ok {
			return NewStructFieldNotFoundError(predicate.Field)
		}

		for i := 0; i < predicate.Comparator.NumberOfArguments(); i++ {
			if params[currentParamIndex].Type != predicate.Comparator.ArgumentTypeFromFieldType(
				structField.Type) {
				return InvalidParamError
			}
			currentParamIndex++
		}
	}

	return nil
}
