// Package utils provides utility functions for encoding and decoding operations
// used throughout the protoc-gen-go-flags project.
package utils

import (
	"encoding/base64"
	"encoding/hex"
)

// MustDecodeBase64 decodes a base64-encoded string and returns the resulting bytes.
// It panics if the input string is not valid base64, following the "Must" prefix convention.
// This function is useful in contexts where the input is expected to be valid base64
// and a panic is appropriate for error handling (such as during initialization or testing).
//
// Parameters:
//   - d: A base64-encoded string to decode
//
// Returns:
//   - []byte: The decoded byte slice
//
// Panics:
//   - If the input string is not valid base64 encoding
func MustDecodeBase64(d string) []byte {
	val, err := base64.StdEncoding.DecodeString(d)
	if err != nil {
		panic(err)
	}
	return val
}

// MustDecodeHex decodes a hex-encoded string and returns the resulting bytes.
// It panics if the input string is not valid hexadecimal, following the "Must" prefix convention.
// This function is useful in contexts where the input is expected to be valid hex
// and a panic is appropriate for error handling (such as during initialization or testing).
//
// Parameters:
//   - d: A hex-encoded string to decode
//
// Returns:
//   - []byte: The decoded byte slice
//
// Panics:
//   - If the input string is not valid hexadecimal encoding
func MustDecodeHex(d string) []byte {
	val, err := hex.DecodeString(d)
	if err != nil {
		panic(err)
	}
	return val
}
