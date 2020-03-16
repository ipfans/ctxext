package ctxext

import (
	"context"
	"sync"
)

type Context struct {
	context.Context

	lock sync.RWMutex
	m    map[string]interface{}
}

func New(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
		m:       map[string]interface{}{},
	}
}

func (c *Context) Set(key string, val interface{}) {
	c.lock.Lock()
	c.m[key] = val
	c.lock.Unlock()
}

func (c *Context) Exists(key string) bool {
	c.lock.Lock()
	_, ok := c.m[key]
	c.lock.Unlock()
	return ok
}

func (c *Context) Value(i interface{}) interface{} {
	key, ok := i.(string)
	if !ok {
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
