# p2ping

简单的 libp2p Ping 程序。

## 使用方法
### 环境
1. 环境变量 `path` 中添加程序路径。
2. CMD 中输入以下命令

### 命令
```bash
p2ping {multiaddr} # 对这个多地址进行 ping 操作
p2ping {peer ID}   # 对这个节点ID进行DHT查询 并ping操作
p2ping -v          # 显示版本号
p2ping explain     # 显示说明
p2ping help        # 显示帮助
```
### 解释
```
1.multiaddr是libp2p节点所使用的一种地址表示方式，可以表达多种协议：
    举例：
    /ip4/1.2.3.4/tcp/1234/p2p/12D3KooWKS71s4iCRVHmdCp1Mg6dJTckiZdhRf77J7dgwJsybvri
        表示 IPv4 地址 1.2.3.4，端口 1234，使用 TCP 协议；
        /p2p/12D3KooWKS71s4iCRVHmdCp1Mg6dJTckiZdhRf77J7dgwJsybvri
        表示 点对点 连接使用此ID的节点；
    /ip6/[::1]/udp/5678/quic-v1/p2p/12D3KooWKS71s4iCRVHmdCp1Mg6dJTckiZdhRf77J7dgwJsybvri
        表示 IPv6 地址 [::1]，端口 5678，使用 UDP 协议上的QUIC-V1协议；
        /p2p/12D3KooWKS71s4iCRVHmdCp1Mg6dJTckiZdhRf77J7dgwJsybvri
        表示 点对点 连接使用此ID的节点；

2.节点ID 是由base58编码的地址，与以太坊地址类似;
    节点ID : 由公钥生成，用于身份验证 与 kad DHT存储标识；
    公钥   : 由私钥生成，用于身份验证；
    生成过程单向，无法反推。
```
