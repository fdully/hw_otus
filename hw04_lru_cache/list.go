package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	// Place your code here
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(interface{}) *listItem
	PushBack(interface{}) *listItem
	Remove(*listItem)
	MoveToFront(*listItem)
}

type listItem struct {
	// Place your code here
	Value      interface{}
	Prev, Next *listItem
}

type list struct {
	// Place your code here
	front, back *listItem
	length      int
}

func NewList() List {
	return &list{}
}

func (l list) Len() int {
	return l.length
}

func (l list) Front() *listItem {
	if l.length == 0 {
		return nil
	}
	return l.front
}

func (l list) Back() *listItem {
	if l.length == 0 {
		return nil
	}
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	item := &listItem{
		Value: v,
	}

	if l.front == nil {
		l.front = item
		l.back = item
		l.length++
		return l.front
	}

	t := l.front
	l.front = item
	l.front.Next = nil
	l.front.Prev = t
	t.Next = l.front
	l.length++
	return l.front
}

func (l *list) PushBack(v interface{}) *listItem {
	item := &listItem{
		Value: v,
	}

	if l.back == nil {
		l.front = item
		l.back = item
		l.length++
		return l.back
	}

	t := l.back
	l.back = item
	l.back.Next = t
	l.back.Prev = nil
	t.Prev = l.back
	l.length++
	return l.back
}

func (l *list) Remove(item *listItem) {
	if l.length == 0 {
		return
	}
	l.remove(item)
}

func (l *list) MoveToFront(item *listItem) {
	if l.length == 0 || l.front == item {
		return
	}

	if item == l.back {
		l.back = item.Next
		l.back.Prev = nil
	} else {
		item.Prev.Next = item.Next
		item.Next.Prev = item.Prev
	}

	t := l.front
	l.front = item
	l.front.Next = nil
	l.front.Prev = t
	t.Next = l.front
}

func (l *list) remove(item *listItem) {
	if item == nil {
		return
	}

	if item == l.back {
		l.back.Prev = nil
		l.back = l.back.Next
		l.length--
		return
	}

	if item == l.front {
		l.front.Next = nil
		l.front = l.front.Prev
		l.length--
		return
	}

	item.Prev.Next = item.Next
	item.Next.Prev = item.Prev
	item.Next = nil
	item.Prev = nil
	l.length--
}
