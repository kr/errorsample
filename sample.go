// Package errorsample samples error values
// uniformly at random
// from an unbounded set of inputs.
// It provides a representative sample
// when the total amount of errors
// is too many to store.
//
// Functions in this package are safe to call concurrently.
package errorsample // import "github.com/kr/errorsample"

import (
	"math/rand"
	"sync"
)

// Set represents an unbounded set of errors.
// Its Sample method returns a sample of bounded size,
// chosen uniformly at random from the set.
//
// Its methods are safe to call concurrently.
// The zero value of Set is a set
// with a capacity of 1, ready to use.
type Set struct {
	mu  sync.Mutex
	n   int
	buf []error // slice header is constant
}

// New returns a new Set
// that samples up to cap errors.
func New(cap int) *Set {
	return &Set{buf: make([]error, cap)}
}

// Empty removes all errors from s.
func (s *Set) Empty() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n = 0
}

// Add adds err to s.
func (s *Set) Add(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.buf) < 1 {
		s.buf = make([]error, 1)
	}

	if s.n < len(s.buf) {
		s.buf[s.n] = err
	} else if i := rand.Intn(s.n); i < len(s.buf) {
		// Sample this item with prob. len(s.buf)/s.n.
		// Replace an existing sample with prob. 1/len(s.buf).
		// See Jeffrey S. Vitter, Random sampling with a reservoir,
		// ACM Trans. Math. Softw. 11 (1985), no. 1, 37–57.
		s.buf[i] = err
	}
	s.n++
}

// Sample returns a uniform random sample
// of n <= cap error values from s.
//
// If fewer than n errors have been added
// since the last call to Empty,
// it returns all errors added so far.
//
// Repeated calls to Sample are not random
// with respect to each other,
// only with respect to
// the sequence of errors added to s.
// In particular, two successive calls to Sample
// with no intervening Add or Empty
// will return the same sample.
func (s *Set) Sample(n int) []error {
	s.mu.Lock()
	defer s.mu.Unlock()

	a := make([]error, min(n, s.n, len(s.buf)))
	for i := range a {
		a[i] = s.buf[i]
	}
	return a
}

// Cap returns the capacity of s.
// Sample will return at most
// this many errors.
func (s *Set) Cap() int {
	return len(s.buf)
}

// Added returns the number of errors added to s
// since the last call to Empty.
func (s *Set) Added() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.n
}

func min(n0 int, n ...int) int {
	for _, ni := range n {
		if ni < n0 {
			n0 = ni
		}
	}
	return n0
}
