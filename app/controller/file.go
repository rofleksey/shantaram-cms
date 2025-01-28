package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
	"shantaram-cms/pkg/database"
	"strconv"
)

type File struct {
	fileService   *service.File
	uploadService *service.Upload
}

func NewFile(
	fileService *service.File,
	uploadService *service.Upload,
) *File {
	return &File{
		fileService:   fileService,
		uploadService: uploadService,
	}
}

func (c *File) Upload(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse multipart form: %v", err),
		})
	}

	fileArr := form.File["file"]
	if len(fileArr) != 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   "invalid number of files",
		})
	}

	file := fileArr[0]

	tempPath, err := c.uploadService.UploadTemp(file)
	if err != nil {
		return fmt.Errorf("failed to upload temp file: %v", err)
	}

	fileInfoArr := form.Value["info"]
	if len(fileInfoArr) != 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   "invalid number of file data",
		})
	}

	fileInfo := fileInfoArr[0]

	var req dao.NewFileRequest

	if err := json.Unmarshal([]byte(fileInfo), &req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   "invalid file info",
		})
	}

	if err := c.fileService.Insert(tempPath, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось загрузить файл: %v", err),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *File) Delete(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("invalid id: %v", err),
		})
	}

	if err := c.fileService.Delete(uint64(id)); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to delete file: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *File) GetAll(ctx *fiber.Ctx) error {
	files, err := c.fileService.GetAll()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to get all files: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[[]database.File]{
		Data: files,
	})
}

func (c *File) Stats(ctx *fiber.Ctx) error {
	stats := c.fileService.Stats()

	return ctx.JSON(dao.SuccessResponse[dao.FileStatsResponse]{
		Data: stats,
	})
}
