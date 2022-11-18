# fvm_go_sdk 命令介绍

1. 版本命令
```bash
$ fvm_go_sdk version
```

输出:
```txt
fvm_go_sdk version v0.1.0+git.6ccc890
```

2. 创建模版命令

执行此命令要求先安装[git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)。
命令执行后，会从github中检索出项目[模版](https://github.com/ipfs-force-community/gofvm-counter), 然后自动运行生成文件，项目编译， 运行预先编辑好的测试。

```bash
fvm_go_sdk new -- hello
```

输出
```txt
$ git clone https://github.com/ipfs-force-community/gofvm-counter.git
正克隆到 'gofvm-counter'...
remote: Enumerating objects: 85, done.
remote: Counting objects: 100% (85/85), done.
remote: Compressing objects: 100% (53/53), done.
remote: Total 85 (delta 36), reused 71 (delta 26), pack-reused 0
接收对象中: 100% (85/85), 175.11 KiB | 914.00 KiB/s, 完成.
处理 delta 中: 100% (36/36), 完成.
$ rm -rf .git
$ mv gofvm-counter hello
module....
$ mkdir client
$ go mod tidy
$ go run main.go
$ go mod tidy
$ fvm_go_sdk build -o hello.wasm
$ fvm_go_sdk test
passed: case counter_install_code
code cid bafk2bzaceb7pkldc5yyqf6onctrd54jlxcvod76yas3ayf64l3lnwsegnjqa2
passed: case counter_create_actor
actor cid f0103
passed: case counter_increase
passed: case counter_get
passed: case counter_increase again
passed: case counter_get_result
```

3. 补丁命令

由于直接使用tinygo生成代码会有编译运行问题，因此sdk需要对本地的tinygo和go代码
进行一些修改，详情见[fvm go_tinygo_patch](https://github.com/ipfs-force-community/go_tinygo_patch)

由于安装方式的不同，可能需要使用sudo运行，需根据实际调整。

```bash
fvm_go_sdk patch
```

输出：
```txt
patching file src/reflect/value.go
Hunk #1 succeeded at 754 (offset 3 lines).
patching file targets/wasi.json
Hunk #2 succeeded at 10 with fuzz 1.
```

4. 编译命令

编译命令使用tinygo生成wasm文件， 然后工具会对这个wasm文件进行编辑，其中包括下面两种处理： 
1. 由于tinygo的init初始化语句是在Start中运行，但是fvm规定的入口是invoke函数，因此需要在invoke执行前把tinygo中init初始化语句运行一遍，这里简单的在invoke函数体中调用了一下Start函数。
2. 在debug模式下，会自动把fd_write的内容输出到fvm的debug log中，如果遇到unreachable这种错误，可以开启debug编译之后在运行，会打印出发生错误的信息。

通常情况下只需要在项目根目录下简单的运行

```bash
fvm_go_sdk build
```

参数说明
* -d：开启debug模式，会重定向程序的输出流到fvm
* -o：指定输出wasm文件的位置
* -w：生成wasm文件对应的wat文件，常用于开发人员查看   

5. 测试命令

测试命令用于测试开发人员编写的合约，默认会读取运行目录下的test.json文件中的用例。并启动一个本地fvm运行合约。

```bash
fvm_go_sdk test --name  <可选｜指定测用例中的一个运行> --  <可选｜读取指定位置的测试文件>  
```

输出:
```txt
passed: case helloworld_install_code
code cid bafk2bzacedh7fvhyw5mcu2otvvaquchonnhq75ecdc5fkiiez37cxm3vteohy
passed: case helloworld_create_actor
actor cid f0103
```