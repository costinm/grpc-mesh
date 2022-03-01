package logs

import (
	"io"
	"io/ioutil"
	"testing"

	"google.golang.org/grpc/grpclog"
	"istio.io/pkg/log"
)

// Based on https://github.com/rs/logbench/blob/master/main_test.go

var tests = map[string]func(out io.Writer, disabled bool) grpclog.LoggerV2{}

func Benchmark(b *testing.B) {
	tests["default"] = func(out io.Writer, disabled bool) grpclog.LoggerV2 {
		return grpclog.Component("default")
	}

	for name, t := range tests {
		b.Run(name, func(b *testing.B) {
			b.Run("Enabled", func(b *testing.B) {
				l := t(ioutil.Discard, false)
				benchGRPC(b, l)
			})
			il := log.RegisterScope("endpoint", "echo serverside", 0)

			//b.Run("Disabled", func(b *testing.B) {
			//	l := t(ioutil.Discard, true)
			//	benchGRPC(b, l)
			//})
		})
	}
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

func benchIstio(b *testing.B, l log.Log) {
	b.Run("Msg", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info(sampleString)
			}
		})
	})
}
