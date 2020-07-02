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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (p Person) MarshalJSON() ([]byte, error) {
	return DynMarshalJSON(reflect.ValueOf(p), p._otherInfo, "_otherInfo")
}

func (p *Person) UnmarshalJSON(data []byte) error {
	return DynUnmarshalJSON(data, reflect.ValueOf(p), &p._otherInfo, "_otherInfo")
}

func (p FavoriteOperatingSystem) MarshalJSON() ([]byte, error) {
	return DynMarshalJSON(reflect.ValueOf(p), p._otherInfo, "_otherInfo")
}

func (p *FavoriteOperatingSystem) UnmarshalJSON(data []byte) error {
	return DynUnmarshalJSON(data, reflect.ValueOf(p), &p._otherInfo, "_otherInfo")
}

func (p FooTagsTest) MarshalJSON() ([]byte, error) {
	return DynMarshalJSON(reflect.ValueOf(p), p._otherInfo, "_otherInfo")
}

func (p *FooTagsTest) UnmarshalJSON(data []byte) error {
	return DynUnmarshalJSON(data, reflect.ValueOf(p), &p._otherInfo, "_otherInfo")
}

func TestDynMarshalJSON(t *testing.T) {
	p := Person{
		ID:            "foobar",
		Name:          "amreo",
		Age:           99,
		AltNames:      []string{"bar", "foo"},
		Certification: nil,
		FavoriteOperatingSystems: []FavoriteOperatingSystem{
			{
				OS:    "Archlinux",
				Since: 2015,
			},
		},
		OptionalMainOperatingSystem: nil,
		OptionalTitle:               nil,
		_otherInfo: map[string]interface{}{
			"Profession": "Gamer",
			"Really":     true,
		},
	}

	expected := `
		{
			"BarID": "foobar",
			"Name": "amreo",
			"Age": 99,
			"AltNames": [ "bar", "foo" ],
			"Certification": null,
			"FavoriteOperatingSystems": [
				{
					"OS": "Archlinux",
					"Since": 2015
				}
			],
			"OptionalMainOperatingSystem": null,
			"OptionalTitle": null,
			"Profession": "Gamer",
			"Really": true
		}
	`

	raw, err := json.Marshal(p)
	require.NoError(t, err)

	assert.JSONEq(t, expected, string(raw))
}

func TestDynMarshalJSONWithTags(t *testing.T) {
	p := FooTagsTest{
		Normal:            4,
		Hidden:            10,
		OnlyOmitEmtpy:     nil,
		OnlyOmitEmtpy2:    ptrStr("Pippo"),
		Renamed:           "Pluto",
		RenamedOmitEmpty:  nil,
		RenamedOmitEmpty2: ptrBool(true),
		_otherInfo: map[string]interface{}{
			"OK": true,
		},
	}

	expected := `
		{
			"Normal": 4,
			"Bar": "Pluto",
			"OnlyOmitEmtpy2": "Pippo",
			"ROE2": true,
			"OK": true
		}
	`

	raw, err := json.Marshal(p)
	require.NoError(t, err)

	assert.JSONEq(t, expected, string(raw))
}

func TestDynUnmarshalJSON(t *testing.T) {
	expected := Person{
		ID:            "foobar",
		Name:          "amreo",
		Age:           99,
		AltNames:      []string{"bar", "foo"},
		Certification: nil,
		FavoriteOperatingSystems: []FavoriteOperatingSystem{
			{
				OS:         "Archlinux",
				Since:      2015,
				_otherInfo: map[string]interface{}{},
			},
		},
		OptionalMainOperatingSystem: nil,
		OptionalTitle:               nil,
		_otherInfo: map[string]interface{}{
			"Profession": "Gamer",
			"Really":     true,
		},
	}

	p := `
		{
			"BarID": "foobar",
			"Name": "amreo",
			"Age": 99,
			"AltNames": [ "bar", "foo" ],
			"Certification": null,
			"FavoriteOperatingSystems": [
				{
					"OS": "Archlinux",
					"Since": 2015
				}
			],
			"OptionalMainOperatingSystem": null,
			"OptionalTitle": null,
			"Profession": "Gamer",
			"Really": true
		}
	`

	var out Person
	require.NoError(t, json.Unmarshal([]byte(p), &out))

	assert.Equal(t, expected, out)
}

func TestDynUnmarshalJSONWithTags(t *testing.T) {
	expected := FooTagsTest{
		Normal:            4,
		Hidden:            0,
		OnlyOmitEmtpy:     nil,
		OnlyOmitEmtpy2:    ptrStr("Pippo"),
		Renamed:           "Pluto",
		RenamedOmitEmpty:  nil,
		RenamedOmitEmpty2: ptrBool(true),
		_otherInfo: map[string]interface{}{
			"OK": true,
		},
	}

	p := `
		{
			"Normal": 4,
			"Bar": "Pluto",
			"OnlyOmitEmtpy2": "Pippo",
			"ROE2": true,
			"OK": true
		}
	`

	var out FooTagsTest
	require.NoError(t, json.Unmarshal([]byte(p), &out))

	assert.Equal(t, expected, out)
}
