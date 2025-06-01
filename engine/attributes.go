package engine

import "sync"

// Attributes contains different mutable attributes
type Attributes struct {
	values map[interface{}]interface{}
	mu     sync.RWMutex
}

// NewAttributes return new attributes
func NewAttributes() *Attributes {
	return &Attributes{
		values: make(map[interface{}]interface{}),
	}
}

// SetAttribute set attribute
func (a *Attributes) SetAttribute(name, value interface{}) {
	a.mu.Lock()
	a.values[name] = value
	a.mu.Unlock()
}

// UpdateAttribute update attribute
func (a *Attributes) UpdateAttribute(name, f func(interface{}) interface{}) {
	a.mu.Lock()
	a.values[name] = f(a.values[name])
	a.mu.Unlock()
}

// GetAttribute return raw attribute
func (a *Attributes) GetAttribute(name interface{}) (value interface{}, exists bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	value, exists = a.values[name]

	return
}

// RemoveAttribute remove given attribute
func (a *Attributes) RemoveAttribute(name interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.values, name)
}

// GetUint32 get attribute and cast to uint32
func (a *Attributes) GetUint32(name interface{}) uint32 {
	v, e := a.GetAttribute(name)
	if !e {
		return 0
	}

	return v.(uint32)
}

// GetUint64 get attribute and cast to uint64
func (a *Attributes) GetUint64(name interface{}) uint64 {
	v, e := a.GetAttribute(name)
	if !e {
		return 0
	}

	return v.(uint64)
}

// GetFloat64 get attribute and cast to float64
func (a *Attributes) GetFloat64(name interface{}) float64 {
	v, e := a.GetAttribute(name)
	if !e {
		return 0
	}

	return v.(float64)
}

// GetInt get attribute and cast to int
func (a *Attributes) GetInt(name interface{}) int {
	v, e := a.GetAttribute(name)
	if !e {
		return 0
	}

	return v.(int)
}

// GetUint8 get attribute and cast to uint8
func (a *Attributes) GetUint8(name interface{}) uint8 {
	v, e := a.GetAttribute(name)
	if !e {
		return 0
	}

	return v.(uint8)
}

// GetBool get attribute and cast to bool
func (a *Attributes) GetBool(name interface{}) bool {
	v, e := a.GetAttribute(name)
	if !e {
		return false
	}

	return v.(bool)
}

// GetString get attribute and cast to string
func (a *Attributes) GetString(name interface{}) string {
	v, e := a.GetAttribute(name)
	if !e {
		return ""
	}

	return v.(string)
}
