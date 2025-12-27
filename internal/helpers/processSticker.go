package helpers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/joho/godotenv"
)

type StickerMetadata struct {
	StickerPackID        string   `json:"sticker-pack-id"`
	StickerPackName      string   `json:"sticker-pack-name"`
	StickerPackPublisher string   `json:"sticker-pack-publisher"`
	Emojis               []string `json:"emojis"`
}

func ProcessSticker(inputData []byte, nameUser string) ([]byte, error) {
	typeFile := http.DetectContentType(inputData)

	switch {
	case typeFile == "image/webp":
		data, _ := convertImageToWebP(inputData)
		return addExifData(data, 123123, nameUser)
	case typeFile == "image/jpeg":
		data, _ := convertImageToWebP(inputData)
		return addExifData(data, 123123, nameUser)
	case typeFile == "image/png":
		data, _ := convertImageToWebP(inputData)
		return addExifData(data, 123123, nameUser)
	case typeFile == "image/gif":
		data, _ := convertGifToWebP(inputData)
		return addExifData(data, 123123, nameUser)
	case strings.HasPrefix(typeFile, "video/"):
		data, _ := convertVideoToGif(inputData)
		return addExifData(data, 123123, nameUser)
	default:
		return nil, fmt.Errorf("tipo de arquivo não suportado: %s", typeFile)
	}
}

func convertGifToWebP(input []byte) ([]byte, error) {
	args := []string{
		"-t", "5",
		"-vf", "scale=512:512:force_original_aspect_ratio=decrease,pad=512:512:(512-iw)/2:(512-ih)/2:color=0x00000000,fps=8",
		"-c:v", "libwebp",
		"-lossless", "0",
		"-compression_level", "6",
		"-q:v", "20",
		"-preset", "picture",
		"-loop", "0",
		"-an",
		"-vsync", "0",
	}
	return runFFmpeg(input, args, ".webp")
}

func convertImageToWebP(input []byte) ([]byte, error) {
	args := []string{
		"-y",
		"-vf", "scale=512x512",
		"-c:v", "libwebp",
		"-lossless", "1",
	}

	return runFFmpeg(input, args, ".webp")
}

func convertVideoToGif(input []byte) ([]byte, error) {
	args := []string{
		"-t", "5",
		"-vf", "fps=10,scale=512:512:flags=lanczos",
		"-f", "gif",
	}
	gifBytes, err := runFFmpeg(input, args, ".gif")
	if err != nil {
		return nil, err
	}
	return convertGifToWebP(gifBytes)
}

func runFFmpeg(data []byte, args []string, outputExt string) ([]byte, error) {
	tmpInput, err := os.CreateTemp("", "input-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	if _, err := tmpInput.Write(data); err != nil {
		return nil, err
	}
	tmpInput.Close()

	tmpOutputPath := filepath.Join(os.TempDir(), "output-"+filepath.Base(tmpInput.Name())+outputExt)
	defer os.Remove(tmpOutputPath)

	finalArgs := append([]string{"-y", "-i", tmpInput.Name()}, args...)
	finalArgs = append(finalArgs, tmpOutputPath)

	cmd := exec.Command("ffmpeg", finalArgs...)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return os.ReadFile(tmpOutputPath)
}

func addExifData(inputData []byte, updateId int64, nameUser string) ([]byte, error) {
	godotenv.Load()

	botName := os.Getenv("BOT_NAME")

	var (
		startingBytes = []byte{0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57, 0x07, 0x00}
		endingBytes   = []byte{0x16, 0x00, 0x00, 0x00}
		b             bytes.Buffer

		currUpdateId = strconv.FormatInt(updateId, 10)
		currPath     = path.Join("downloads", currUpdateId)
		inputPath    = path.Join(currPath, "input_exif.webm")
		outputPath   = path.Join(currPath, "output_exif.webp")
		exifDataPath = path.Join(currPath, "raw.exif")
	)

	_, err := b.Write(startingBytes)
	if err != nil {
		return nil, err
	}

	jsonData := map[string]interface{}{
		"sticker-pack-id":        "orvit.viitorags.com.github.",
		"sticker-pack-name":      nameUser,
		"sticker-pack-publisher": fmt.Sprintf("%s IA", botName),
		"emojis":                 []string{"😀"},
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	jsonLength := (uint32)(len(jsonBytes))

	lenBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuffer, jsonLength)

	if _, err := b.Write(lenBuffer); err != nil {
		return nil, err
	}
	if _, err := b.Write(endingBytes); err != nil {
		return nil, err
	}
	if _, err := b.Write(jsonBytes); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(currPath, os.ModePerm); err != nil {
		return nil, err
	}
	defer os.RemoveAll(currPath)

	if err := os.WriteFile(inputPath, inputData, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.WriteFile(exifDataPath, b.Bytes(), os.ModePerm); err != nil {
		return nil, err
	}

	if err := os.WriteFile(exifDataPath, b.Bytes(), os.ModePerm); err != nil {
		return nil, err
	}

	cmd := exec.Command("webpmux",
		"-set", "exif",
		exifDataPath, inputPath,
		"-o", outputPath,
	)

	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run webpmux command",
			fmt.Errorf("%v", err),
		)
		return nil, err
	}

	return os.ReadFile(outputPath)
}
