package tinydb

// Iterator TinyRecs convert to iterator.
func (arr TinyRecsArr) Iterator() *TinyRecsIter {
	return &TinyRecsIter{0, arr}
}

// Next Get TinyRecsIter Next item.
func (i *TinyRecsIter) Next() TinyRecs {
	if i.index >= len(i.item) {
		return nil
	}
	i.index++
	return i.item[i.index-1]
}

// HasNext TinyRecsIter has next item.
func (i *TinyRecsIter) HasNext() bool {
	return i.index < len(i.item)
}
