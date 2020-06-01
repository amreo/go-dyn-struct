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
)

// DynMarshalJSON return the JSON encoding of the dynamic struct _struct
// _struct contains the reflect.Value of the struct
// extraFields is the map that contains the extra fields
// extraFieldsName is the name of the field in the struct that contains the extra fields
func DynMarshalJSON(_struct reflect.Value, extraFields map[string]interface{}, extraFieldsName string) ([]byte, error) {
	// out is the map that will be marshalled
	out := make(map[string]interface{})

	if _struct.Kind() == reflect.Ptr {
		_struct = _struct.Elem()
	}

	// add each field except extraFieldsName into out
	typ := _struct.Type()
	for i := 0; i < _struct.NumField(); i++ {
		fi := typ.Field(i)
		if fi.Name != extraFieldsName {
			out[fi.Name] = _struct.Field(i).Interface()
		}
	}

	// add the missing extra fields
	for v, k := range extraFields {
		out[v] = k
	}

	return json.Marshal(out)
}

// DynUnmarshalJSON parses the JSON encoded data and store the result into ptrStruct. The fields that aren't part of the struct are set inside extraFieldsPtr
// data contains the JSON encoded rappresentation of the data
// ptrStruct contains a reflect.Value pointer to the struct
// extraFieldsPtr is the pointer to the extraFields map
func DynUnmarshalJSON(data []byte, ptrStruct reflect.Value, extraFieldsPtr *map[string]interface{}) error {
	// initialize the map that contains the extra fields
	*extraFieldsPtr = make(map[string]interface{})

	// get the list of key/value pairs of the map
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(data, &objmap)
	if err != nil {
		return err
	}

	// for each key/value pair set it to a field of struct or add it to extraFields
	for k, v := range objmap {
		field := ptrStruct.Elem().FieldByName(k)

		if field.IsValid() {
			// the field k is part of the struct, so the value will be set inside
			err = json.Unmarshal(objmap[k], field.Addr().Interface())
			if err != nil {
				return err
			}
		} else {
			// the field k is not part of the struct, so the kv will be added to extraFields
			var out interface{}
			err = json.Unmarshal(v, &out)
			if err != nil {
				return err
			}
			(*extraFieldsPtr)[k] = out
		}
	}

	return nil
}
