package pfdomain

import "pyrorhythm.dev/fn"

type Data[T any] struct {
	Data *T `json:"data"`
}

func (d Data[T]) Get() *T {
	return d.Data
}

type Item[T any] struct {
	Item *T `json:"item"`
}

func (d Item[T]) Get() *T {
	return d.Item
}

type MatchedItem[T any] struct {
	Item          *Data[T] `json:"item"`
	MatchedFields []any    `json:"matchedFields"`
}

func (m MatchedItem[T]) Get() *T {
	if m.Item == nil {
		return nil
	}
	return m.Item.Data
}

type MatchedList[T any] struct {
	Items []*MatchedItem[T] `json:"items"`
}

func (l MatchedList[T]) GetMatched() []*MatchedItem[T] {
	return l.Items
}

func (l MatchedList[T]) Get() []*T {
	return fn.Map(l.Items, func(i *MatchedItem[T]) *T { return i.Get() })
}

type WrappedItem[T any] struct {
	Typename string `json:"__typename"`
	Data     *T     `json:"data"`
}

func (w WrappedItem[T]) Get() *T {
	return w.Data
}

type WrappedList[T any] struct {
	Items []*WrappedItem[T] `json:"items"`
}

func (l WrappedList[T]) GetWrapped() []*WrappedItem[T] {
	return l.Items
}

func (l WrappedList[T]) Get() []*T {
	return fn.Map(l.Items, func(i *WrappedItem[T]) *T { return i.Get() })
}

type ItemList[T any] struct {
	Items []*T `json:"items"`
}

type ItemCountList[T any] struct {
	ItemList[T]
	TotalCount
}

func (i ItemList[T]) Get() []*T {
	return i.Items
}

type ItemV2List[T any] struct {
	Items []*T `json:"itemsV2"`
}

func (i ItemV2List[T]) Get() []*T {
	return i.Items
}

type TotalCount struct {
	TotalCount int `json:"totalCount"`
}

func (c TotalCount) Count() int {
	return c.TotalCount
}
