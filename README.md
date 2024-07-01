# Little Interpreter

## Specification

### Integer types

```
u8 u16 u32 u64
i8 i16 i32 i64
```

The value of an n-bit integer is n bits wide and represented using two's
complement arithmetic. Explicit conversions are required when different integer
types are mixed in an expression or assignment. When converting between integer
types, if the value is a signed integer, it is sign extended to implicit
infinite precision; otherwise it is zero extended. It is then truncated to fit
in the result type's size.

### Operators

```
+  -  *  /  %
&  |  ^
<<  >>
==  !=  <  <=  >  >=
&&  ||
```

The operand types must be identical unless the operation involves shifts.

The shift operators shift the left operand by the shift count specified by the right operand, which must be non-negative. If the shift count is negative at run time, a run-time panic occurs. The shift operators implement arithmetic shifts if the left operand is a signed integer and logical shifts if it is an unsigned integer. There is no upper limit on the shift count.

Comparison operators compare two operands and yield a result of the u8 type.

For two integer values x and y, the integer quotient q = x / y and remainder r = x % y satisfy the following relationships:

```
x = q*y + r  and  |r| < |y|
```

with x / y truncated towards zero.

The one exception to this rule is that if the dividend x is the most negative value for the interger type of x, the quotient q = x / -1 is equal to x (and r = 0) due to two's-complement integer overflow.
If the divisor is zero at run time, a run-time panic occurs. 