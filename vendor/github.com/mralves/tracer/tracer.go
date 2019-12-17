package tracer

import (
	"github.com/v2pro/plz/gls"
	"sync"
)

type IntSet map[int64]bool


var tracer = Tracer{
	Locker: &sync.RWMutex{},
	ActiveTransactions: map[string]IntSet{},
	TransactionParticipants: map[int64]string{},
}

func (set IntSet) Add(x int64) bool {
	exists := set.Contains(x)
	set[x] = true
	return !exists
}

func (set IntSet) Contains(x int64) bool {
	_, contains := set[x]
	return contains
}

type Tracer struct {
	sync.Locker
	ActiveTransactions      map[string]IntSet
	TransactionParticipants map[int64]string
}

func (tracer *Tracer) BeginTransaction(transaction string) {
	tracer.Lock()
	defer tracer.Unlock()
	if _, exists := tracer.ActiveTransactions[transaction]; !exists {
		tracer.ActiveTransactions[transaction] = IntSet{}
	}
	goID := gls.GoID()
	tracer.ActiveTransactions[transaction].Add(goID)
	tracer.TransactionParticipants[goID] = transaction
}

func (tracer *Tracer) CommitTransaction(transaction string) {
	tracer.Lock()
	defer tracer.Unlock()
	for participant := range tracer.ActiveTransactions[transaction] {
		delete(tracer.TransactionParticipants, participant)
	}
	delete(tracer.ActiveTransactions, transaction)
}

func (tracer *Tracer) GetActiveTransaction() string {
	tracer.Lock()
	defer tracer.Unlock()
	goID := gls.GoID()
	return tracer.TransactionParticipants[goID]
}