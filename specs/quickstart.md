# Quick Start

## Requirement

1. Install [git](https://github.com/git-guides/install-git) 
2. Install [Go](https://go.dev/doc/install) version 1.16.x/1.17.x
3. Install `TinyGo` - fvm [release](https://github.com/ipfs-force-community/tinygo/tags)
4. Install [go-fvm-sdk](https://github.com/ipfs-force-community/go-fvm-sdk/releases); Then rename it to `go-fvm-sdk-tools`

Add above tools to your ```PATH``` environment.
```bash
export PATH=$PATH:<dir to go-fvm-sdk>:<tinygo dir/bin>:<go dir>/bin
```

## Create an actor project

```sh
go-fvm-sdk-tools new -- mycounter
```

If all goes well, a template actor project will be genereated for you. Write new actors at your will in the directory structure created for you. **Note: after modify contract code, you will need to re-run the generate command**

## Generator

Generator takes care of following aspects of the actor code...

1. Marshal and unmarshal code for actor state structs
2. Entry code for the actor
3. Actor client code which you can use to interact with filecoin

```sh
cd gen && go run main.go
```

## Compile

Compile your actor with the following command.
```sh
go-fvm-sdk-tools build  # execute at project root
```

## Test

```sh
{
  "accounts": [ # mock accounts, used for the send test message
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
      "name": "counter",                  # test name
      "binary": "../gofvm-counter.wasm",  # Path to the binary that is generated during compile step
      "constructor": "",                  # contructor parameters
      "cases": [
        {                                 # execute specify method that defined in actor
          "name": "increase",             # test name
          "method_num": 2,                # which actor method to run
          "params": "1832"                # actor method parameter
          "send_from": 0,                 # caller of this test message
          "expect_code": 0,               # expect code if fail
          "expect_message": "",           # expect message if fail
          "return_data": "",              # check return_data if any
        },
      ]
    }
  ]
}
```

Run the test.
```sh
go-fvm-sdk-tools test -- <directory for test file>
```
