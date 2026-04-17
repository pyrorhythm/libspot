package pfdomain

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
