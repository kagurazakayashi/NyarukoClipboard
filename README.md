![icon](ico/icon.ico)

# NyarukoClipboard

剪贴板点对点同步

## 使用方法

`NyarukoClipboard [参数]`

- `-h` 显示帮助
- `-s`: 服务器模式，填写[服务监听的地址和端口](#连接地址的写法)
- `-c`: 客户端模式，填写[目标服务器地址和端口](#连接地址的写法)
- `-nr`: 禁止接收(仅连接)
- `-ns`: 禁止发送(仅连接)
- `-r`: 剪贴板检查间隔(毫秒，默认 `1000`(1秒) )

`-s` 和 `-c` 必须有其中一项，且只能有一项，填哪个代表使用哪个模式(服务器/客户端)。如果省略，默认以**服务器**模式启动监听 `tcp://:7976` 。

### 连接地址的写法

`[tcp/udp]://[host]:[port]`

- `[tcp/udp]`: 建议使用 tcp
- `[host]`: 作为服务器时的监听地址 / 作为客户端时的连接地址
- `[port]`: 端口号
- 示例: `tcp://:7976`, `tcp://192.168.1.2:7976`

默认为服务器模式，监听不限 IP 的 7976 端口。

## 编译

- `go get`
- `build.bat` (Windows 环境，其他环境命令类似)

## LICENSE

Copyright (c) 2024 KagurazakaYashi NyarukoClipboard is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
