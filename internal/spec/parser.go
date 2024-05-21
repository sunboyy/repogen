package spec

import (
	"go/types"
	"strconv"

	"github.com/fatih/camelcase"
	"github.com/sunboyy/repogen/internal/code"
)

// ParseInterfaceMethod returns repository method spec from declared interface
// method.
func ParseInterfaceMethod(pkg *types.Package, namedStruct *types.Named,
	method *types.Func) (MethodSpec, error) {

	parser := interfaceMethodParser{
		NamedStruct:      namedStruct,
		UnderlyingStruct: namedStruct.Underlying().(*types.Struct),
		MethodName:       method.Name(),
		Signature:        method.Type().(*types.Signature),
	}

	return parser.Parse()
}

type interfaceMethodParser struct {
	NamedStruct      *types.Named
	UnderlyingStruct *types.Struct
	MethodName       string
	Signature        *types.Signature
}

func (p interfaceMethodParser) Parse() (MethodSpec, error) {
	operation, err := p.parseMethod()
	if err != nil {
		return MethodSpec{}, err
	}

	return MethodSpec{
		Name:      p.MethodName,
		Signature: p.Signature,
		Operation: operation,
	}, nil
}

func (p interfaceMethodParser) parseMethod() (Operation, error) {
	methodNameTokens := camelcase.Split(p.MethodName)
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
	mode, err := p.extractInsertReturns(p.Signature.Results())
	if err != nil {
		return nil, err
	}

	if err := p.validateContextParam(); err != nil {
		return nil, err
	}

	pointerType := types.NewPointer(p.NamedStruct)
	if mode == QueryModeOne && !types.Identical(p.Signature.Params().At(1).Type(), pointerType) {
		return nil, ErrInvalidParam
	}

	arrayType := types.NewSlice(pointerType)
	if mode == QueryModeMany && !types.Identical(p.Signature.Params().At(1).Type(), arrayType) {
		return nil, ErrInvalidParam
	}

	return InsertOperation{
		Mode: mode,
	}, nil
}

func (p interfaceMethodParser) extractInsertReturns(returns *types.Tuple) (QueryMode, error) {
	if returns.Len() != 2 {
		return "", NewOperationReturnCountUnmatchedError(2)
	}

	if !types.Identical(returns.At(1).Type(), code.TypeError) {
		return "", NewUnsupportedReturnError(returns.At(1).Type(), 1)
	}

	switch t := returns.At(0).Type().(type) {
	case *types.Interface:
		if t.Empty() {
			return QueryModeOne, nil
		}

	case *types.Slice:
		interfaceType, ok := t.Elem().(*types.Interface)
		if ok && interfaceType.Empty() {
			return QueryModeMany, nil
		}
	}

	return "", NewUnsupportedReturnError(returns.At(0).Type(), 0)
}

func (p interfaceMethodParser) parseFindOperation(tokens []string) (Operation, error) {
	mode, err := p.extractModelOrSliceReturns(p.Signature.Results())
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

	if err := p.validateQueryOnlyParams(querySpec); err != nil {
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
	fields, ok := resolveStructField(p.UnderlyingStruct, t)
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

func (p interfaceMethodParser) extractModelOrSliceReturns(returns *types.Tuple) (QueryMode, error) {
	if returns.Len() != 2 {
		return "", NewOperationReturnCountUnmatchedError(2)
	}

	if !types.Identical(returns.At(1).Type(), code.TypeError) {
		return "", NewUnsupportedReturnError(returns.At(1).Type(), 1)
	}

	switch t := returns.At(0).Type().(type) {
	case *types.Pointer:
		if types.Identical(t.Elem(), p.NamedStruct) {
			return QueryModeOne, nil
		}

	case *types.Slice:
		pointerType, ok := t.Elem().(*types.Pointer)
		if ok {
			if types.Identical(pointerType.Elem(), p.NamedStruct) {
				return QueryModeMany, nil
			}
		}
	}

	return "", NewUnsupportedReturnError(returns.At(0).Type(), 0)
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
	mode, err := p.extractIntOrBoolReturns(p.Signature.Results())
	if err != nil {
		return nil, err
	}

	querySpec, err := p.parseQuery(tokens, 1)
	if err != nil {
		return nil, err
	}

	if err := p.validateQueryOnlyParams(querySpec); err != nil {
		return nil, err
	}

	return DeleteOperation{
		Mode:  mode,
		Query: querySpec,
	}, nil
}

func (p interfaceMethodParser) parseCountOperation(tokens []string) (Operation, error) {
	if err := p.validateCountReturns(p.Signature.Results()); err != nil {
		return nil, err
	}

	querySpec, err := p.parseQuery(tokens, 1)
	if err != nil {
		return nil, err
	}

	if err := p.validateQueryOnlyParams(querySpec); err != nil {
		return nil, err
	}

	return CountOperation{
		Query: querySpec,
	}, nil
}

func (p interfaceMethodParser) validateCountReturns(returns *types.Tuple) error {
	if returns.Len() != 2 {
		return NewOperationReturnCountUnmatchedError(2)
	}

	if !types.Identical(returns.At(0).Type(), code.TypeInt) {
		return NewUnsupportedReturnError(returns.At(0).Type(), 0)
	}

	if !types.Identical(returns.At(1).Type(), code.TypeError) {
		return NewUnsupportedReturnError(returns.At(1).Type(), 1)
	}

	return nil
}

func (p interfaceMethodParser) extractIntOrBoolReturns(returns *types.Tuple) (QueryMode, error) {
	if returns.Len() != 2 {
		return "", NewOperationReturnCountUnmatchedError(2)
	}

	if !types.Identical(returns.At(1).Type(), code.TypeError) {
		return "", NewUnsupportedReturnError(returns.At(1).Type(), 1)
	}

	basicType, ok := returns.At(0).Type().(*types.Basic)
	if ok {
		if types.Identical(basicType, code.TypeBool) {
			return QueryModeOne, nil
		}
		if types.Identical(basicType, code.TypeInt) {
			return QueryModeMany, nil
		}
	}

	return "", NewUnsupportedReturnError(returns.At(0).Type(), 0)
}

func (p interfaceMethodParser) validateQueryOnlyParams(querySpec QuerySpec) error {
	if err := p.validateContextParam(); err != nil {
		return err
	}

	if err := p.validateQueryFromParams(p.Signature.Params(), 1, querySpec); err != nil {
		return err
	}

	return nil
}

func (p interfaceMethodParser) validateContextParam() error {
	if p.Signature.Params().Len() == 0 || p.Signature.Params().At(0).Type().String() != "context.Context" {
		return ErrContextParamRequired
	}
	return nil
}

func (p interfaceMethodParser) validateQueryFromParams(params *types.Tuple, startIndex int, querySpec QuerySpec) error {
	if params.Len()-startIndex != querySpec.NumberOfArguments() {
		return ErrInvalidParam
	}

	currentParamIndex := startIndex
	for _, predicate := range querySpec.Predicates {
		if (predicate.Comparator == ComparatorTrue || predicate.Comparator == ComparatorFalse) &&
			!types.Identical(predicate.FieldReference.ReferencedField().Var.Type(), code.TypeBool) {
			return NewIncompatibleComparatorError(predicate.Comparator,
				predicate.FieldReference.ReferencedField())
		}

		for i := 0; i < predicate.Comparator.NumberOfArguments(); i++ {
			requiredType := predicate.Comparator.ArgumentTypeFromFieldType(
				predicate.FieldReference.ReferencedField().Var.Type(),
			)

			if !types.Identical(params.At(currentParamIndex).Type(), requiredType) {
				return NewArgumentTypeNotMatchedError(predicate.FieldReference.ReferencingCode(), requiredType,
					params.At(currentParamIndex).Type())
			}
			currentParamIndex++
		}
	}

	return nil
}

func (p interfaceMethodParser) parseQuery(queryTokens []string, paramIndex int) (QuerySpec, error) {
	queryParser := queryParser{
		UnderlyingStruct: p.UnderlyingStruct,
	}
	return queryParser.parseQuery(queryTokens, paramIndex)
}
