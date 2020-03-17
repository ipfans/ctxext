package ctxext

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestContext_Set(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name  string
		args  args
		wantM map[string]interface{}
	}{
		{
			"string",
			args{
				"str",
				"111",
			},
			map[string]interface{}{"str": "111"},
		},
		{
			"int",
			args{
				"str",
				111,
			},
			map[string]interface{}{"str": 111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := New(nil)
			ctx.Set(tt.args.key, tt.args.val)
			if !reflect.DeepEqual(ctx.m, tt.wantM) {
				t.Errorf("Set() = %v, want %v", ctx.m, tt.wantM)
			}
		})
	}
}

func TestContext_Value(t *testing.T) {
	type fields struct {
		ctx context.Context
		key interface{}
		val interface{}
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			"Without value",
			fields{
				context.TODO(),
				nil,
				nil,
			},
			args{
				"not_exists_data",
			},
			nil,
		},
		{
			"Value only in map",
			fields{
				context.TODO(),
				"data",
				"123",
			},
			args{
				"data",
			},
			"123",
		},
		{
			"Value only in context chain",
			fields{
				context.WithValue(context.TODO(), "data", "123"),
				nil,
				nil,
			},
			args{
				"data",
			},
			"123",
		},
		{
			"Value in map and context chain",
			fields{
				context.WithValue(context.TODO(), "data", "456"),
				"data",
				"123",
			},
			args{
				"data",
			},
			"123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.ctx)
			if tt.fields.key != nil {
				c.Set(tt.fields.key.(string), tt.fields.val)
			}
			if got := c.Value(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_ValueWithWrap(t *testing.T) {
	type fields struct {
		ctx context.Context
		key interface{}
		val interface{}
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			"Without value",
			fields{
				context.TODO(),
				nil,
				nil,
			},
			args{
				"not_exists_data",
			},
			nil,
		},
		{
			"Value only in map",
			fields{
				context.TODO(),
				"data",
				"123",
			},
			args{
				"data",
			},
			"123",
		},
		{
			"Value in map and insider",
			fields{
				context.WithValue(context.TODO(), "data", "456"),
				"data",
				"123",
			},
			args{
				"data",
			},
			"123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.ctx)
			if tt.fields.key != nil {
				c.Set(tt.fields.key.(string), tt.fields.val)
			}
			ctx, _ := context.WithCancel(c)
			if got := ctx.Value(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_Exists(t *testing.T) {
	type fields struct {
		ctx context.Context
		key interface{}
		val interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"Context not exists.",
			fields{
				context.TODO(),
				nil,
				nil,
			},
			args{
				"not_exists",
			},
			false,
		},
		{
			"Exists in map",
			fields{
				context.TODO(),
				"data",
				"123",
			},
			args{
				"data",
			},
			true,
		},
		{
			"Exists only in chain",
			fields{
				context.WithValue(context.TODO(), "data", "123"),
				nil,
				nil,
			},
			args{
				"data",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := New(tt.fields.ctx)
			if tt.fields.key != nil {
				ctx.Set(tt.fields.key.(string), tt.fields.val)
			}
			if !reflect.DeepEqual(ctx.Exists(tt.args.key), tt.want) {
				t.Errorf("Exists() = %v, want %v", ctx.m, tt.want)
			}
		})
	}
}

func TestCopy_Exception(t *testing.T) {
	ctx := New(context.TODO())
	ctx.Set("data", "123")
	newctx := Copy(ctx)
	newctx.Set("data", "456")
	if !reflect.DeepEqual(ctx.Value("data"), "123") {
		t.Errorf("Old context Value() = %v, want %v", ctx.Value("data"), "123")
		return
	}
	if !reflect.DeepEqual(newctx.Value("data"), "456") {
		t.Errorf("Old context Value() = %v, want %v", ctx.Value("data"), "456")
		return
	}

	// map will be modified.
	m := map[string]string{
		"aaa": "bbb",
	}
	ctx.Set("map", m)
	newctx = Copy(ctx)
	mm := newctx.Value("map").(map[string]string)
	mm["aaa"] = "ccc"
	oldv := ctx.Value("map").(map[string]string)["aaa"]
	if !reflect.DeepEqual(oldv, "ccc") {
		t.Errorf("Old context Value() = %v, want %v", oldv, "ccc")
	}
	newv := newctx.Value("map").(map[string]string)["aaa"]
	if !reflect.DeepEqual(newv, "ccc") {
		t.Errorf("Old context Value() = %v, want %v", newv, "ccc")
	}
}

func TestCopy1(t *testing.T) {
	ctx := New(context.TODO())
	ctx.Set("data", "123")
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *Context
	}{
		{
			"nil context",
			args{
				nil,
			},
			&Context{
				context.TODO(),
				sync.RWMutex{},
				map[string]interface{}{},
			},
		},
		{
			"context wrap",
			args{
				ctx,
			},
			&Context{
				ctx,
				sync.RWMutex{},
				map[string]interface{}{
					"data": "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.args.ctx
			if got := Copy(ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}
