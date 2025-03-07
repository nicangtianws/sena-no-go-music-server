# ffmpeg
- 编译仅支持图像处理和音频处理的ffmpeg
- 版本：4.4
- 依赖
```sh
sudo apt install libpng-dev libjpeg-dev libvorbis-dev libmp3lame-dev nasm pkg-config
```
- 配置
```sh
./configure \
  --prefix=/usr/local/ffmpeg-custom \
  --enable-static \
  --disable-shared \
  --disable-doc \
  --disable-avdevice \
  --disable-swscale \
  --disable-postproc \
  --disable-network \
  --disable-ffprobe \
  --disable-everything \
  --enable-gpl \
  --enable-nonfree \
  --enable-swresample \
  --enable-protocol=file \
  --enable-libmp3lame \
  --enable-libvorbis \
  --enable-zlib \
  --enable-encoder=aac,flac,libmp3lame,libvorbis,pcm_s16le,png,jpeg2000 \
  --enable-decoder=aac,flac,mp3,libvorbis,pcm_s16le,png,jpeg2000 \
  --enable-muxer=mp3,flac,wav,ogg,image2 \
  --enable-demuxer=mp3,flac,wav,ogg,image2 \
  --enable-parser=aac,flac,mpegaudio,vorbis,png,jpeg2000 \
  --enable-bsf=aac_adtstoasc \
  --enable-filter=aresample \
  --enable-small \
  --extra-cflags="-static" \
  --extra-ldflags="-static" \
  --pkg-config-flags="--static"
```
- 构建
```sh
make clean
make -j$(nproc)
sudo make install
```

- 清除旧的配置
```sh
make distclean
```