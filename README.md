# 简介
基于go、ffmpeg的音乐服务器
## 构建
### 开发
```sh
git clone https://github.com/nicangtianws/sena-no-go-music-server senanomusic
cd senanomusic
make dev
```

### 部署
#### 打包
```sh
make build
```
#### 运行
```sh
./target/senanomusic -env=./target/.env
```

### 使用docker镜像部署
#### go编译基础镜像
基于alpine:3.21
apk源替换为中科大源
内置基础编译环境
内置taglib
[builder](docker/builder-Dockerfile)
```sh
docker build -f script/builder-Dockerfile -t golang-alpine-custom .
```

#### 服务运行基础镜像
基于alpine:3.21

#### 打包镜像
```sh
docker build -t senanomusic .
```

#### 运行
```sh
docker run -d --name senanomusic \
-p 8000:8000 \
-v /data/senanomusic/music/:/data/senanomusic/music/ \
senanomusic
```

## ffmpeg
- 由于格式转换依赖于ffmpeg，所以编译一个仅支持简单图像处理和音频处理的ffmpeg
- 版本：7.1.1
### 依赖
```sh
sudo apt install pkg-config
```
- [ffmpeg依赖编译脚本](script/compile-ffmpeg-deps.sh)
```sh
sh script/compile-ffmpeg-deps.sh
```

### 配置
```sh
./configure \
  --prefix=/opt/ffmpeg \
  --extra-cflags="-I/data/senanomusic/ffmpeg-deps/build/include" \
  --extra-ldflags="-L/data/senanomusic/ffmpeg-deps/build/lib" \
  --pkg-config-flags="--static" \
  --enable-static \
  --disable-shared \
  --enable-pic \
  --disable-all \
  --enable-ffmpeg \
  --enable-avcodec \
  --enable-avformat \
  --enable-avutil \
  --enable-avfilter \
  --enable-gpl \
  --enable-nonfree \
  --enable-libmp3lame \
  --enable-libvorbis \
  --enable-libfdk_aac \
  --enable-zlib \
  --enable-swscale \
  --enable-swresample \
  --enable-encoder=libmp3lame,flac,libvorbis,wavpack,pcm_s16le,libfdk_aac,png,zlib,jpeg2000,mjpeg \
  --enable-decoder=libmp3lame,flac,libvorbis,wavpack,pcm_s16le,libfdk_aac,png,zlib,jpeg2000,mjpeg \
  --enable-demuxer=mp3,wav,flac,ogg,image2,mjpeg \
  --enable-muxer=mp3,wav,flac,ogg,image2,mjpeg \
  --enable-parser=aac,flac,vorbis,mjpeg,jpeg2000 \
  --enable-protocol=file \
  --enable-filter=scale,format,aresample \
  --disable-network \
  --disable-iconv \
  --disable-bzlib \
  --disable-lzma \
  --disable-vaapi \
  --disable-vdpau \
  --disable-ffplay \
  --disable-ffprobe
```
- 查看支持的编码器，decoder到filter同理
```sh
./configure --list-encoders
```

### 构建
```sh
make clean
make -j$(nproc)
sudo make install
```

- 清除旧的配置
```sh
make distclean
```