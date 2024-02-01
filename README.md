# Little Interpreter

## Examples

```
func fib(n)
    if n <= 1
        return n
    end

    return fib(n - 1) + fib(n - 2)
end

func main()
    return fib(9)
end
```
