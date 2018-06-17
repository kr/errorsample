package errorsample_test

import (
	"errors"
	"fmt"

	"github.com/kr/errorsample"
)

func Example() {
	set := errorsample.New(20)
	set.Add(errors.New("first"))  // add errors
	set.Add(errors.New("second")) // lots of errors
	set.Add(errors.New("third"))  // lots and LOTS of errors
	errs := make([]error, set.Cap())
	errs = errs[:set.Sample(errs)]
	fmt.Println("our sample:", errs)
	// Output:
	// our sample: [first second third]
}
