# Dict Passing 在 Haskell/Rust 中的实现：Ad-hoc 多态的技术细节

## 动机与背景

Dict passing（字典传递）是 Haskell、Rust 等语言实现类型类（type classes）或特性（traits）这类 ad-hoc 多态的核心技术。其核心动机是：

1. **解决参数化多态的局限性**：
    - 传统参数化多态（如 Java 泛型）无法表达"这个类型必须支持某些操作"的约束
    - 需要一种方式在编译时保证类型满足特定接口，同时保持运行时灵活性

2. **实现零成本抽象**：
    - 不需要虚函数表的开销
    - 保持静态分派的性能优势

3. **支持类型类推导**：
    - 编译器自动推断并传递所需的类型类实例
    - 减少用户显式传递样板代码

4. **模块化设计**：
    - 实现与接口分离
    - 允许事后扩展类型的功能

## 技术实现细节

### 1. 基本原理

Dict passing 的核心思想是将类型类的约束转换为隐式的"字典"参数传递：

```haskell
-- Haskell 代码
showList :: Show a => [a] -> String

-- 会被脱糖(desugar)为：
showList :: ShowDict a -> [a] -> String
```

### 2. 编译器内部表示

#### 2.1 字典结构

对于每个类型类，编译器生成一个字典结构：

```rust
// 对应 Haskell 的 `Show a` 类
struct ShowDict<A> {
show: fn(&A) -> String,
showsPrec: fn(i32, &A) -> String,
// ... 其他方法
}
```

#### 2.2 实例解析

当遇到类型类实例时：

```haskell
instance Show Bool where
show True = "True"
show False = "False"
```

编译器生成：

```rust
const boolShowDict: ShowDict<bool> = ShowDict {
show: |x| match x {
true => "True".to_string(),
false => "False".to_string(),
},
// ... 其他方法
};
```

### 3. 编译过程

#### 3.1 类型检查阶段

1. 收集所有类型类约束
2. 构建约束求解环境
3. 解析约束到具体字典

#### 3.2 脱糖转换

将类型类约束转换为显式字典参数：

```haskell
-- 原代码
printTwice :: Show a => a -> IO ()
printTwice x = do
print x
print x

-- 脱糖后
printTwice :: ShowDict a -> a -> IO ()
printTwice dict x = do
print (dict.show x)
print (dict.show x)
```

#### 3.3 字典传递优化

1. **全局实例**：对单例类型类（如 `Show Int`），直接使用全局静态字典
2. **参数化传递**：多态函数携带额外的字典参数
3. **字典融合**：合并多个相关字典减少传递开销

### 4. Rust 的特例化实现

Rust 的 trait 系统使用类似的机制，但有一些特殊处理：

```rust
// Rust trait
trait Show {
fn show(&self) -> String;
}

// 编译后实际上生成
struct ShowVtable<T> {
show: fn(&T) -> String,
// 其他方法...
}

// 泛型函数
fn print_twice<T: Show>(x: T) {
println!("{}", x.show());
println!("{}", x.show());
}

// 脱糖为
fn print_twice<T>(x: T, vtable: &ShowVtable<T>) {
(vtable.show)(&x);
(vtable.show)(&x);
}
```

### 5. 高级优化技术

1. **字典消解**：对于单态化（monomorphized）的泛型函数，直接内联字典内容
2. **静态方法解析**：当具体类型已知时，直接静态分派方法调用
3. **层次化字典**：对继承的类型类（如 `Eq a => Ord a`），构建字典的层次结构
4. **惰性字典生成**：只在需要时才生成字典实例

### 6. 与虚函数表的区别

Dict passing 与传统的 OO 虚函数表关键区别：

1. **传递方式**：
    - 虚函数表：通过对象指针间接引用
    - 字典传递：通过额外参数显式传递

2. **组合性**：
    - 虚函数表：固定单一表结构
    - 字典传递：可以灵活组合多个字典

3. **多参数支持**：
    - 虚函数表：只支持单一接收者
    - 字典传递：可以处理多参数类型类（如 `Add`）

## 性能特点

1. **优点**：
    - 静态分派带来的优化空间
    - 无运行时类型检查开销
    - 更好的内联可能性

2. **缺点**：
    - 代码膨胀（多态函数的多个实例化版本）
    - 字典参数的传递开销（通常被优化掉）

## 总结

Dict passing 是一种将 ad-hoc 多态转换为显式参数传递的编译技术，它保持了语言的表达力同时不牺牲性能。这种技术在 Haskell 和 Rust 中的成功应用证明了其在系统编程和函数式编程领域的价值。现代编译器通过多种优化手段，使得这种抽象几乎达到了零成本。