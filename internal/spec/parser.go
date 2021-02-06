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
	methodNameTokens := camelcase.Split(p.Method.Name)
	switch methodNameTokens[0] {
	case "Insert":
		return p.parseInsertMethod(methodNameTokens[1:])
	case "Find":
		return p.parseFindMethod(methodNameTokens[1:])
	case "Update":
		return p.parseUpdateMethod(methodNameTokens[1:])
	case "Delete":
		return p.parseDeleteMethod(methodNameTokens[1:])
	case "Count":
		return p.parseCountMethod(methodNameTokens[1:])
	}
	return MethodSpec{}, UnknownOperationError
}

func (p interfaceMethodParser) parseInsertMethod(tokens []string) (MethodSpec, error) {
	mode, err := p.extractInsertReturns(p.Method.Returns)
	if err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateContextParam(); err != nil {
		return MethodSpec{}, err
	}

	pointerType := code.PointerType{ContainedType: p.StructModel.ReferencedType()}
	if mode == QueryModeOne && p.Method.Params[1].Type != pointerType {
		return MethodSpec{}, InvalidParamError
	}

	arrayType := code.ArrayType{ContainedType: pointerType}
	if mode == QueryModeMany && p.Method.Params[1].Type != arrayType {
		return MethodSpec{}, InvalidParamError
	}

	return p.createMethodSpec(InsertOperation{
		Mode: mode,
	}), nil
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

func (p interfaceMethodParser) parseFindMethod(tokens []string) (MethodSpec, error) {
	if len(tokens) == 0 {
		return MethodSpec{}, UnsupportedNameError
	}

	mode, err := p.extractModelOrSliceReturns(p.Method.Returns)
	if err != nil {
		return MethodSpec{}, err
	}

	querySpec, err := parseQuery(tokens, 1)
	if err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateContextParam(); err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[1:], querySpec); err != nil {
		return MethodSpec{}, err
	}

	return p.createMethodSpec(FindOperation{
		Mode:  mode,
		Query: querySpec,
	}), nil
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

func (p interfaceMethodParser) parseUpdateMethod(tokens []string) (MethodSpec, error) {
	if len(tokens) == 0 {
		return MethodSpec{}, UnsupportedNameError
	}

	mode, err := p.extractIntOrBoolReturns(p.Method.Returns)
	if err != nil {
		return MethodSpec{}, err
	}

	paramIndex := 1
	var fields []UpdateField
	var aggregatedToken string
	for i, token := range tokens {
		if token == "By" || token == "All" {
			tokens = tokens[i:]
			break
		} else if token != "And" {
			aggregatedToken += token
		} else if len(aggregatedToken) == 0 {
			return MethodSpec{}, InvalidUpdateFieldsError
		} else {
			fields = append(fields, UpdateField{Name: aggregatedToken, ParamIndex: paramIndex})
			paramIndex++
			aggregatedToken = ""
		}
	}
	if len(aggregatedToken) == 0 {
		return MethodSpec{}, InvalidUpdateFieldsError
	}
	fields = append(fields, UpdateField{Name: aggregatedToken, ParamIndex: paramIndex})

	querySpec, err := parseQuery(tokens, 1+len(fields))
	if err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateContextParam(); err != nil {
		return MethodSpec{}, err
	}

	for _, field := range fields {
		structField, ok := p.StructModel.Fields.ByName(field.Name)
		if !ok {
			return MethodSpec{}, StructFieldNotFoundError
		}

		if structField.Type != p.Method.Params[field.ParamIndex].Type {
			return MethodSpec{}, InvalidParamError
		}
	}

	if err := p.validateQueryFromParams(p.Method.Params[len(fields)+1:], querySpec); err != nil {
		return MethodSpec{}, err
	}

	return p.createMethodSpec(UpdateOperation{
		Fields: fields,
		Mode:   mode,
		Query:  querySpec,
	}), nil
}

func (p interfaceMethodParser) parseDeleteMethod(tokens []string) (MethodSpec, error) {
	if len(tokens) == 0 {
		return MethodSpec{}, UnsupportedNameError
	}

	mode, err := p.extractIntOrBoolReturns(p.Method.Returns)
	if err != nil {
		return MethodSpec{}, err
	}

	querySpec, err := parseQuery(tokens, 1)
	if err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateContextParam(); err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[1:], querySpec); err != nil {
		return MethodSpec{}, err
	}

	return p.createMethodSpec(DeleteOperation{
		Mode:  mode,
		Query: querySpec,
	}), nil
}

func (p interfaceMethodParser) parseCountMethod(tokens []string) (MethodSpec, error) {
	if len(tokens) == 0 {
		return MethodSpec{}, UnsupportedNameError
	}

	if err := p.validateCountReturns(p.Method.Returns); err != nil {
		return MethodSpec{}, err
	}

	querySpec, err := parseQuery(tokens, 1)
	if err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateContextParam(); err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateQueryFromParams(p.Method.Params[1:], querySpec); err != nil {
		return MethodSpec{}, err
	}

	return p.createMethodSpec(CountOperation{
		Query: querySpec,
	}), nil
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
			return StructFieldNotFoundError
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

func (p interfaceMethodParser) createMethodSpec(operation Operation) MethodSpec {
	return MethodSpec{
		Name:      p.Method.Name,
		Params:    p.Method.Params,
		Returns:   p.Method.Returns,
		Operation: operation,
	}
}
