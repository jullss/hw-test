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
	count int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.count == 0 {
		l.front = newItem
		l.back = newItem
	} else {
		l.front.Prev = newItem
		newItem.Next = l.front
		newItem.Prev = nil
		l.front = newItem
	}

	l.count++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.count == 0 {
		l.front = newItem
		l.back = newItem
	} else {
		l.back.Next = newItem
		newItem.Prev = l.back
		l.back = newItem
	}

	l.count++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.count--

	i.Prev = nil
	i.Next = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if l.count <= 1 {
		return
	}

	if i.Prev == nil {
		return
	}

	if i.Next == nil {
		i.Prev.Next = nil
		l.back = i.Prev
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.front.Prev = i
	i.Next = l.front
	i.Prev = nil
	l.front = i
}

func NewList() List {
	return new(list)
}
