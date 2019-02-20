// +build gofuzz

package sdp

import "reflect"

// Fuzz implements a randomized fuzz test of the sdp
// parser using go-fuzz.
//
// To run the fuzzer, first download go-fuzz:
// `go get github.com/dvyukov/go-fuzz/...`
//
// Then build the testing package:
// `go-fuzz-build github.com/pions/sdp`
//
// And run the fuzzer on the corpus:
// `go-fuzz -bin=sdp-fuzz.zip -workdir=fuzzer`
func Fuzz(data []byte) int {
	var sdes SessionDescription
	if err := sdes.Unmarshal(string(data)); err != nil {
		return 0
	}
	out := sdes.Marshal()
	if !reflect.DeepEqual(out, data) {
		panic("not equal")
	}

	return 1
}
