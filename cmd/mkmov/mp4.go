package main

import (
	"encoding/binary"
	"image"
	"os"

	"github.com/HexHacks/goffer/pkg/global"

	"github.com/mshafiee/mp4"
	"github.com/mshafiee/mp4/box"
)

const (
	MP4Colors  = 3 // rgb
	MP4PerComp = 4 // 4 bytes per component
)

var (
	// version = 0
	// creationtime = 0
	// modificationtine = 0
	// flags = [3]byte{0, 0, 0}
	MP4VideoSampleSize = uint32(global.Width * global.Height * MP4PerComp * MP4Colors)
	MP4VideoTimescale  = uint32(global.FPS) // units / sec
	MP4VideoDuration   = uint32(global.FrameCount)
	MP4VideoOpColor    = [3]uint16{4, 4, 4} // 32bit ?
	MP4VideoHdlr       = &box.HdlrBox{
		//Version     byte
		//Flags       [3]byte
		PreDefined: 0,
		// HandlerType can be: "vide" (video track), "soun" (audio track),
		// "hint" (hint track), "meta" (timed Metadata track),
		// "auxv" (auxiliary video track).
		HandlerType: "vide",
		Name:        "peyron_video",
	}
)

func mp4OutputFile(filename string, imgs []image.Image) {
	strct := mp4Struct(imgs)

	ofile, err := os.Create(filename)
	panicOn(err)

	defer ofile.Close()

	err = strct.Encode(ofile)
	panicOn(err)
}

func mp4Identity() []byte {
	zero := make([]byte, 3)
	one := make([]byte, 3)
	binary.BigEndian.PutUint32(one, 1)

	i := 0
	out := make([]byte, 36) // 3x3 matrix
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if x == y {
				out[i] = one[0]
				out[i+1] = one[1]
				out[i+2] = one[2]
			} else {
				out[i] = zero[0]
				out[i+1] = zero[1]
				out[i+2] = zero[2]
			}

			i = i + 3
		}
	}

	if i != 36 {
		panic("huston, we have that problem")
	}

	return out
}

func mp4Struct(imgs []image.Image) *mp4.MP4 {

	return &mp4.MP4{
		Ftyp: &box.FtypBox{ // filetype
			MajorBrand:       "peyron_out",
			MinorVersion:     []byte{0, 0, 1},
			CompatibleBrands: []string{},
		},
		Moov: &box.MoovBox{ // metadata
			Mvhd: &box.MvhdBox{
				// version, flags, creation and modification (skip)
				Timescale:   MP4VideoTimescale, // units / sec
				Duration:    MP4VideoDuration,
				NextTrackId: 0, // audio?
				Rate:        0, // audio?
				Volume:      0, // audio.
			},
			Iods: &box.IodsBox{}, // optional
			Trak: []*box.TrakBox{
				&box.TrakBox{ // video track
					Tkhd: &box.TkhdBox{
						// version, flags, creation and modification (skip)
						TrackId:        0,
						Layer:          0,
						AlternateGroup: 0,
						Volume:         0,
						Duration:       MP4VideoDuration,
						Matrix:         mp4Identity(),
						Width:          box.Fixed32(global.Width),
						Height:         box.Fixed32(global.Height),
					},
					Mdia: &box.MdiaBox{
						Mdhd: &box.MdhdBox{
							// version, flags, creation and modification (skip)
							Timescale: MP4VideoTimescale,
							Duration:  MP4VideoDuration,
							Language:  0,
						},
						Hdlr: MP4VideoHdlr, // 1
						Minf: &box.MinfBox{
							Vmhd: &box.VmhdBox{
								// version, flags
								GraphicsMode: 0,
								OpColor:      MP4VideoOpColor,
							},
							Smhd: &box.SmhdBox{}, // version, flags, balance
							Stbl: &box.StblBox{ // TODO: function for this?
								Stsd: &box.StsdBox{}, // version, flags
								Stts: &box.SttsBox{
									//the number of consecutive samples having the same duration
									SampleCount: []uint32{MP4VideoDuration},
									//duration in time units
									SampleTimeDelta: []uint32{MP4VideoTimescale},
								},
								Stss: &box.StssBox{
									SampleNumber: []uint32{},
								},
								Stsc: &box.StscBox{
									// all chunks starting at this index up to the next
									// first chunk have the same sample count/description
									FirstChunk: []uint32{0},
									// number of samples in the chunk
									SamplesPerChunk: []uint32{MP4VideoDuration},
									// description (see the sample description box - stsd)
									SampleDescriptionID: []uint32{0},
								},
								Stsz: &box.StszBox{
									SampleUniformSize: MP4VideoSampleSize,
									SampleNumber:      0,
									SampleSize:        []uint32{},
								},
								Stco: &box.StcoBox{
									ChunkOffset: mp4ChunkOffsets(len(imgs)),
								},
								Ctts: &box.CttsBox{}, // optional
							},
							Dinf: &box.DinfBox{
								Dref: &box.DrefBox{}, // flags, version
							},
							Hdlr: MP4VideoHdlr, // 2
						},
					},
					Edts: &box.EdtsBox{
						Elst: &box.ElstBox{
							SegmentDuration:   []uint32{}, // should be uint32 (version 0)
							MediaTime:         []uint32{}, // should be int32 (version 0)
							MediaRateInteger:  []uint16{},
							MediaRateFraction: []uint16{}, // should be int16
						},
					},
				},
			},
			Udta: &box.UdtaBox{ // user data
				Meta: &box.MetaBox{},
			},
		},
		Mdat: &box.MdatBox{}, // optional
	}
}

func mp4ChunkOffsets(count int) []uint32 {
	out := make([]uint32, count)
	for i := range out {
		out[i] = uint32(i) * MP4VideoSampleSize
	}
	return out
}
