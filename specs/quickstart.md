# Quick Start

## Requirement

1. install git on the machine
2. install go version 1.16.x/1.17.x [install golang](https://go.dev/doc/install)
3. install tinygo  [fvm tinygo release](https://github.com/ipfs-force-community/tinygo/tags)
4. install go-fvm-sdk [gofvm tool](https://github.com/ipfs-force-community/go-fvm-sdk/releases) rename tool to ```go-fvm-sdk-tools```

add above tools to your ```PATH```environment
```azure
export PATH=$PATH:<dir to go-fvm-sdk>:<tinygo dir/bin>:<go dir>/bin
```

## Checkout actor example

```sh
go-fvm-sdk-tools new -- mycounter
```

If everything is ok, you can get the a simple project, and then you can try your own idea in this example. but **After modify contract colde, you need to re-run the generate command**

## Generate code

the generate tool produce three kinds of code
1. marshal/unmarshal code of (actor state)/(method input and output)
2. the entry code for actor
3. actor client code ,you can use this client to interact with filecoin

```sh
cd gen && go run main.go
```

## Compile

this command 
```sh
go-fvm-sdk-tools build  #execute this in project root path
```

## Test

```sh
{
  "accounts": [ #pre-made accounts, used for the send test message
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
      "name": "counter",                  #test name
      "binary": "../gofvm-counter.wasm",  #test binary path that generate by compile step
      "constructor": "",                  #contructor parameters
      "cases": [
        {                                 #execute specify method that defined in actor
          "name": "increase",             #test name
          "method_num": 2,                #which method to run in actor
          "params":"1832"                 #actor method parameter
          "send_from":0,                  #caller of this test message
          "expect_code":0,                #expect code if fail
          "expect_message":"",            #expect message if fail
          "return_data":"",               #check return_data if any
        },
      ]
    }
  ]
}
```

execute test command
```sh
go-fvm-sdk-tools test -- <directory for test file>
```
