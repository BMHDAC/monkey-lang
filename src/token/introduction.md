### The monkey language has these features:
  - C-like syntax
  - variable binding
  - intergers and boolean
  - athrimetic expression
  - built-in functions
  - first-class and higher-order functions
  - closures
  - a string data structure
  - an array data structure
  - a hash data structrure

Here is an example program in the monkey lang
```cpp
  let age = 1;
  let name = "Monkey";
  let result = 10 * (20 / 2);
  let array = [1,2,3,4,5,6];
  let my_love = { name: "Chan", age: 22};
  array[1] = 2;
  my_love.name = "Chan";
  //This is how you do a built-in function
  let my_function = fn(a,b) { return a + b; };
  //Or this
  let my_function = fn(a,b) { a + b;};
  my_function(1,2) // => 3

  // Control flow in go
  let fibonaci = fn(x) {
    if (x == 0) {
      0
    } else {
      if (x == 1) {
        1
      } else {
        fibonaci(x - 1) + fibonaci(x - 2);
      }
    }
  };

  // Higher order function in Monkey
  let twice = fn(f,x) {
    return f(f(x))
  }
  let add_two = fn(x) {
    return x + 2;
  }
```
