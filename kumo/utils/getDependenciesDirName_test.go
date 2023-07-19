package utils

import "testing"

func TestGetDependenciesDirName(t *testing.T) {
	expectedResult := "deps"

	result := GetDependenciesDirName()

	if result != expectedResult {
		t.Errorf("GetDependenciesDirName() failed, expected %v, got %v", expectedResult, result)
	}
}