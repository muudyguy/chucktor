package tree


type TreeInterface interface {
	Search(value interface{}) interface{}
	Add(value interface{})
}

type Comparable interface {
	Bigger(in Comparable) bool
	Equals(in Comparable) bool
}


