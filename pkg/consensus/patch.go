package consensus

import (
	"errors"
	"math"
)

const xDegreePatch float64 = 10
const yDegreePatch float64 = 10

var worldPatches = createMapPatches(xDegreePatch, yDegreePatch)

//GeoPatch represents a geographic section on the earth based on latitude and longitude
type GeoPatch struct {
	patchid int
	left    float64
	right   float64
	top     float64
	bottom  float64
}

//NewGeoPatch creates a new geo patch using its ID
func NewGeoPatch(id int) (GeoPatch, error) {
	for i, p := range worldPatches {
		if i == id {
			return p, nil
		}
	}
	return GeoPatch{}, errors.New("patch id doesn't exist")
}

//ID returns the geo patch ID
func (p GeoPatch) ID() int {
	return p.patchid
}

func createMapPatches(xDegree float64, yDegree float64) []GeoPatch {

	geoPatches := make([]GeoPatch, 0)
	i := 0
	for x := -180.0; x < 180.0; x += xDegree {
		for y := -90.0; y < 90.0; y += yDegree {
			geoPatches = append(geoPatches, GeoPatch{
				patchid: i,
				left:    x,
				right:   x + 10,
				bottom:  y,
				top:     y + 10,
			})
			i++
		}
	}
	return geoPatches
}

//ComputeGeoPatch identifies a geographic patch from a given geographic peer position
func ComputeGeoPatch(lat float64, lon float64) (p GeoPatch) {
	for _, patch := range worldPatches {

		if lon >= patch.left && lon <= patch.right && lat >= patch.bottom && lat <= patch.top {
			p = patch
			break
		}
	}

	return p
}

//availablePatches returns all available patches from the all network
func availablePatches(r NodeReader) (patches []GeoPatch, err error) {
	onlineNodes, err := r.Reachables()
	if err != nil {
		return nil, err
	}

	for _, n := range onlineNodes {
		var found bool
		var i int
		for !found && i < len(patches) {
			if patches[i].patchid == n.patch.patchid {
				found = true
			}
			i++
		}
		if !found {
			patches = append(patches, n.patch)
		}
	}
	return
}

//validationRequiredPatchNumber returns the required number of patches for transaction validation
func validationRequiredPatchNumber(nbRequiredValidation int, r NodeReader) (int, error) {
	availablePatches, err := availablePatches(r)
	if err != nil {
		return 0, err
	}
	return int(math.Min(float64(len(availablePatches)), float64(nbRequiredValidation))), nil
}

//storageRequiredPatchNumber returns the required number of patches for transaction storage
func storageRequiredPatchNumber(r NodeReader) (int, error) {
	availablePatches, err := availablePatches(r)
	if err != nil {
		return 0, err
	}
	return len(availablePatches), nil
}
