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
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DynMarshalBSON return the BSON encoding of the dynamic struct _struct
// _struct contains the reflect.Value of the struct
// extraFields is the map that contains the extra fields
// extraFieldsName is the name of the field in the struct that contains the extra fields
func DynMarshalBSON(_struct reflect.Value, extraFields map[string]interface{}, extraFieldsName string) ([]byte, error) {
	// out is the document that will be marshalled
	out := make(bson.D, 0)

	if _struct.Kind() == reflect.Ptr {
		_struct = _struct.Elem()
	}

	// add each field except extraFieldsName into out
	typ := _struct.Type()
	for i := 0; i < _struct.NumField(); i++ {
		fi := typ.Field(i)

		if fi.Name != extraFieldsName {
			val, _ := fi.Tag.Lookup("bson")
			fi, err := buildFieldInfo(fi.Name, _struct.Field(i), val)
			if err != nil {
				return nil, err
			}

			if !fi.omitted && (!fi.omitEmpty || !fi.fieldValue.IsZero()) {
				out = append(out, bson.E{Key: fi.actualFieldName, Value: fi.fieldValue.Interface()})
			}
		}
	}

	// add the missing extra fields
	tempList := make(bson.D, 0)
	for v, k := range extraFields {
		tempList = append(tempList, bson.E{Key: v, Value: k})
	}
	sort.Slice(tempList, func(i, j int) bool {
		return strings.Compare(tempList[i].Key, tempList[j].Key) < 0
	})

	out = append(out, tempList...)

	return bson.Marshal(out)
}

// DynUnmarshalBSON parses the BSON encoded data and store the result into ptrStruct. The fields that aren't part of the struct are set inside extraFieldsPtr
// data contains the BSON encoded rappresentation of the data
// ptrStruct contains a reflect.Value pointer to the struct
// extraFieldsPtr is the pointer to the extraFields map
func DynUnmarshalBSON(data []byte, ptrStruct reflect.Value, extraFieldsPtr *map[string]interface{}, extraFieldsName string) error {
	// initialize the map that contains the extra fields
	*extraFieldsPtr = make(map[string]interface{})

	// get the list of key/value pairs of the document
	var document map[string]bson.RawValue
	err := bson.Unmarshal(data, &document)
	if err != nil {
		return err
	}

	// create a map of every struct fields
	structFields := make(map[string]fieldInfo)

	typ := reflect.Indirect(ptrStruct).Type()
	for i := 0; i < typ.NumField(); i++ {
		fi := typ.Field(i)

		if fi.Name != extraFieldsName {
			val, _ := fi.Tag.Lookup("bson")
			info, err := buildFieldInfo(fi.Name, ptrStruct.Elem().Field(i), val)
			if err != nil {
				return err
			}

			structFields[info.actualFieldName] = info
		}
	}

	var othersList bson.D = []primitive.E{}

	// for each key/value pair set it to a field of struct or add it to othersList
	for k, v := range document {
		field := structFields[k]

		if field.fieldValue.IsValid() {
			if !field.omitted {
				// the field k is part of the struct, so the value will be set inside
				if v.Type == bson.TypeNull && field.fieldValue.Type().Kind() == reflect.Ptr && field.fieldValue.Type().Elem().Kind() == reflect.Struct {
					nilValue := reflect.Zero(field.fieldValue.Type())
					field.fieldValue.Set(nilValue)
				} else {
					err = v.Unmarshal(field.fieldValue.Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		} else {
			// the field k is not part of the struct, so the kv is added to othersList
			othersList = append(othersList, bson.E{Key: k, Value: v})
		}
	}

	// marshal otherList to BSON.
	tempOthersListRaw, err := bson.Marshal(othersList)
	if err != nil {
		return err
	}

	// unmarshal it to extraFieldsPtr
	err = bson.Unmarshal(tempOthersListRaw, extraFieldsPtr)
	if err != nil {
		return err
	}

	return nil
}
