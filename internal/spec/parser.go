package spec

import (
	"errors"
	"fmt"

	"github.com/fatih/camelcase"
	"github.com/sunboyy/repogen/internal/code"
)

// ParseRepositoryInterface returns repository spec from declared repository interface
func ParseRepositoryInterface(structModel code.Struct, intf code.Interface) (RepositorySpec, error) {
	parser := repositoryInterfaceParser{
		StructModel: structModel,
		Interface:   intf,
	}

	return parser.Parse()
}

type repositoryInterfaceParser struct {
	StructModel code.Struct
	Interface   code.Interface
}

func (p repositoryInterfaceParser) Parse() (RepositorySpec, error) {
	repositorySpec := RepositorySpec{
		InterfaceName: p.Interface.Name,
	}

	for _, method := range p.Interface.Methods {
		methodSpec, err := p.parseMethod(method)
		if err != nil {
			return RepositorySpec{}, err
		}
		repositorySpec.Methods = append(repositorySpec.Methods, methodSpec)
	}

	return repositorySpec, nil
}

func (p repositoryInterfaceParser) parseMethod(method code.Method) (MethodSpec, error) {
	methodNameTokens := camelcase.Split(method.Name)
	switch methodNameTokens[0] {
	case "Find":
		return p.parseFindMethod(method, methodNameTokens[1:])
	}
	return MethodSpec{}, errors.New("method name not supported")
}

func (p repositoryInterfaceParser) parseFindMethod(method code.Method, tokens []string) (MethodSpec, error) {
	if len(tokens) == 0 {
		return MethodSpec{}, errors.New("method name not supported")
	}

	mode, err := p.extractFindReturns(method.Returns)
	if err != nil {
		return MethodSpec{}, err
	}

	querySpec, err := p.parseQuery(tokens)
	if err != nil {
		return MethodSpec{}, err
	}

	if querySpec.NumberOfArguments()+1 != len(method.Params) {
		return MethodSpec{}, errors.New("method parameter not supported")
	}

	return MethodSpec{
		Name:    method.Name,
		Params:  method.Params,
		Returns: method.Returns,
		Operation: FindOperation{
			Mode:  mode,
			Query: querySpec,
		},
	}, nil
}

func (p repositoryInterfaceParser) extractFindReturns(returns []code.Type) (QueryMode, error) {
	if len(returns) != 2 {
		return "", errors.New("method return not supported")
	}

	if returns[1] != code.SimpleType("error") {
		return "", errors.New("method return not supported")
	}

	pointerType, ok := returns[0].(code.PointerType)
	if ok {
		simpleType := pointerType.ContainedType
		if simpleType == code.SimpleType(p.StructModel.Name) {
			return QueryModeOne, nil
		}
		return "", fmt.Errorf("invalid return type %s", pointerType.Code())
	}

	arrayType, ok := returns[0].(code.ArrayType)
	if ok {
		pointerType, ok := arrayType.ContainedType.(code.PointerType)
		if ok {
			simpleType := pointerType.ContainedType
			if simpleType == code.SimpleType(p.StructModel.Name) {
				return QueryModeMany, nil
			}
			return "", fmt.Errorf("invalid return type %s", pointerType.Code())
		}
	}

	return "", errors.New("method return not supported")
}

func (p repositoryInterfaceParser) parseQuery(tokens []string) (QuerySpec, error) {
	if len(tokens) == 0 {
		return QuerySpec{}, errors.New("method name not supported")
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
		return QuerySpec{}, errors.New("method name not supported")
	}

	var operator Operator
	var predicates []Predicate
	var aggregatedToken predicateToken
	for _, token := range tokens {
		if token != "And" && token != "Or" {
			aggregatedToken = append(aggregatedToken, token)
		} else if token == "And" && operator != OperatorOr {
			operator = OperatorAnd
			predicates = append(predicates, aggregatedToken.ToPredicate())
			aggregatedToken = predicateToken{}
		} else if token == "Or" && operator != OperatorAnd {
			operator = OperatorOr
			predicates = append(predicates, aggregatedToken.ToPredicate())
			aggregatedToken = predicateToken{}
		} else {
			return QuerySpec{}, errors.New("method name contains ambiguous query")
		}
	}
	if len(aggregatedToken) == 0 {
		return QuerySpec{}, errors.New("method name not supported")
	}
	predicates = append(predicates, aggregatedToken.ToPredicate())

	return QuerySpec{Operator: operator, Predicates: predicates}, nil
}
