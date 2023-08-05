package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/ZHJBOBO/multipath-quic-go/example/sample/utils"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/mux"
	lmq "github.com/libp2p/go-libp2p-mpquic-transport"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	//"main/utils"
	"os"
	"strconv"
	"strings"
	"time"
)
const BUFFERSIZE int64 = 1024
func download(stream mux.MuxedStream, savePath string, addr string) {

	fmt.Println("addr: ", addr)
	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	stream.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	fmt.Println("file size received: ", fileSize)

	stream.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("file name received: ", fileName, "time: ", time.Now().Format("2006-01-02 15:04:05.000"))

	newFile, err := os.Create(savePath + "/" + fileName)
	utils.HandleError(err)

	var receivedBytes int64
	start := time.Now()

	for {
		//fmt.Println(fileName, " 0")
		if (fileSize-receivedBytes) < BUFFERSIZE && (fileSize-receivedBytes) > 0 {

			recv, err := io.CopyN(newFile, stream, (fileSize - receivedBytes))
			if err != nil {
				fmt.Println(err)
			}
			receivedBytes += recv
			fmt.Printf("\033[2K\rReceived: %d / %d", receivedBytes, fileSize)
			break
		}
		//fmt.Println(fileName, "1")
		cn, err := io.CopyN(newFile, stream, BUFFERSIZE)
		if err != nil && err != io.EOF {
			fmt.Println("error_addr:"+addr+"time:"+time.Now().Format("2006-01-02 15:04:05.000"), "copysize:", cn)
			fmt.Println("streamID is ", stream)
			fmt.Println(err)
		}
		//fmt.Println(fileName, "2")
		receivedBytes += cn

		fmt.Printf("\033[2K\rReceived: %d / %d", receivedBytes, fileSize)

		if receivedBytes == fileSize || err == io.EOF {
			break
		}
		//fmt.Println(fileName, "3")
	}
	elapsed := strconv.FormatInt(time.Now().Sub(start).Nanoseconds(), 10)
	fmt.Println("\nTransfer took: ", elapsed)

	fmt.Println("\n\nReceived file completely!")
	newFile.Close()

	//var content string = fileName + "," + strconv.FormatInt(fileSize, 10) + "," + elapsed + "," + start.Format("2006-01-02 15:04:05.000") + "," + time.Now().Format("2006-01-02 15:04:05.000") + "\n"
	//_, err = fd.WriteString(content)
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func main() {
	savePath := "D:\\testdata"
	fmt.Println("Saving file to: ", savePath)

	//addr_l := os.Args[1] + ":4242"
	//localAddrs := []string{addr_l}

	fmt.Println("Server started! Waiting for streams from client...")

	////回写数据到csv
	//var header string = "file,file_size,transfer_time(ns),start_time,end_time" + "\n"
	//
	//fd, err := os.OpenFile("receive.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	//if err != nil {
	//	log.Println("open file error :", err)
	//	return
	//}
	//n, err := fd.WriteString(header)
	//if err == nil && n < len(header) {
	//	return
	//}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	//fmt.Println(rsaKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	key, err := ic.UnmarshalRsaPrivateKey(x509.MarshalPKCS1PrivateKey(rsaKey))
	publickey:=key.GetPublic()
	publickeyMarshal,err:=ic.MarshalPublicKey(publickey)
	fmt.Println("publickeyMarshal:",publickeyMarshal)
	for _,v:= range publickeyMarshal{
		fmt.Print(v,",")
	}
	//fmt.Println(key)
	//fmt.Println(x509.MarshalPKCS1PrivateKey(rsaKey))
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := lmq.NewTransport(key, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/4242/quic")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Trying to listen to: ", localAddr.String())
	listener, err := t.Listen(localAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(listener)
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("conn created: ", conn)
	for {

		stream, err := conn.AcceptStream()
		if err != nil {
			fmt.Println("acceptstream error!")
		}
		fmt.Println("stream created: ", stream)
		go download(stream, savePath, localAddr.String())
	}

	select {}
	return

}
