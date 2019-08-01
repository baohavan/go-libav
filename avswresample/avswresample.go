package avswresample

// #include "libswresample/swresample.h"
// #include "libavutil/avutil.h"
// #define MAX_AUDIO_FRAME_SIZE 192000 // 1 second of 48khz 32bit audio
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

func (swr *SwrContext) SwrConvert(frame *avutil.Frame, out_buffer uintptr) error {
	errCode := C.swr_convert((*C.SwrContext)(unsafe.Pointer(swr.CAVSwrContext)), (**C.uchar)(unsafe.Pointer(out_buffer)),
		(C.int)(C.MAX_AUDIO_FRAME_SIZE),
		(**C.uint8_t)(frame.ExtendedData()), (C.int)(frame.NumberOfSamples()))
	if (int)(errCode) < 0 {
		return avutil.NewErrorFromCode((avutil.ErrorCode)(errCode))
	}
	return nil
}

func AllocBuffer(outputCtx *avcodec.Context, frame *avutil.Frame) uintptr {
	var ptr uintptr
	ptr = uintptr(unsafe.Pointer(C.av_malloc((C.uint)(outputCtx.Channels()) * (C.uint)(C.MAX_AUDIO_FRAME_SIZE))))
	if ptr == 0 {
		return 0
	}
	return ptr
}

func FreeBuffer(ptr uintptr) {
	C.av_free(unsafe.Pointer(ptr))
}
