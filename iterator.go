package tinydb

func (arr TinyRecsArr) Iterator() *TinyRecsIter {
	return &TinyRecsIter{0, arr}
}

func (i *TinyRecsIter) Next() TinyRecs {
	if i.index >= len(i.item) {
		return nil
	}
	i.index++
	return i.item[i.index-1]
}

func (i *TinyRecsIter) HasNext() bool {
	return i.index < len(i.item)
}
