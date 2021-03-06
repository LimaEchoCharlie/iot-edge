/*
 * Copyright 2020-2021 ForgeRock AS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package thing

import (
	"reflect"
	"strings"
	"testing"
)

func TestJSONContent_GetString(t *testing.T) {
	key := "key1"
	content := make(JSONContent)
	tests := []struct {
		name     string
		value    interface{}
		err      error
		expected string
	}{
		{name: "error", value: 1, err: readError{key: key}},
		{name: "ok", value: "one", expected: "one"},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			content[key] = subtest.value
			a, err := content.GetString(key)
			if subtest.err != nil {
				if subtest.err != err {
					t.Errorf("Expected a read error; got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if a != subtest.expected {
				t.Errorf("expected %v; got %v", subtest.expected, a)
			}
		})
	}
}

func TestJSONContent_GetNumber(t *testing.T) {
	key := "key1"
	content := make(JSONContent)
	tests := []struct {
		name     string
		value    interface{}
		err      error
		expected float64
	}{
		{name: "error", value: "one", err: readError{key: key}},
		{name: "ok", value: 1.2, expected: 1.2},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			content[key] = subtest.value
			a, err := content.GetNumber(key)
			if subtest.err != nil {
				if subtest.err != err {
					t.Errorf("Expected a read error; got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if a != subtest.expected {
				t.Errorf("expected %v; got %v", subtest.expected, a)
			}
		})
	}
}

func TestJSONContent_GetBool(t *testing.T) {
	key := "key1"
	content := make(JSONContent)
	tests := []struct {
		name     string
		value    interface{}
		err      error
		expected bool
	}{
		{name: "error", value: "one", err: readError{key: key}},
		{name: "ok", value: true, expected: true},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			content[key] = subtest.value
			a, err := content.GetBool(key)
			if subtest.err != nil {
				if subtest.err != err {
					t.Errorf("Expected a read error; got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if a != subtest.expected {
				t.Errorf("expected %v; got %v", subtest.expected, a)
			}
		})
	}
}

func TestJSONContent_GetStringArray(t *testing.T) {
	key := "key1"
	content := make(JSONContent)
	tests := []struct {
		name     string
		value    interface{}
		err      error
		expected []string
	}{
		{name: "error", value: "one", err: readError{key: key}},
		{name: "single-type", value: []interface{}{"one", "two"}, expected: []string{"one", "two"}},
		{name: "mixed-type", value: []interface{}{"one", "two", 3}, expected: []string{"one", "two"}},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			content[key] = subtest.value
			a, err := content.GetStringArray(key)
			if subtest.err != nil {
				if subtest.err != err {
					t.Errorf("Expected a read error; got %v", err)
				}
				return
			}
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(a, subtest.expected) {
				t.Errorf("expected %v; got %v", subtest.expected, a)
			}
		})
	}
}

func TestIntrospectionResponse_Active(t *testing.T) {
	introspection := IntrospectionResponse{
		Content: make(JSONContent),
	}
	tests := []struct {
		name   string
		active bool
	}{
		{name: "true", active: true},
		{name: "false", active: false},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			introspection.Content["active"] = subtest.active
			active, _ := introspection.Active()
			if active != subtest.active {
				t.Errorf("expected %v; got %v", subtest.active, active)
			}
		})
	}
}

func TestIntrospectionResponse_Scopes(t *testing.T) {
	introspection := IntrospectionResponse{
		Content: make(JSONContent),
	}
	tests := []struct {
		name   string
		scopes []string
	}{
		{name: "missing", scopes: nil},
		{name: "empty", scopes: []string{}},
		{name: "single", scopes: []string{"one"}},
		{name: "multiple", scopes: []string{"one", "two"}},
	}
	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			var expected []string
			if subtest.scopes == nil {
				expected = nil
				delete(introspection.Content, "scope")
			} else {
				expected = subtest.scopes
				introspection.Content["scope"] = strings.Join(expected, " ")
			}
			scopes, _ := introspection.Scope()
			if !reflect.DeepEqual(expected, scopes) {
				t.Errorf("expected %v; got %v", expected, scopes)
			}
		})
	}

}
