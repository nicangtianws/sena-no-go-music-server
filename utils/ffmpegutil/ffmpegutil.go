package ffmpegutil

import (
	"gin-jwt/utils/mylog"
	"os/exec"
)

func ConvertTo44kOGG(inputPath, outputPath string) error {
	// 构造 FFmpeg 命令
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath, // 输入文件
		"-vn",          // 不处理封面
		"-ar", "44100", // 设置采样率为 44.1kHz
		"-b:a", "320k", // 设置音频码率为 320k
		"-c:a", "libvorbis", // 使用 libvorbis 编码器
		outputPath, // 输出文件
	)

	mylog.LOG.Info().Msg(cmd.String())

	// 执行命令
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ConvertCover(inputPath, outputPath string) error {
	// 构造 FFmpeg 命令
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath, // 输入文件
		"-an",
		"-vcodec", "copy",
		outputPath, // 输出文件
	)

	mylog.LOG.Info().Msg(cmd.String())

	// 执行命令
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
