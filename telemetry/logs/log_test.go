package logs

import (
	"io"
	"testing"

	"google.golang.org/grpc/grpclog"
	"istio.io/pkg/log"
)

// Based on https://github.com/rs/logbench/blob/master/main_test.go

var tests = map[string]func(out io.Writer, disabled bool) grpclog.LoggerV2{}

func Benchmark(b *testing.B) {
	b.Run("grpclog", func(b *testing.B) {
		l := grpclog.Component("default")
		benchGRPC(b, l)
	})
	il := log.RegisterScope("endpoint", "echo serverside", 0)

	b.Run("istio", func(b *testing.B) {
		benchIstio(b, il)
	})
}

var sampleString = "some string with a somewhat realistic length"

func benchGRPC(b *testing.B, l grpclog.LoggerV2) {
	b.Run("Msg", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info(sampleString)
			}
		})
	})
}

func benchIstio(b *testing.B, l *log.Scope) {
	b.Run("Msg", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info(sampleString)
			}
		})
	})
}
