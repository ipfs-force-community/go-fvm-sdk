解释系统调用。

fvm具备的系统调用

系统调用函数的一个大致结构 第一个参数是返回值， 后面追加的是参数， 函数返回调用返回码，通常0代表成功，其他代表失败

如何导入fvm系统调用函数

如果在golang中导入fvm系统函数。参数或者返回值有这么几种情况

1. 简单值传递（数值型）
2. 传递字符串  获取stringHeader 传递postion+字符串长度
3. 传递slice  获取sliceHeader 传递postion+切片长度
4. 传递结构体  结构体转换成slice在传递
5. 返回字符串  目前还不存在
6. 返回slice  构建一个slice buf(长度预设或者通过stat接口获取)，获取其指针和位置传递给fvm
7. 返回结构体  直接传递结构体指针（结构体字段顺序和类型必须严格匹配）

库处理结构

引用的库无论是系统库还是外部库都不能引用和系统（常见的有文件，网络，time）相关的内容，如果有引用，就需要改造这个库并移除系统相关的依赖

1. 系统库的处理办法， tinygo里面增加这部分自定义系统库。fmt目前采用这种做法， 删除了Println Printf这类输出流打印
2. 外部库，常见于加密库，和文件流相关的库， fork过来，删除相关代码，在合约项目中replace掉，目前filecoin的一些库采用这种做法

cbor生成:

目前cbor生成方式，需要待生成包能够完成编译，但是游戏引用了fvm的系统调用，在正常的环境下无法完成编译。 只能复制结构体出来生成

合约的一个基本结构：

State: 保存合约的状态信息，记录合约的数据
Invoke: 导出合约调用入口，Invoke通常写法就是一个switch case， 通过传入的方法序号决定执行哪一段逻辑。其中方法1是用于初始化合约的，必须定义出来。
参数和返回值： 合约参数及返回值都通过sdk.GetBlock sdk.Put来获取和传递。，Invoke的参数是个uint32，通过这个数值调用sdk.GetBlock. Invoke的返回值是个Uint32, 把要返回的数据通过sdk.Put放到fvm里面得到一个序号，外边会通过这个序号到blockregistry里面获取返回的数据


合约部署

lotus chain install-actor <wasm-path>

```go
	params, err := actors.SerializeParams(&init8.InstallParams{
			Code: code,
		})
	msg := &types.Message{
			To:     builtin.InitActorAddr,
			From:   fromAddr,
			Value:  big.Zero(),
			Method: 3,
			Params: params,
		}

```

lotus chain create-actor

```go
    params, err := actors.SerializeParams(&init8.ExecParams{
			CodeCID:           codeCid,
			ConstructorParams: cparams,
	})
	msg := &types.Message{
			To:     builtin.InitActorAddr,
			From:   fromAddr,
			Value:  big.Zero(),
			Method: 2,
			Params: params,
		}
```

louts chain invoke 

```go
    msg := &types.Message{
        To:     addr,
        From:   fromAddr,
        Value:  big.Zero(),
        Method: abi.MethodNum(method),
        Params: params,
    }
```
