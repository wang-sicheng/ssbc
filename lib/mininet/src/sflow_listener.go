package main

import (
	"bytes"
	"fmt"
	"github.com/Cistern/sflow"
	"net"
)

func sFlowParser(buffer []byte) {
	reader := bytes.NewReader(buffer)
	d := sflow.NewDecoder(reader)
	dgram, err := d.Decode()
	if err != nil {
		fmt.Println("Decode err",err)

		return
	}
	fmt.Println(dgram)
	for _, sample := range dgram.Samples {

		fmt.Println(sample)

	}
}

func sFlowListener(AppState *app_state, SFlowConfig sflow_config) (err error) {
	defer AppState.Wait.Done()

	var udp_addr = fmt.Sprintf("[%s]:%d", SFlowConfig.Address, SFlowConfig.Port)

	DebugLogger.Println("Binding sFlow listener to", udp_addr)

	//UDPAddr, err := net.ResolveUDPAddr("udp", udp_addr)
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 6343,
	})
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}
	defer conn.Close()
	for AppState.Running{
		data := make([]byte, 10000)
		read, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("read data ", err)
			return err
		}

		fmt.Println(read, remoteAddr)

		//fmt.Printf("%s\n", data)
		sFlowParser(data)
		//send_data := []byte("hi client!")
		//_, err = conn.WriteToUDP(send_data, remoteAddr)
		//if err != nil {
		//	return err
		//	fmt.Println("send fail!", err)
		//}
	}
	//var buffer []byte
	//for AppState.Running {
	//	/*
	//	  Normally read would block, but we want to be able to break this
	//	  loop gracefuly. So add read timeout and every 0.1s check if it is
	//	  time to finish
	//	*/
	//	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	//
	//	var read, _, err = conn.ReadFromUDP(buffer)
	//	if err!=nil{
	//		fmt.Println(err)
	//	}
	//	if read > 0  {
	//		fmt.Println("in")
	//		sFlowParser(buffer)
	//	}
	//
	//}



	return nil
}
