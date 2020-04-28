package switches

import (
	"time"
	"fmt"
	"os"
	"github.com/brian-armstrong/gpio"
)
const FWD_IO    uint = 36
const REV_IO  uint = 37
const UP_IO   uint = 38
const DOWN_IO uint = 39
const HOME_IO     uint = 6

const SAMPLERATE uint =  20
const CAPTURETIME  uint = 1
const BUMPTIME  uint = 8
const HOLDTIME  uint = 2

func readswitch(pin gpio.Pin)(value uint){
	value,_ = pin.Read()
	value = 1 - value
	return
}
var rtimer uint
var rbutton uint
const RBBUMP uint = 8
func processRB(pin1,pin2 gpio.Pin)(result uint){
	button := readswitch(pin1) * 2 + readswitch(pin2)
	result = 0
	if rtimer == 0 {
		if button > 0 {
			rtimer++
		}
	} else if rtimer < CAPTURETIME {
		rtimer++
	} else if rtimer == CAPTURETIME {
		if button > 0 {
			rbutton = button
			rtimer++
		} else {
			rtimer = 0
		}
	} else if rtimer < BUMPTIME {
		if button > 0 {
			rtimer++
		} else {
			result = (RBBUMP+rbutton)
			rtimer = 0
		}
	} else if rtimer == BUMPTIME {
		result = (RBBUMP+rbutton+3)
		rtimer++
	} else if rtimer < (BUMPTIME + HOLDTIME - 1) {
		rtimer++
	} else {
		if button > 0{
			rtimer=BUMPTIME
		} else {
			result = (RBBUMP+7)
			rtimer=0
		}
	}
	return		
}

var otimer uint
var obutton uint
const ODOBUMP uint = 16
func processOdo(pin1,pin2 gpio.Pin)(result uint){
	button := readswitch(pin1) * 2 + readswitch(pin2)
	result = 0;
	if otimer == 0 {
		if button > 0 {
			otimer++
		}
	} else if otimer < CAPTURETIME {
		otimer++
	} else if otimer == CAPTURETIME {
		if button > 0 {
			obutton = button
			otimer++
		} else {
			otimer = 0
		}
	} else if otimer < BUMPTIME {
		if button > 0 {
			otimer++
		} else {
			result = (ODOBUMP+obutton)
			otimer = 0
		}
	} else if otimer == BUMPTIME {
		result = (ODOBUMP+obutton+3)
		otimer++
	} else if otimer < (BUMPTIME + HOLDTIME - 1) {
		otimer++;
	} else {
		if button > 0{
			otimer=BUMPTIME;
		} else {
			result = (ODOBUMP+7)
			otimer=0;
		}
	}
	return		
}


var htimer uint
var hbutton uint
const HOMEBUMP uint = 0
func processHome(pin gpio.Pin)(result uint){
	button := readswitch(pin);
	result = 0;
	if htimer == 0 {
		if button > 0 {
			htimer++
		}
	} else if htimer < CAPTURETIME {
		htimer++
	} else if htimer == CAPTURETIME {
		if button > 0 {
			hbutton = button;
			htimer++;
		} else {
			htimer = 0
		}
	} else if htimer < BUMPTIME {
		if button > 0 {
			htimer++
		} else {
			result = (HOMEBUMP+1)
			htimer = 0
		}
	} else if htimer == BUMPTIME {
		result = (HOMEBUMP+4);
		htimer++
	} else if htimer < (BUMPTIME + HOLDTIME - 1) {
		htimer++;
	} else {
		if button > 0{
			htimer=BUMPTIME;
		} else {
			result = (HOMEBUMP+7);
			htimer=0;
		}
	}
	return		
}

func Startswitches() {
	export, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
        if err != nil {
                fmt.Printf("switches.Startswitchs: failed to open /sys/class/gpio/export file for writing\n")
                fmt.Printf("switches.Startswitchs: gpio not available\n")
                return
        }
	export.Close()

	fwdpin := gpio.NewInput(FWD_IO)
	revpin := gpio.NewInput(REV_IO)
	uppin := gpio.NewInput(UP_IO)
	downpin := gpio.NewInput(DOWN_IO)
	homepin := gpio.NewInput(HOME_IO)
	htimer = 0;
	hbutton = 0;
	
	t := time.NewTicker(50 * time.Millisecond)
	for _ = range t.C {
	
		r := processRB(fwdpin,revpin)
		o := processOdo(uppin,downpin)
		h := processHome(homepin)
		if(r != 0) { 
			fmt.Println(r) 
		}
		if(o != 0) { 
			fmt.Println(o) 
		}
		if(h != 0) { 
			fmt.Println(h) 
		}

	}
	return
}

