package SecretGopher

// set type is a type alias of `map[interface{}]struct{}`
type set map[interface{}]struct{}

// add Adds an element to the set
func (s set) add(elem interface{}) {
	s[elem] = struct{}{}
}

// addAll Adds a list of elements to the set
func (s set) addAll(elem ...interface{}) {
	for e := range elem {
		s[e] = struct{}{}
	}
}

// remove Removes an element from the set
func (s set) remove(elem interface{}) {
	delete(s, elem)
}

// clear Removes an element from the set
func (s set) clear() {
	for e := range s {
		s.remove(e)
	}
}

// has Returns a boolean value describing if the element exists in the set
func (s set) has(elem interface{}) bool {
	_, ok := s[elem]
	return ok
}
