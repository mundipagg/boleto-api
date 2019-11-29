package tracer

import (
	"sync"
	"time"
)

const (
	Fatal         = uint8(0)
	Alert         = uint8(1)
	Critical      = uint8(2)
	Error         = uint8(3)
	Warning       = uint8(4)
	Notice        = uint8(5)
	Informational = uint8(6)
	Debug         = uint8(7)
	Unset         = uint8(255)
)

var LevelNames = []string{
	"FATAL",
	"ALERT",
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFORMATIONAL",
	"DEBUG",
}

func RegisterWriter(writer Writer) {
	DefaultContext.RegisterWriter(writer)
}

func GetLogger(owner string, ctx ...Context) Logger {
	ctx = append(ctx, DefaultContext)
	return ctx[0].GetLogger(owner)
}

type Logger interface {
	Debug(message string, args ...interface{})
	D(message string, args ...interface{})
	Info(message string, args ...interface{})
	I(message string, args ...interface{})
	Notice(message string, args ...interface{})
	N(message string, args ...interface{})
	Warn(message string, args ...interface{})
	W(message string, args ...interface{})
	Error(message string, args ...interface{})
	E(message string, args ...interface{})
	Critical(message string, args ...interface{})
	C(message string, args ...interface{})
	Alert(message string, args ...interface{})
	A(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	F(message string, args ...interface{})
	Trace(transactionId ...string) Logger
	Commit()
	RegisterWriter(writer Writer)
	MinimumLevel(level uint8)
	GetMinimumLevel() uint8
}

type logger struct {
	sync.Locker
	Context
	transactionId string
	owner         string
}

func (l *logger) Debug(message string, args ...interface{}) {
	l.log(Debug, message, args)
}

func (l *logger) D(message string, args ...interface{}) {
	l.log(Debug, message, args)
}

func (l *logger) Info(message string, args ...interface{}) {
	l.log(Informational, message, args)
}

func (l *logger) I(message string, args ...interface{}) {
	l.log(Informational, message, args)
}

func (l *logger) Notice(message string, args ...interface{}) {
	l.log(Notice, message, args)
}

func (l *logger) N(message string, args ...interface{}) {
	l.log(Notice, message, args)
}

func (l *logger) Warn(message string, args ...interface{}) {
	l.log(Warning, message, args)
}

func (l *logger) W(message string, args ...interface{}) {
	l.log(Warning, message, args)
}

func (l *logger) Error(message string, args ...interface{}) {
	l.log(Error, message, args)
}

func (l *logger) E(message string, args ...interface{}) {
	l.log(Error, message, args)
}

func (l *logger) Critical(message string, args ...interface{}) {
	l.log(Critical, message, args)
}

func (l *logger) C(message string, args ...interface{}) {
	l.log(Critical, message, args)
}

func (l *logger) Alert(message string, args ...interface{}) {
	l.log(Alert, message, args)
}

func (l *logger) A(message string, args ...interface{}) {
	l.log(Alert, message, args)
}

func (l *logger) Fatal(message string, args ...interface{}) {
	l.log(Fatal, message, args)
}

func (l *logger) F(message string, args ...interface{}) {
	l.log(Fatal, message, args)
}

func (l *logger) Trace(transactionId ...string) Logger {
	if l.transactionId == "" {
		transactionId = append(transactionId, tracer.GetActiveTransaction())
	} else {
		transactionId = append(transactionId, l.transactionId)
	}
	tId := transactionId[0]
	tracer.BeginTransaction(tId)
	return &logger{
		Locker:        &sync.RWMutex{},
		Context:       l.Context,
		owner:         l.owner,
		transactionId: tId,
	}
}

func (l *logger) Commit() {
	tracer.CommitTransaction(tracer.GetActiveTransaction())
	l.transactionId = ""
}

func (l *logger) log(level uint8, message string, args []interface{}) {
	transactionId := tracer.GetActiveTransaction()
	if l.transactionId != "" {
		transactionId = l.transactionId
	}
	if level > l.GetMinimumLevel() {
		return
	}
	entry := Entry{
		Owner:         l.owner,
		Level:         level,
		Message:       message,
		Args:          args,
		Time:          time.Now(),
		TransactionId: transactionId,
		StackTrace:    GetStackTrace(3),
	}
	for _, writer := range l.GetWriters() {
		writer.Write(entry)
	}
}
