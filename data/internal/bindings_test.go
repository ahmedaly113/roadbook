package internal

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

const example = `
<gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="OpenRally.org" xmlns:openrally="http://www.openrally.org/xmlschemas/GpxExtensions/v0.2-DRAFT">
	<wpt lat='1' lon='2'>
		<extensions>
			<openrally:distance>2.71</openrally:distance>
			<openrally:tulip>svg tulip</openrally:tulip>
			<openrally:notes>svg note</openrally:notes>
		</extensions>
	</wpt>
	<wpt lat='3' lon='4'>
		<extensions>
			<openrally:distance>3.14</openrally:distance>
			<openrally:tulip>svg tulip</openrally:tulip>
			<openrally:notes>svg note</openrally:notes>
		</extensions>
	</wpt>
</gpx>`

func TestRead(t *testing.T) {
	v := GPX{}
	err := xml.Unmarshal([]byte(example), &v)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(v.Waypoints), "expected 2 waypoints")
	assert.Equal(t, "1", v.Waypoints[0].Latitude, "first waypoint lat")
	assert.Equal(t, "2", v.Waypoints[0].Longitude, "first waypoint long")
	assert.Equal(t, "2.71", v.Waypoints[0].Distance, "first waypoint distance")
	assert.Equal(t, "svg tulip", v.Waypoints[0].Tulip, "first waypoint tulip")
	assert.Equal(t, "svg note", v.Waypoints[0].Notes, "first waypoint note")
	assert.Equal(t, "3", v.Waypoints[1].Latitude, "second waypoint lat")
	assert.Equal(t, "4", v.Waypoints[1].Longitude, "second waypoint long")
	assert.Equal(t, "3.14", v.Waypoints[1].Distance, "second waypoint distance")
	assert.Equal(t, "svg tulip", v.Waypoints[1].Tulip, "second waypoint tulip")
	assert.Equal(t, "svg note", v.Waypoints[1].Notes, "second waypoint note")
}

const svgExample = `<svg xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'><image height='396px' width='550px' xlink:href='data:image/png;base64,GOODSTUFFHERE'/></svg>`

func TestCheatBase64(t *testing.T) {
	s, _ := CheatBase64(svgExample)
	assert.Equal(t, "GOODSTUFFHERE", s)
}
