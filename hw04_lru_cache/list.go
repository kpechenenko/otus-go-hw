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

// list связный список.
type list struct {
	front  *ListItem // Первый элемент
	back   *ListItem // Последний элемент
	length int       // Текущее количество элементов
}

// Len получить количество элементов в списке.
func (l *list) Len() int {
	return l.length
}

// Front получить ссылку на первый элемент списка.
func (l *list) Front() *ListItem {
	return l.front
}

// Back получить ссылку на последний элемент списка.
func (l *list) Back() *ListItem {
	return l.back
}

// PushFront добавить значение в начало списка.
func (l *list) PushFront(v interface{}) *ListItem {
	node := &ListItem{Value: v}
	if l.length == 0 {
		l.back = node
	} else {
		node.Next = l.front
		l.front.Prev = node
	}
	l.front = node
	l.length++
	return node
}

// PushBack добавить значение в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	node := &ListItem{Value: v}
	if l.length == 0 {
		l.front = node
		l.back = node
	} else {
		node.Prev = l.back
		l.back.Next = node
		l.back = node
	}
	l.length++
	return node
}

// Remove удалить элемент из списка.
func (l *list) Remove(node *ListItem) {
	if node == nil {
		return
	}
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if node == l.front {
		l.front = node.Next
	}
	if node == l.back {
		l.back = node.Prev
	}
	l.length--
}

// MoveToFront переместить элемент в начало списка.
func (l *list) MoveToFront(node *ListItem) {
	l.Remove(node)
	l.PushFront(node.Value)
}

// NewList создать новый список.
func NewList() List {
	return &list{front: nil, back: nil, length: 0}
}
