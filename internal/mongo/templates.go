package mongo

import (
	"github.com/sunboyy/repogen/internal/spec"
)

const constructorBody = `	return &{{.ImplStructName}}{
		collection: collection,
	}`

type constructorBodyData struct {
	ImplStructName string
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
	Sorts      []findSort
}

type findSort struct {
	BsonTag  string
	Ordering spec.Ordering
}

func (s findSort) OrderNum() int {
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
