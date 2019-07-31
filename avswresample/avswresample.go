package avswresample

// #include "libswresample/swresample.h"
// #include "libavutil/avutil.h"
// #cgo pkg-config: libswresample libavutil
import "C"
import (
	"errors"
	"github.com/baohavan/go-libav/avcodec"
	"github.com/baohavan/go-libav/avutil"
	"unsafe"
)

type SwrContext struct {
	CAVSwrContext uintptr
}

func NewSwrContext(inputCtx *avcodec.Context, outputCtx *avcodec.Context) (*SwrContext, error) {
	swrCtxOut := SwrContext{}
	outputChannels, _ := avutil.FindDefaultChannelLayout(outputCtx.Channels())
	inputChannels, _ := avutil.FindDefaultChannelLayout(inputCtx.Channels())
	swrCtxOut.CAVSwrContext = uintptr(unsafe.Pointer(C.swr_alloc_set_opts((*C.SwrContext)(C.NULL),
		(C.int64_t)(outputChannels),
		(C.enum_AVSampleFormat)(outputCtx.SampleFormat()), (C.int)(outputCtx.SampleRate()),
		(C.int64_t)(inputChannels),
		(C.enum_AVSampleFormat)(inputCtx.SampleFormat()), (C.int)(inputCtx.SampleRate()),
		0, C.NULL)))

	if swrCtxOut.CAVSwrContext == 0 {
		return nil, errors.New("Could not allocate swresample context\n")
	}

	return &swrCtxOut, nil
}

func (swr *SwrContext) Init() error {
	if (int)(C.swr_init((*C.SwrContext)(unsafe.Pointer(swr.CAVSwrContext)))) < 0 {
		return errors.New("Could not init swresample context\n")
	}
	return nil
}

func (swr *SwrContext) Free() {
	C.swr_free((**C.SwrContext)(unsafe.Pointer(&swr.CAVSwrContext)))
}

func (swr *SwrContext) SwrConvert(frame *avutil.Frame) error {
	frame, err := avutil.NewFrame()
	if err != nil {
		return err
	}
	defer frame.Free()

	errCode := C.swr_convert((*C.SwrContext)(unsafe.Pointer(swr.CAVSwrContext)), (**C.uchar)(frame.ExtendedData()), (C.int)(frame.NumberOfSamples()),
		(**C.uint8_t)(frame.ExtendedData()), (C.int)(frame.NumberOfSamples()))
	if (int)(errCode) < 0 {
		return avutil.NewErrorFromCode((avutil.ErrorCode)(errCode))
	}
	return nil
}

func SampleAlloc(outputCtx *avcodec.Context, frame *avutil.Frame) uintptr {
	var result uintptr

	var ptr uintptr
	ptr = uintptr(unsafe.Pointer(C.calloc((C.uint)(outputCtx.Channels()), C.sizeof_uint8_t)))
	if ptr == 0 {
		return 0
	}

	if (int)(C.av_samples_alloc((**C.uint8_t)(unsafe.Pointer(ptr)), (*C.int)(C.NULL),
		(C.int)(outputCtx.Channels()),
		(C.int)(frame.NumberOfSamples()),
		(C.enum_AVSampleFormat)(outputCtx.SampleFormat()), 0)) == 0 {
		C.av_freep(unsafe.Pointer(ptr))
		C.free(unsafe.Pointer(ptr))
		return 0
	}
	result = (uintptr)(unsafe.Pointer(ptr))
	return result
}

func FreeSample(ptr uintptr) {
	C.av_freep(unsafe.Pointer(ptr))
	C.free(unsafe.Pointer(ptr))
}
