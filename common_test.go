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

type Person struct {
	Name                     string
	Age                      int
	AltNames                 []string
	Certification            []string
	FavoriteOperatingSystems []FavoriteOperatingSystem
	_otherInfo               map[string]interface{}
}

type FavoriteOperatingSystem struct {
	OS    string
	Since int

	_otherInfo map[string]interface{}
}
