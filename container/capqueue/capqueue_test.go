package capqueue_test

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/gohornet/hornet/pkg/model/mselection/container/capqueue"
	"github.com/stretchr/testify/assert"
)

const testCapacity = 10

func TestNew(t *testing.T) {
	q := New(testCapacity)
	assert.Equal(t, 0, q.Len())
	assert.Equal(t, testCapacity, q.Cap())
}

func TestCapQueue_Max(t *testing.T) {
	q := New(testCapacity)
	assert.Panics(t, func() { _, _ = q.Max() })

	q.Add("1", 1)
	maxKey, maxValue := q.Max()
	assert.Equal(t, "1", maxKey)
	assert.Equal(t, 1, maxValue)
}

func TestCapQueue_Add(t *testing.T) {
	q := New(testCapacity)
	for i := 1; i <= testCapacity+1; i++ {
		q.Add(fmt.Sprint(i), i)
	}
	assert.Equal(t, testCapacity, q.Len())

	_, max := q.Max()
	assert.Equal(t, max, testCapacity+1)
}

func TestCapQueue_Delete(t *testing.T) {
	q := New(testCapacity)
	for i := 1; i <= testCapacity; i++ {
		q.Add(fmt.Sprint(i), i)
	}

	assert.False(t, q.Delete("not contained"))

	for i := testCapacity - 1; i >= 0; i-- {
		maxKey, _ := q.Max()
		assert.True(t, q.Delete(maxKey))
		assert.Equal(t, i, q.Len())
	}
}

func TestCapQueue_Value(t *testing.T) {
	q := New(testCapacity)
	for i := 1; i <= testCapacity; i++ {
		q.Add(fmt.Sprint(i), i)
	}

	assert.Zero(t, q.Value("not contained"))

	for i := 1; i <= testCapacity; i++ {
		assert.Equal(t, i, q.Value(fmt.Sprint(i)))
	}
}

func BenchmarkCapQueue_Add(b *testing.B) {
	q := New(b.N)
	// prepare random adds
	data := make([]int, b.N)
	for i := range data {
		data[i] = rand.Int()
	}
	b.ResetTimer()

	for i := range data {
		q.Add("", data[i])
	}
}

func BenchmarkCapQueue_FullAdd(b *testing.B) {
	// create a queue full of random values
	q := New(b.N)
	for i := 0; i < b.N; i++ {
		v := rand.Intn(b.N)
		q.Add(fmt.Sprint(v), v)
	}
	// prepare random adds
	data := make([]int, b.N)
	for i := range data {
		data[i] = rand.Int()
	}
	b.ResetTimer()

	for i := range data {
		q.Add("", data[i])
	}
}

func BenchmarkCapQueue_Delete(b *testing.B) {
	// create a full queue
	q := New(b.N)
	for i := 0; i < b.N; i++ {
		q.Add(fmt.Sprint(i), i)
	}
	// prepare deletes in random order
	data := make([]string, b.N)
	for i := range data {
		data[i] = fmt.Sprint(i)
	}
	rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
	b.ResetTimer()

	for i := range data {
		q.Delete(data[i])
	}
}
