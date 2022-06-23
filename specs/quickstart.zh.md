# Quick Start

## 要求

1. 安装git
2. 安装go  版本要求1.16.x/1.17.x [install golang](https://go.dev/doc/install)
3. 安装tinygo  [fvm tinygo release](https://github.com/ipfs-force-community/tinygo/tags)
4. 安装go-fvm-sdk-tools [gofvm tool](https://github.com/ipfs-force-community/go-fvm-sdk/releases) 下载文件重命名为```go-fvm-sdk-tools```

需要把上面三个工具加入到PATH环境变量里面
```azure
export PATH=$PATH:<dir to go-fvm-sdk>:<tinygo dir/bin>:<go dir>/bin
```

## 创建合约项目

```sh
go-fvm-sdk-tools new -- mycounter
```

如果一切正常，你可以看到一个简单的go actor的项目，可以在里面添加自己的内容。**增加自己的内容后，需要重新运行生成命令**

## 生成代码

目前生成工具主要生成三种内容
1. 合约结构体及输入输出结构体所需要的序列化反序列化文件
2. 合约入口文件
3. 合约客户端代码，可以使用该客户端代码和合约进行交互

```sh
cd gen && go run main.go
```

## 合约编译

```sh
go-fvm-sdk-tools build  #在项目根目录运行
```
## 合约测试

```sh
{
  "accounts": [ #预制账号，用于做测试消息的from
    {
      "priv_key": "6c3b9aa767f785b537c0d8ba5fa54677e6a6e281320dfbb27c889b8fa460670f",
      "address": "f1m674sjwmga36qi3wkowt3wozwpahrkdlvd4tpci",
      "balance": 10000
    },
    {
      "priv_key": "b10da48cea4c09676b8e0efcd806941465060736032bb898420d0863dca72538",
      "address": "f1dwyrbh74hr5nwqv2gjedjyvgphxxkffxug4rkkq",
      "balance": 10000
    },
    {
      "priv_key": "c51b8a31c98b9fe13065b485c9f8658c194c430843570ccac2720a3b30b47adb",
      "address": "f15o3zaqettjmmarblwzjr66lwddsi6rbtxjwzngq",
      "balance": 10000
    }
  ],
  "contracts":[
    {
      "name": "counter",                  #测试名
      "binary": "../gofvm-counter.wasm",  #运行的合约路径
      "constructor": "",                  #构造函数参数
      "cases": [
        {                                 #执行的合约方法测试
          "name": "increase",             #测试名称
          "method_num": 2,                #测试方法
          "params":"1832"                 #执行的方法参数
          "send_from":0,                  #消息的发起人
          "expect_code":0,                #期望的消息退出码
          "expect_message":"",            #期望的错误返回
          "return_data":"",               #期望的消息返回数据
        },
      ]
    }
  ]
}
```

运行合约测试命令
```sh
go-fvm-sdk-tools test -- <存放test文件的目录>
```
