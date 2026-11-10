package winres

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"testing"
	"unicode/utf16"
)

func testPNG(t *testing.T, size int) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	img.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 255}) // top-left red, rest transparent
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// object is a parsed .syso: its resource section and relocation offsets.
type object struct {
	section []byte
	relocs  []uint32
}

func parseObject(t *testing.T, b []byte) object {
	t.Helper()
	if m := binary.LittleEndian.Uint16(b); m != machineAMD64 {
		t.Fatalf("machine = %#x", m)
	}
	sh := b[20:]
	if string(sh[:5]) != ".rsrc" {
		t.Fatalf("section name %q", sh[:8])
	}
	size := binary.LittleEndian.Uint32(sh[16:])
	ptr := binary.LittleEndian.Uint32(sh[20:])
	rptr := binary.LittleEndian.Uint32(sh[24:])
	n := binary.LittleEndian.Uint16(sh[32:])
	o := object{section: b[ptr : ptr+size]}
	for i := range uint32(n) {
		r := b[rptr+10*i:]
		if typ := binary.LittleEndian.Uint16(r[8:]); typ != 3 {
			t.Fatalf("relocation type %d", typ)
		}
		o.relocs = append(o.relocs, binary.LittleEndian.Uint32(r))
	}
	return o
}

// get returns the data of resource typ/id (first language).
func (o object) get(t *testing.T, typ, id uint32) []byte {
	t.Helper()
	d := o.section
	u32 := func(off uint32) uint32 { return binary.LittleEndian.Uint32(d[off:]) }
	find := func(dir, key uint32) uint32 {
		n := uint32(binary.LittleEndian.Uint16(d[dir+14:]))
		for i := range n {
			e := dir + 16 + 8*i
			if key == 0 || u32(e) == key {
				return u32(e+4) &^ 0x80000000
			}
		}
		t.Fatalf("resource %d/%d not found", typ, id)
		return 0
	}
	entry := find(find(find(0, typ), id), 0)
	off, size := u32(entry), u32(entry+4)
	return d[off : off+size]
}

func TestBuild(t *testing.T) {
	syso, err := Build(Resources{
		Icons:           [][]byte{testPNG(t, 16), testPNG(t, 256), testPNG(t, 32)},
		Version:         "v2.1.3-4-gabc",
		ProductName:     "DLC Desktop App",
		FileDescription: "DLC Desktop App",
		OriginalName:    "dlc-desktop-app.exe",
	}, "amd64")
	if err != nil {
		t.Fatal(err)
	}
	o := parseObject(t, syso)
	if len(o.relocs) != 3+3 { // 3 icons, group, version, manifest
		t.Errorf("%d relocations, want 6", len(o.relocs))
	}

	grp := o.get(t, rtGroupIcon, 1)
	if n := binary.LittleEndian.Uint16(grp[4:]); n != 3 {
		t.Fatalf("group has %d icons", n)
	}
	wantSizes := []byte{0, 32, 16} // largest first; 0 means 256
	for i, want := range wantSizes {
		e := grp[6+14*i:]
		id := binary.LittleEndian.Uint16(e[12:])
		img := o.get(t, rtIcon, uint32(id))
		if e[0] != want || binary.LittleEndian.Uint32(e[8:]) != uint32(len(img)) {
			t.Errorf("entry %d: size %d, bytes %d (resource has %d)", i, e[0], binary.LittleEndian.Uint32(e[8:]), len(img))
		}
		isPNG := bytes.HasPrefix(img, []byte("\x89PNG"))
		if isPNG != (want == 0) {
			t.Errorf("entry %d: PNG = %v", i, isPNG)
		}
	}

	// 16x16 DIB: header height is doubled; the top-left red pixel is in the
	// last pixel row (bottom-up) and opaque, so its mask bit is clear.
	dibData := o.get(t, rtIcon, 3)
	if h := int32(binary.LittleEndian.Uint32(dibData[8:])); h != 32 {
		t.Errorf("DIB height %d, want 32", h)
	}
	px := dibData[40+15*16*4:]
	if !bytes.Equal(px[:4], []byte{0, 0, 255, 255}) {
		t.Errorf("top-left pixel BGRA = %v", px[:4])
	}
	mask := dibData[40+16*16*4:]
	if mask[15*4]&0x80 != 0 || mask[15*4]&0x40 == 0 {
		t.Errorf("mask row for top line = %08b", mask[15*4])
	}

	ver := o.get(t, rtVersion, 1)
	if l := binary.LittleEndian.Uint16(ver); int(l) != len(ver) {
		t.Errorf("VS_VERSIONINFO length %d, resource %d", l, len(ver))
	}
	fixed := ver[40:] // 6-byte header + "VS_VERSION_INFO\0" (32) + 2 padding
	if sig := binary.LittleEndian.Uint32(fixed); sig != 0xFEEF04BD {
		t.Fatalf("fixed info signature %#x", sig)
	}
	if ms, ls := binary.LittleEndian.Uint32(fixed[8:]), binary.LittleEndian.Uint32(fixed[12:]); ms != 2<<16|1 || ls != 3<<16 {
		t.Errorf("file version %#x %#x", ms, ls)
	}
	for _, s := range []string{"ProductVersion\x00v2.1.3-4-gabc", "FileVersion\x00\x002.1.3.0", "OriginalFilename\x00dlc-desktop-app.exe"} {
		if !bytes.Contains(ver, utf16le(s)) {
			t.Errorf("version info lacks %q", s)
		}
	}

	if m := o.get(t, rtManifest, 1); !bytes.Contains(m, []byte("asInvoker")) {
		t.Error("manifest missing")
	}
}

func utf16le(s string) []byte {
	var b []byte
	for _, u := range utf16.Encode([]rune(s)) {
		b = binary.LittleEndian.AppendUint16(b, u)
	}
	return b
}

func TestBuildErrors(t *testing.T) {
	cases := map[string]Resources{
		"no icons":   {},
		"not square": {Icons: [][]byte{pngOf(t, 16, 8)}},
		"too big":    {Icons: [][]byte{pngOf(t, 512, 512)}},
		"not png":    {Icons: [][]byte{[]byte("nope")}},
	}
	for name, r := range cases {
		if _, err := Build(r, "amd64"); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, err := Build(Resources{Icons: [][]byte{testPNG(t, 16)}}, "386"); err == nil {
		t.Error("386: no error")
	}
}

func pngOf(t *testing.T, w, h int) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewNRGBA(image.Rect(0, 0, w, h))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestParseVersion(t *testing.T) {
	for in, want := range map[string][4]uint16{
		"v2.0.0":              {2, 0, 0, 0},
		"2.1":                 {2, 1, 0, 0},
		"v1.2.3.4":            {1, 2, 3, 4},
		"v2.0.1-3-gabc-dirty": {2, 0, 1, 0},
		"fd9ce4f":             {0, 0, 0, 0},
		"4e2a1b0":             {0, 0, 0, 0}, // a commit hash, not a version
		"dev":                 {0, 0, 0, 0},
		"v99999.1":            {0, 1, 0, 0},
	} {
		if got := ParseVersion(in); got != want {
			t.Errorf("ParseVersion(%q) = %v, want %v", in, got, want)
		}
	}
}
