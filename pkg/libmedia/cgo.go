package libmedia

/*
#cgo CFLAGS: -Wall -Wextra -I/opt/app/gweb/third_party/libmedia/include -I/opt/app/gweb/third_party/libmedia/deps/ffmpeg-7.1/include
#cgo CXXFLAGS: -std=c++20 -O2
#cgo LDFLAGS: -L/opt/app/gweb/build/libmedia -lmedia -L/opt/app/gweb/third_party/libmedia/deps/ffmpeg-7.1/lib -lavformat -lavcodec -lavutil -lswscale -lswresample
*/
import "C"
