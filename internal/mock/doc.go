// Package mock splits up the functions and types found in the exported
// package mock, into individual internal packages, so that other packages can
// still use some mocks without import cycles.
//
// While having so many small packages isn't pretty, it's still better than
// having to maintain code duplicates in multiple packages, that exist due to
// import cycles.
package mock
