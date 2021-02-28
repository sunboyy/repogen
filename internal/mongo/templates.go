package mongo

import (
	"fmt"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
)

const constructorTemplate = `
import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (r *{{.StructName}}) {{.MethodName}}({{.Parameters}}){{.Returns}} {
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

const insertOneTemplate = `	result, err := r.collection.InsertOne(arg0, arg1)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil`

const insertManyTemplate = `	var entities []interface{}
	for _, model := range arg1 {
		entities = append(entities, model)
	}
	result, err := r.collection.InsertMany(arg0, entities)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil`

type mongoFindTemplateData struct {
	EntityType string
	QuerySpec  querySpec
	Sorts      []sort
}

type sort struct {
	BsonTag  string
	Ordering spec.Ordering
}

func (s sort) OrderNum() int {
	if s.Ordering == spec.OrderingAscending {
		return 1
	}
	return -1
}

const findOneTemplate = `	var entity {{.EntityType}}
	if err := r.collection.FindOne(arg0, bson.M{
{{.QuerySpec.Code}}
	}, options.FindOne().SetSort(bson.M{
{{range $index, $element := .Sorts}}		"{{$element.BsonTag}}": {{$element.OrderNum}},
{{end}}	})).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil`

const findManyTemplate = `	cursor, err := r.collection.Find(arg0, bson.M{
{{.QuerySpec.Code}}
	}, options.Find().SetSort(bson.M{
{{range $index, $element := .Sorts}}		"{{$element.BsonTag}}": {{$element.OrderNum}},
{{end}}	}))
	if err != nil {
		return nil, err
	}
	var entities []*{{.EntityType}}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`

type mongoUpdateTemplateData struct {
	Update    update
	QuerySpec querySpec
}

const updateOneTemplate = `	result, err := r.collection.UpdateOne(arg0, bson.M{
{{.QuerySpec.Code}}
	}, bson.M{
{{.Update.Code}}
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err`

const updateManyTemplate = `	result, err := r.collection.UpdateMany(arg0, bson.M{
{{.QuerySpec.Code}}
	}, bson.M{
{{.Update.Code}}
	})
	if err != nil {
		return 0, err
	}
	return int(result.MatchedCount), err`

type mongoDeleteTemplateData struct {
	QuerySpec querySpec
}

const deleteOneTemplate = `	result, err := r.collection.DeleteOne(arg0, bson.M{
{{.QuerySpec.Code}}
	})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil`

const deleteManyTemplate = `	result, err := r.collection.DeleteMany(arg0, bson.M{
{{.QuerySpec.Code}}
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`

type mongoCountTemplateData struct {
	QuerySpec querySpec
}

const countTemplate = `	count, err := r.collection.CountDocuments(arg0, bson.M{
{{.QuerySpec.Code}}
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`
