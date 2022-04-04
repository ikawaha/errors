errors
===

The errors package is a libraries that adds context to errors.

## contexter

The contexter adds key/value context and stack traces to the error.

Setting a value to an error:
```
err := errors.New("some error")
err = contextre.WithValue(err, "hello", "goodbye")
```

Getting the value associated with this error for key:
```
val, ok := contexter.Value(err, "hello")
fmt.Println(val, ok)

OUTPUT:
goodbye, true
```

Logging stack trace to an error:
```
err = contexter.WithStackTrace(err)
```

Retrieving the record from error:
```
trace, ok := contexter.StackTrace(err)
fmt.Println(trace)

OUTPUT:
error caused by fn1:
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:12 github.com/ikawaha/verbose_test.fn1
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:16 github.com/ikawaha/verbose_test.fn2
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:20 github.com/ikawaha/verbose_test.TestWithStackTrace
    /usr/local/opt/go/libexec/src/testing/testing.go:1439 testing.tRunner
    /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1571 runtime.goexit
```

## linker

The linker allows multiple errors to be linked (embedded) together and treated as a single error. It can also expand linked errors back into multiple errors.

Linking errors:
```
lhs := errors.New("lhs")
rhs := errors.New("rhs")
e := linker.Append(lhs, rhs)

fmt.Println("errors.Is(e, lhs)=", errors.Is(e, lhs))
fmt.Println("errors.Is(e, rhs)=", errors.Is(e, rhs))
OUTPUT:
errors.Is(e, lhs)= true
errors.Is(e, rhs)= true
```

Yield multiple errors:
```
errs := linker.Yield(e)
fmt.Println("linker.Yield(e)=", errs)
OUTPUT:
linker.Yield(e)== [lhs rhs]
```

---
MIT