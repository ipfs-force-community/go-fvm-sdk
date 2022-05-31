# gc_conservative analysis

## 内存布局

## 内存分配

## 内存回收

### markglobal

### markStack
p:代表上一个帧数据的指针
s:代表帧上保存的slot数量，这表示在这个栈上面分配了多少个内存指针

本gc中，每个栈前面都会插入一些描述栈帧信息的数据， 所有的栈帧指针数据之间是一个链表。markStack的过程是从最后的帧位置开始，通过这个链表，检查每一个栈帧。

每个帧的结构是 ptr|slotnum|ptr1,ptr2 ...ptrn|

p1 s1 <ptr>|p2 s2 <ptr>|p3 s3 <ptr>|....|pn sn <ptr>|
其中pn -> p(n-1) -> p(n-2) -> ····· -> p3 -> p2 -> p1

栈帧数据结构， parent指向上一个栈帧的信息，numSlots表示

```go
type stackChainObject struct {
	parent   *stackChainObject
	numSlots uintptr
	... //alloc ptr
}
```

这里stackChainObject结构体重只有前两个字段，后面的分配的指针是通过编译器在所有发生内存分配的函数中插入的一些指令处理的。
```go
var funcMap3   []int32

//go:export shim_map_size
func SetMapSize(oldFuncNumber int32) {
	funcMap = make([]int32, oldFuncNumber)
}

```

转换出来的wasm代码如下
```wat
 (func $shim_map_size (type 4) (param i32)
    (local i32 i32 i32)   ;; p l1 l2 l3
    global.get $__stack_pointer
    i32.const 16
    i32.sub
    local.tee 1
    global.set $__stack_pointer     ;;l1= sp-16 sp=sp-16    创建stackobject

    local.get 1
    i64.const 1
    i64.store offset=4 align=4      ;;addr(sp(4-12)) = 1

    i32.const 0
    i32.load offset=65556
    local.set 2                     ;;  l2 = i32.load(65556)  中间变量保存老的stackobject

    i32.const 0
    local.get 1
    i32.store offset=65556          ;; i32.store(65556, l1)  全局最新 stackobject指针指向新的stackobject

    local.get 1
    local.get 2
    i32.store                       ;; i32.store(l1, l2)    新stackobject的parent字段指向老的stackobject

    block  ;; label = @1
      local.get 0
      i32.const 1073741824
      i32.lt_u
      br_if 0 (;@1;)
      unreachable
      unreachable
    end

    local.get 0
    i32.const 2
    i32.shl                         ;;size = 2^method_number
    call $runtime.alloc
    local.set 3                     ;;l3 = alloc(size)
    
    i32.const 0
    local.get 2
    i32.store offset=65556          ;;i32(65556) = l2  退出的时候

    local.get 1
    i32.const 8
    i32.add
    local.get 3
    i32.store                       ;; addr(l1+8) = l3  设置stackobject中的指针字段

    i32.const 0
    local.get 0
    i32.store offset=65568          ;;设置slice长度

    i32.const 0
    local.get 0
    i32.store offset=65564          ;;设置slice容量

    i32.const 0
    local.get 3
    i32.store offset=65560          ;;设置slice指针

    local.get 1
    i32.const 16
    i32.add
    global.set $__stack_pointer)    ;;退出栈
```

markStack函数会比遍历整个链表，案后通过ptr和slot在遍历每个帧上保存的数据指针，所有在栈栈上的指针指向的内存部分会被标识成在使用。


