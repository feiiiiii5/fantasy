package jsonrepair

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRepairJSONPreservesValidStringEndingInBackslash(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{name: "one trailing backslash", input: `{"a": "val\\"}`},
		{name: "two trailing backslashes", input: `{"a":"val\\\\"}`},
		{name: "backslashes before content", input: `{"a": "x\\\\b"}`},
		{name: "trailing backslash in array object", input: `[{"a": "\\"}]`},
		{name: "single backslash value", input: `{"a": "\\"}`},
		{name: "multiple trailing backslashes after text", input: `{"a":"ends \\\\"}`},
		{name: "command ending in backslash", input: `{"cmd": "grep -r 'x' src\\"}`},
		{name: "key ending in backslash", input: `{"a\\":1}`},
		{name: "single backslash key", input: `{"\\":1}`},
		{name: "key ending in two backslashes", input: `{"a\\\\":1}`},
		{name: "backslash inside key", input: `{"x\\\\y":1}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !json.Valid([]byte(tc.input)) {
				t.Fatal("regression input must be valid JSON")
			}

			got, err := RepairJSON(tc.input)
			if err != nil {
				t.Fatalf("RepairJSON returned an error for valid JSON: %v", err)
			}

			var wantValue, gotValue any
			if err := json.Unmarshal([]byte(tc.input), &wantValue); err != nil {
				t.Fatalf("parse input: %v", err)
			}
			if err := json.Unmarshal([]byte(got), &gotValue); err != nil {
				t.Fatalf("parse repaired JSON %q: %v", got, err)
			}
			if !reflect.DeepEqual(gotValue, wantValue) {
				t.Fatalf("RepairJSON changed the decoded value: input=%q repaired=%q want=%#v got=%#v", tc.input, got, wantValue, gotValue)
			}
		})
	}
}
