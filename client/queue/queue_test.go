package queue

import "testing"

func TestPop(t *testing.T) {
	q := Queue{}

	q.Push("test_string")

	test_string := q.Pop()
	if test_string != "test_string" {
		t.Errorf("got %q, wanted %q", test_string, "test_string")
	}
}

func TestEmpty(t *testing.T) {
	q := Queue{}

	if !q.Empty() {
		t.Error("Queue shall be empty")
	}

	q.Push("test_string")
	if q.Empty() {
		t.Error("Queue shall not be empty")
	}

	q.Pop()
	if !q.Empty() {
		t.Error("Queue shall be empty")
	}
}

func TestPush(t *testing.T) {
	q := Queue{}

	q.Push("test_string1")
	q.Push("test_string2")

	test_string := q.Pop()
	if test_string != "test_string1" {
		t.Errorf("got %q, wanted %q", test_string, "test_string1")
	}
}
