package test

import (
	"sports_service/barrage/client"
	"testing"
)

func BenchmarkClientConn(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.ClientConn()
		}
	})
}
