package gps

import(
	"time"
	"sync"
	"github.com/ahmedaly113/roadbook/model"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
)

const ascii_0 byte = 48
const ascii_a byte = 97
const ascii_A byte = 65
const ascii_dollar byte = 36
const ascii_star byte = 42
 

func asciihextobyte(hex1 byte, hex2 byte) (value byte) {

	value = 0
	if hex1 >= ascii_a {
		value = (hex1 - ascii_a + 10) * 16
	} else if hex1 >= ascii_A {
		value = (hex1 - ascii_A + 10) * 16
	} else {
		value = (hex1 - ascii_0) * 16
	}

	if hex2 >= ascii_a {
		value += (hex2 - ascii_a + 10)
	} else if hex2 >= ascii_A {
		value += (hex2 - ascii_A + 10)
	} else {
		value += (hex2 - ascii_0)
	}
	return
}

func validatechecksum(sentence []byte, length int) (valid bool){
	var checksum byte = 0
	valid = false
	for i := 1; i < length-3;i++ {
		if sentence[i] == ascii_star {
			//fmt.Printf("%x %s\n",checksum, sentence[i:i+3])
			checksum ^= asciihextobyte(sentence[i+1],sentence[i+2])
			if checksum == 0 {
				valid = true;
			}
			return
		} else {
			checksum ^= sentence[i];
		}
	}
	return
}

func Gpsreader(state *model.Model,stateMu *sync.Mutex,devicename string) {
	if devicename == "/dev/null" { //no gps input intended, exit thread quietly
		return  
	}
	options := serial.OpenOptions{
		PortName: devicename,
		BaudRate: 19200,
		DataBits: 8,
		StopBits: 1,
		InterCharacterTimeout: 0,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	if err != nil { //unable to open serial port or access ioctl, return with diagnostic
		fmt.Println("gps:Gpsreader() Unable to open serial port:",devicename)
		fmt.Println("gps:Gpsreader():", err)
		return
	}
	defer port.Close()

	//19200 baud is about 1920bytes per sec.  96 bytes per 50mS
	const buflength int = 200
	t := time.NewTicker(50 * time.Millisecond)

	bufread := make([]byte,buflength)
	bufassemble := make([]byte,buflength)
	bufassembleindex := 0
	bufparse := make([]byte,buflength)

	for _ = range t.C {
		n,_ := port.Read(bufread)
//		fmt.Println(n,string(b[:n]))
		for i:=0; i<n; i++ {
			if bufread[i] == ascii_dollar {
				copy(bufparse,bufassemble)
				bufassemble[0] = bufread[i]
				for j:=1; j<buflength;j++{
					bufassemble[j]=0;
				}
				bufassembleindex = 1
				//fmt.Println("buf",string(so))
				if validatechecksum(bufparse,buflength) {
					_,di,sp,he,_ := Parsenmea(string(bufparse))
					stateMu.Lock()
					state.Heading = float32(he)
					state.Speed = float32(sp)
					state.Distance += float32(di)
					stateMu.Unlock()
				}
			} else {
				bufassemble[bufassembleindex] = bufread[i]
				if bufassembleindex < buflength-1 { 
					bufassembleindex++ 
				} 
			}
		}	
	}
}
