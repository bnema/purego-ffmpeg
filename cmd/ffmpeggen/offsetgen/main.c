// cmd/ffmpeggen/offsetgen/main.c
// Offset generator for FFmpeg struct fields.
// Platform: Linux amd64 (offsets are architecture and version specific)
// Target: FFmpeg 7.x (libavcodec 62, libavformat 62, libavutil 60)
//
// Usage:
//   cc -o offsetgen main.c $(pkg-config --cflags --libs libavformat libavcodec libavutil)
//   ./offsetgen
//
// Update the values in cmd/ffmpeggen/internal/overrides/overrides.go
// whenever you target a different FFmpeg version or architecture.

#include <stdio.h>
#include <stddef.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>
#include <libavutil/hwcontext.h>

#define OFFSET(type, field) printf("  %s.%s = %zu\n", #type, #field, offsetof(type, field))

int main() {
    printf("// FFmpeg struct offsets for overrides.go\n");
    printf("// Compile: cc -o offsetgen main.c $(pkg-config --cflags --libs libavformat libavcodec libavutil)\n\n");

    printf("AVFormatContext:\n");
    OFFSET(AVFormatContext, nb_streams);
    OFFSET(AVFormatContext, streams);
    OFFSET(AVFormatContext, duration);
    OFFSET(AVFormatContext, bit_rate);

    printf("\nAVStream:\n");
    OFFSET(AVStream, index);
    OFFSET(AVStream, codecpar);
    OFFSET(AVStream, time_base);
    OFFSET(AVStream, duration);
    OFFSET(AVStream, nb_frames);

    printf("\nAVCodecContext:\n");
    OFFSET(AVCodecContext, codec_type);
    OFFSET(AVCodecContext, codec_id);
    OFFSET(AVCodecContext, time_base);
    OFFSET(AVCodecContext, width);
    OFFSET(AVCodecContext, height);
    OFFSET(AVCodecContext, pix_fmt);
    OFFSET(AVCodecContext, sample_rate);
    OFFSET(AVCodecContext, sample_fmt);
    OFFSET(AVCodecContext, hw_device_ctx);
    OFFSET(AVCodecContext, hw_frames_ctx);

    printf("\nAVFrame:\n");
    OFFSET(AVFrame, data);
    OFFSET(AVFrame, linesize);
    OFFSET(AVFrame, width);
    OFFSET(AVFrame, height);
    OFFSET(AVFrame, nb_samples);
    OFFSET(AVFrame, format);
    OFFSET(AVFrame, pts);
    OFFSET(AVFrame, pkt_dts);
    OFFSET(AVFrame, sample_rate);
    OFFSET(AVFrame, hw_frames_ctx);
    OFFSET(AVFrame, ch_layout);

    printf("\nAVCodecParameters:\n");
    OFFSET(AVCodecParameters, codec_type);
    OFFSET(AVCodecParameters, codec_id);
    OFFSET(AVCodecParameters, width);
    OFFSET(AVCodecParameters, height);
    OFFSET(AVCodecParameters, sample_rate);
    OFFSET(AVCodecParameters, format);

    return 0;
}
