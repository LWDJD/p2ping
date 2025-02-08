package main

import (
	"context"
	// "errors"
	"os"
    "fmt"
	// "strconv"
	"time"
	// "encoding/hex"
	// cry "github.com/libp2p/go-libp2p/core/crypto"
    "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	// b58 "github.com/mr-tron/base58/base58"
	"github.com/libp2p/go-libp2p/core/host"
	// "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	multiaddr "github.com/multiformats/go-multiaddr"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	// dht "github.com/libp2p/go-libp2p-kad-dht"
)

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

	if len(os.Args) > 1 {
		if len(os.Args) > 2{
			fmt.Print("P2Ping Error：未知参")
		}else{
			fmt.Println("\n正在 Ping "+os.Args[1]+" 具有 32 字节的数据:")
			err := Ping(node,os.Args[1],4)
			if err != nil{
				fmt.Println("P2Ping Error：",err,"\n")
			}
		}
	} else {
		fmt.Print("P2Ping Error：无参")
	}

    // 打印节点的侦听地址
    // fmt.Println("监听地址:", node.Addrs(), "\n")
	// fmt.Println("节点ID:" , node.ID())
	// nodeID,err := cry.MarshalPrivateKey(node.Peerstore().PrivKey(node.ID()))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("节点私钥:" ,b58.Encode(nodeID))

	

}

// 关闭节点
func Stop(node host.Host){

    if err := node.Close(); err != nil {
        panic(err)
    }
	fmt.Println("程序退出！！")
	os.Exit(0)
}
/* Ping功能
 * node 节点
 * addrStr 地址
*/
func Ping(node host.Host,addrStr string,quantity int) error{
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		return err
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}
	if err := node.Connect(context.Background(), *peer); err != nil {
		return err
	}
	ch := ping.Ping(context.Background(),node,peer.ID)
	var times []time.Duration
	for i := 0; i < quantity; i++ {
		res := <-ch
		if res.RTT.Nanoseconds()==0{
			fmt.Println("请求失败。")
		}else {
			fmt.Println("来自", peer.ID, "的回复: 时间", res.RTT)
		}
		
		times = append(times,res.RTT)
		time.Sleep(1 * time.Second)
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
	fmt.Println(peer.ID,"的 Ping 统计信息:\n    数据包: 已发送 =",quantity,"，已接收 =",stat,"，丢失 =",quantity-stat," (",((float64(quantity-stat) / float64(quantity)) * 100),"% 丢失)，")
	average = average / time.Duration(q)
	fmt.Println("往返行程的估计时间:\n    最短 = ",times[s],"，最长 = ",times[l],"，平均 = ",average)
	return nil
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