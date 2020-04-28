package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/ahmedaly113/roadbook/data"
	"github.com/ahmedaly113/roadbook/gps"
	"github.com/ahmedaly113/roadbook/model"
	"github.com/ahmedaly113/roadbook/safe"
	"github.com/ahmedaly113/roadbook/switches"
	"github.com/ahmedaly113/roadbook/ui"
	"github.com/vharitonsky/iniflags"
)

var stateMu sync.Mutex
var state *model.Model
var dataDirPtr *string

func main() {
	// DEFAULTS are for the PRODUCTION HARDWARE
	// override in your own INI file (e.g. ahmedaly113.ini)
	serialportPtr := flag.String("port", "/dev/ttymxc4", "serial device")
	filenamePtr := flag.String("file", "", "playback log file")
	enableGpioPtr := flag.Bool("gpio", true, "enable/disable gpio hardware")
	dataDirPtr = flag.String("data_dir", "", "application data dir")
	currentRoadbookPtr := flag.String("current_roadbook", "", "current roadbook file")

	iniflags.Parse()

	// Load current book. All access to state that is potentially shared by
	// concurrent goroutines (even reads) must be guarded by a sync.Mutex.
	stateMu.Lock()
	if *currentRoadbookPtr != "" {
		// state = data.LoadDummy()
		file, err := os.Open(*currentRoadbookPtr)
		if err != nil {
			panic(err)
		}
		bytes, err := ioutil.ReadAll(bufio.NewReader(file))
		if err != nil {
			panic(err)
		}
		s, err := data.FromGPX(bytes)
		if err != nil {
			panic(err)
		}
		state = s
	} else {
		state = &model.Model{
			Book:     []model.Waypoint{},
			Distance: float32(0),
			Speed:    float32(0),
			Heading:  float32(0),
			Idx:      0,
		}
	}

	// This is for SpeedZone UI Testing
	state.IsSpeedZone = true
	state.SpeedLimit = 50

	stateMu.Unlock()

	// Run simulator in lieu of real sensors. Simulator is launched as a goroutine,
	// allowing it to execute concurrently with the rendering.
	//go gps.Simulator(state,&stateMu)
	if *filenamePtr != "" {
		go safe.Do("gps Playback", func() { gps.Playback(state, &stateMu, *filenamePtr, false) })
	} else {
		go safe.Do("gps Reader", func() { gps.Gpsreader(state, &stateMu, *serialportPtr) })
	}

	if *enableGpioPtr {
		go safe.Do("switches", func() { switches.Startswitches() })
	}
	// Setup all screen elements, using defer to tear down when main returns.
	quit := ui.Init()
	defer ui.Done()

	// Render forever (simulator is already running, no further logic necessary)
	for {
		// look for quit
		ui.ProcessEvents()

		// render a frame
		stateMu.Lock()
		delay := ui.Render(state)
		stateMu.Unlock()

		// sleep outside render to minimize mutex lock time
		time.Sleep(delay)

		// non-blocking check of the quit channel
		select {
		case <-quit:
			// exit program when the channel is closed
			// (closed channels always return, open but empty channels block)
			return
		default:
			// no signal yet, move on
		}
	}
}
