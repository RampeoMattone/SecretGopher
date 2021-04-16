package SecretGopher

// Set type is a type alias of `map[interface{}]struct{}`
type Set map[interface{}]struct{}

// add Adds an element to the set
func (s Set) add(elem interface{}) {
	s[elem] = struct{}{}
}

// addAll Adds a list of elements to the set
func (s Set) addAll(elem ...interface{}) {
	for e := range elem {
		s[e] = struct{}{}
	}
}

// remove Removes an element from the set
func (s Set) remove(elem interface{}) {
	delete(s, elem)
}

// clear Removes an element from the set
func (s Set) clear() {
	for e := range s {
		s.remove(e)
	}
}

// has Returns a boolean value describing if the element exists in the set
func (s Set) has(elem interface{}) bool {
	_, ok := s[elem]
	return ok
}
