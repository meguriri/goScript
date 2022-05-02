package main

import (
	"archive/zip"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-gomail/gomail"
)

var (
	fileNames []string
	zipName   string
	mailTitle string
	Sender    string
	Receiver  []string
	password  string
)

func GetFiles(folder string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("获取目录文件失败: " + err.Error())
		return
	}
	fmt.Println("获取" + folder + "目录下文件: ")
	for _, v := range files {
		if v.Name() == "mail.exe" {
			continue
		}
		if v.IsDir() {
			GetFiles(folder + "/" + v.Name())
		} else {
			fmt.Println(folder + "/" + v.Name())
			fileNames = append(fileNames, folder+"/"+v.Name())
		}
	}
}

func Zip(files []string) {

	archive, err := os.Create(zipName + ".zip")
	if err != nil {
		fmt.Println("创建压缩文件失败: " + err.Error())
		return
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	for _, v := range files {
		f1, _ := os.Open(v)
		w1, err := zipWriter.Create(v)
		if err != nil {
			fmt.Println("创建zipWriter失败: " + err.Error())
			return
		}
		if _, err := io.Copy(w1, f1); err != nil {
			fmt.Println("写入压缩文件失败: " + err.Error())
			return
		}
		f1.Close()
	}

	zipWriter.Close()
}

func SetMessage() {
	fmt.Println("请输入发送方邮箱: ")
	fmt.Scan(&Sender)
	fmt.Println("请输入接收方邮箱: (可输入多个接受邮箱，以#为输入结束)")
	for {
		rec := ""
		fmt.Scan(&rec)
		if rec == "#" {
			break
		}
		Receiver = append(Receiver, rec)
	}
	fmt.Println("请输入邮件标题: (无空格！)")
	fmt.Scan(&mailTitle)
	fmt.Println("请输入压缩文件名称: ")
	fmt.Scan(&zipName)
	fmt.Println("请输入发送方邮箱授权码: ")
	fmt.Scan(&password)
}

func PostMail() {
	m := gomail.NewMessage()
	m.Attach("./" + zipName + ".zip")
	m.SetHeader("From", Sender)       //发送者
	m.SetHeader("To", Receiver...)    //接受者
	m.SetHeader("Subject", mailTitle) // 邮件标题

	//这里第一个参数为服务器地址，第二个为端口号，第三个为发送者邮箱号
	d := gomail.NewDialer("smtp.qq.com", 465, Sender, password) //

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("发送失败: ", err.Error())
	} else {
		fmt.Println("发送成功! 发送压缩文件为" + "./" + zipName + ".zip")
	}
	os.Remove("./" + zipName + ".zip")
}

func main() {
	SetMessage()
	GetFiles(".")
	Zip(fileNames)
	PostMail()
	time.Sleep(time.Second * 3)
}
