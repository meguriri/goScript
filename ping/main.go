package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_PG = 65500 //数据包最大字节数
)

//ICMP首部
type ICMP struct {
	Type       uint8  //类型
	Code       uint8  //代码
	Checksum   uint16 //检验和
	Identifier uint16 //标识符
	SeqenceNum uint16 //序号
}

var (
	originBytes []byte // 数据包
)

func init() {
	originBytes = make([]byte, MAX_PG)
}

//ICMP校验和
func CheckSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length int = len(data)
		index  int = 0
	)

	for length > 1 { //2个字节凑一行二进制加法（模2和），保留溢出所以用32位加法
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	if length > 0 {
		sum += uint32(data[index]) << 8 //剩余不够凑16个字节的放到高八位，后面补0
	}

	rt = uint16(sum) + uint16(sum>>16) //低16位与高16位相加

	return ^rt //按位取反
}

func Ping(domain string, length, Count int) {

	var (
		icmp     ICMP
		laddr    = net.IPAddr{IP: net.ParseIP("0.0.0.0")} // 得到本机的IP地址结构
		raddr, _ = net.ResolveIPAddr("ip", domain)        // 解析域名得到 IP 地址结构
	)

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	icmp = ICMP{8, 0, 0, 0, 0} //类型：8，代码：0 ，为ping命令的请求

	var buffer bytes.Buffer //大端构造ICMP报文
	binary.Write(&buffer, binary.BigEndian, icmp)
	binary.Write(&buffer, binary.BigEndian, originBytes[0:length])
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], CheckSum(b)) //放入首部校验和

	fmt.Printf("\n正在 Ping %s [%s] 具有 %d字节的数据:\n", domain, raddr.String(), length)

	recv := make([]byte, 1024)

	for i := 1; i <= Count; i++ {
		//向目标地址发送二进制报文包
		if _, err := conn.Write(b); err != nil {
			fmt.Println("请求超时。")
			continue
		}
		// 记录当前得时间
		t_start := time.Now()
		conn.SetReadDeadline((time.Now().Add(time.Second * 3)))

		len, err := conn.Read(recv)
		if err != nil {
			fmt.Println("请求超时。")
			continue
		}

		t_end := time.Now()
		dur := t_end.Sub(t_start).Milliseconds()
		var TTL = recv[8]

		fmt.Printf("来自 %s 的回复: 顺序号 = %d 时间 = %dms TTL = %d IP包总长度 = %d 实际字节数 = %d ECHO数据包字节数 = %d\n", raddr.String(), i, dur, TTL, len, len-20, len-28)
		time.Sleep(time.Second)
	}
}
func main() {
	var length, count = 32, 4
	reader := bufio.NewReader(os.Stdin)
	res, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("reader.ReadLine() error:", err)
	}

	s := strings.Split(string(res), " ")
	for i, v := range s {
		if v == "-l" {
			length, _ = strconv.Atoi(s[i+1])
			if length < 0 {
				fmt.Println("错误: 最小值为0")
				return
			}
		}
		if v == "-n" {
			count, _ = strconv.Atoi(s[i+1])
		}
	}
	Ping(s[1], length, count)
}
