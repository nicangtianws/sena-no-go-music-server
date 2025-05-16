package audiofileutil

import "errors"

type AudioFileType int

const (
	OGG AudioFileType = iota
	FLAC
	AAC
	WAV
	MP3
)

// 字符串转音频枚举
func GetAudioFileTypeByStr(typeStr string) (AudioFileType, error) {
	switch typeStr {
	case "OGG":
		return OGG, nil
	case "FLAC":
		return FLAC, nil
	case "AAC":
		return AAC, nil
	case "WAV":
		return WAV, nil
	case "MP3":
		return MP3, nil
	default:
		return OGG, errors.New("not supported type")
	}
}

// 根据mime类型获取音频类型
func GetAudioFileType(mimeType string) (AudioFileType, error) {
	switch mimeType {
	case "audio/ogg":
		return OGG, nil
	case "audio/flac":
		return FLAC, nil
	case "audio/aac":
		return AAC, nil
	case "audio/wav":
		return WAV, nil
	case "audio/x-wav":
		return WAV, nil
	case "audio/vnd.wave":
		return WAV, nil
	case "audio/wave":
		return WAV, nil
	case "audio/mpeg":
		return MP3, nil
	case "audio/x-mpeg":
		return MP3, nil
	case "audio/mp3":
		return MP3, nil
	default:
		return OGG, errors.New("not supported type")
	}
}

// 获取音频类型对应的mime类型及别称
func GetMimeTypeList(audioType AudioFileType) ([]string, error) {
	switch audioType {
	case OGG:
		return []string{"audio/ogg"}, nil
	case FLAC:
		return []string{"audio/flac"}, nil
	case AAC:
		return []string{"audio/aac"}, nil
	case WAV:
		return []string{"audio/wav", "audio/x-wav", "audio/vnd.wave", "audio/wave"}, nil
	case MP3:
		return []string{"audio/mpeg", "audio/x-mpeg", "audio/mp3"}, nil
	default:
		return []string{}, errors.New("not supported type")
	}
}

// 获取音频类型对应的mime类型
func GetMimeType(audioType AudioFileType) (string, error) {
	switch audioType {
	case OGG:
		return "audio/ogg", nil
	case FLAC:
		return "audio/flac", nil
	case AAC:
		return "audio/aac", nil
	case WAV:
		return "audio/wav", nil
	case MP3:
		return "audio/mpeg", nil
	default:
		return "", errors.New("not supported type")
	}
}
