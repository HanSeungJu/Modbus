package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	//"reflect"
	//"strings"

)
func Scrclr(){
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

}

//MBClient config
type MBClient struct {
	IP      string
	Port    int

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
func NewClient(IP string, port int) *MBClient {
	m := &MBClient{}
	m.IP = IP
	m.Port = port


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
	os.Exit(12)
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
func Qurry(conn net.Conn, pdu []byte) ([]byte, error) {
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
	//conn.SetReadDeadline(time.Now().Add(timeout))
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
func (m *MBClient) ReadCoil(id uint8, addr uint16, leng uint16) ([]int, error) {
	pdu := []byte{id, 0x01, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	res, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return []int{}, err
	}
	//convert
	result := []int{}
	bc := res[2]
	for i := 0; i < int(bc); i++ {
		for j := 0; j < 8; j++ {
			if (res[3+i] & (byte(1) << byte(j))) != 0 {
				result = append(result, 1)
			} else {
				result = append(result, 0)
			}
		}
	}
	result = result[:leng]
	return result, nil
}

//ReadCoilIn mdbus function 2 qurry and return []uint16
func (m *MBClient) ReadCoilIn(id uint8, addr uint16, leng uint16) ([]int, error) {

	pdu := []byte{id, 0x02, byte(addr >> 8), byte(addr), byte(leng >> 8), byte(leng)}

	//write
	res, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return []int{}, err
	}

	//convert
	result := []int{}
	bc := res[2]
	for i := 0; i < int(bc); i++ {
		for j := 0; j < 8; j++ {
			if (res[3+i] & (byte(1) << byte(j))) != 0 {
				result = append(result, 1)
			} else {
				result = append(result, 0)
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
	res, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
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
	res, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
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
	_, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
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
	_, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
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
		pb, _ := strconv.ParseBool(data[i])
		pa,_ := strconv.Atoi(data[i])
		fmt.Println("alias  ", addr, ": ", pa)
		if pb {
			tbuf |= byte(1 << uint(i%8))
		}

		if (i+1)%8 == 0 || i == len(data)-1 {
			pdu = append(pdu, tbuf)
			tbuf = 0
		}
		addr++
	}
	//write
	_, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil
		}
		return err
	}
	return nil
}
//WriteRegs mdbus function 16(0x10) qurry and return []uint16


func (m *MBClient) WriteRegs(id uint8, addr uint16, data []string)  error {


	//var data []byte
	pdu := []byte{id, 0x10, byte(addr >> 8), byte(addr), byte(len(data) >> 8), byte(len(data)), byte(len(data)) * 2}
	for i := 0; i < len(data); i++ {
		pi, _ := strconv.ParseUint(data[i], 10, 16)
		fmt.Println("alias  ",addr, ": ", pi)
		pdu = append(pdu, byte(pi>>8))
		pdu = append(pdu, byte(pi))
		//fmt.Println(pdu)
		addr++
	}


	//write
	_, err := Qurry(m.Conn, pdu)
	if err != nil {
		if err.Error() == Disconnect {
			fmt.Println("\n\n\n\n\n\n@@@Disconnect error@@@\n\n\n\n\n\n")
			m.Close()
			m.Conn = nil


		}
		return err
	}
	return nil
}

func Error(){
	println("-------------------------------------")
	fmt.Println("error")
	println("-------------------------------------")
	fmt.Println("You entered it iolncorrectly.")
	println("-------------------------------------\n\n\n")
	fmt.Print("1 : Back\n","2 : Main menu\n\n\n")
	fmt.Print("Select number Enter:")
}
func Continue(){
	println("\n\n\n-------------------------------------")
	fmt.Println("1 : Back","\n2 : Main menu")
	println("-------------------------------------\n\n\n")
	fmt.Print("Select number Enter:")
}

func main() {
	//var arr []string
	var num int
	var add uint16
	var leng uint16
	//TCP coonnetion
	mbc := NewClient("192.168.0.222", 502)
	mbc.Open()
	/*err := mbc.Open()
	if err != nil {
		log.Println("disconnect",err)

	}*/

	defer mbc.Close()

	data, _ := mbc.ReadReg(1, 0, 10)
	log.Println(data)

	input := bufio.NewScanner(os.Stdin)
Main:
	for {
		Scrclr()
		println("-------------------------------------")
		fmt.Println("Main Menu")
		println("-------------------------------------")
		fmt.Print("1 : Output Coils\n","2 : Input Coils\n","3 : Input Registers\n","4 : Holding Registers\n\n\n\n\n")
		fmt.Print("Select number Enter:")
		fmt.Scanln(&num)

		if num == 1{
			OutputCoils:
			for {
				Scrclr()
				println("-------------------------------------")
				fmt.Println("Output Coils")
				println("-------------------------------------")
				fmt.Print("1 : Write Coils\n","2 : Read Coils\n","3 : Go back\n\n\n")
				fmt.Print("Select number Enter:")
				fmt.Scanln(&num)

				if num == 1{

					Scrclr()
					var arr []string
					var leng int

					println("-------------------------------------")
					fmt.Println("Write Coils")
					println("-------------------------------------")
					fmt.Print("\n\n\nStart address: ")
					fmt.Scanln(&add)
					fmt.Print("length values:")
					fmt.Scanln(&leng)

					fmt.Print("Enter data values: ")
					input.Scan()
					err, _ := strconv.Atoi(input.Text())
					arr = strings.Split(input.Text(), " ")
						if err == 1 || err == 0  && leng == len(arr) && add < 100{
							mbc.WriteCoils(1, add, arr)
							Continue()
							fmt.Scanln(&num)
							if num == 1 {
								continue OutputCoils
							}else if num == 2{
								continue Main
							}
						} else {
							Scrclr()
							Error()
							fmt.Scanln(&num)
							if num == 1{
								continue OutputCoils
							}else if num == 2{
								continue Main
							}
						}
				}
				if num == 2{
					Scrclr()
					println("-------------------------------------")
					fmt.Println("Read Coils")
					println("-------------------------------------")
					fmt.Print("\n\n\nStart address: ")
					fmt.Scanln(&add)
					fmt.Print("length values:")
					fmt.Scanln(&leng)

					data, _ := mbc.ReadCoil(1, add, leng)
					fmt.Println("Data values : ", data)
					Continue()
					fmt.Scanln(&num)
					if num == 1 {
						continue OutputCoils
						} else if num == 2 {
							continue Main
						}
					}
				if num == 3{
					continue Main
				} else{
					fmt.Println("Please enter again")
					continue OutputCoils
				}
			}
		}
		if num == 2{
			InputCoils:
			for {
				Scrclr()
				println("-------------------------------------")
				fmt.Println("Input Coils")
				println("-------------------------------------")
				fmt.Print("\n\n\nStart address: ")
				fmt.Scanln(&add)
				fmt.Print("length values:")
				fmt.Scanln(&leng)

				data, _ := mbc.ReadCoilIn(1, add, leng)
				fmt.Println("Data value:",data)
				Continue()
				fmt.Scanln(&num)
				if num == 1 {
					continue InputCoils
				} else if num == 2 {
					continue Main
				}
			}
		}
		if num == 3{
			InputRegisters:
			for {
				Scrclr()
				println("-------------------------------------")
				fmt.Println("Input Registers")
				println("-------------------------------------")
				fmt.Print("\n\n\nStart address: ")
				fmt.Scanln(&add)
				fmt.Print("length values:")
				fmt.Scanln(&leng)
				data, _ := mbc.ReadRegIn(1, add, leng)
				fmt.Println("Data values : ",data)
				Continue()
				fmt.Scanln(&num)
				if num == 1 {
					continue InputRegisters
				} else if num == 2 {
					continue Main
				}
			}
		}

		if num == 4 {
		HoldingRegisters:
			for {

				Scrclr()
				println("-------------------------------------")
				fmt.Println("Holding Registers")
				println("-------------------------------------")
				fmt.Print("1 : Write Registers\n","2 : Read Registers\n","3 : Go back\n\n\n")
				fmt.Print("Select number Enter:")
				fmt.Scanln(&num)
				if num == 1 {
					Scrclr()
					var arr []string
					var leng int
					println("-------------------------------------")
					fmt.Println("Write Registers")
					println("-------------------------------------")
					fmt.Print("\n\n\nStart address: ")
					fmt.Scanln(&add)
					fmt.Print("length values:")
					fmt.Scanln(&leng)
					fmt.Print("Enter Data values: ")
					input.Scan()
					err, _ := strconv.Atoi(input.Text())

					if err > 65535 {
						Scrclr()
						fmt.Println("\n[Max Excess error]")
						Error()
						fmt.Scanln(&num)
						if num == 1{
							continue HoldingRegisters
						}else if num == 2{
							continue Main
						}
					}
						arr = strings.Split(input.Text()," ")
						//bbb := strconv.Itoa(leng)
						if leng == len(arr){
							mbc.WriteRegs(1, add, arr)
							Continue()
							fmt.Scanln(&num)
							if num == 1 {
								continue HoldingRegisters
							} else if num == 2 {
								continue Main
							}
						}else{
							Scrclr()
							fmt.Println("\n[Entered incorrectly length values]")
							Error()
							fmt.Scanln(&num)
							if num == 1{
								continue HoldingRegisters
							}else if num == 2{
								continue Main
							}
						}

				}
				if num == 2{
					println("-------------------------------------")
					fmt.Println("Read Registers")
					println("-------------------------------------")
					fmt.Print("\n\n\nStart address: ")
					fmt.Scanln(&add)
					fmt.Print("length values:")
					fmt.Scanln(&leng)
					data, _ := mbc.ReadReg(1, add, leng)
					fmt.Println("Data values : ",data)
					Continue()
					fmt.Scanln(&num)
					if num == 1 {
						continue HoldingRegisters
					} else if num == 2 {
						continue Main
					}
				}
				if num == 3{
					continue Main
				} else{
					continue HoldingRegisters
				}
			}
		}
	}
}
	/*readData := make([]byte, 3)
	  readData[0] = byte(200 >> 8)   // (High Byte)
	  readData[1] = byte(200 & 0xff) // (Low Byte)
	  readData[2] = 0x01
	  Rx, rerr := TCPRead(m.Conn, 300, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00, readData, trace)
	  if  rerr != nil{
	          log.Println(rerr)
	  }
	  log.Println(Rx)*/

