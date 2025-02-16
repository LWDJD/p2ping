package main

import (
	"context"
	"errors"
	"os"
    "fmt"
	"strings"
	// "strconv"
	"time"
	// "encoding/hex"
	// cry "github.com/libp2p/go-libp2p/core/crypto"
    "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	b58 "github.com/mr-tron/base58/base58"
	"github.com/libp2p/go-libp2p/core/host"
	// "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	multiaddr "github.com/multiformats/go-multiaddr"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

var versions string = "1.1.0"
func main(){
	// fmt.Println("开始运行\n")

	//设置节点私钥
	// pk,err := b58.Decode("23jhTfQKFxeuxYVYhMZKrLeMX5bFLfavf2TsUeR1Qg4J3rieF4mrN9tYWUoGAXxX8kBKpLX9JQ5csCY7TijCLKVakrxzt")
	// if err != nil {
    //     panic(err)
    // }
	// privKey,err := cry.UnmarshalPrivateKey(pk)
	// if err != nil {
    //     panic(err)
    // }




	// 启动 libp2p 节点
    node, err := libp2p.New()
    if err != nil {
        panic(err)
    }

	switch len(os.Args) {
		case 1:
			fmt.Println("P2Ping Error：无参")
			help()
		case 2:
			switch strings.ToLower(os.Args[1]){
				case "-v":
					fmt.Println("\np2ping版本：v"+versions+"\n")
				case "help":
					help()
				case "explain":
					explain()
				default:
					err := Ping(node,os.Args[1],4)
					if err != nil{
						fmt.Println("P2Ping Error：",err,"\n")
						fmt.Println(`输入 p2ping help 获得帮助`)
					}
			}
		default:
			fmt.Println("P2Ping Error：未知参")
			help()
	}

	// if len(os.Args) > 1 {
	// 	if len(os.Args) > 2{
	// 		fmt.Println("P2Ping Error：未知参")
	// 	}else{
	// 		fmt.Println("\n正在 Ping "+os.Args[1]+" 具有 32 字节的数据:")
	// 		err := Ping(node,os.Args[1],4)
	// 		if err != nil{
	// 			fmt.Println("P2Ping Error：",err,"\n")
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("P2Ping Error：无参")
	// }

    // 打印节点的侦听地址
    // fmt.Println("监听地址:", node.Addrs(), "\n")
	// fmt.Println("节点ID:" , node.ID())
	// nodeID,err := cry.MarshalPrivateKey(node.Peerstore().PrivKey(node.ID()))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("节点私钥:" ,b58.Encode(nodeID))

	Stop(node)
}

// 关闭节点
func Stop(node host.Host){

    if err := node.Close(); err != nil {
        panic(err)
    }
	// fmt.Println("程序退出！！")
	os.Exit(0)
}
/* Ping功能
 * node 节点
 * addrStr 地址
*/
func Ping(node host.Host,addrStr string,quantity int) error{
	var addrInfo peerstore.AddrInfo 
	compare := "/p2p/"
	if peerID,err := IDFromString(addrStr); err==nil||addrStr[0:5]==compare[0:5]{

		if addrStr[0:5]==compare[0:5] {
			// fmt.Println("记录点S")
			addr, err := multiaddr.NewMultiaddr(addrStr)
			if err != nil {
				return err
			}
			// fmt.Println("记录点L")
			addrInfoP, err := peerstore.AddrInfoFromP2pAddr(addr)
			if err != nil {
				return err
			}
			addrInfo = *addrInfoP
			peerID = addrInfo.ID
		}
		
		dhtCtx := context.Background()

		// 创建 kad-dht 实例
		dht, err := dht.New(dhtCtx, node , dht.BootstrapPeers(dht.GetDefaultBootstrapPeerAddrInfos()...))
		if err != nil {
			return(err)
		}

		// 启动 DHT
		err = dht.Bootstrap(dhtCtx)
		if err != nil {
			return(err)
		}
		fmt.Println("DHT正在预热...")
		time.Sleep(8 * time.Second)
		connectedPeers := node.Network().Peers()
		fmt.Printf("DHT预热结束，已连接节点数: %d\n", len(connectedPeers))
		fmt.Println("正在查找DHT中的节点连接地址...")
		addrInfo,err = dht.FindPeer(dhtCtx, peerID)

		// fmt.Println(addrInfo)
		if len(addrInfo.Addrs)==0{
			return errors.New("未查询到"+peerID.String()+"的连接地址")
		}
		
	}
	

	if addrInfo.ID=="" {
		// fmt.Println("记录点A")
		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		// fmt.Println("记录点B")
		addrInfoP, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		addrInfo = *addrInfoP
	}

	
	
	// fmt.Println("记录点C")
	if err := node.Connect(context.Background(), addrInfo); err != nil {
		return err
	}
	// fmt.Println("记录点D")
	ch := ping.Ping(context.Background(),node,addrInfo.ID)
	// fmt.Println("记录点E")
	var times []time.Duration
	fmt.Println("\n正在 Ping "+os.Args[1]+" 具有 32 字节的数据:")
	for i := 0; i < quantity; i++ {
		time.Sleep(1 * time.Second)
		res := <-ch
		if res.RTT.Nanoseconds()==0{
			fmt.Println("请求失败。")
		}else {
			fmt.Println("来自", addrInfo.ID, "的回复: 时间", res.RTT)
		}
		
		times = append(times,res.RTT)
	}

	var stat int = quantity
	var s int
	var l int
	var q int
	var average time.Duration
	for i :=0; i < len(times); i++{
		if times[i].Nanoseconds()==0{
			stat--
		}else {
			average+=times[i]
			q++
			if times[s].Nanoseconds()>times[i].Nanoseconds(){
				s = i	
			}
		}
		if times[l].Nanoseconds()<times[i].Nanoseconds(){
			l = i
		}
	}
	fmt.Println()
	fmt.Println(addrInfo.ID,"的 Ping 统计信息:\n    数据包: 已发送 =",quantity,"，已接收 =",stat,"，丢失 =",quantity-stat," (",((float64(quantity-stat) / float64(quantity)) * 100),"% 丢失)，")
	if average.Nanoseconds() != 0{
		average = average / time.Duration(q)
		fmt.Println("往返行程的估计时间:\n    最短 = ",times[s],"，最长 = ",times[l],"，平均 = ",average,"\n")
	}else{
		fmt.Println()
	}
	return nil
}

func help(){
	fmt.Println(`
p2ping 帮助：
p2ping {multiaddr} : 对这个多地址进行ping操作
p2ping -v          : 显示版本号
p2ping explain     : 显示说明
p2ping help        : 显示帮助
`)
}

func explain(){
	fmt.Println(`
解释：
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

2.节点ID 是由base58编码的地址，与比特币地址类似;
    节点ID : 由公钥生成，用于身份验证 与 kad DHT存储标识；
    公钥   : 由私钥生成，用于身份验证；
    生成过程单向，无法反推。
`)
}
// func Status(node host.Host) error{
// 	if node ==nil{
// 		return errors.New("Status(node host.Host)节点不能为空")
// 	}
// 	// 获取所有已连接节点的 PeerID
// 	connectedPeers := node.Network().Peers()
// 	fmt.Printf("已连接节点数: %d\n", len(connectedPeers))
// 	return nil
// }

/* 转换ID为ID
 * 
 * IDString base58编码的peerID
*/
func IDFromString(IDString string) (peerstore.ID, error) {
	// 将 Base58 编码的字符串解码为字节切片
	decoded, err := b58.Decode(IDString)
	if err != nil {
		return "", err
	}

	// 从字节切片创建 peer.ID
	id, err := peerstore.IDFromBytes(decoded)
	if err != nil {
		return "", err
	}

	return id, nil
}

