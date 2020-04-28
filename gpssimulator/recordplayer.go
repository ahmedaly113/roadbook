package main


import(
	"flag"
	"time"
	"os"
	"bufio"
	"github.com/jacobsa/go-serial/serial"
//	"fmt"
)
func trap(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	serialportPtr := flag.String("p","/dev/ttyUSB1","serial device")
	filenamePtr := flag.String("f","../assets/samples/home-to-office-route/ico_gps_log.txt","log file name")

	flag.Parse()

	repeat := false

	options := serial.OpenOptions{
		PortName: *serialportPtr,
		BaudRate: 19200,
		DataBits: 8,
		StopBits: 1,
		InterCharacterTimeout: 0,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	trap(err)
	defer port.Close()
	
	t := time.NewTicker(200 * time.Millisecond)
	f,err := os.Open(*filenamePtr)
	trap(err)
	defer func() {
		err = f.Close()
		trap(err)
	}()

	s := bufio.NewScanner(f)

	for _ = range t.C {
		if s.Scan() {
			//fmt.Println(s.Text())
			port.Write([]byte(s.Text()))
			port.Write([]byte("\n\r"))
		} else {
			trap(s.Err())
			//EOF returns false from Scan() but does not set an error
			if repeat {
				f.Seek(0,0)
				s = bufio.NewScanner(f)
			} else {
				return
			}
		}
	}

}

