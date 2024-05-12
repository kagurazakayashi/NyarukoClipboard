![icon](ico/icon.ico)

# NyarukoClipboard

使剪贴板中的文本和图片在两台电脑间点对点同步，不经过第三方服务器，支持使用证书进行验证和加密。

## 使用方法

`NyarukoClipboard [参数]`

- `-h` 显示帮助
- `-s`: 服务器模式，填写[服务监听的地址和端口](#连接地址的写法)
- `-c`: 客户端模式，填写[目标服务器地址和端口](#连接地址的写法)
- `-e`: X.509 证书文件，支持 CN 限定监听域名，确保连接方身份和数据安全。双方应使用同样的证书。
- `-k`: X.509 证书私钥(仅作为服务端时使用)
- `-nr`: 禁止接收(单向传输用)
- `-ns`: 禁止发送(单向传输用)
- `-r`: 剪贴板检查间隔(毫秒，默认 `1000`(1秒) )：每隔多久检查一下剪贴板中的内容是否有变化，如果有变化则收集并发送。
- `-cd`: 剪贴板读取冷却时间(毫秒，默认 `2000`(2秒) )：在发送一次剪贴板内容后，多长时间不再次发送。防止剪贴板短时间内大量变化导致过于频繁的网络传输和目标剪贴板内容替换。
- `-v`: 显示调试信息，显示剪贴板监听和网络操作的详细过程。

### 注意事项

- `-s` 和 `-c` 必须只能有其中一项，填哪个代表使用哪个模式(服务器/客户端)。如果省略，默认以**服务器**模式启动监听 `tcp://:7976` 。
- 图片复制有尺寸和大小上限，通常复制过大的内容到剪贴板时不进行同步。
- 最大单次传输内容为 `10,000,000 bytes` (10 MB)

### 连接地址的写法

`[tcp/udp]://[host]:[port]`

- `[tcp/udp]`: 建议使用 tcp
- `[host]`: 作为服务器时的监听地址 / 作为客户端时的连接地址
- `[port]`: 端口号
- 示例: `tcp://:7976`, `tcp://192.168.1.2:7976`

默认为服务器模式，监听不限 IP 的 7976 端口。

### 证书创建

1. 创建一个 `.cnf` 文件，内容参考 `testcert/san.cnf` 。
2. 在 `[alt_names]` 设置服务端监听的域名或 IP 地址。域名以 `DNS.1, DNS.2, DNS...` 开头，IP 地址以 `IP.1, IP.2, IP...` ，可以指定多个。
3. 使用命令 `openssl req -new -x509 -days 365 -nodes -config san.cnf -keyout key.pem -out cert.pem` 完成创建。

- `365` 为证书有效期（天）。
- `san.cnf` 为刚才创建的 `cnf` 文件。
- `key.pem` 是要保存的新私钥文件。
- `cert.pem` 是要保存的新证书文件。

## 编译

- `go get`
- `build.bat` (Windows 环境，其他环境命令类似)

根据需要，可以通过调整 `main.go` 中的 `bufSize` 常量来决定单次可以传输的最大数据量。

## LICENSE

Copyright (c) 2024 KagurazakaYashi NyarukoClipboard is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
