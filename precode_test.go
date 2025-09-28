package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func contains(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}

func TestTaskStructTags(t *testing.T) {
	typ := reflect.TypeOf(Task{})

	tests := map[string]string{
		"ID":           "id",
		"Description":  "description",
		"Note":         "note",
		"Applications": "applications",
	}

	for field, want := range tests {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("field %q not found on Task", field)
		}
		if got := f.Tag.Get("json"); got != want {
			t.Errorf("wrong json tag for %s: got %q, want %q", field, got, want)
		}
	}
}

func TestTasksFixtureBasic(t *testing.T) {
	if len(tasks) != 2 {
		t.Fatalf("unexpected tasks length: got %d, want %d", len(tasks), 2)
	}

	t1, ok := tasks["1"]
	if !ok {
		t.Fatalf("task '1' not found")
	}
	if t1.ID != "1" || t1.Description == "" || t1.Note == "" {
		t.Errorf("task '1' has invalid content: %+v", t1)
	}
	if !contains(t1.Applications, "git") {
		t.Errorf("task '1' Applications missing 'git': %v", t1.Applications)
	}

	t2, ok := tasks["2"]
	if !ok {
		t.Fatalf("task '2' not found")
	}
	if t2.ID != "2" || t2.Description == "" || t2.Note == "" {
		t.Errorf("task '2' has invalid content: %+v", t2)
	}
	if !contains(t2.Applications, "Postman") {
		t.Errorf("task '2' Applications missing 'Postman': %v", t2.Applications)
	}
}

func TestTaskJSONMarshalling(t *testing.T) {
	tk := tasks["1"]

	data, err := json.Marshal(tk)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	required := []string{"id", "description", "note", "applications"}
	for _, key := range required {
		if _, ok := m[key]; !ok {
			t.Errorf("missing JSON key %q", key)
		}
	}
	if len(m) != len(required) {
		t.Errorf("unexpected number of JSON keys: got %d, want %d", len(m), len(required))
	}
}
