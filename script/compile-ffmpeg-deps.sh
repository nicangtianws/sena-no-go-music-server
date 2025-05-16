# 创建目录
mkdir -p /data/senanomusic/ffmpeg-deps/build && cd /data/senanomusic/ffmpeg-deps

# 依赖下载
# zlib必须使用该版本
wget -O zlib-1.2.11.tar.gz https://github.com/madler/zlib/archive/refs/tags/v1.2.11.tar.gz
wget https://www.ijg.org/files/jpegsrc.v9e.tar.gz
wget https://sourceforge.net/projects/libpng/files/libpng16/1.6.43/libpng-1.6.43.tar.xz
wget https://github.com/xiph/flac/releases/download/1.4.3/flac-1.4.3.tar.xz
wget https://downloads.xiph.org/releases/ogg/libogg-1.3.5.tar.xz
wget https://downloads.xiph.org/releases/vorbis/libvorbis-1.3.7.tar.xz
wget https://sourceforge.net/projects/lame/files/lame/3.100/lame-3.100.tar.gz
wget https://sourceforge.net/projects/opencore-amr/files/fdk-aac/fdk-aac-2.0.3.tar.gz

# 解压
tar -zxf zlib-1.2.11.tar.gz
tar -zxf jpegsrc.v9e.tar.gz
tar -zxf lame-3.100.tar.gz
tar -zxf fdk-aac-2.0.3.tar.gz

tar -Jxf libpng-1.6.43.tar.xz
tar -Jxf flac-1.4.3.tar.xz
tar -Jxf libogg-1.3.5.tar.xz
tar -Jxf libvorbis-1.3.7.tar.xz

# 编译 zlib
cd /data/senanomusic/ffmpeg-deps/zlib-1.2.11
./configure --static --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译 libjpeg
cd /data/senanomusic/ffmpeg-deps/jpeg-9e
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译 libpng (依赖 zlib)
cd /data/senanomusic/ffmpeg-deps/libpng-1.6.43
CFLAGS="-I/data/senanomusic/ffmpeg-deps/build/include" LDFLAGS="-L/data/senanomusic/ffmpeg-deps/build/lib" \
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译 libogg
cd /data/senanomusic/ffmpeg-deps/libogg-1.3.5
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译 libvorbis (依赖 libogg)
cd /data/senanomusic/ffmpeg-deps/libvorbis-1.3.7
CFLAGS="-I/data/senanomusic/ffmpeg-deps/build/include" LDFLAGS="-L/data/senanomusic/ffmpeg-deps/build/lib" \
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译 libflac
cd /data/senanomusic/ffmpeg-deps/flac-1.4.3
CFLAGS="-I/data/senanomusic/ffmpeg-deps/build/include" LDFLAGS="-L/data/senanomusic/ffmpeg-deps/build/lib" \
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译 libmp3lame
cd /data/senanomusic/ffmpeg-deps/lame-3.100
CFLAGS="-I/data/senanomusic/ffmpeg-deps/build/include" LDFLAGS="-L/data/senanomusic/ffmpeg-deps/build/lib" \
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install

# 编译fdk_acc
cd /data/senanomusic/ffmpeg-deps/fdk-aac-2.0.3
./configure --enable-static --disable-shared --prefix=/data/senanomusic/ffmpeg-deps/build
make -j$(nproc) && make install