// go-dyn-struct
// Copyright (C) 2020  Andrea Laisa

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// © 2020 GitHub, Inc.

package godynstruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID                          string `bson:"FooID" json:"BarID"`
	Name                        string
	Age                         int
	AltNames                    []string
	Certification               []string
	FavoriteOperatingSystems    []FavoriteOperatingSystem
	OptionalMainOperatingSystem *FavoriteOperatingSystem
	OptionalTitle               *string
	_otherInfo                  map[string]interface{}
}

type FavoriteOperatingSystem struct {
	OS    string
	Since int

	_otherInfo map[string]interface{}
}

type FooTagsTest struct {
	Normal            uint
	Renamed           string  `json:"Bar" bson:"Bar"`
	Hidden            int     `json:"-" bson:"-"`
	RenamedOmitEmpty  *bool   `json:"ROE,omitempty" bson:"ROE,omitempty"`
	OnlyOmitEmtpy     *string `json:",omitempty" bson:",omitempty"`
	RenamedOmitEmpty2 *bool   `json:"ROE2,omitempty" bson:"ROE2,omitempty"`
	OnlyOmitEmtpy2    *string `json:",omitempty" bson:",omitempty"`
	_otherInfo        map[string]interface{}
}

func ptrStr(str string) *string {
	return &str
}

func ptrBool(val bool) *bool {
	return &val
}

func TestBuildFieldInfo(t *testing.T) {
	var ftt FooTagsTest
	fttValue := reflect.ValueOf(ftt)

	outInfo, err := buildFieldInfo("Normal", fttValue, "")
	assert.NoError(t, err)
	assert.Equal(t, fieldInfo{
		actualFieldName: "Normal",
		fieldValue:      fttValue,
		omitted:         false,
		omitEmpty:       false,
	}, outInfo)

	outInfo, err = buildFieldInfo("Renamed", fttValue, "Bar")
	assert.NoError(t, err)
	assert.Equal(t, fieldInfo{
		actualFieldName: "Bar",
		fieldValue:      fttValue,
		omitted:         false,
		omitEmpty:       false,
	}, outInfo)

	outInfo, err = buildFieldInfo("Hidden", fttValue, "-")
	assert.NoError(t, err)
	assert.Equal(t, fieldInfo{
		actualFieldName: "",
		fieldValue:      fttValue,
		omitted:         true,
		omitEmpty:       false,
	}, outInfo)

	outInfo, err = buildFieldInfo("RenamedOmitEmpty", fttValue, "ROE,omitempty")
	assert.NoError(t, err)
	assert.Equal(t, fieldInfo{
		actualFieldName: "ROE",
		fieldValue:      fttValue,
		omitted:         false,
		omitEmpty:       true,
	}, outInfo)

	outInfo, err = buildFieldInfo("OnlyOmitEmtpy", fttValue, ",omitempty")
	assert.NoError(t, err)
	assert.Equal(t, fieldInfo{
		actualFieldName: "OnlyOmitEmtpy",
		fieldValue:      fttValue,
		omitted:         false,
		omitEmpty:       true,
	}, outInfo)
}
