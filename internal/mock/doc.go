// Package mock splits up the exposed functions and types found in the exported
// package mock, into individual packages.
// This is done so that other packages that are imported by some mocks, can use
// other mocks without creating import cycles.
//
// While having so many single-file packages isn't pretty, it's still better
// than having to maintain code duplicates in multiple packages, that were
// placed there to prevent import cycles.
package mock
