package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Find PatchID
	Given a latitude and logitude
	When I want to get the patch id
	Then I get the wanted result.
*/
func TestPatchId(t *testing.T) {
	lat1 := 0.0
	lon1 := 0.0
	p := ComputeGeoPatch(lat1, lon1)
	assert.NotEqual(t, 0, p.patchid)

	//position of Eiffel Tower, Paris
	lat2 := 48.8583728827653310
	lon2 := 2.2944796085357666
	//position of triumphal arch, Paris
	lat3 := 48.873804445573874
	lon3 := 2.2950267791748047

	p2 := ComputeGeoPatch(lat2, lon2)
	p3 := ComputeGeoPatch(lat3, lon3)
	assert.Equal(t, p2.patchid, p3.patchid)

	//position of statue of liberty, New york
	lat4 := 40.689039
	lon4 := -74.044396
	//position of Clock Habib Bourguiba, Tunis
	lat5 := 36.800236
	lon5 := 10.186422

	p4 := ComputeGeoPatch(lat4, lon4)
	p5 := ComputeGeoPatch(lat5, lon5)
	assert.NotEqual(t, p4.patchid, p5.patchid)
}
