package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"profile-service/api/presenter"
	"profile-service/pkg/entities"
	"profile-service/pkg/userfile"

	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

// UploadUserFile is handler/controller which upload file
// @Summary      Upload user file
// @Description  Upload user file
// @Tags         Upload File
// @Accept       multipart/form-data
// @Produce      json
// @Security BearerAuth
// @Param        file  formData  file  true  "User File"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/upload-file [post]
func UploadUserFile(userfileService userfile.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil user_id dari JWT

		// 4. Ambil file baru dari form-data
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(presenter.ErrorResponse("file is required"))
		}
		//validate file
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fmt.Println(ext)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return c.Status(fiber.StatusBadRequest).SendString("Hanya file gambar yang diperbolehkan")
		}

		//validasi ukuran
		// Validasi ukuran file (max 100 KB)
		const maxSize = 100 * 1024 // 100 KB
		if file.Size > maxSize {
			return c.Status(fiber.StatusBadRequest).SendString("Ukuran file maksimal 100KB")
		}

		dddr, _ := RandomString(10)

		// bikin folder user -> uploads/<email>/
		dirPath := filepath.Join("uploads", dddr)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(presenter.ErrorResponse("failed to create dir: " + err.Error()))
		}

		savePath := filepath.Join(dirPath, file.Filename)

		// 5. Simpan file ke folder
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(presenter.ErrorResponse("failed to save file: " + err.Error()))
		}

		// 6. Simpan metadata user_file
		userFile := entities.File{
			FileUri:          savePath,
			FileThumbnailUri: savePath,
		}

		result, err := userfileService.UploadUserFile(&userFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(fiber.Map{
			"fileId":           result.ID,
			"fileUri":          result.FileUri,
			"fileThumbnailUri": result.FileThumbnailUri,
		})
	}
}
