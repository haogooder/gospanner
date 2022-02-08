package middleware

import (
    "context"
    "github.com/haogooder/gospanner/mongo/field"
    "github.com/haogooder/gospanner/mongo/hook"
    "github.com/haogooder/gospanner/mongo/operator"
    "github.com/haogooder/gospanner/mongo/validator"
)

// callback define the callback function type
type callback func(ctx context.Context, doc interface{}, opType operator.OpType, opts ...interface{}) error

// middlewareCallback the register callback slice
// some callbacks initial here without Register() for order
var middlewareCallback = []callback{
    hook.Do,
    field.Do,
    validator.Do,
}

// Register register callback into middleware
func Register(cb callback) {
    middlewareCallback = append(middlewareCallback, cb)
}

// Do call every registers
// The doc is always the document to operate
func Do(ctx context.Context, doc interface{}, opType operator.OpType, opts ...interface{}) error {
    for _, cb := range middlewareCallback {
        if err := cb(ctx, doc, opType, opts...); err != nil {
            return err
        }
    }
    return nil
}
