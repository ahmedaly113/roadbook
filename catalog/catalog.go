// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package catalog

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p := messageKeyToIndex[key]
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"de": &dictionary{index: deIndex, data: deData},
		"en": &dictionary{index: enIndex, data: enData},
		"es": &dictionary{index: esIndex, data: esData},
		"fr": &dictionary{index: frIndex, data: frData},
		"pt": &dictionary{index: ptIndex, data: ptData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Decimal Point": 3,
	"Display Prefs": 1,
	"Exit":          4,
	"Go Back":       5,
	"Language":      0,
	"Prefs1":        6,
	"Prefs2":        7,
	"Prefs3":        8,
	"Routes":        2,
}

var deIndex = []uint32{ // 10 elements
	0x00000000, 0x00000008, 0x0000001f, 0x00000026,
	0x00000034, 0x0000003c, 0x00000048, 0x00000057,
	0x00000066, 0x00000075,
} // Size: 64 bytes

const deData string = "" + // Size: 117 bytes
	"\x02Sprache\x02Einstellungen anzeigen\x02Routen\x02Decimal Point\x02Been" +
	"den\x02Geh zurück\x02Einstellungen1\x02Einstellungen2\x02Einstellungen3"

var enIndex = []uint32{ // 10 elements
	0x00000000, 0x00000009, 0x00000017, 0x0000001e,
	0x0000002c, 0x00000031, 0x00000039, 0x00000040,
	0x00000047, 0x0000004e,
} // Size: 64 bytes

const enData string = "" + // Size: 78 bytes
	"\x02Language\x02Display Prefs\x02Routes\x02Decimal Point\x02Exit\x02Go B" +
	"ack\x02Prefs1\x02Prefs2\x02Prefs3"

var esIndex = []uint32{ // 10 elements
	0x00000000, 0x00000007, 0x00000020, 0x00000026,
	0x00000034, 0x0000003b, 0x00000043, 0x00000051,
	0x0000005f, 0x0000006d,
} // Size: 64 bytes

const esData string = "" + // Size: 109 bytes
	"\x02Idioma\x02Preferencias de pantalla\x02Rutas\x02Decimal Point\x02Sali" +
	"da\x02Regresa\x02Preferencias1\x02Preferencias2\x02Preferencias3"

var frIndex = []uint32{ // 10 elements
	0x00000000, 0x0000000a, 0x00000024, 0x00000031,
	0x0000003f, 0x00000046, 0x00000050, 0x0000005e,
	0x0000006c, 0x0000007a,
} // Size: 64 bytes

const frData string = "" + // Size: 122 bytes
	"\x02La langue\x02Préférences d'affichage\x02Itinéraires\x02Decimal Point" +
	"\x02Sortie\x02Retourner\x02Préférence1\x02Préférence2\x02Préférence3"

var ptIndex = []uint32{ // 10 elements
	0x00000000, 0x00000008, 0x00000024, 0x0000002a,
	0x00000038, 0x0000003f, 0x00000045, 0x00000053,
	0x00000061, 0x0000006f,
} // Size: 64 bytes

const ptData string = "" + // Size: 111 bytes
	"\x02Língua\x02Preferências de exibição\x02Rotas\x02Decimal Point\x02Saíd" +
	"a\x02Volte\x02Preferência1\x02Preferência2\x02Preferência3"

	// Total table size 857 bytes (0KiB); checksum: 760D7A1B
