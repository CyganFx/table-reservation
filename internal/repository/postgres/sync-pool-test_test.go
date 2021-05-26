package postgres

import "testing"

func BenchmarkWithPoolReseting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithPoolReseting()
	}
}

func BenchmarkWithPoolAssigning(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithPoolAssigning()
	}
}

func BenchmarkWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithoutPool()
	}
}
