package generator

import "testing"

func BenchmarkDefaultGenerator_Generate(b *testing.B) {
	generator := NewLineGenerator(100)
	for i := 0; i < b.N; i++ {
		_ = generator.Generate()
	}
}
