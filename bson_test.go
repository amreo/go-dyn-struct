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
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (p Person) MarshalBSON() ([]byte, error) {
	return DynMarshalBSON(reflect.ValueOf(p), p._otherInfo, "_otherInfo")
}

func (p *Person) UnmarshalBSON(data []byte) error {
	return DynUnmarshalBSON(data, reflect.ValueOf(p), &p._otherInfo)
}

func (p FavoriteOperatingSystem) MarshalBSON() ([]byte, error) {
	return DynMarshalBSON(reflect.ValueOf(p), p._otherInfo, "_otherInfo")
}

func (p *FavoriteOperatingSystem) UnmarshalBSON(data []byte) error {
	return DynUnmarshalBSON(data, reflect.ValueOf(p), &p._otherInfo)
}

func TestDynMarshalBSON(t *testing.T) {
	p1 := Person{
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
		// OptionalMainOperatingSystem: nil,
		OptionalTitle: nil,
		_otherInfo: map[string]interface{}{
			"Profession": "Gamer",
			"Really":     true,
		},
	}

	p2 := bson.D{
		bson.E{Key: "Name", Value: "amreo"},
		bson.E{Key: "Age", Value: 99},
		bson.E{Key: "AltNames", Value: bson.A{
			"bar",
			"foo",
		}},
		bson.E{Key: "Certification", Value: nil},
		bson.E{Key: "FavoriteOperatingSystems", Value: bson.A{
			bson.D{
				bson.E{Key: "OS", Value: "Archlinux"},
				bson.E{Key: "Since", Value: 2015},
			},
		}},
		bson.E{Key: "OptionalMainOperatingSystem", Value: nil},
		bson.E{Key: "OptionalTitle", Value: nil},
		bson.E{Key: "Profession", Value: "Gamer"},
		bson.E{Key: "Really", Value: true},
	}

	raw1, err := bson.Marshal(p1)
	require.NoError(t, err)

	raw2, err := bson.Marshal(p2)
	require.NoError(t, err)

	assert.Equal(t, raw1, raw2)
}

func TestDynUnmarshalBSON(t *testing.T) {
	p := bson.D{
		bson.E{Key: "Name", Value: "amreo"},
		bson.E{Key: "Age", Value: 99},
		bson.E{Key: "AltNames", Value: bson.A{
			"bar",
			"foo",
		}},
		bson.E{Key: "Certification", Value: nil},
		bson.E{Key: "FavoriteOperatingSystems", Value: bson.A{
			bson.D{
				bson.E{Key: "OS", Value: "Archlinux"},
				bson.E{Key: "Since", Value: 2015},
			},
		}},
		bson.E{Key: "OptionalMainOperatingSystem", Value: nil},
		bson.E{Key: "OptionalTitle", Value: nil},
		bson.E{Key: "Profession", Value: "Gamer"},
		bson.E{Key: "Really", Value: true},
	}

	expected := Person{
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

	raw, err := bson.Marshal(p)
	require.NoError(t, err)

	var out Person
	require.NoError(t, bson.Unmarshal(raw, &out))

	assert.Equal(t, expected, out)
}
