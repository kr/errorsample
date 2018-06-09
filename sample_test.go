package errorsample

import (
	"errors"
	"reflect"
	"testing"
)

var err0 = errors.New("err0")

func TestNew(t *testing.T) {
	got := New(1)
	want := &Set{buf: make([]error, 1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("New(1) = %v, want %v", got, want)
	}
}

func TestReset(t *testing.T) {
	set := &Set{n: 1}
	set.Reset()
	want := 0
	if set.n != want {
		t.Errorf("after Reset, set.n = %d, want %d", set.n, want)
	}
}

func TestAdd0(t *testing.T) {
	set := &Set{buf: make([]error, 2)}
	set.Add(err0)
	if set.n != 1 {
		t.Errorf("after Add, n = %d, want 1", set.n)
	}
	if set.buf[0] != err0 {
		t.Errorf("after Add, buf[0] = %v, want err0", set.buf[0])
	}
	if set.buf[1] != nil {
		t.Errorf("after Add, buf[1] = %v, want nil", set.buf[1])
	}
}

func TestAddN(t *testing.T) {
	set := &Set{buf: make([]error, 2), n: 2}
	set.Add(err0)
	if set.n != 3 {
		t.Fatalf("after Add, set.n = %d, want 3", set.n)
	}
	if (set.buf[0] == nil) == (set.buf[1] == nil) {
		t.Errorf("after Add, exactly one buffer cell must be nil")
	}
	if (set.buf[0] == err0) == (set.buf[1] == err0) {
		t.Errorf("after Add, exactly one buffer cell must be err0")
	}
}

func TestSample(t *testing.T) {
	cases := []struct{ cap, param, added int }{
		{1, 2, 3},
		{3, 1, 2},
		{2, 3, 1},
	}
	for _, test := range cases {
		set := &Set{buf: make([]error, test.cap), n: test.added}
		got := set.Sample(test.param)
		if len(got) != 1 {
			t.Errorf("case %+v len(got) = %d, want 1", test, len(got))
		}
	}
}

func TestCap(t *testing.T) {
	set := &Set{buf: make([]error, 2)}
	got := set.Cap()
	if got != 2 {
		t.Errorf("set.Cap() = %d, want 2", got)
	}
}

func TestAdded(t *testing.T) {
	set := &Set{n: 2}
	got := set.Added()
	if got != 2 {
		t.Errorf("set.Added() = %d, want 2", got)
	}
}

func TestZero(t *testing.T) {
	var set Set
	gotCap := set.Cap()
	if gotCap != 0 {
		t.Errorf("Set{}.Cap() = %d, want 0", gotCap)
	}
	gotAdded := set.Added()
	if gotAdded != 0 {
		t.Errorf("Set{}.Added() = %d, want 0", gotAdded)
	}
	gotSample := set.Sample(5)
	wantSample := []error{}
	if !reflect.DeepEqual(gotSample, wantSample) {
		t.Errorf("Set{}.Sample() = %v, want %v", gotSample, wantSample)
	}

	set.Add(errors.New("a"))
	set.Add(errors.New("a"))
	set.Add(errors.New("a"))

	gotAdded = set.Added()
	if gotAdded != 3 {
		t.Errorf("Set{}.Added() = %d, want 3", gotAdded)
	}
	gotSample = set.Sample(5)
	wantSample = []error{errors.New("a")}
	if !reflect.DeepEqual(gotSample, wantSample) {
		t.Errorf("Set{}.Sample() = %v, want %v", gotSample, wantSample)
	}
}
