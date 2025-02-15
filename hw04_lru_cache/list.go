package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.first == nil {
		l.first = newItem
		l.last = newItem
		l.len = 1
		return newItem
	}

	firstItem := l.first
	l.first = newItem
	newItem.Next = firstItem
	firstItem.Prev = newItem
	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.first == nil {
		l.first = newItem
		l.last = newItem
		l.len++
		return newItem
	}

	lastItem := l.last
	l.last = newItem
	newItem.Prev = lastItem
	lastItem.Next = newItem
	l.len++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	prevItem := i.Prev
	nextItem := i.Next
	prevItem.Next = nextItem
	nextItem.Prev = prevItem
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.first {
		return
	}

	prevItem := i.Prev
	nextItem := i.Next
	prevItem.Next = nextItem
	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	firstItem := l.first
	i.Next = firstItem
	firstItem.Prev = i
	i.Prev = nil
	l.first = i
}
