package gps
import(
//	"fmt"
	"strings"
	"strconv"
	"math"
)

var lvalid bool
var llat float64
var llon float64

func dmmtof(dmm,hemisphere string) (coord float64) {
	//0123456789
	//dddmm.mmmmm
	//ddmm.mmmmm
	//dmm.mmmmm
	pointindex := strings.Index(dmm,".")
	//fmt.Println(dmm,hemisphere,pointindex)
	coordd,_ := strconv.ParseFloat(dmm[0:pointindex-2],64)
	coordm,_ := strconv.ParseFloat(dmm[pointindex-2:],64)
	coord = coordd + coordm / 60.0;
	
	if hemisphere == "S" || hemisphere == "W" {
		coord *= -1.0
	}
	return

}

//time (hhmmss) distance (km) speed (km/h) heading (*) active
func Parsenmea(sentence string) (time int,distance,speed,heading float64, active bool) {

	if strings.Count(sentence,",") < 7 { return }

	thisstrings := strings.Split(sentence,",")
	if thisstrings[0] == "$GPRMC" && thisstrings[2] == "A" {
		timestring := strings.Split(thisstrings[1],".")
		t,_ := strconv.ParseInt(timestring[0],10,32)
		time = int(t)
		//fmt.Println(time)
	
		//work with lat/lon in radians, not degrees
		tlat := dmmtof(thisstrings[3],thisstrings[4]) * math.Pi / 180.0
		tlon := dmmtof(thisstrings[5],thisstrings[6]) * math.Pi / 180.0
		//fmt.Println(tlat*180/math.Pi,tlon*180/math.Pi)


		if(lvalid){
			dlat := tlat-llat;
			
			dlon := tlon-llon;

			//handle crossing +/- 180* longitude
			if dlon > math.Pi { 
				dlon -= 2.0 * math.Pi 
			}
			if dlon < -math.Pi { 
				dlon += 2.0 * math.Pi 
			}

			oh1 := math.Cos(tlat)
			oh2 := math.Cos(2.0 * tlat)
			dx := dlon * oh1 * 6383.616723
			dy := dlat * (6367.399725 - 32.432276 * oh2)
			distance = math.Sqrt(dx*dx + dy*dy)
		} else {
			distance = 0.0
		}
		llat = tlat
		llon = tlon
		lvalid = true

		speed,_  =  strconv.ParseFloat(thisstrings[7],64)
		speed *= 1.852  //convert from knots to km/h
		heading,_ = strconv.ParseFloat(thisstrings[8],64)	

		active = true
		
	} else {
		distance = 0
		speed = 0
		heading = 0
		active = false
	}
	return 
}
