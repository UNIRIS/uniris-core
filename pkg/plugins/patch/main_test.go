package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Generate a patch  from coordinates
	Given a latitude and logitude
	When I want to get the patch id
	Then I get a patch with an ID
*/
func TestFindGeoPatchByCoord(t *testing.T) {
	lat1 := 0.0
	lon1 := 0.0
	p := FindGeoPatchByCoord(lat1, lon1)
	assert.NotEqual(t, 0, p.(GeoPatch).ID())

	//position of Eiffel Tower, Paris
	lat2 := 48.8583728827653310
	lon2 := 2.2944796085357666
	//position of triumphal arch, Paris
	lat3 := 48.873804445573874
	lon3 := 2.2950267791748047

	p2 := FindGeoPatchByCoord(lat2, lon2)
	p3 := FindGeoPatchByCoord(lat3, lon3)
	assert.Equal(t, p2.(GeoPatch).ID(), p3.(GeoPatch).ID())

	//position of statue of liberty, New york
	lat4 := 40.689039
	lon4 := -74.044396
	//position of Clock Habib Bourguiba, Tunis
	lat5 := 36.800236
	lon5 := 10.186422

	p4 := FindGeoPatchByCoord(lat4, lon4)
	p5 := FindGeoPatchByCoord(lat5, lon5)
	assert.NotEqual(t, p4.(GeoPatch).ID(), p5.(GeoPatch).ID())
}

/*
Scenario: Get the available patches from a network topology
	Given 5 nodes in 3 patches
	When I want to get the available patches
	Then I get 3 patches
*/
func TestGetAvailablePatches(t *testing.T) {

	patches := []GeoPatch{
		geoPatch{id: 1},
		geoPatch{id: 2},
		geoPatch{id: 3},
		geoPatch{id: 2},
		geoPatch{id: 1},
	}

	assert.Len(t, availablePatches(patches), 3)
}

/*
Scenario: Get the number of required patches for transaction validation
	Given a 3 available patches and 5 minimum transaction validations
	When I went to get the number of required number patches
	Then I get 3
*/
func TestValidationRequiredPatchNumber(t *testing.T) {
	patches := []interface{}{
		geoPatch{id: 1},
		geoPatch{id: 2},
		geoPatch{id: 3},
		geoPatch{id: 2},
		geoPatch{id: 1},
	}

	nbPatches, err := ValidationRequiredPatchNumber(5, patches)
	assert.Nil(t, err)
	assert.Equal(t, 3, nbPatches)
}

/*
Scenario: Get the number of required patches for transaction storage
	Given a 3 available patches
	When I went to get the number of required number patches
	Then I get 3
*/
func TestStorageRequiredPatchNumber(t *testing.T) {
	patches := []interface{}{
		geoPatch{id: 1},
		geoPatch{id: 2},
		geoPatch{id: 3},
		geoPatch{id: 2},
		geoPatch{id: 1},
	}

	nbPatches, err := StorageRequiredPatchNumber(patches)
	assert.Nil(t, err)
	assert.Equal(t, 3, nbPatches)
}
