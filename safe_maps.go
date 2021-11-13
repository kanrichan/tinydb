package main

import "sync"

type TinyTabsMap struct {
	*sync.RWMutex
	item TinyTabs
}

type TinyDocsMap struct {
	TinyTabsMap
	table string
}

type TinyRecsMap struct {
	TinyDocsMap
	id int
}

type TinyTabs map[string]TinyDocs
type TinyDocs map[int]TinyRecs
type TinyRecs map[string]TinyData
type TinyData interface{}

func (m *TinyTabsMap) Put(table string, value TinyDocs) {
	m.Lock()
	defer m.Unlock()
	if m.item == nil {
		m.item = TinyTabs{}
	}
	m.item[table] = value
}

func (m *TinyTabsMap) Get(table string) TinyDocs {
	m.RLock()
	defer m.RUnlock()
	if m.item == nil {
		return nil
	}
	return m.item[table]
}

func (m *TinyTabsMap) Docs(table string) TinyDocsMap {
	return TinyDocsMap{*m, table}
}

func (m *TinyDocsMap) Put(id int, value TinyRecs) {
	m.Lock()
	defer m.Unlock()
	m.item[m.table][id] = value
}

func (m *TinyDocsMap) Get(id int) TinyRecs {
	m.RLock()
	defer m.RUnlock()
	return m.item[m.table][id]
}

func (m *TinyDocsMap) Recs(id int) TinyRecsMap {
	return TinyRecsMap{*m, id}
}

func (m *TinyRecsMap) Put(key string, value TinyData) {
	m.Lock()
	defer m.Unlock()
	m.item[m.table][m.id][key] = value
}

func (m *TinyRecsMap) Get(key string) TinyData {
	m.RLock()
	defer m.RUnlock()
	return m.item[m.table][m.id][key]
}
