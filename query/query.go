package query

import (
	"bytes"
	"strconv"
)

type Query struct {
	buf *bytes.Buffer
}

func NewQuery() *Query {
	b := make([]byte, 0, 1024)

	return &Query{
		buf: bytes.NewBuffer(b),
	}
}

func (q *Query) separate() {
	if q.buf.Len() > 0 {
		q.buf.WriteByte('&')
	}
}

func (q *Query) WriteString(key string, value string) {
	q.separate()
	q.buf.WriteString(key)
	q.buf.WriteByte('=')
	q.buf.WriteString(value)
}

func (q *Query) WriteInt(key string, value int64) {
	s := strconv.FormatInt(value, 10)
	q.WriteString(key, s)
}

func (q *Query) WriteUint(key string, value uint) {
	s := strconv.FormatUint(uint64(value), 10)
	q.WriteString(key, s)
}

func (q *Query) WriteBool(key string, value bool) {
	s := strconv.FormatBool(value)
	q.WriteString(key, s)
}

func (q *Query) WriteFloat(key string, value float64) {
	s := strconv.FormatFloat(value, 'f', 2, 64)
	q.WriteString(key, s)
}

func (q *Query) AppendQuery(query *Query) {
	q.separate()
	query.buf.WriteTo(q.buf)
}

func (q *Query) String() string {
	return q.buf.String()
}

func (q *Query) Bytes() []byte {
	return q.buf.Bytes()
}
