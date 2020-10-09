package test

import (
	"testing"
	"sports_service/server/barrage/client"
)

func BenchmarkClientConn(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.ClientConn()
		}
	})
}



