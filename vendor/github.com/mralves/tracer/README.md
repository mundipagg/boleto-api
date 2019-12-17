# Tracer
Easy to use and extend log library

## How to install
Using go get (not recommended):
```bash
go get github.com/mralves/tracer
```

Using [dep](github.com/golang/dep) (recommended):
```bash
dep ensure --add github.com/stretchr/testify@<version>
```

## How to use

Below follows a simple example of how to use this lib:

```go
package main

import (
	"github.com/mralves/tracer"
	"os"
	"time"
)

func inner() {
	logger := tracer.GetLogger("moduleA.inner")
	logger.Info("don't know which transaction is this")
	logger.Info("but this log in this transaction")
	logger = logger.Trace()
	go func() {
		logger.Info("this is also inside the same transaction")
		func() {
			logger := tracer.GetLogger("moduleA.inner.nested")
			logger.Info("but not this one...")

		}()
	}()
}

func main() {
	logger := tracer.GetLogger("moduleA")
	tracer.RegisterWriter(tracer.NewFileWriter(os.Stdout, tracer.SimpleFormatter("message='@message' transaction=@transactionId\n", nil)))
	logger.Info("logging in transaction 'A'", "A")
	logger.Info("logging in transaction 'B'", "B")
	logger.Info("logging in transaction 'B'", "B")
	logger.Info("logging in transaction 'A'", "A")
	logger.Info("logging in transaction 'A'", "A")
	logger = logger.Trace("C") // now all logs on this logger will be on the transaction C
	logger.Info("logging in transaction 'C'")
	logger.Info("logging in transaction 'C'", "A")
	inner()
	time.Sleep(50 * time.Millisecond)
}
```

