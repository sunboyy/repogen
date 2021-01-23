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
	case "Find":
		return p.parseFindMethod(methodNameTokens[1:])
	}
	return MethodSpec{}, UnknownOperationError
}

func (p interfaceMethodParser) parseFindMethod(tokens []string) (MethodSpec, error) {
	if len(tokens) == 0 {
		return MethodSpec{}, UnsupportedNameError
	}

	mode, err := p.extractFindReturns(p.Method.Returns)
	if err != nil {
		return MethodSpec{}, err
	}

	querySpec, err := p.parseQuery(tokens)
	if err != nil {
		return MethodSpec{}, err
	}

	if err := p.validateMethodSignature(querySpec); err != nil {
		return MethodSpec{}, err
	}

	return MethodSpec{
		Name:    p.Method.Name,
		Params:  p.Method.Params,
		Returns: p.Method.Returns,
		Operation: FindOperation{
			Mode:  mode,
			Query: querySpec,
		},
	}, nil
}

func (p interfaceMethodParser) extractFindReturns(returns []code.Type) (QueryMode, error) {
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

func (p interfaceMethodParser) parseQuery(tokens []string) (QuerySpec, error) {
	if len(tokens) == 0 {
		return QuerySpec{}, InvalidQueryError
	}

	if len(tokens) == 1 && tokens[0] == "All" {
		return QuerySpec{}, nil
	}

	if tokens[0] == "One" {
		tokens = tokens[1:]
	}
	if tokens[0] == "By" {
		tokens = tokens[1:]
	}

	if tokens[0] == "And" || tokens[0] == "Or" {
		return QuerySpec{}, InvalidQueryError
	}

	var operator Operator
	var predicates []Predicate
	var aggregatedToken predicateToken
	for _, token := range tokens {
		if token != "And" && token != "Or" {
			aggregatedToken = append(aggregatedToken, token)
		} else if len(aggregatedToken) == 0 {
			return QuerySpec{}, InvalidQueryError
		} else if token == "And" && operator != OperatorOr {
			operator = OperatorAnd
			predicates = append(predicates, aggregatedToken.ToPredicate())
			aggregatedToken = predicateToken{}
		} else if token == "Or" && operator != OperatorAnd {
			operator = OperatorOr
			predicates = append(predicates, aggregatedToken.ToPredicate())
			aggregatedToken = predicateToken{}
		} else {
			return QuerySpec{}, InvalidQueryError
		}
	}
	if len(aggregatedToken) == 0 {
		return QuerySpec{}, InvalidQueryError
	}
	predicates = append(predicates, aggregatedToken.ToPredicate())

	return QuerySpec{Operator: operator, Predicates: predicates}, nil
}

func (p interfaceMethodParser) validateMethodSignature(querySpec QuerySpec) error {
	contextType := code.ExternalType{PackageAlias: "context", Name: "Context"}
	if len(p.Method.Params) == 0 || p.Method.Params[0].Type != contextType {
		return ContextParamRequiredError
	}

	if querySpec.NumberOfArguments()+1 != len(p.Method.Params) {
		return InvalidParamError
	}

	currentParamIndex := 1
	for _, predicate := range querySpec.Predicates {
		structField, ok := p.StructModel.Fields.ByName(predicate.Field)
		if !ok {
			return StructFieldNotFoundError
		}

		if p.Method.Params[currentParamIndex].Type != predicate.Comparator.ArgumentTypeFromFieldType(
			structField.Type) {
			return InvalidParamError
		}

		currentParamIndex++
	}

	return nil
}
