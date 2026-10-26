package jsonrepair

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRepairJSONPreservesKeysStartingWithObjectDelimiter(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{name: "key starting with comma", input: `{"a":"x",",b":1}`},
		{name: "key starting with closing brace", input: `{"a":"x","}":1}`},
		{name: "comma key with whitespace after separator", input: `{"a":"x", ",b":1}`},
		{name: "lone comma key", input: `{"a":"x",",":1}`},
		{name: "nested comma key", input: `{"outer":{"a":"x",",b":1}}`},
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
