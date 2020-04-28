// Package internal provides XML bindings for mapping model structs to/from XML. Making the package
// internal allows the default XML binding behavior (which requires Exported symbols) without polluting
// namespace with structs that aren't intended to be used elsewhere.
//
// Go note: "internal" is a magic name: https://docs.google.com/document/d/1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU/edit
package internal

import (
	"encoding/xml"
	"errors"
	"strings"
)

// <gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="OpenRally.org" xmlns:openrally="http://www.openrally.org/xmlschemas/GpxExtensions/v0.2-DRAFT">
// <wpt lat='0' lon='0'>
// <extensions>
// <openrally:distance>0.00</openrally:distance>
// <openrally:tulip>svg tulip</openrally:tulip>
// <openrally:notes>svg note</openrally:notes>
// </extensions>
// </wpt>

type WPT struct {
	Latitude  string `xml:"lat,attr"`
	Longitude string `xml:"lon,attr"`
	Distance  string `xml:"extensions>distance"`
	Tulip     string `xml:"extensions>tulip"`
	Notes     string `xml:"extensions>notes"`
	DZ		  *struct{} `xml:"extensions>dz"`
	FZ		  *struct{} `xml:"extensions>fz"`
}

type GPX struct {
	XMLName   xml.Name `xml:"gpx"`
	Creator   string   `xml:"creator,attr"`
	Waypoints []WPT    `xml:"wpt"`
}

func CheatBase64(svg string) (string, error) {
	anchor := "data:image/png;base64,"
	i := strings.Index(svg, anchor)
	if i <= 0 {
		return "", errors.New("base64 PNG not found in SVG")
	}
	right := svg[i+len(anchor):]

	i = strings.Index(right, "'")
	if i <= 0 {
		return "", errors.New("unterminated base64 string in SVG")
	}

	return right[:i], nil
}
