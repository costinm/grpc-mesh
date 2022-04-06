// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_zap

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

// ReplaceGrpcLogger sets the given zap.Logger as a gRPC-level logger.
// This should be called *before* any other initialization, preferably from init() functions.
// Deprecated: use ReplaceGrpcLoggerV2.
//func ReplaceGrpcLogger(logger *zap.Logger) {
//	zgl := &zapGrpcLogger{logger.With(SystemField, zap.Bool("grpc_log", true))}
//	grpclog.SetLogger(zgl)
//}

// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *zapGrpcLoggerV2) InfoDepth(depth int, args ...interface{}) {
	// TODO: keep a cache of logger for different depths
	//msg := ""
	//if len(args) > 0 {
	//	if m, ok := (args[0]).(string); ok {
	//			msg = m
	//	}
	//}
	//ce := l.logger.Check(zap.InfoLevel, msg)
	//
	//ce.Write()

	// When called from componentLogger: first part is component, second is string
	l1 := l.logger.WithOptions(zap.AddCallerSkip(depth - 2))
	l1.Sugar().Info(args...)

	//l.s.Infow(args[0].(string), args[1:]...)
}

// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *zapGrpcLoggerV2) WarningDepth(depth int, args ...interface{}) {
	l.s.Warn(args...)
	//	l.s.Warnw(args[0].(string), args[1:]...)
}

// ErrorDepth logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *zapGrpcLoggerV2) ErrorDepth(depth int, args ...interface{}) {
	l.s.Error(args...)
	//	l.s.Errorw(args[0].(string), args[1:]...)
}

// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *zapGrpcLoggerV2) FatalDepth(depth int, args ...interface{}) {
	l.s.Fatalw(args[0].(string), args[1:]...)
}

// ReplaceGrpcLoggerV2 replaces the grpc_log.LoggerV2 with the provided logger.
// It should be called before any gRPC functions.
func ReplaceGrpcLoggerV2(logger *zap.Logger) {
	ReplaceGrpcLoggerV2WithVerbosity(logger, 0)
}

// ReplaceGrpcLoggerV2WithVerbosity replaces the grpc_.LoggerV2 with the provided logger and verbosity.
// It should be called before any gRPC functions.
func ReplaceGrpcLoggerV2WithVerbosity(logger *zap.Logger, verbosity int) {
	zgl := &zapGrpcLoggerV2{
		logger:    logger.With(SystemField, zap.Bool("grpc_log", true)),
		verbosity: verbosity,
	}
	zgl.s = zgl.logger.Sugar()
	grpclog.SetLoggerV2(zgl)
}

// SetGrpcLoggerV2 replaces the grpc_log.LoggerV2 with the provided logger.
// It can be used even when grpc infrastructure was initialized.
//func SetGrpcLoggerV2(settable grpc_logsettable.SettableLoggerV2, logger *zap.Logger) {
//	SetGrpcLoggerV2WithVerbosity(settable, logger, 0)
//}
//
//// SetGrpcLoggerV2WithVerbosity replaces the grpc_.LoggerV2 with the provided logger and verbosity.
//// It can be used even when grpc infrastructure was initialized.
//func SetGrpcLoggerV2WithVerbosity(settable grpc_logsettable.SettableLoggerV2, logger *zap.Logger, verbosity int) {
//	zgl := &zapGrpcLoggerV2{
//		logger:    logger.With(SystemField, zap.Bool("grpc_log", true)),
//		verbosity: verbosity,
//	}
//	settable.Set(zgl)
//}

type zapGrpcLoggerV2 struct {
	logger    *zap.Logger
	verbosity int
	s         *zap.SugaredLogger
}

func (l *zapGrpcLoggerV2) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Infoln(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *zapGrpcLoggerV2) Warning(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Warningln(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Warningf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *zapGrpcLoggerV2) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Errorln(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *zapGrpcLoggerV2) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Fatalln(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *zapGrpcLoggerV2) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *zapGrpcLoggerV2) V(level int) bool {
	return l.verbosity <= level
}
