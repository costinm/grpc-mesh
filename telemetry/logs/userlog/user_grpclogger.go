package userlog

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

type userGrpcLoggerV2 struct {
	logger    *zap.Logger
	verbosity int
	s         *zap.SugaredLogger
}

type LogEvent struct {
	Time      uint64
	Level     int
	Verbosity int

	Source string
	Msg    string
	Args   []interface{}
}

func (l *userGrpcLoggerV2) Subscribe() {

}

func (l *userGrpcLoggerV2) Publish() {

}

// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *userGrpcLoggerV2) InfoDepth(depth int, args ...interface{}) {
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
func (l *userGrpcLoggerV2) WarningDepth(depth int, args ...interface{}) {
	l.s.Warn(args...)
	//	l.s.Warnw(args[0].(string), args[1:]...)
}

// ErrorDepth logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *userGrpcLoggerV2) ErrorDepth(depth int, args ...interface{}) {
	l.s.Error(args...)
	//	l.s.Errorw(args[0].(string), args[1:]...)
}

// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Println.
func (l *userGrpcLoggerV2) FatalDepth(depth int, args ...interface{}) {
	l.s.Fatalw(args[0].(string), args[1:]...)
}

func ReplaceGrpcLogger() {
	zgl := &userGrpcLoggerV2{}
	zgl.s = zgl.logger.Sugar()
	grpclog.SetLoggerV2(zgl)
}

func (l *userGrpcLoggerV2) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Infoln(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *userGrpcLoggerV2) Warning(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Warningln(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Warningf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *userGrpcLoggerV2) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Errorln(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *userGrpcLoggerV2) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Fatalln(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *userGrpcLoggerV2) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *userGrpcLoggerV2) V(level int) bool {
	return l.verbosity <= level
}
