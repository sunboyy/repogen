package generator_test

import (
	"errors"
	"os"
	"testing"

	"github.com/sunboyy/repogen/internal/generator"
	"github.com/sunboyy/repogen/internal/testutils"
)

const (
	validStructModelName   = "User"
	validRepoInterfaceName = "UserRepositoryIntegration"
)

func TestGenerateRepository_Success(t *testing.T) {
	expectedBytes, err := os.ReadFile("../../test/generator_test_expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	expectedCode := string(expectedBytes)

	code, err := generator.GenerateRepositoryImpl(testutils.Pkg, validStructModelName, validRepoInterfaceName)

	if err != nil {
		t.Fatal(err)
	}
	if err := testutils.ExpectMultiLineString(expectedCode, code); err != nil {
		t.Error(err)
	}
}

func TestGenerateRepositoryImpl_StructNotFound(t *testing.T) {
	_, err := generator.GenerateRepositoryImpl(testutils.Pkg, "UnknownModel", validRepoInterfaceName)

	expectedError := generator.ErrStructNotFound
	if !errors.Is(err, expectedError) {
		t.Errorf("\nExpected = %+v\nReceived = %+v", expectedError, err)
	}
}

func TestGenerateRepositoryImpl_ModelNameNotStruct(t *testing.T) {
	_, err := generator.GenerateRepositoryImpl(testutils.Pkg, "UserRepositoryFind", validRepoInterfaceName)

	expectedError := generator.ErrNotNamedStruct
	if !errors.Is(err, expectedError) {
		t.Errorf("\nExpected = %+v\nReceived = %+v", expectedError, err)
	}
}

func TestGenerateRepositoryImpl_InterfaceNotFound(t *testing.T) {
	_, err := generator.GenerateRepositoryImpl(testutils.Pkg, validStructModelName, "UnknownRepository")

	expectedError := generator.ErrInterfaceNotFound
	if !errors.Is(err, expectedError) {
		t.Errorf("\nExpected = %+v\nReceived = %+v", expectedError, err)
	}
}

func TestGenerateRepositoryImpl_RepoNameNotInterface(t *testing.T) {
	_, err := generator.GenerateRepositoryImpl(testutils.Pkg, validStructModelName, "User")

	expectedError := generator.ErrNotInterface
	if !errors.Is(err, expectedError) {
		t.Errorf("\nExpected = %+v\nReceived = %+v", expectedError, err)
	}
}
