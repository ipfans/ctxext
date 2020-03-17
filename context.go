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

func New(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
		m:       map[string]interface{}{},
	}
}

func Copy(ctx context.Context) *Context {
	c := New(ctx)
	val := ctx.Value(ctxKey)
	if val != nil {
		c.m = val.(map[string]interface{})
	}
	return c
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
		if _, ok := i.(contextKey); ok {
			c.lock.RLock()
			defer c.lock.RUnlock()
			newMap := make(map[string]interface{}, len(c.m))
			for k := range c.m {
				newMap[k] = c.m[k]
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
