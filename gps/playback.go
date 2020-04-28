package gps


import(
	"time"
	"sync"
	"github.com/ahmedaly113/roadbook/model"
	"os"
	"bufio"
//	"fmt"
)
func trap(err error) {
	if err != nil {
		panic(err)
	}
}

func Playback(state *model.Model,stateMu *sync.Mutex,filename string, repeat bool) {
	t := time.NewTicker(200 * time.Millisecond)
	f,err := os.Open(filename)
	trap(err)
	defer func() {
		err = f.Close()
		trap(err)
	}()

	s := bufio.NewScanner(f)

	for _ = range t.C {
		if s.Scan() {
			//fmt.Println(s.Text())
			_,di,sp,he,_ := Parsenmea(s.Text())
			stateMu.Lock()
			state.Heading = float32(he)
			state.Speed = float32(sp)
			state.Distance += float32(di)
			stateMu.Unlock()
			
		} else {
			trap(s.Err())
			//EOF returns false from Scan() but does not set an error
			if repeat {
				f.Seek(0,0)
				s = bufio.NewScanner(f)
			}
		}
	}

}

