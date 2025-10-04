package internal

import "fmt"

type CompareFunc[T any] = func(a PriorityItem[T], b PriorityItem[T]) bool

type PriorityItem[T any] struct {
	Value    T
	Priority int
}

type MinHeap[T any] struct {
	items    []PriorityItem[T]
	size     int
	capacity int
	compare  CompareFunc[T]
}

type PriorityQueue[T any] struct {
	heap *MinHeap[T]
}

func DefaultMinHeapComparator[T any](a PriorityItem[T], b PriorityItem[T]) bool {
	return a.Priority < b.Priority
}

func MaxHeapComparator[T any](a PriorityItem[T], b PriorityItem[T]) bool {
	return a.Priority > b.Priority
}

func NewMinHeapWithComparator[T any](capacity int, compareFn CompareFunc[T]) *MinHeap[T] {
	return &MinHeap[T]{
		items:    make([]PriorityItem[T], capacity),
		size:     0,
		capacity: capacity,
		compare:  compareFn,
	}
}

func NewMinHeap[T any](capacity int, compareFn CompareFunc[T]) *MinHeap[T] {
	return NewMinHeapWithComparator(capacity, DefaultMinHeapComparator[T])
}

func (mh *MinHeap[T]) HeapifyDown(index int) {
	for mh.hasChildren(index) {
		higherPriorityChildIndex := mh.getHigherPriorityChildIndex(index)
		if mh.compare(mh.items[higherPriorityChildIndex], mh.items[index]) {
			mh.items[index], mh.items[higherPriorityChildIndex] = mh.items[higherPriorityChildIndex], mh.items[index]
			index = higherPriorityChildIndex
		} else {
			break
		}
	}
}

func (mh *MinHeap[T]) hasChildren(index int) bool {
	return ((2 * index) + 1) < mh.size // at least left child is there
}

func (mh *MinHeap[T]) getHigherPriorityChildIndex(index int) int {
	leftChildIndex := (2 * index) + 1
	rightChildIndex := (2 * index) + 2
	highestPriorityIndex := leftChildIndex

	if rightChildIndex < mh.size && mh.compare(mh.items[rightChildIndex], mh.items[leftChildIndex]) {
		highestPriorityIndex = rightChildIndex
	}

	return highestPriorityIndex
}

func (mh *MinHeap[T]) ExtractMin() (PriorityItem[T], error) {
	if mh.size == 0 {
		return PriorityItem[T]{}, fmt.Errorf("heap is empty")
	}
	min := mh.items[0]
	mh.items[0] = mh.items[mh.size-1]
	mh.HeapifyDown(0)
	mh.size--
	return min, nil
}

func (mh *MinHeap[T]) HeapifyUp(index int) {
	for index > 0 {
		parentIndex := mh.getParentIndex(index)
		if mh.compare(mh.items[index], mh.items[parentIndex]) {
			mh.items[index], mh.items[parentIndex] = mh.items[parentIndex], mh.items[index]
			index = parentIndex
		} else {
			break
		}

	}
}

func (mh *MinHeap[T]) getParentIndex(index int) int {
	return (index - 1) / 2
}

func (mh *MinHeap[T]) Insert(item PriorityItem[T]) error {
	if mh.size >= mh.capacity {
		return fmt.Errorf("heap is full")
	}

	mh.items[mh.size] = item
	mh.HeapifyUp(mh.size)
	mh.size++
	return nil
}

func (mh *MinHeap[T]) IsEmpty() bool {
	return mh.size == 0
}

func (mh *MinHeap[T]) Size() int {
	return mh.size
}

func (mh *MinHeap[T]) Peek() (PriorityItem[T], error) {
	if mh.size == 0 {
		return PriorityItem[T]{}, fmt.Errorf("heap is empty")
	}
	return mh.items[0], nil
}

func NewPriorityQueueWithComparator[T any](capacity int, compareFn CompareFunc[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: NewMinHeapWithComparator(capacity, compareFn),
	}
}

func NewPriorityQueue[T any](capacity int) *PriorityQueue[T] {
	return NewPriorityQueueWithComparator(capacity, DefaultMinHeapComparator[T])
}

func NewMaxHeapPriorityQueue[T any](capacity int) *PriorityQueue[T] {
	return NewPriorityQueueWithComparator(capacity, MaxHeapComparator[T])
}

func (pq *PriorityQueue[T]) Enqueue(value T, priority int) error {
	item := PriorityItem[T]{
		Value:    value,
		Priority: priority,
	}
	return pq.heap.Insert(item)
}

func (pq *PriorityQueue[T]) Dequeue() (T, error) {
	item, err := pq.heap.ExtractMin()
	var zero T
	if err != nil {
		return zero, err
	}

	return item.Value, nil
}

func (pq *PriorityQueue[T]) Peek() (T, error) {
	item, err := pq.heap.Peek()
	var zero T
	if err != nil {
		return zero, err
	}

	return item.Value, nil
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.heap.IsEmpty()
}

func (pq *PriorityQueue[T]) Size() int {
	return pq.heap.Size()
}
