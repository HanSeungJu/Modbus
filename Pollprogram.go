package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	      //"reflect"
	//"strings"
	"time"
)

//MBClient config
type MBClient struct {
	IP      string
	Port    int
	Timeout time.Duration
	Conn    net.Conn
}

//state show for error
const (
	Init        = "Init"
	ModbusError = "ModbusError"
	Ok          = "Ok"
	Disconnect  = "Disconnect"
)

// NewClient creates a new Modbus Client config.
func NewClient(IP string, port int, timeout time.Duration ) *MBClient {
	m := &MBClient{}
	m.IP = IP
	m.Port = port
	m.Timeout = timeout

	return m
}

//Open modbus tcp connetion
func (m *MBClient) Open() error {
	addr := m.IP + ":" + strconv.Itoa(m.Port)
	// var err error
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Panicln(err)
	}
	m.Conn = conn

	return nil
}

//Close modbus tcp connetion
func (m *MBClient) Close() {
	if m.Conn != nil {
		m.Conn.Close()
	}
}

//IsConnected for check modbus connetection
func (m *MBClient) IsConnected() bool {
	if m.Conn != nil {
		return true
	}
	return false
}

//Qurry make a modbus tcp qurry
func Qurry(conn net.Conn, timeout time.Duration, pdu []byte) ([]byte, error) {
	if conn == nil {
		return []byte{}, fmt.Errorf(Disconnect)
	}
	header := []byte{0, 0, 0, 0, byte(len(pdu) << 10), byte(len(pdu))}
	wbuf := append(header, pdu...)
	//write
	_, err := conn.Write([]byte(wbuf))
	if err != nil {
		return nil, fmt.Errorf(Disconnect)
	}

	//read
	rbuf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(timeout))
	len, err := conn.Read(rbuf)
	if err != nil {
		return nil, fmt.Errorf(Disconnect)
	}
	if len < 10 {
		return nil, fmt.Errorf(ModbusError)
	}
	return rbuf[6:len], nil
}

//ReadCoil mdbus function 1 qurry and return []uint16
func (m *MBClient) ReadCoil(id uint8, addr uint16, leng uint16) ([]bool, error) {
	pdu := []byte{id, 0x01, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return []bool{}, err
	}
	//convert
	result := []bool{}
	bc := res[2]
	for i := 0; i < int(bc); i++ {
		for j := 0; j < 8; j++ {
			if (res[3+i] & (byte(1) << byte(j))) != 0 {
				result = append(result, true)
			} else {
				result = append(result, false)
			}
		}
	}
	result = result[:leng]
	return result, nil
}

//ReadCoilIn mdbus function 2 qurry and return []uint16
func (m *MBClient) ReadCoilIn(id uint8, addr uint16, leng uint16) ([]bool, error) {

	pdu := []byte{id, 0x02, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return []bool{}, err
	}

	//convert
	result := []bool{}
	bc := res[2]
	for i := 0; i < int(bc); i++ {
		for j := 0; j < 8; j++ {
			if (res[3+i] & (byte(1) << byte(j))) != 0 {
				result = append(result, true)
			} else {
				result = append(result, false)
			}
		}
	}
	result = result[:leng]

	return result, nil
}

//ReadReg mdbus function 3 qurry and return []uint16
func (m *MBClient) ReadReg(id uint8, addr uint16, leng uint16) ([]uint16, error) {

	pdu := []byte{id, 0x03, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return []uint16{}, err
	}
	//convert
	result := []uint16{}
	for i := 0; i < int(leng); i++ {
		var b uint16
		b = uint16(res[i*2+3]) << 8
		b |= uint16(res[i*2+4])
		result = append(result, b)
	}

	return result, nil
}

//ReadRegIn mdbus function 4 qurry and return []uint16
func (m *MBClient) ReadRegIn(id uint8, addr uint16, leng uint16) ([]uint16, error) {

	pdu := []byte{id, 0x04, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return []uint16{}, err
	}

	//convert
	result := []uint16{}
	for i := 0; i < int(leng); i++ {
		var b uint16
		b = uint16(res[i*2+3]) << 8
		b |= uint16(res[i*2+4])
		result = append(result, b)
	}

	return result, nil
}

//WriteCoil mdbus function 5 qurry and return []uint16
func (m *MBClient) WriteCoil(id uint8, addr uint16, data bool) error {

	var pdu = []byte{}
	if data == true {
		pdu = []byte{id, 0x5, byte(addr >> 8), byte(addr), 0xff, 0x00}
	} else {
		pdu = []byte{id, 0x5, byte(addr >> 8), byte(addr), 0x00, 0x00}
	}

	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return err
	}

	return nil
}

//WriteReg mdbus function 6 qurry and return []uint16
func (m *MBClient) WriteReg(id uint8, addr uint16, data uint16) error {

	pdu := []byte{id, 0x06, byte(addr >> 8), byte(addr), byte(data >> 8), byte(data)}

	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return err
	}

	return nil
}

//WriteCoils mdbus function 15(0x0f) qurry and return []uint16
func (m *MBClient) WriteCoils(id uint8, addr uint16, data []string) error {
	pdu := []byte{}
	if len(data)%8 == 0 {
		pdu = []byte{id, 0x0f, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data) / 8)}
	} else {
		pdu = []byte{id, 0x0f, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data)/8) + 1}
	}
	var tbuf byte
	for i := 0; i < len(data); i++ {
		pb, _ := strconv.ParseBool(data[i]);
		if pb {
			tbuf |= byte(1 << uint(i%8))
		}

		if (i+1)%8 == 0 || i == len(data)-1 {
			pdu = append(pdu, tbuf)
			tbuf = 0
		}
	}

	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return err
	}
	return nil
}
//WriteRegs mdbus function 16(0x10) qurry and return []uint16


func (m *MBClient) WriteRegs(id uint8, addr uint16, data []string)  error {
	fmt.Println("id:",id,""," addr:",addr,"")


	//var data []byte
	pdu := []byte{id, 0x10, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data)) * 2}
	for i := 0; i < len(data); i++ {
		pi, _ := strconv.ParseUint(data[i], 10, 16);
		if pi == 0{
			fmt.Println("")
			main()
		}
		fmt.Println("alias  ",i, ": ", pi)
		pdu = append(pdu, byte(pi>>8))
		pdu = append(pdu, byte(pi))
		//fmt.Println(pdu)
	}
	println("안녕하세용")

	//write
	_, err := Qurry(m.Conn, m.Timeout, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			m.Close()
			m.Conn = nil
		}
		return err
	}
	return nil
}

func main() {
	//var arr []string
	var num int
	var alias uint16
	var leng uint16
	//TCP coonnetion
	mbc := NewClient("192.168.0.222", 502, time.Second)

	err := mbc.Open()
	if err != nil {
		log.Panicln(err)
		os.Exit(12)
	}

	defer mbc.Close()

	data, _ := mbc.ReadReg(1, 0, 10)
	log.Println(data)
	input := bufio.NewScanner(os.Stdin)
loop :
	for {
		fmt.Print("\n\n\n\n\n\n##Modbus protocol##\n")
		fmt.Print("1 : Output Coils\n","2 : Input Coils\n","3 : Input Registers\n","4 : Holding Registers\n\n\n\n\n")
		fmt.Print("Select number Enter:")
		fmt.Scanln(&num,"\n\n\n\n\n\n")
		if num == 1{
		loop2:
			for {
				fmt.Print("\n\nno.1    : Write Coils \n","no.2    : Read Coils\n","Any key : Go back\n\n\n")
				fmt.Print("Select number Enter:")
				fmt.Scanln(&num)
				if num == 1{
					var arr []string
					fmt.Println("\n\n\n\n@@Output Coils@@\nWrite Coils\n\n\n")
					fmt.Print("\nalias values: ")
					fmt.Scanln(&alias)
					fmt.Print("leng values:")
					fmt.Scanln(&leng)
					for i := 0; i < int(leng); i++ {
						fmt.Print("Enter data values: ")
						input.Scan()
						err, _ := strconv.Atoi(input.Text())
						if err == 1 || err == 0{
							arr = append(arr,input.Text())
							mbc.WriteCoils(1, alias, arr)
						} else{
							fmt.Println("\n\n\nPlease enter again")
							continue loop2
						}
					}

				}
				if num == 2{
					fmt.Println("\n\n\n\n@@Output Coils@@\nRead Coils\n\n\n")
					fmt.Print("\nalias values: ")
					fmt.Scanln(&alias)
					fmt.Print("leng values:")
					fmt.Scanln(&leng)
					data, _ := mbc.ReadCoil(1, alias, leng)
					fmt.Println("Data values : ",data)
					continue loop
				}else{
					fmt.Println("Go back")
					continue loop
				}
			}
		}
		if num == 2{

			for {
				fmt.Print("\n\nno.1    : Read Coils Input\n","Any key : Go back\n\n\n")
				fmt.Print("Select number Enter:")
				fmt.Scanln(&num)
				if num == 1{
					fmt.Println("\n\n\n\n@@Input Coils@@\nRead Coils Input\n\n\n")
					fmt.Print("\nalias values: ")
					fmt.Scanln(&alias)
					fmt.Print("leng values:")
					fmt.Scanln(&leng)
					data, _ := mbc.ReadCoilIn(1, alias, leng)
					fmt.Println("Data value:",data,"\n\n")
				}else{
					fmt.Println("Go back")
					continue loop
				}
			}
		}
		if num == 3{

			for {
				fmt.Print("\n\nno.1    : Read Register Input\n","Any key : Go back\n\n\n")
				fmt.Print("Select number Enter:")
				fmt.Scanln(&num)
				if num == 1{
					fmt.Println("\n\n\n\n@@Input Registers@@\nRead Register Input\n\n\n")
					fmt.Print("\nalias values: ")
					fmt.Scanln(&alias)
					fmt.Print("leng values:")
					fmt.Scanln(&leng)
					data, _ := mbc.ReadRegIn(1, alias, leng)
					fmt.Println("Data values : ",data)
				}else{
					fmt.Println("Go back")
					continue loop
				}
			}
		}

		if num == 4 {
		loop4:
			for {
				fmt.Print("\n\nno.1    : Write Registers\n", "no.2    : Read Registers\n","Any key : Go back\n\n\n")
				fmt.Print("Select number Enter:")
				fmt.Scanln(&num)

				if num == 1 {
					var arr []string
					fmt.Println("\n\n\n\n@@Holding Registers@@\nWrite Registers\n\n\n")
					fmt.Print("\nalias values: ")
					fmt.Scanln(&alias)
					fmt.Print("leng values:")
					fmt.Scanln(&leng)
					for i := 0; i < int(leng); i++ {
						fmt.Print("Enter Data values: ")
						input.Scan()
						err, _ := strconv.Atoi(input.Text())

						if err > 6555 {
							fmt.Println("max Excess error")
							continue loop
						}
						arr = append(arr, input.Text())
						mbc.WriteRegs(1, alias, arr)
					}
				}
				if num == 2{
					fmt.Println("\n\n\n\n@@Holding Registers@@\nRead Registers\n\n\n")
					fmt.Print("\nalias 입력: ")
					fmt.Scanln(&alias)
					fmt.Print("leng 입력:")
					fmt.Scanln(&leng)
					data, _ := mbc.ReadReg(1, alias, leng)
					fmt.Println("Data values : ",data)
					continue loop4
				}else{
					fmt.Println("Go back")
					continue loop
				}
			}
		}
	}
}


