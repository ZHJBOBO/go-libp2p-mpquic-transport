package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/ZHJBOBO/multipath-quic-go/example/sample/utils"
	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/transport"
	libp2pmpquic "github.com/libp2p/go-libp2p-mpquic-transport"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	//lmq "github.com/libp2p/go-libp2p-mpquic-transport"
	ma "github.com/multiformats/go-multiaddr"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const BUFFERSIZE int64 = 1024

func upload(path string, file fs.FileInfo, conn transport.CapableConn, fd *os.File) {

	var file_path string = filepath.Join(path, file.Name())
	stream, err := conn.OpenStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	starttime := time.Now()
	var repeat_count int = 1
	var avg_elapsed int64 = 0
	for i := 1; i <= repeat_count; i++ {
		avg_elapsed += sendFile(stream, file_path)
	}

	avg_elapsed = avg_elapsed / int64(repeat_count)
	stream.Close()
	var content string = file_path + "," + strconv.FormatInt(file.Size(), 10) + "," + strconv.FormatInt(avg_elapsed, 10) + "," + starttime.Format("2006-01-02 15:04:05.000") + "," + time.Now().Format("2006-01-02 15:04:05.000") + "," + strconv.FormatInt(time.Now().Sub(starttime).Nanoseconds(), 10) + "\n"
	_, err = fd.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("stream.Close()")

}

func sendFile(stream mux.MuxedStream, fileToSend string) int64 {

	file, err := os.Open(fileToSend)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := utils.FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := utils.FillString(fileInfo.Name(), 64)

	//fmt.Println("Sending filename and filesize!")
	stream.Write([]byte(fileSize))
	stream.Write([]byte(fileName))

	var sentBytes int64
	start := time.Now()

	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		n, err := file.Read(sendBuffer)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			return -1
		}

		// 如果已经读到文件结尾，退出循环
		if err != nil && err.Error() == "EOF" && n == 0 {
			break
		}

		// 将读取到的数据发送到远程服务器
		n1, err := stream.Write(sendBuffer[:n])
		if err != nil {
			fmt.Println(err)
			fmt.Println("size n:", n, " n1:", n1)
			return -2
		}
		sentBytes += int64(n)
		fmt.Printf("\033[2K\rSent: %d / %d", sentBytes, fileInfo.Size())
	}

	end := time.Now()
	elapsed := end.Sub(start).Nanoseconds()
	fmt.Printf("\n file,%s,transfer time(ns),%s,size %s", fileToSend, string(elapsed), string(sentBytes))

	file.Close()

	return elapsed
}

func main() {
	bgcontext := context.Background()
	//回写数据到csv
	var header string = "file,file_size,transfer_time(ns),start_time,end_time,since_time" + "\n"

	f0_result, err := os.OpenFile("send.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	n, err := f0_result.WriteString(header)
	if err == nil && n < len(header) {
		return
	}

	//发送文件处理
	dataset_path := "C:\\Users\\24201\\Desktop\\libp2p-test-fbc\\libp2p-test\\send"
	files, err := ioutil.ReadDir(dataset_path)
	path, _ := filepath.Abs(dataset_path)
	if err != nil {
		log.Fatal(err)
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	key, err := ic.UnmarshalRsaPrivateKey(x509.MarshalPKCS1PrivateKey(rsaKey))
	fmt.Println(key.GetPublic())
	t, err := libp2pmpquic.NewTransport(key, nil, nil)
	remoteAddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/4242/quic")
	fmt.Println("Trying to connect to: ", remoteAddr.String())

	conn, err := t.Dial(bgcontext, remoteAddr, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("conn created: ", conn)

	resultFileList := []*os.File{f0_result}

	for _, file := range files {
		upload(path, file, conn, resultFileList[0])
	}

	if err0 := f0_result.Close(); err == nil {
		err = err0
	}

	return

}
