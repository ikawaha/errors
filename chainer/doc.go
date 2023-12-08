// Package chainer offers a set of tools for linking errors in Go applications.
// Differing from errors created using Wrap or Join, chainer treats errors as the head
// of an error list. This method keeps the error chain succinct, while still holding
// each error's detailed information. It allows for precise root cause analysis
// without losing the individual errors' context.
package chainer
