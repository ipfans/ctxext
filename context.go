package ctxext

import (
	"context"
	"sync"
)

type contextKey struct{}

type Context struct {
	context.Context

	lock sync.RWMutex
	m    map[string]interface{}
}

var ctxKey = contextKey{}

// New generate *Context based on context.Context
func New(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Context{
		Context: ctx,
		m:       map[string]interface{}{},
	}
}

// Copy or new *Context based on context.Context
func Copy(ctx context.Context) *Context {
	if ctx == nil {
		return New(nil)
	}
	c := New(ctx)
	val := ctx.Value(ctxKey)
	if val != nil {
		c.m = val.(map[string]interface{})
	}
	return c
}

// Set value to store.
// WARN: Take care if you are storing map/slice etc.
func (c *Context) Set(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.m[key] = val
}

// Exists check key status.
func (c *Context) Exists(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.m[key]
	return ok
}

// Value for context.Context interface.
func (c *Context) Value(i interface{}) interface{} {
	key, ok := i.(string)
	if !ok {
		if _, ok := i.(contextKey); ok {
			c.lock.RLock()
			defer c.lock.RUnlock()
			newMap := make(map[string]interface{}, len(c.m))
			for k, v := range c.m {
				newMap[k] = v
			}
			return newMap
		}
		return c.Context.Value(i)
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	v, ok := c.m[key]
	if !ok {
		return c.Context.Value(i)
	}
	return v
}
