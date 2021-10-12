# goproxy 更新日志

## 2021.10.12

内网穿透：支持多client模式，transfer支持穿透某个client的port，只需要在新建transfer时指定machine参数

### 查看client列表

1）`cmd` 方式
```bash
➜  proxy-test git:(main) ✗ go run cmd/main.go server list
{
    "items": [
        "127.0.0.1:53245"
    ]
}
```

2）`api` 方式
```bash
➜  proxy-test git:(main) ✗ curl --location --request GET 'http://localhost:12351/client'
{"code":200,"msg":"ok","data":{"items":["127.0.0.1:53245"]}}
```

### 添加transfer支持选择machine

**如果没有指定machine，系统将随机选择存在的client，如果系统没有则会挂起**

1）`cmd` 方式
```bash
➜  proxy-test git:(main) ✗ go run cmd/main.go server add -l :2222 -r :22 -m 127.0.0.1:53245
200 ok <nil>
```

2）`api` 方式
```bash
➜  proxy-test git:(main) ✗ curl --location --request POST 'http://localhost:12351/transfer/tcp' \
--header 'Content-Type: application/json' \
--data-raw '{
    "laddr": [
        {
            "ip": "",
            "port": 5222
        }
    ],
    "raddr": {
        "port": 22
    },
    "machine": "127.0.0.1:53245"
}'
{"code":200,"msg":"ok","data":null}
```

### 查看transfer

1）`cmd` 方式
```bash
➜  proxy-test git:(main) ✗ go run cmd/main.go server list -t                               
[
    {
        "id": 0,
        "laddr": {
            "port": 2222
        },
        "network": 1,
        "raddr": {
            "port": 22
        },
        "machine": "127.0.0.1:53245"
    },
    {
        "id": 1,
        "laddr": {
            "port": 5222
        },
        "network": 1,
        "raddr": {
            "port": 22
        },
        "machine": "127.0.0.1:53245"
    }
]
```

2）`api` 方式
```bash
➜  proxy-test git:(main) ✗ curl --location --request GET 'http://127.0.0.1:12351/transfer' | json_pp
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   284  100   284    0     0   138k      0 --:--:-- --:--:-- --:--:--  138k
{
   "msg" : "ok",
   "code" : 200,
   "data" : [
      {
         "id" : 0,
         "raddr" : {
            "port" : 22
         },
         "laddr" : {
            "port" : 2222
         },
         "network" : 1,
         "machine": "127.0.0.1:53245"
      },
      {
         "id": 1,
         "laddr": {
            "port": 5222
         },
         "network": 1,
         "raddr": {
            "port": 22
         },
         "machine": "127.0.0.1:53245"
      }
   ]
}
```