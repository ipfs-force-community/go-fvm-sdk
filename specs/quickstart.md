# Quick Start

## Requirement

1. Install [git](https://github.com/git-guides/install-git) 
2. Install [Go](https://go.dev/doc/install) version 1.16.x/1.17.x
3. Install [TinyGo](https://tinygo.org/getting-started/install/)
4. Install [go-fvm-sdk](https://github.com/ipfs-force-community/go-fvm-sdk/releases); Then rename it to `go-fvm-sdk-tools`

Add go-fvm-sdk-tools tools to your ```PATH``` environment.
```bash
export PATH=$PATH:<dir to go-fvm-sdk>
```

## Patch your local environment

```bash
go-fvm-sdk-tools patch
```
this command change your local go/tinygo std package,  this may cause other code not work properly. more details refer link [patch](https://github.com/ipfs-force-community/go_tinygo_patch)

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
grep -v '#' > test.json << EOF
{
  # mock accounts, used for the send test message
  "accounts": [ 
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
      # test name
      "name": "counter",       
      # Path to the binary that is generated during compile step           
      "binary": "../gofvm-counter.wasm",  
      # contructor parameters
      "constructor": "",                 
      "cases": [
       # execute specify method that defined in actor
        {           
          # test name
          "name": "increase",    
          # which actor method to run         
          "method_num": 2, 
          # actor method parameter              
          "params": "1832"   
          # caller of this test message             
          "send_from": 0,           
          # expect code if fail      
          "expect_code": 0,       
          # expect message if fail        
          "expect_message": "",      
          # check return_data if any     
          "return_data": "",              
        },
      ]
    }
  ]
}
EOF
```

Run the test.
```sh
go-fvm-sdk-tools test
```
