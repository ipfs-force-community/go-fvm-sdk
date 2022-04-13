# Introduction
目前filecoin在积极推进基于wasm的虚拟机项目，其中主要的的一个重要原因是 wasm是个中间语言，能够把其他语言翻译成wasm执行代码，这样可以快速应用不同平台上现有成果。

go-fvm-sdk整体上会实现成三个部分，

1. 利用现有方案或新实现一个能够把tinygo语言翻译成wasm的工具。这里采用tinygo的原因是目前条件下tinygo支持wasi模式，而go语言的转换功能是基于浏览器运行的假设下完成的，因此无法翻译出干净的wasm代码。

2. 在fvm系统调用的协议上，实现一个go语言版本的SDK，其中主要包含的内容是fvm所有系统调用的go版本接口实现。

3. 根据实际的支持情况，提供一些常用的go语言类库，其中可能是从go系统库中重新实现出来的兼容fvm的版本，也可能是从filecoin项目中整合出来的一个常用功能类库

go-fvm-sdk的理想目标就是能够比较简单的把go上的一个项目迁移到fvm上面。同时能够在fvm上运行，并能够正常的和系统actor和其他actor（包括go/rust生成的）正常交互。

# Motivations

1. go语言语法简单，易于学习，使用。具备更高的开发效率

2. go语言基础库完善，具有大部分区块链领域常见的类库和算法的实现

3. filecoin中使用的大部分算法在在golang中都有对应的实现

4. go/tinygo在当前的实现中已经支持了go语言转换wasm的能力

# Goals

1. go语言能否生成干净的wasm代码，并能够在fvm中正常运行

2. 完全实现fvm要求的SDK系统调用部分

3. go语言生成wasm能够在fvm和其他合约（包括原生合约/rust生成的合约/go生成的合约）正常交互

4. 为了提高go-fvm-sdk代码的执行效率，探索一种在gc环境下使用不安全的编程方式来手动管理内存

# Risks

1. go生成的wasm在和系统合约或者其他合约交互的时候，参数能否正确传递，返回值能否正常获取

2. 由于go的系统库可能会涉及到系统调用部分，那么在合约中使用系统库的时候是否有常用系统库无法在fvm上运行。

3. 不同语言（go/rust)翻译出来的wasm在处理内存上可能存在不同的方法，这是否会导致不能互相兼容。

4. go生成的代码在fvm上运行时，内存方面是否存在什么限制，由于gc的存在，执行效率和内存是否存在严重问题。比如不合理的内存增长，过多运算指令。

5. 如果go-fvm-sdk使用gc模式时候，那么gc在扩充托管内存空间的时候，是否会把rust-wasm的内存空间损坏。

# Refrence

[tinygo转换wasm的例子](https://tinygo.org/docs/guides/webassembly/)

[tinygo编译命令](https://tinygo.org/docs/reference/usage/important-options/)

[tinygo conservative博客讲解](https://aykevl.nl/2020/09/gc-tinygo)

[tinygo conservative具体实现](https://github.com/tinygo-org/tinygo/blob/master/src/runtime/gc_conservative.go)

[wasm标准指令](https://webassembly.github.io/spec/core/syntax/instructions.html)
