package task

import "testing"

func TestShortID(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "full id", input: "T3k7x-example", want: "T3k7x"},
		{name: "short id", input: "E1a1a", want: "E1a1a"},
		{name: "other", input: "not-an-id", want: "not-an-id"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := ShortID(tc.input); got != tc.want {
				t.Fatalf("ShortID(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestResolveTaskID(t *testing.T) {
	tasks := map[string]*Task{
		"T1a1a-foo": {ID: "T1a1a-foo"},
		"T2b2b-bar": {ID: "T2b2b-bar"},
	}

	resolved, err := ResolveTaskID(tasks, "T1a1a-foo")
	if err != nil {
		t.Fatalf("ResolveTaskID(full) error: %v", err)
	}
	if resolved != "T1a1a-foo" {
		t.Fatalf("ResolveTaskID(full) = %q, want %q", resolved, "T1a1a-foo")
	}

	resolved, err = ResolveTaskID(tasks, "T2b2b")
	if err != nil {
		t.Fatalf("ResolveTaskID(short) error: %v", err)
	}
	if resolved != "T2b2b-bar" {
		t.Fatalf("ResolveTaskID(short) = %q, want %q", resolved, "T2b2b-bar")
	}

	resolved, err = ResolveTaskID(tasks, "tasks/T1a1a-foo/task.md")
	if err != nil {
		t.Fatalf("ResolveTaskID(path) error: %v", err)
	}
	if resolved != "T1a1a-foo" {
		t.Fatalf("ResolveTaskID(path) = %q, want %q", resolved, "T1a1a-foo")
	}

	_, err = ResolveTaskID(tasks, "T9z9z")
	if err == nil {
		t.Fatalf("ResolveTaskID(missing) expected error")
	}
}

func TestIsValidTaskID(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "valid short id min length", input: "T1a1a", want: true},
		{name: "valid short id max length", input: "T1a1abc", want: true},
		{name: "valid full id", input: "T3k7x-example", want: true},
		{name: "valid full id with numbers", input: "E2k7x-123-test", want: true},
		{name: "empty string", input: "", want: false},
		{name: "whitespace only", input: "   ", want: false},
		{name: "lowercase letter start", input: "t1a1a", want: false},
		{name: "too short", input: "T1a1", want: false},
		{name: "too long", input: "T1a1abcd", want: false},
		{name: "special characters", input: "T1@1a", want: false},
		{name: "no hyphen in full format", input: "T3k7xexample", want: false},
		{name: "random text", input: "not-an-id", want: false},
		{name: "number start", input: "13k7x-example", want: false},
		{name: "valid with whitespace", input: "  T3k7x-example  ", want: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidTaskID(tc.input); got != tc.want {
				t.Fatalf("IsValidTaskID(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestResolveTaskIDAmbiguous(t *testing.T) {
	tasks := map[string]*Task{
		"T1a1a-foo": {ID: "T1a1a-foo"},
		"T1a1a-bar": {ID: "T1a1a-bar"},
	}

	_, err := ResolveTaskID(tasks, "T1a1a")
	if err == nil {
		t.Fatalf("ResolveTaskID(ambiguous) expected error")
	}
}
