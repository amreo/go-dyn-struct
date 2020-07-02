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
	"errors"
	"reflect"
	"strings"
)

type fieldInfo struct {
	actualFieldName string
	fieldValue      reflect.Value
	omitted         bool
	omitEmpty       bool
}

func buildFieldInfo(fieldName string, fieldValue reflect.Value, tags string) (fieldInfo, error) {
	out := fieldInfo{
		fieldValue: fieldValue,
	}

	if tags == "-" {
		out.omitted = true
		return out, nil
	}

	out.actualFieldName = fieldName
	parts := strings.Split(tags, ",")
	if parts[0] != "" {
		out.actualFieldName = parts[0]
	}

	for _, part := range parts[1:] {
		switch part {
		case "omitempty":
			out.omitEmpty = true
		default:
			return fieldInfo{}, errors.New("Unrecognized part in field tags " + tags)
		}
	}

	return out, nil
}
