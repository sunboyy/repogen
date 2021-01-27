package mongo

import (
	"fmt"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

const constructorTemplate = `
import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func New{{.InterfaceName}}(collection *mongo.Collection) {{.InterfaceName}} {
	return &{{.ImplStructName}}{
		collection: collection,
	}
}

type {{.ImplStructName}} struct {
	collection *mongo.Collection
}
`

type mongoConstructorTemplateData struct {
	InterfaceName  string
	ImplStructName string
}

const methodTemplate = `
func (r *{{.StructName}}) {{.MethodName}}(ctx context.Context, {{.Parameters}}){{.Returns}} {
{{.Implementation}}
}
`

type mongoMethodTemplateData struct {
	StructName     string
	MethodName     string
	ParamTypes     []code.Type
	ReturnTypes    []code.Type
	Implementation string
}

func (data mongoMethodTemplateData) Parameters() string {
	var params []string
	for i, paramType := range data.ParamTypes {
		params = append(params, fmt.Sprintf("arg%d %s", i, paramType.Code()))
	}
	return strings.Join(params, ", ")
}

func (data mongoMethodTemplateData) Returns() string {
	if len(data.ReturnTypes) == 0 {
		return ""
	}

	if len(data.ReturnTypes) == 1 {
		return fmt.Sprintf(" %s", data.ReturnTypes[0].Code())
	}

	var returns []string
	for _, returnType := range data.ReturnTypes {
		returns = append(returns, returnType.Code())
	}
	return fmt.Sprintf(" (%s)", strings.Join(returns, ", "))
}

type mongoFindTemplateData struct {
	EntityType string
	QuerySpec  querySpec
}

const findOneTemplate = `	var entity {{.EntityType}}
	if err := r.collection.FindOne(ctx, bson.M{
{{.QuerySpec.Code}}
	}).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil`

const findManyTemplate = `	cursor, err := r.collection.Find(ctx, bson.M{
{{.QuerySpec.Code}}
	})
	if err != nil {
		return nil, err
	}
	var entities []*{{.EntityType}}
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return entities, nil`

type mongoDeleteTemplateData struct {
	QuerySpec querySpec
}

const deleteOneTemplate = `	result, err := r.collection.DeleteOne(ctx, bson.M{
{{.QuerySpec.Code}}
	})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil`

const deleteManyTemplate = `	result, err := r.collection.DeleteMany(ctx, bson.M{
{{.QuerySpec.Code}}
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`
