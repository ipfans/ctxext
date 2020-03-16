package ctxext

import (
	"context"
	"reflect"
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
			ctx := New(context.TODO())
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
	type args struct{
		key string
	}
	tests := []struct {
		name  string
		fields fields
		args  args
		want bool
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
