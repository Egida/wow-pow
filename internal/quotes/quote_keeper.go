package quotes

import (
	"math/rand"
)

type QuoteKeeper struct {
	quotes []string
}

func New(quotes []string) *QuoteKeeper {
	return &QuoteKeeper{
		quotes: quotes,
	}
}

func (q *QuoteKeeper) Quote() string {
	i := rand.Intn(len(q.quotes))

	return q.quotes[i]
}
