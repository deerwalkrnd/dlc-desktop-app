// Package winres writes Windows resources (icon, version information and
// application manifest) as a COFF object file (.syso) that the Go linker
// embeds when it builds a Windows executable.
//
// It replaces third-party tools such as go-winres or rsrc, which keeps the
// project dependency-free. Only what the app needs is supported: one icon
// group, one version resource and one manifest, for amd64 or arm64.
package winres

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"regexp"
	"slices"
	"strconv"
	"unicode/utf16"
)

// Resource types and the language used for every resource.
const (
	rtIcon      = 3
	rtGroupIcon = 14
	rtVersion   = 16
	rtManifest  = 24

	langEnUS    = 0x0409
	codepageUTF = 0x04B0 // the "Unicode" code page used in version info
)

// Resources describes what to embed.
type Resources struct {
	// Icons are PNG files, one per size, each square and at most 256 pixels.
	Icons [][]byte
	// Version is a version string such as "v2.0.0" or "2.1.0-3-gabc123".
	// Its leading numbers become the numeric file and product version.
	Version string

	ProductName     string
	FileDescription string
	CompanyName     string
	Copyright       string
	OriginalName    string // e.g. "dlc-desktop-app.exe"
}

// Machine types for COFF headers.
const (
	machineAMD64 = 0x8664
	machineARM64 = 0xAA64
)

// Build returns the .syso contents for goarch ("amd64" or "arm64").
func Build(r Resources, goarch string) ([]byte, error) {
	var machine uint16
	var relocType uint16
	switch goarch {
	case "amd64":
		machine, relocType = machineAMD64, 3 // IMAGE_REL_AMD64_ADDR32NB
	case "arm64":
		machine, relocType = machineARM64, 2 // IMAGE_REL_ARM64_ADDR32NB
	default:
		return nil, fmt.Errorf("unsupported GOARCH %q", goarch)
	}

	var tree resourceTree
	group, err := addIcons(&tree, r.Icons)
	if err != nil {
		return nil, err
	}
	tree.add(rtGroupIcon, 1, group)
	tree.add(rtVersion, 1, versionInfo(r))
	tree.add(rtManifest, 1, []byte(manifest))

	section, relocs := tree.encode()
	return coff(machine, relocType, section, relocs), nil
}

// --- icons -----------------------------------------------------------------

// addIcons stores each image as an RT_ICON (PNG for 256 px, a 32-bit DIB for
// smaller sizes, which every part of Windows can draw) and returns the
// RT_GROUP_ICON directory that lists them.
func addIcons(tree *resourceTree, pngs [][]byte) ([]byte, error) {
	if len(pngs) == 0 {
		return nil, fmt.Errorf("no icon images")
	}
	type icon struct {
		size int
		data []byte
	}
	icons := make([]icon, 0, len(pngs))
	for i, p := range pngs {
		img, err := png.Decode(bytes.NewReader(p))
		if err != nil {
			return nil, fmt.Errorf("icon %d: %w", i, err)
		}
		b := img.Bounds()
		if b.Dx() != b.Dy() || b.Dx() < 1 || b.Dx() > 256 {
			return nil, fmt.Errorf("icon %d: must be square and at most 256 pixels, got %dx%d", i, b.Dx(), b.Dy())
		}
		data := p
		if b.Dx() < 256 {
			data = dib(img)
		}
		icons = append(icons, icon{b.Dx(), data})
	}
	// Largest first, as icon editors write them.
	slices.SortStableFunc(icons, func(a, b icon) int { return b.size - a.size })

	var g bytes.Buffer
	le(&g, uint16(0), uint16(1), uint16(len(icons))) // GRPICONDIR
	for i, ic := range icons {
		id := uint16(i + 1)
		tree.add(rtIcon, id, ic.data)
		dim := uint8(ic.size)
		if ic.size == 256 {
			dim = 0 // 0 means 256
		}
		// GRPICONDIRENTRY
		le(&g, dim, dim, uint8(0), uint8(0), uint16(1), uint16(32), uint32(len(ic.data)), id)
	}
	return g.Bytes(), nil
}

// dib encodes img as the bitmap form of an icon image: a BITMAPINFOHEADER
// with doubled height, bottom-up 32-bit BGRA pixels, then a 1-bit AND mask.
func dib(img image.Image) []byte {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	rgba := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	var buf bytes.Buffer
	maskStride := (w + 31) / 32 * 4
	le(&buf, uint32(40), int32(w), int32(2*h), uint16(1), uint16(32), uint32(0),
		uint32(w*h*4+maskStride*h), int32(0), int32(0), uint32(0), uint32(0))
	for y := h - 1; y >= 0; y-- {
		for x := range w {
			c := rgba.NRGBAAt(x, y)
			buf.Write([]byte{c.B, c.G, c.R, c.A})
		}
	}
	mask := make([]byte, maskStride)
	for y := h - 1; y >= 0; y-- {
		clear(mask)
		for x := range w {
			if rgba.NRGBAAt(x, y).A == 0 {
				mask[x/8] |= 0x80 >> (x % 8)
			}
		}
		buf.Write(mask)
	}
	return buf.Bytes()
}

// --- version information ---------------------------------------------------

var versionNumbers = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:$|[-+])`)

// ParseVersion returns the four numeric parts of a version string; anything
// that doesn't start with a number (such as a bare commit hash) is 0.0.0.0.
func ParseVersion(s string) [4]uint16 {
	var v [4]uint16
	m := versionNumbers.FindStringSubmatch(s)
	if m == nil {
		return v
	}
	for i, part := range m[1:] {
		if n, err := strconv.ParseUint(part, 10, 16); err == nil {
			v[i] = uint16(n)
		}
	}
	return v
}

// versionInfo builds a VS_VERSIONINFO resource.
func versionInfo(r Resources) []byte {
	v := ParseVersion(r.Version)
	ms := uint32(v[0])<<16 | uint32(v[1])
	ls := uint32(v[2])<<16 | uint32(v[3])

	var fixed bytes.Buffer
	le(&fixed,
		uint32(0xFEEF04BD), uint32(0x00010000), // signature, struct version
		ms, ls, ms, ls, // file version, product version
		uint32(0x3F), uint32(0), // flags mask, flags
		uint32(0x00040004),              // VOS_NT_WINDOWS32
		uint32(1),                       // VFT_APP
		uint32(0), uint32(0), uint32(0), // subtype, date
	)

	numeric := fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
	strings := []struct{ key, value string }{
		{"CompanyName", r.CompanyName},
		{"FileDescription", r.FileDescription},
		{"FileVersion", numeric},
		{"InternalName", r.ProductName},
		{"LegalCopyright", r.Copyright},
		{"OriginalFilename", r.OriginalName},
		{"ProductName", r.ProductName},
		{"ProductVersion", r.Version},
	}
	var table []block
	for _, s := range strings {
		if s.value != "" {
			table = append(table, textBlock(s.key, s.value))
		}
	}
	translation := make([]byte, 4)
	binary.LittleEndian.PutUint16(translation, langEnUS)
	binary.LittleEndian.PutUint16(translation[2:], codepageUTF)

	root := block{key: "VS_VERSION_INFO", value: fixed.Bytes(), valueLen: uint16(fixed.Len()), children: []block{
		{key: "StringFileInfo", text: true, children: []block{
			{key: fmt.Sprintf("%04X%04X", langEnUS, codepageUTF), text: true, children: table},
		}},
		{key: "VarFileInfo", text: true, children: []block{
			{key: "Translation", value: translation, valueLen: 4},
		}},
	}}
	return root.encode()
}

// block is one node of a version resource: wLength, wValueLength, wType, a
// NUL-terminated UTF-16 key, padding, the value, padding, then children.
type block struct {
	key      string
	value    []byte
	valueLen uint16 // in bytes for binary values, in UTF-16 units for text
	text     bool
	children []block
}

func textBlock(key, value string) block {
	v := utf16z(value)
	return block{key: key, value: v, valueLen: uint16(len(v) / 2), text: true}
}

func (b block) encode() []byte {
	var buf bytes.Buffer
	wType := uint16(0)
	if b.text {
		wType = 1
	}
	le(&buf, uint16(0), b.valueLen, wType) // length is patched below
	buf.Write(utf16z(b.key))
	pad4(&buf)
	buf.Write(b.value)
	for _, c := range b.children {
		pad4(&buf)
		buf.Write(c.encode())
	}
	out := buf.Bytes()
	binary.LittleEndian.PutUint16(out, uint16(len(out)))
	return out
}

// --- manifest --------------------------------------------------------------

// The manifest declares Windows 10/11 support (so Windows doesn't apply
// compatibility shims), runs without elevation, and opts in to long paths,
// since lesson file names are long and libraries can be nested deeply.
const manifest = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
    </application>
  </compatibility>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <longPathAware xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">true</longPathAware>
    </windowsSettings>
  </application>
</assembly>
`

// --- resource directory and COFF -------------------------------------------

type resource struct {
	typ, id uint16
	data    []byte
}

// resourceTree holds resources with integer types and IDs, all in one
// language, and lays them out as the three-level PE resource directory.
type resourceTree struct{ items []resource }

func (t *resourceTree) add(typ, id uint16, data []byte) {
	t.items = append(t.items, resource{typ, id, data})
}

// encode returns the .rsrc section and the offsets of the data-entry fields
// that need an image-relative (RVA) relocation.
func (t *resourceTree) encode() (section []byte, relocs []uint32) {
	items := slices.Clone(t.items)
	slices.SortFunc(items, func(a, b resource) int {
		if a.typ != b.typ {
			return int(a.typ) - int(b.typ)
		}
		return int(a.id) - int(b.id)
	})
	var types []uint16
	byType := map[uint16][]resource{}
	for _, it := range items {
		if len(byType[it.typ]) == 0 {
			types = append(types, it.typ)
		}
		byType[it.typ] = append(byType[it.typ], it)
	}

	const dirSize, entrySize, dataEntrySize = 16, 8, 16
	// Layout: root dir, one name-level dir per type, one language-level dir
	// per resource, then data entries, then the data itself.
	off := uint32(dirSize + entrySize*len(types))
	nameDir := map[uint16]uint32{}
	for _, typ := range types {
		nameDir[typ] = off
		off += uint32(dirSize + entrySize*len(byType[typ]))
	}
	langDir := make([]uint32, len(items))
	for i := range items {
		langDir[i] = off
		off += dirSize + entrySize
	}
	dataEntry := make([]uint32, len(items))
	for i := range items {
		dataEntry[i] = off
		off += dataEntrySize
	}
	dataOff := make([]uint32, len(items))
	for i, it := range items {
		off = (off + 7) &^ 7
		dataOff[i] = off
		off += uint32(len(it.data))
	}

	buf := make([]byte, (off+7)&^7)
	put16 := func(at uint32, v uint16) { binary.LittleEndian.PutUint16(buf[at:], v) }
	put32 := func(at uint32, v uint32) { binary.LittleEndian.PutUint32(buf[at:], v) }
	dir := func(at uint32, ids int) uint32 { put16(at+14, uint16(ids)); return at + dirSize }
	const subdir = 0x80000000

	e := dir(0, len(types))
	for _, typ := range types {
		put32(e, uint32(typ))
		put32(e+4, subdir|nameDir[typ])
		e += entrySize
	}
	i := 0
	for _, typ := range types {
		e := dir(nameDir[typ], len(byType[typ]))
		for _, it := range byType[typ] {
			put32(e, uint32(it.id))
			put32(e+4, subdir|langDir[i])
			e += entrySize
			lang := dir(langDir[i], 1)
			put32(lang, langEnUS)
			put32(lang+4, dataEntry[i])
			i++
		}
	}
	for i, it := range items {
		put32(dataEntry[i], dataOff[i]) // made image-relative by the relocation
		put32(dataEntry[i]+4, uint32(len(it.data)))
		copy(buf[dataOff[i]:], it.data)
		relocs = append(relocs, dataEntry[i])
	}
	return buf, relocs
}

// coff wraps the section in an object file with one section, the relocations
// and a single static symbol for the section that they refer to.
func coff(machine, relocType uint16, section []byte, relocs []uint32) []byte {
	const headerSize, sectionHeaderSize, relocSize = 20, 40, 10
	dataPtr := uint32(headerSize + sectionHeaderSize)
	relocPtr := dataPtr + uint32(len(section))
	symPtr := relocPtr + uint32(relocSize*len(relocs))

	var buf bytes.Buffer
	// IMAGE_FILE_HEADER
	le(&buf, machine, uint16(1), uint32(0), symPtr, uint32(1), uint16(0), uint16(0))
	// IMAGE_SECTION_HEADER
	buf.WriteString(".rsrc\x00\x00\x00")
	le(&buf, uint32(0), uint32(0), uint32(len(section)), dataPtr, relocPtr, uint32(0),
		uint16(len(relocs)), uint16(0),
		uint32(0x40000040)) // IMAGE_SCN_CNT_INITIALIZED_DATA | IMAGE_SCN_MEM_READ
	buf.Write(section)
	for _, r := range relocs {
		le(&buf, r, uint32(0), relocType) // offset, symbol index 0, type
	}
	// Symbol ".rsrc": value 0, section 1, type 0, storage class STATIC.
	buf.WriteString(".rsrc\x00\x00\x00")
	le(&buf, uint32(0), int16(1), uint16(0), uint8(3), uint8(0))
	le(&buf, uint32(4)) // empty string table
	return buf.Bytes()
}

// --- helpers ---------------------------------------------------------------

func le(buf *bytes.Buffer, values ...any) {
	for _, v := range values {
		_ = binary.Write(buf, binary.LittleEndian, v)
	}
}

func pad4(buf *bytes.Buffer) {
	for buf.Len()%4 != 0 {
		buf.WriteByte(0)
	}
}

func utf16z(s string) []byte {
	units := append(utf16.Encode([]rune(s)), 0)
	out := make([]byte, 2*len(units))
	for i, u := range units {
		binary.LittleEndian.PutUint16(out[2*i:], u)
	}
	return out
}
