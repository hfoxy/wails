package parser

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wailsapp/wails/v3/internal/flags"
)

//go:embed testdata
var testdata embed.FS

func getFile(filename string) string {
	// get the file from the testdata FS
	file, err := fs.ReadFile(testdata, filename)
	if err != nil {
		return ""
	}
	return string(file)
}

func TestGenerateBindings(t *testing.T) {

	tests := []struct {
		name          string
		dir           string
		want          map[string]map[string]string
		useIDs        bool
		useTypescript bool
	}{
		{
			name: "complex_json",
			dir:  "testdata/complex_json",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_json/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "complex_json",
			dir:  "testdata/complex_json",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_json/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "complex_json - Typescript - CallByID",
			dir:  "testdata/complex_json",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_json/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "complex_json - Typescript - CallByName",
			dir:  "testdata/complex_json",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_json/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "complex_method",
			dir:  "testdata/complex_method",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_method/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "complex_method",
			dir:  "testdata/complex_method",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_method/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "complex_method - Typescript - CallByID",
			dir:  "testdata/complex_method",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_method/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "complex_method - Typescript - CallByName",
			dir:  "testdata/complex_method",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/complex_method/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "enum",
			dir:  "testdata/enum",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "enum",
			dir:  "testdata/enum",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "enum - Typescript - CallByID",
			dir:  "testdata/enum",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "enum - Typescript - CallByName",
			dir:  "testdata/enum",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "enum_from_imported_package",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum_from_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
			},
		},
		{
			name: "enum_from_imported_package",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum_from_imported_package/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "enum_from_imported_package - Typescript - CallByID",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum_from_imported_package/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "enum_from_imported_package - Typescript - CallByName",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/enum_from_imported_package/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "function_single",
			dir:  "testdata/function_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "function_single",
			dir:  "testdata/function_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "function single - Typescript - CallByID",
			dir:  "testdata/function_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "function single - Typescript - CallByName",
			dir:  "testdata/function_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "function_multiple_files",
			dir:  "testdata/function_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_multiple_files/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "function_multiple_files",
			dir:  "testdata/function_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_multiple_files/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "function_multiple_files - Typescript - CallByID",
			dir:  "testdata/function_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_multiple_files/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "function_multiple_files - Typescript - CallByName",
			dir:  "testdata/function_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_multiple_files/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "function_from_imported_package - CallByName",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "function_from_imported_package - CallById",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "function_from_imported_package - Typescript - CallByID",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.ts"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "function_from_imported_package - Typescript - CallByName",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.name.ts"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "function_from_nested_imported_package",
			dir:  "testdata/function_from_nested_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
				"services/other": {
					"OtherService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/services/other/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "function_from_nested_imported_package",
			dir:  "testdata/function_from_nested_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/main/GreetService.js"),
				},
				"services/other": {
					"OtherService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/services/other/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "function_from_nested_imported_package - Typescript - CallByID",
			dir:  "testdata/function_from_nested_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/main/GreetService.ts"),
				},
				"services/other": {
					"OtherService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/services/other/OtherService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "function_from_nested_imported_package - Typescript - CallByName",
			dir:  "testdata/function_from_nested_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/main/GreetService.name.ts"),
				},
				"services/other": {
					"OtherService": getFile("testdata/function_from_nested_imported_package/frontend/bindings/services/other/OtherService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "struct_literal_multiple",
			dir:  "testdata/struct_literal_multiple",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/GreetService.js"),
					"OtherService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "struct_literal_multiple",
			dir:  "testdata/struct_literal_multiple",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/GreetService.name.js"),
					"OtherService": getFile("testdata/struct_literal_multiple/frontend/bindings/main/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "function_from_imported_package",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "function_from_imported_package",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_from_imported_package/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/function_from_imported_package/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "variable_single",
			dir:  "testdata/variable_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "variable_single",
			dir:  "testdata/variable_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "variable_single - Typescript - CallByID",
			dir:  "testdata/variable_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "variable_single - Typescript - CallByName",
			dir:  "testdata/variable_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "variable_single_from_function",
			dir:  "testdata/variable_single_from_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_function/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "variable_single_from_function",
			dir:  "testdata/variable_single_from_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_function/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "variable_single_from_function - Typescript - CallByID",
			dir:  "testdata/variable_single_from_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_function/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "variable_single_from_function - Typescript - CallByName",
			dir:  "testdata/variable_single_from_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_function/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "variable_single_from_other_function",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_other_function/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "variable_single_from_other_function",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_other_function/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "variable_single_from_other_function - Typescript - CallByID",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_other_function/frontend/bindings/main/GreetService.ts"),
				},
				"services": {
					"OtherService": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/OtherService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "variable_single_from_other_function - Typescript - CallByName",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/variable_single_from_other_function/frontend/bindings/main/GreetService.name.ts"),
				},
				"services": {
					"OtherService": getFile("testdata/variable_single_from_other_function/frontend/bindings/services/OtherService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "struct_literal_single",
			dir:  "testdata/struct_literal_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "struct_literal_single",
			dir:  "testdata/struct_literal_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "struct_literal_single - Typescript - CallByID",
			dir:  "testdata/struct_literal_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_single/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "struct_literal_single - Typescript - CallByName",
			dir:  "testdata/struct_literal_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_single/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "struct_literal_multiple_other",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/GreetService.name.js"),
				},
				"services": {
					"OtherService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "struct_literal_multiple_other",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/GreetService.js"),
				},
				"services": {
					"OtherService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "struct_literal_multiple_other - Typescript - CallByID",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/GreetService.ts"),
				},
				"services": {
					"OtherService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/OtherService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "struct_literal_multiple_other - Typescript - CallByName",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/main/GreetService.name.ts"),
				},
				"services": {
					"OtherService": getFile("testdata/struct_literal_multiple_other/frontend/bindings/services/OtherService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "struct_literal_non_pointer_single",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "struct_literal_non_pointer_single",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "struct_literal_non_pointer_single - Typescript - CallByID",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "struct_literal_non_pointer_single - Typescript - CallByName",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_non_pointer_single/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "struct_literal_multiple_files",
			dir:  "testdata/struct_literal_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/GreetService.name.js"),
					"OtherService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/OtherService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "struct_literal_multiple_files",
			dir:  "testdata/struct_literal_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/GreetService.js"),
					"OtherService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/OtherService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "struct_literal_multiple_files - Typescript - CallByID",
			dir:  "testdata/struct_literal_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/GreetService.ts"),
					"OtherService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/OtherService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "struct_literal_multiple_files - Typescript - CallByName",
			dir:  "testdata/struct_literal_multiple_files",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/GreetService.name.ts"),
					"OtherService": getFile("testdata/struct_literal_multiple_files/frontend/bindings/main/OtherService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
		{
			name: "function_single_context",
			dir:  "testdata/function_single_context",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single_context/frontend/bindings/main/GreetService.js"),
				},
			},
			useIDs: true,
		},
		{
			name: "function_single_context",
			dir:  "testdata/function_single_context",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single_context/frontend/bindings/main/GreetService.name.js"),
				},
			},
			useIDs: false,
		},
		{
			name: "function single - Typescript - CallByID",
			dir:  "testdata/function_single_context",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single_context/frontend/bindings/main/GreetService.ts"),
				},
			},
			useIDs:        true,
			useTypescript: true,
		},
		{
			name: "function single - Typescript - CallByName",
			dir:  "testdata/function_single_context",
			want: map[string]map[string]string{
				"main": {
					"GreetService": getFile("testdata/function_single_context/frontend/bindings/main/GreetService.name.ts"),
				},
			},
			useIDs:        false,
			useTypescript: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run parser on directory
			absDir, err := filepath.Abs(tt.dir)
			if err != nil {
				t.Errorf("filepath.Abs() error = %v", err)
				return
			}

			options := &flags.GenerateBindingsOptions{
				ModelsFilename:  "models",
				UseNames:        !tt.useIDs,
				TS:              tt.useTypescript,
				OutputDirectory: "frontend/bindings",
			}

			project, err := ParseProject([]string{absDir}, options)
			if err != nil {
				t.Errorf("ParseProject() error = %v", err)
				return
			}

			// Generate Bindings
			got, err := project.GenerateBindings(options)
			if err != nil {
				t.Fatalf("GenerateBindings() error = %v", err)
			}

			// check if bindings are missing
			for dirName, structDetails := range tt.want {
				for name := range structDetails {
					_, ok := got[dirName][name]
					if !ok {
						t.Errorf("GenerateBindings() missing binding = %v/%v", dirName, name)
						continue
					}
				}
			}

			for dirName, structDetails := range got {
				// iterate the struct names in structDetails
				for name, binding := range structDetails {
					expected, ok := tt.want[dirName][name]
					if !ok {
						t.Errorf("GenerateBindings() unexpected binding = %v/%v", dirName, name)
						continue
					}

					// compare the binding

					// convert all line endings to \n
					binding = convertLineEndings(binding)
					expected = convertLineEndings(expected)

					if diff := cmp.Diff(expected, binding); diff != "" {
						originalFilename := name
						if !tt.useIDs {
							originalFilename += ".name"
						}
						outFileName := originalFilename + ".got"
						if tt.useTypescript {
							originalFilename += ".ts"
							outFileName += ".ts"
						} else {
							originalFilename += ".js"
							outFileName += ".js"
						}

						originalFile := filepath.Join(tt.dir, options.OutputDirectory, dirName, originalFilename)
						// Check if file exists
						if _, err := os.Stat(originalFile); err != nil {
							outFileName = originalFilename
						}

						outFile := filepath.Join(tt.dir, options.OutputDirectory, dirName, outFileName)
						err = os.WriteFile(outFile, []byte(binding), 0644)
						if err != nil {
							t.Errorf("os.WriteFile() error = %v", err)
							continue
						}

						t.Errorf("GenerateBindings() mismatch (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
