// +build ffmpeg33

package avcodec

//#include <libavutil/avutil.h>
//#include <libavcodec/avcodec.h>
//
//
//static void go_avcodec_parameters_free(void *pParam) {
//	avcodec_parameters_free((AVCodecParameters**)(&pParam));
//}
// #cgo pkg-config: libavcodec libavutil
import "C"

import (
	"unsafe"

	"github.com/baohavan/go-libav/avutil"
)

type CodecParameters struct {
	//CAVCodecParameters *C.AVCodecParameters
	CAVCodecParameters uintptr
}

func NewCodecParameters() (*CodecParameters, error) {
	cPkt := uintptr(unsafe.Pointer(C.avcodec_parameters_alloc()))
	if cPkt == 0 {
		return nil, ErrAllocationError
	}
	return NewCodecParametersFromC(cPkt), nil
}

func NewCodecParametersFromC(cPSD uintptr) *CodecParameters {
	return &CodecParameters{CAVCodecParameters: cPSD}
}

func (cParams *CodecParameters) Free() {
	C.go_avcodec_parameters_free(unsafe.Pointer(cParams.CAVCodecParameters))
}

func (ctx *Context) CopyTo(dst *Context) error {
	// added in lavc 57.33.100
	parameters, err := NewCodecParameters()
	if err != nil {
		return err
	}
	defer parameters.Free()
	cParams := (*C.AVCodecParameters)(unsafe.Pointer(parameters.CAVCodecParameters))
	code := C.avcodec_parameters_from_context(cParams, ctx.CodeContext())
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	code = C.avcodec_parameters_to_context(dst.CodeContext(), cParams)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

func (ctx *Context) DecodeVideo(pkt *Packet, frame *avutil.Frame) (bool, int, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	cPkt := (*C.AVPacket)(unsafe.Pointer(&pkt.CAVPacket))
	C.avcodec_send_packet(ctx.CodeContext(), cPkt)
	code := C.avcodec_receive_frame(ctx.CodeContext(), cFrame)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		code = 0
	}
	return code == 0, int(code), err
}

func (ctx *Context) DecodeAudio(pkt *Packet, frame *avutil.Frame) (bool, int, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	cPkt := (*C.AVPacket)(unsafe.Pointer(&pkt.CAVPacket))
	C.avcodec_send_packet(ctx.CodeContext(), cPkt)
	code := C.avcodec_receive_frame(ctx.CodeContext(), cFrame)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		code = 0
	}
	return code == 0, int(code), err
}

func (ctx *Context) EncodeVideo(pkt *Packet, frame *avutil.Frame) (bool, error) {
	var cGotFrame C.int
	var cFrame *C.AVFrame
	if frame != nil {
		cFrame = (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	}
	cPkt := (*C.AVPacket)(unsafe.Pointer(&pkt.CAVPacket))
	code := C.avcodec_send_frame(ctx.CodeContext(), cFrame)
	C.avcodec_receive_packet(ctx.CodeContext(), cPkt)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return (cGotFrame != (C.int)(0)), err
}

func (ctx *Context) EncodeAudio(pkt *Packet, frame *avutil.Frame) (bool, error) {
	var cGotFrame C.int
	var cFrame *C.AVFrame
	if frame != nil {
		cFrame = (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	}
	cPkt := (*C.AVPacket)(unsafe.Pointer(&pkt.CAVPacket))
	code := C.avcodec_send_frame(ctx.CodeContext(), cFrame)
	C.avcodec_receive_packet(ctx.CodeContext(), cPkt)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return (cGotFrame != (C.int)(0)), err
}
