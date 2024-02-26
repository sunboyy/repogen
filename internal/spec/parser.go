package spec

import (
	"strconv"

	"github.com/fatih/camelcase"
	"github.com/sunboyy/repogen/internal/code"
)

// ParseInterfaceMethod returns repository method spec from declared interface
// method.
func ParseInterfaceMethod(structs map[string]code.Struct, structModel code.Struct,
	method code.Method) (MethodSpec, error) {

	parser := interfaceMethodParser{
		fieldResolver: fieldResolver{
			Structs: structs,
		},
		StructModel: structModel,
		Method:      method,
	}

	return parser.Parse()
}

type interfaceMethodParser struct {
	fieldResolver fieldResolver
	StructModel   code.Struct
	Method        code.Method
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
		return nil, ErrInvalidParam
	}

	arrayType := code.ArrayType{ContainedType: pointerType}
	if mode == QueryModeMany && p.Method.Params[1].Type != arrayType {
		return nil, ErrInvalidParam
	}

	return InsertOperation{
		Mode: mode,
	}, nil
}

func (p interfaceMethodParser) extractInsertReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", NewOperationReturnCountUnmatchedError(2)
	}

	if returns[1] != code.TypeError {
		return "", NewUnsupportedReturnError(returns[1], 1)
	}

	switch t := returns[0].(type) {
	case code.InterfaceType:
		if len(t.Methods) == 0 {
			return QueryModeOne, nil
		}

	case code.ArrayType:
		interfaceType, ok := t.ContainedType.(code.InterfaceType)
		if ok && len(interfaceType.Methods) == 0 {
			return QueryModeMany, nil
		}
	}

	return "", NewUnsupportedReturnError(returns[0], 0)
}

func (p interfaceMethodParser) parseFindOperation(tokens []string) (Operation, error) {
	mode, err := p.extractModelOrSliceReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	limit, tokens, err := p.parseFindTop(tokens)
	if err != nil {
		return nil, err
	}
	if mode == QueryModeOne && limit != 0 {
		return nil, ErrLimitOnFindOne
	}

	queryTokens, sortTokens := p.splitQueryAndSortTokens(tokens)

	querySpec, err := p.parseQuery(queryTokens, 1)
	if err != nil {
		return nil, err
	}

	sorts, err := p.parseSort(sortTokens)
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
		Sorts: sorts,
		Limit: limit,
	}, nil
}

func (p interfaceMethodParser) parseFindTop(tokens []string) (int, []string,
	error) {

	if len(tokens) >= 1 && tokens[0] == "Top" {
		if len(tokens) < 2 {
			return 0, nil, ErrLimitAmountRequired
		}

		limit, err := strconv.Atoi(tokens[1])
		if err != nil {
			return 0, nil, ErrLimitAmountRequired
		}

		if limit <= 0 {
			return 0, nil, ErrLimitNonPositive
		}
		return limit, tokens[2:], nil
	}

	return 0, tokens, nil
}

func (p interfaceMethodParser) parseSort(rawTokens []string) ([]Sort, error) {
	if len(rawTokens) == 0 {
		return nil, nil
	}

	sortTokens, ok := splitByAnd(rawTokens[2:])
	if !ok {
		return nil, NewInvalidSortError(rawTokens)
	}

	var sorts []Sort
	for _, token := range sortTokens {
		sort, err := p.parseSortToken(token)
		if err != nil {
			return nil, err
		}
		sorts = append(sorts, sort)
	}

	return sorts, nil
}

func (p interfaceMethodParser) parseSortToken(t []string) (Sort, error) {
	if len(t) > 1 && t[len(t)-1] == "Asc" {
		return p.createSort(t[:len(t)-1], OrderingAscending)
	}
	if len(t) > 1 && t[len(t)-1] == "Desc" {
		return p.createSort(t[:len(t)-1], OrderingDescending)
	}
	return p.createSort(t, OrderingAscending)
}

func (p interfaceMethodParser) createSort(t []string, ordering Ordering) (Sort, error) {
	fields, ok := p.fieldResolver.ResolveStructField(p.StructModel, t)
	if !ok {
		return Sort{}, NewStructFieldNotFoundError(t)
	}

	return Sort{
		FieldReference: fields,
		Ordering:       ordering,
	}, nil
}

func (p interfaceMethodParser) splitQueryAndSortTokens(tokens []string) ([]string, []string) {
	var queryTokens []string
	var sortTokens []string

	for i, token := range tokens {
		if len(tokens) > i && token == "Order" && tokens[i+1] == "By" {
			sortTokens = tokens[i:]
			break
		} else {
			queryTokens = append(queryTokens, token)
		}
	}

	return queryTokens, sortTokens
}

func (p interfaceMethodParser) extractModelOrSliceReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", NewOperationReturnCountUnmatchedError(2)
	}

	if returns[1] != code.TypeError {
		return "", NewUnsupportedReturnError(returns[1], 1)
	}

	switch t := returns[0].(type) {
	case code.PointerType:
		pointerType := code.PointerType{ContainedType: p.StructModel.ReferencedType()}
		if t == pointerType {
			return QueryModeOne, nil
		}

	case code.ArrayType:
		pointerType := code.PointerType{ContainedType: p.StructModel.ReferencedType()}
		arrayType := code.ArrayType{ContainedType: pointerType}
		if t == arrayType {
			return QueryModeMany, nil
		}
	}

	return "", NewUnsupportedReturnError(returns[0], 0)
}

func splitByAnd(tokens []string) ([][]string, bool) {
	var updateFieldTokens [][]string
	var aggregatedToken []string

	for _, token := range tokens {
		if token != "And" {
			aggregatedToken = append(aggregatedToken, token)
		} else if len(aggregatedToken) == 0 {
			return nil, false
		} else {
			updateFieldTokens = append(updateFieldTokens, aggregatedToken)
			aggregatedToken = nil
		}
	}
	if len(aggregatedToken) == 0 {
		return nil, false
	}
	updateFieldTokens = append(updateFieldTokens, aggregatedToken)

	return updateFieldTokens, true
}

func (p interfaceMethodParser) splitUpdateAndQueryTokens(tokens []string) ([]string, []string) {
	var updateTokens []string
	var queryTokens []string

	for i, token := range tokens {
		if token == "By" || token == "All" {
			queryTokens = tokens[i:]
			break
		} else {
			updateTokens = append(updateTokens, token)
		}
	}

	return updateTokens, queryTokens
}

func (p interfaceMethodParser) parseDeleteOperation(tokens []string) (Operation, error) {
	mode, err := p.extractIntOrBoolReturns(p.Method.Returns)
	if err != nil {
		return nil, err
	}

	querySpec, err := p.parseQuery(tokens, 1)
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

	querySpec, err := p.parseQuery(tokens, 1)
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
		return NewOperationReturnCountUnmatchedError(2)
	}

	if returns[0] != code.TypeInt {
		return NewUnsupportedReturnError(returns[0], 0)
	}

	if returns[1] != code.TypeError {
		return NewUnsupportedReturnError(returns[1], 1)
	}

	return nil
}

func (p interfaceMethodParser) extractIntOrBoolReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", NewOperationReturnCountUnmatchedError(2)
	}

	if returns[1] != code.TypeError {
		return "", NewUnsupportedReturnError(returns[1], 1)
	}

	simpleType, ok := returns[0].(code.SimpleType)
	if ok {
		if simpleType == code.TypeBool {
			return QueryModeOne, nil
		}
		if simpleType == code.TypeInt {
			return QueryModeMany, nil
		}
	}

	return "", NewUnsupportedReturnError(returns[0], 0)
}

func (p interfaceMethodParser) validateContextParam() error {
	contextType := code.ExternalType{PackageAlias: "context", Name: "Context"}
	if len(p.Method.Params) == 0 || p.Method.Params[0].Type != contextType {
		return ErrContextParamRequired
	}
	return nil
}

func (p interfaceMethodParser) validateQueryFromParams(params []code.Param, querySpec QuerySpec) error {
	if querySpec.NumberOfArguments() != len(params) {
		return ErrInvalidParam
	}

	var currentParamIndex int
	for _, predicate := range querySpec.Predicates {
		if (predicate.Comparator == ComparatorTrue || predicate.Comparator == ComparatorFalse) &&
			predicate.FieldReference.ReferencedField().Type != code.TypeBool {
			return NewIncompatibleComparatorError(predicate.Comparator,
				predicate.FieldReference.ReferencedField())
		}

		for i := 0; i < predicate.Comparator.NumberOfArguments(); i++ {
			requiredType := predicate.Comparator.ArgumentTypeFromFieldType(
				predicate.FieldReference.ReferencedField().Type,
			)

			if params[currentParamIndex].Type != requiredType {
				return NewArgumentTypeNotMatchedError(predicate.FieldReference.ReferencingCode(), requiredType,
					params[currentParamIndex].Type)
			}
			currentParamIndex++
		}
	}

	return nil
}

func (p interfaceMethodParser) parseQuery(queryTokens []string, paramIndex int) (QuerySpec, error) {
	queryParser := queryParser{
		fieldResolver: p.fieldResolver,
		StructModel:   p.StructModel,
	}
	return queryParser.parseQuery(queryTokens, paramIndex)
}
