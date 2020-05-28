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
	return DynUnmarshalJSON(data, reflect.ValueOf(p), &p._otherInfo)
}

func TestDynMarshalJSON(t *testing.T) {
	p := Person{
		Name: "amreo",
		Age:  99,
		_otherInfo: map[string]interface{}{
			"Profession": "Gamer",
			"Really":     true,
		},
	}

	expected := `
		{
			"Name": "amreo",
			"Age": 99,
			"Profession": "Gamer",
			"Really": true
		}
	`

	raw, err := json.Marshal(p)
	require.NoError(t, err)

	assert.JSONEq(t, expected, string(raw))
}

func TestDynUnmarshalJSON(t *testing.T) {
	expected := Person{
		Name: "amreo",
		Age:  99,
		_otherInfo: map[string]interface{}{
			"Profession": "Gamer",
			"Really":     true,
		},
	}

	p := `
		{
			"Name": "amreo",
			"Age": 99,
			"Profession": "Gamer",
			"Really": true
		}
	`

	var out Person
	require.NoError(t, json.Unmarshal([]byte(p), &out))

	assert.Equal(t, expected, out)
}
