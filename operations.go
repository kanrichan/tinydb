package main

import "errors"

func Delete(field string) func(map[string]interface{}) error {
	transform := func(doc map[string]interface{}) error {
		if _, ok := doc[field]; !ok {
			return errors.New("key not found")
		}
		delete(doc, field)
		return nil
	}
	return transform
}

func Add(field string, n int) func(map[string]interface{}) error {
	transform := func(doc map[string]interface{}) error {
		if _, ok := doc[field]; !ok {
			return errors.New("key not found")
		}
		if _, ok := doc[field].(int); !ok {
			return errors.New("val is not int")
		}
		doc[field] = doc[field].(int) + n
		return nil
	}
	return transform
}

func Subtract(field string, n int) func(map[string]interface{}) error {
	transform := func(doc map[string]interface{}) error {
		if _, ok := doc[field]; !ok {
			return errors.New("key not found")
		}
		if _, ok := doc[field].(int); !ok {
			return errors.New("val is not int")
		}
		doc[field] = doc[field].(int) - n
		return nil
	}
	return transform
}

func Set(field string, val interface{}) func(map[string]interface{}) error {
	transform := func(doc map[string]interface{}) error {
		if _, ok := doc[field]; !ok {
			return errors.New("key not found")
		}
		doc[field] = val
		return nil
	}
	return transform
}

func Increment(field string, n int) func(map[string]interface{}) error {
	transform := func(doc map[string]interface{}) error {
		if _, ok := doc[field]; !ok {
			return errors.New("key not found")
		}
		if _, ok := doc[field].(int); !ok {
			return errors.New("val is not int")
		}
		doc[field] = doc[field].(int) + 1
		return nil
	}
	return transform
}

func Decrement(field string, n int) func(map[string]interface{}) error {
	transform := func(doc map[string]interface{}) error {
		if _, ok := doc[field]; !ok {
			return errors.New("key not found")
		}
		if _, ok := doc[field].(int); !ok {
			return errors.New("val is not int")
		}
		doc[field] = doc[field].(int) - 1
		return nil
	}
	return transform
}
