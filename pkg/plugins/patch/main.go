package main

import (
	"errors"
	"math"
)

const xDegreePatch float64 = 10
const yDegreePatch float64 = 10

var worldPatches = createMapPatches(xDegreePatch, yDegreePatch)

//GeoPatch represents a geographic section on the earth based on latitude and longitude
type GeoPatch interface {
	ID() int
}

type geoPatch struct {
	id     int
	left   float64
	right  float64
	top    float64
	bottom float64
}

//NewGeoPatch creates a new geo patch using its ID
func NewGeoPatch(id int) (interface{}, error) {
	for i, p := range worldPatches {
		if i == id {
			return p, nil
		}
	}
	return geoPatch{}, errors.New("patch id doesn't exist")
}

//ID returns the geo patch ID
func (p geoPatch) ID() int {
	return p.id
}

func createMapPatches(xDegree float64, yDegree float64) []geoPatch {

	geoPatches := make([]geoPatch, 0)
	i := 0
	for x := -180.0; x < 180.0; x += xDegree {
		for y := -90.0; y < 90.0; y += yDegree {
			geoPatches = append(geoPatches, geoPatch{
				id:     i,
				left:   x,
				right:  x + 10,
				bottom: y,
				top:    y + 10,
			})
			i++
		}
	}
	return geoPatches
}

//FindGeoPatchByCoord retreives a geographic patch from a given geographic peer position
func FindGeoPatchByCoord(lat float64, lon float64) (p interface{}) {
	for _, patch := range worldPatches {

		if lon >= patch.left && lon <= patch.right && lat >= patch.bottom && lat <= patch.top {
			p = patch
			break
		}
	}

	return p
}

//ValidationRequiredPatchNumber returns the required number of patches for transaction validation
func ValidationRequiredPatchNumber(nbRequiredValidation int, reachablesNodesPatch []interface{}) (int, error) {

	patches := make([]GeoPatch, 0)
	for _, p := range reachablesNodesPatch {

		patch, ok := p.(GeoPatch)
		if !ok {
			return 0, errors.New("node patch is not valid GeoPatch")
		}

		patches = append(patches, patch)
	}

	return int(math.Min(float64(len(availablePatches(patches))), float64(nbRequiredValidation))), nil
}

//StorageRequiredPatchNumber returns the required number of patches for transaction storage
func StorageRequiredPatchNumber(reachablesNodesPatch []interface{}) (int, error) {
	patches := make([]GeoPatch, 0)
	for _, p := range reachablesNodesPatch {

		patch, ok := p.(GeoPatch)
		if !ok {
			return 0, errors.New("node patch is not valid GeoPatch")
		}

		patches = append(patches, patch)
	}

	return len(availablePatches(patches)), nil
}

//availablePatches returns all available patches from the all network
func availablePatches(reachablesNodesPatch []GeoPatch) (patches []GeoPatch) {
	for _, p := range reachablesNodesPatch {
		var found bool
		var i int
		for !found && i < len(patches) {
			if patches[i].ID() == p.ID() {
				found = true
			}
			i++
		}
		if !found {
			patches = append(patches, p)
		}
	}
	return
}
