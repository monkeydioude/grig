package model

type IndexBuilder interface {
	SetIndex(int)
	GetIndex() int
	GetParent() IndexBuilder
	SetParent(IndexBuilder)
	GetName() string
}
