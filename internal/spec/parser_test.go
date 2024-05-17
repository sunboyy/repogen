package spec_test

import (
	"errors"
	"go/types"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestParseInterfaceMethod_Insert(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInsert").Type().Underlying().(*types.Interface)

	expectedOperations := []spec.Operation{
		// InsertMany
		spec.InsertOperation{
			Mode: spec.QueryModeMany,
		},
		// InsertOne
		spec.InsertOperation{
			Mode: spec.QueryModeOne,
		},
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if method.Name() != actualSpec.Name {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Name(), actualSpec.Name)
			}
			if !types.Identical(method.Type(), actualSpec.Signature) {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Type(), actualSpec.Signature)
			}
			if !reflect.DeepEqual(expectedOperations[i], actualSpec.Operation) {
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOperations[i], actualSpec.Operation)
			}
		})
	}
}

func TestParseInterfaceMethod_Find(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryFind").Type().Underlying().(*types.Interface)

	expectedOperations := []spec.Operation{
		// FindAll
		spec.FindOperation{
			Mode: spec.QueryModeMany,
		},
		// FindByAgeBetween
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorBetween,
					ParamIndex: 1,
				},
			}},
		},
		// FindByAgeGreaterThan
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorGreaterThan,
					ParamIndex: 1,
				},
			}},
		},
		// FindByAgeGreaterThanEqual
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorGreaterThanEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByAgeLessThan
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorLessThan,
					ParamIndex: 1,
				},
			}},
		},
		// FindByAgeLessThanEqual
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorLessThanEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByCity
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByCityAndGender
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{
				Operator: spec.OperatorAnd,
				Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 1,
					},
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 2,
					},
				},
			},
		},
		// FindByCityIn
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorIn,
					ParamIndex: 1,
				},
			}},
		},
		// FindByCityNot
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorNot,
					ParamIndex: 1,
				},
			}},
		},
		// FindByCityNotIn
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorNotIn,
					ParamIndex: 1,
				},
			}},
		},
		// FindByCityOrGender
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{
				Operator: spec.OperatorOr,
				Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 1,
					},
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 2,
					},
				},
			},
		},
		// FindByCityOrderByAge
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
			Sorts: []spec.Sort{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Ordering: spec.OrderingAscending,
				},
			},
		},
		// FindByCityOrderByAgeAsc
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
			Sorts: []spec.Sort{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Ordering: spec.OrderingAscending,
				},
			},
		},
		// FindByCityOrderByAgeDesc
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
			Sorts: []spec.Sort{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Ordering: spec.OrderingDescending,
				},
			},
		},
		// FindByCityOrderByCityAndAgeDesc
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
			Sorts: []spec.Sort{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Ordering: spec.OrderingAscending,
				},
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Ordering: spec.OrderingDescending,
				},
			},
		},
		// FindByCityOrderByNameFirst
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
			Sorts: []spec.Sort{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
						testutils.FindStructFieldByName(testutils.TypeNameStruct, "First"),
					},
					Ordering: spec.OrderingAscending,
				},
			},
		},
		// FindByEnabledFalse
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Enabled"),
					},
					Comparator: spec.ComparatorFalse,
					ParamIndex: 1,
				},
			}},
		},
		// FindByEnabledTrue
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Enabled"),
					},
					Comparator: spec.ComparatorTrue,
					ParamIndex: 1,
				},
			}},
		},
		// FindByID
		spec.FindOperation{
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByNameFirst
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
						testutils.FindStructFieldByName(testutils.TypeNameStruct, "First"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByPhoneNumber
		spec.FindOperation{
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "PhoneNumber"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByReferrerExists
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Referrer"),
					},
					Comparator: spec.ComparatorExists,
					ParamIndex: 1,
				},
			}},
		},
		// FindByReferrerID
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Referrer"),
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// FindByReferrerNotExists
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Referrer"),
					},
					Comparator: spec.ComparatorNotExists,
					ParamIndex: 1,
				},
			}},
		},
		// FindTop5ByGenderOrderByAgeDesc
		spec.FindOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
			Sorts: []spec.Sort{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Ordering: spec.OrderingDescending,
				},
			},
			Limit: 5,
		},
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if method.Name() != actualSpec.Name {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Name(), actualSpec.Name)
			}
			if !types.Identical(method.Type(), actualSpec.Signature) {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Type(), actualSpec.Signature)
			}
			if !reflect.DeepEqual(expectedOperations[i], actualSpec.Operation) {
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOperations[i], actualSpec.Operation)
			}
		})
	}
}

func TestParseInterfaceMethod_Update(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryUpdate").Type().Underlying().(*types.Interface)

	expectedOperations := []spec.Operation{
		// UpdateAgeIncByID
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorInc,
				},
			},
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 2,
				},
			}},
		},
		// UpdateByID
		spec.UpdateOperation{
			Update: spec.UpdateModel{},
			Mode:   spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 2,
				},
			}},
		},
		// UpdateConsentHistoryPushByID
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ConsentHistory"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorPush,
				},
			},
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 2,
				},
			}},
		},
		// UpdateEnabledAndConsentHistoryPushByID
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Enabled"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorSet,
				},
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ConsentHistory"),
					},
					ParamIndex: 2,
					Operator:   spec.UpdateOperatorPush,
				},
			},
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 3,
				},
			}},
		},
		// UpdateGenderAndCityByID
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorSet,
				},
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					ParamIndex: 2,
					Operator:   spec.UpdateOperatorSet,
				},
			},
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 3,
				},
			}},
		},
		// UpdateGenderByAge
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorSet,
				},
			},
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 2,
				},
			}},
		},
		// UpdateGenderByID
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorSet,
				},
			},
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 2,
				},
			}},
		},
		// UpdateNameFirstByID
		spec.UpdateOperation{
			Update: spec.UpdateFields{
				spec.UpdateField{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
						testutils.FindStructFieldByName(testutils.TypeNameStruct, "First"),
					},
					ParamIndex: 1,
					Operator:   spec.UpdateOperatorSet,
				},
			},
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 2,
				},
			}},
		},
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if method.Name() != actualSpec.Name {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Name(), actualSpec.Name)
			}
			if !types.Identical(method.Type(), actualSpec.Signature) {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Type(), actualSpec.Signature)
			}
			if !reflect.DeepEqual(expectedOperations[i], actualSpec.Operation) {
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOperations[i], actualSpec.Operation)
			}
		})
	}
}

func TestParseInterfaceMethod_Delete(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryDelete").Type().Underlying().(*types.Interface)

	expectedOperations := []spec.Operation{
		// DeleteAll
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
		},
		// DeleteByAgeBetween
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorBetween,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByAgeGreaterThan
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorGreaterThan,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByAgeGreaterThanEqual
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorGreaterThanEqual,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByAgeLessThan
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorLessThan,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByAgeLessThanEqual
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
					},
					Comparator: spec.ComparatorLessThanEqual,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByCity
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByCityAndGender
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{
				Operator: spec.OperatorAnd,
				Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 1,
					},
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 2,
					},
				},
			},
		},
		// DeleteByCityIn
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorIn,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByCityNot
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
					},
					Comparator: spec.ComparatorNot,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByCityOrGender
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{
				Operator: spec.OperatorOr,
				Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 1,
					},
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 2,
					},
				},
			},
		},
		// DeleteByID
		spec.DeleteOperation{
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByNameFirst
		spec.DeleteOperation{
			Mode: spec.QueryModeMany,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
						testutils.FindStructFieldByName(testutils.TypeNameStruct, "First"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
		// DeleteByPhoneNumber
		spec.DeleteOperation{
			Mode: spec.QueryModeOne,
			Query: spec.QuerySpec{Predicates: []spec.Predicate{
				{
					FieldReference: spec.FieldReference{
						testutils.FindStructFieldByName(testutils.TypeUserStruct, "PhoneNumber"),
					},
					Comparator: spec.ComparatorEqual,
					ParamIndex: 1,
				},
			}},
		},
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if method.Name() != actualSpec.Name {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Name(), actualSpec.Name)
			}
			if !types.Identical(method.Type(), actualSpec.Signature) {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Type(), actualSpec.Signature)
			}
			if !reflect.DeepEqual(expectedOperations[i], actualSpec.Operation) {
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOperations[i], actualSpec.Operation)
			}
		})
	}
}

func TestParseInterfaceMethod_Count(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryCount").Type().Underlying().(*types.Interface)

	expectedOperations := []spec.Operation{
		// CountAll
		spec.CountOperation{
			Query: spec.QuerySpec{},
		},
		// CountByGender
		spec.CountOperation{
			Query: spec.QuerySpec{
				Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 1,
					},
				},
			},
		},
		// CountByNameFirst
		spec.CountOperation{
			Query: spec.QuerySpec{
				Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{
							testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
							testutils.FindStructFieldByName(testutils.TypeNameStruct, "First"),
						},
						Comparator: spec.ComparatorEqual,
						ParamIndex: 1,
					},
				},
			},
		},
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if method.Name() != actualSpec.Name {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Name(), actualSpec.Name)
			}
			if !types.Identical(method.Type(), actualSpec.Signature) {
				t.Errorf("Expected = %+v\nReceived = %+v", method.Type(), actualSpec.Signature)
			}
			if !reflect.DeepEqual(expectedOperations[i], actualSpec.Operation) {
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOperations[i], actualSpec.Operation)
			}
		})
	}
}

func TestParseInterfaceMethod_InvalidOperation(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInvalidOperation").Type().Underlying().(*types.Interface)
	method := repoIntf.Method(0)

	_, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

	expectedError := spec.NewUnknownOperationError("Search")
	if !errors.Is(err, expectedError) {
		t.Errorf("\nExpected = %+v\nReceived = %+v", expectedError, err)
	}
}

func TestParseInterfaceMethod_Insert_Invalid(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInvalidInsert").Type().Underlying().(*types.Interface)

	expectedErrors := []error{
		// Insert1
		spec.NewOperationReturnCountUnmatchedError(2),
		// Insert2
		spec.NewUnsupportedReturnError(types.NewPointer(testutils.TypeUserNamed), 0),
		// Insert3
		spec.NewUnsupportedReturnError(repoIntf.Method(2).Type().(*types.Signature).Results().At(0).Type(), 0),
		// Insert4
		spec.NewUnsupportedReturnError(repoIntf.Method(3).Type().(*types.Signature).Results().At(1).Type(), 1),
		// Insert5
		spec.ErrContextParamRequired,
		// Insert6
		spec.ErrInvalidParam,
		// Insert7
		spec.ErrInvalidParam,
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err.Error() != expectedErrors[i].Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", expectedErrors[i], err)
			}
		})
	}
}

func TestParseInterfaceMethod_Find_Invalid(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInvalidFind").Type().Underlying().(*types.Interface)

	expectedErrors := []error{
		// Find
		spec.ErrQueryRequired,
		// FindAll
		spec.NewOperationReturnCountUnmatchedError(2),
		// FindAllOrderByAgeAnd
		spec.NewInvalidSortError([]string{"Order", "By", "Age", "And"}),
		// FindAllOrderByAgeAndAndGender
		spec.NewInvalidSortError([]string{"Order", "By", "Age", "And", "And", "Gender"}),
		// FindAllOrderByAndAge
		spec.NewInvalidSortError([]string{"Order", "By", "And", "Age"}),
		// FindAllOrderByCountry
		spec.NewStructFieldNotFoundError([]string{"Country"}),
		// FindByAge
		spec.ErrContextParamRequired,
		// FindByAndGender
		spec.NewInvalidQueryError([]string{"And", "Gender"}),
		// FindByCity
		spec.ErrInvalidParam,
		// FindByCityIn
		spec.NewArgumentTypeNotMatchedError("City", types.NewSlice(code.TypeString), code.TypeString),
		// FindByCountry
		spec.NewStructFieldNotFoundError([]string{"Country"}),
		// FindByGender
		spec.NewArgumentTypeNotMatchedError("Gender", testutils.TypeGenderNamed, code.TypeString),
		// FindByGenderAnd
		spec.NewInvalidQueryError([]string{"Gender", "And"}),
		// FindByGenderAndAndCity
		spec.NewInvalidQueryError([]string{"Gender", "And", "And", "City"}),
		// FindByGenderAndCityOrAge
		spec.NewInvalidQueryError([]string{"Gender", "And", "City", "Or", "Age"}),
		// FindByGenderFalse
		spec.NewIncompatibleComparatorError(spec.ComparatorFalse,
			testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender")),
		// FindByGenderTrue
		spec.NewIncompatibleComparatorError(spec.ComparatorTrue,
			testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender")),
		// FindByID
		spec.NewUnsupportedReturnError(testutils.TypeUserNamed, 0),
		// FindByNameMiddle
		spec.NewStructFieldNotFoundError([]string{"Name", "Middle"}),
		// FindTop
		spec.ErrLimitAmountRequired,
		// FindTop0All
		spec.ErrLimitNonPositive,
		// FindTop5All
		spec.ErrLimitOnFindOne,
		// FindTopAll
		spec.ErrLimitAmountRequired,
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err.Error() != expectedErrors[i].Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", expectedErrors[i], err)
			}
		})
	}
}

func TestParseInterfaceMethod_Update_Invalid(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInvalidUpdate").Type().Underlying().(*types.Interface)

	expectedErrors := []error{
		// UpdateAgeAndAndGenderByID
		spec.ErrInvalidUpdateFields,
		// UpdateAgeByGender
		spec.ErrContextParamRequired,
		// UpdateAgeByID
		spec.NewOperationReturnCountUnmatchedError(2),
		// UpdateAgeByIDAndUsernameOrGender
		spec.NewInvalidQueryError([]string{"ID", "And", "Username", "Or", "Gender"}),
		// UpdateByGender
		spec.ErrInvalidUpdateFields,
		// UpdateByID
		spec.ErrInvalidUpdateFields,
		// UpdateCity
		spec.ErrQueryRequired,
		// UpdateCityByID
		spec.NewUnsupportedReturnError(code.TypeFloat64, 0),
		// UpdateCityIncByID
		spec.NewIncompatibleUpdateOperatorError(spec.UpdateOperatorInc, spec.FieldReference{
			testutils.FindStructFieldByName(testutils.TypeUserStruct, "City"),
		}),
		// UpdateConsentHistoryPushByID
		spec.NewArgumentTypeNotMatchedError("ConsentHistory",
			testutils.TypeConsentHistoryNamed, types.NewSlice(testutils.TypeConsentHistoryNamed)),
		// UpdateCountryByGender
		spec.NewStructFieldNotFoundError([]string{"Country"}),
		// UpdateEnabledAll
		spec.ErrInvalidUpdateFields,
		// UpdateEnabledByCity
		spec.NewArgumentTypeNotMatchedError("City", code.TypeString, code.TypeInt),
		// UpdateEnabledByGender
		spec.NewArgumentTypeNotMatchedError("Enabled", code.TypeBool, code.TypeInt),
		// UpdateEnabledByID
		spec.NewUnsupportedReturnError(code.TypeBool, 1),
		// UpdateGenderPushByID
		spec.NewIncompatibleUpdateOperatorError(spec.UpdateOperatorPush, spec.FieldReference{
			testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
		}),
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err.Error() != expectedErrors[i].Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", expectedErrors[i], err)
			}
		})
	}
}

func TestParseInterfaceMethod_Delete_Invalid(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInvalidDelete").Type().Underlying().(*types.Interface)

	expectedErrors := []error{
		// Delete
		spec.ErrQueryRequired,
		// DeleteAll
		spec.NewOperationReturnCountUnmatchedError(2),
		// DeleteByAge
		spec.NewUnsupportedReturnError(code.TypeFloat64, 0),
		// DeleteByAndGender
		spec.NewInvalidQueryError([]string{"And", "Gender"}),
		// DeleteByCity
		spec.NewUnsupportedReturnError(code.TypeBool, 1),
		// DeleteByCityIn
		spec.NewArgumentTypeNotMatchedError("City", types.NewSlice(code.TypeString), code.TypeString),
		// DeleteByCountry
		spec.NewStructFieldNotFoundError([]string{"Country"}),
		// DeleteByEnabled
		spec.ErrInvalidParam,
		// DeleteByGender
		spec.ErrContextParamRequired,
		// DeleteByGenderAnd
		spec.NewInvalidQueryError([]string{"Gender", "And"}),
		// DeleteByGenderAndAndCity
		spec.NewInvalidQueryError([]string{"Gender", "And", "And", "City"}),
		// DeleteByGenderAndCityOrAge
		spec.NewInvalidQueryError([]string{"Gender", "And", "City", "Or", "Age"}),
		// DeleteByPhoneNumber
		spec.NewArgumentTypeNotMatchedError("PhoneNumber", code.TypeString, code.TypeInt),
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err.Error() != expectedErrors[i].Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", expectedErrors[i], err)
			}
		})
	}
}

func TestParseInterfaceMethod_Count_Invalid(t *testing.T) {
	repoIntf := testutils.Pkg.Scope().Lookup("UserRepositoryInvalidCount").Type().Underlying().(*types.Interface)

	expectedErrors := []error{
		// Count
		spec.ErrQueryRequired,
		// CountAll
		spec.NewOperationReturnCountUnmatchedError(2),
		// CountBy
		spec.NewInvalidQueryError([]string{"By"}),
		// CountByAge
		spec.NewUnsupportedReturnError(code.TypeInt64, 0),
		// CountByCity
		spec.NewUnsupportedReturnError(code.TypeBool, 1),
		// CountByCountry
		spec.NewStructFieldNotFoundError([]string{"Country"}),
		// CountByEnabled
		spec.ErrInvalidParam,
		// CountByGender
		spec.ErrContextParamRequired,
		// CountByPhoneNumber
		spec.NewArgumentTypeNotMatchedError("PhoneNumber", code.TypeString, code.TypeInt),
	}

	for i := 0; i < repoIntf.NumMethods(); i++ {
		method := repoIntf.Method(i)

		t.Run(method.Name(), func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(testutils.Pkg, testutils.TypeUserNamed, method)

			if err.Error() != expectedErrors[i].Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", expectedErrors[i], err)
			}
		})
	}
}
