package service

import (
	"fmt"
	"github.com/ricochet2200/go-disk-usage/du"
	"os"
	"path/filepath"
	"shantaram-cms/app/dao"
	"shantaram-cms/pkg/database"
)

type File struct {
	db            *database.Database
	uploadService *Upload
}

func NewFile(
	db *database.Database,
	uploadService *Upload,
) *File {
	return &File{
		db:            db,
		uploadService: uploadService,
	}
}

func (s *File) GetAll() ([]database.File, error) {
	pages, err := s.db.GetAllFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get all files: %v", err)
	}

	return pages, nil
}

func (s *File) Delete(id uint64) error {
	file, err := s.db.GetFileByID(id)
	if err != nil {
		return fmt.Errorf("failed to get file by id %d: %v", id, err)
	}

	filePath := filepath.Join("data", "uploads", file.Path)

	if err := s.db.DeleteFile(id); err != nil {
		return fmt.Errorf("failed to delete file by id %d: %v", id, err)
	}

	_ = os.Remove(filePath)

	return nil
}

func (s *File) Insert(tempPath string, info dao.NewFileRequest) error {
	stats, err := os.Stat(tempPath)
	if err != nil {
		return fmt.Errorf("failed to get file stats: %v", err)
	}

	size := stats.Size()

	uploadPath, err := s.uploadService.SaveToUploads(tempPath)
	if err != nil {
		return fmt.Errorf("failed to save file to uploads: %v", err)
	}

	path := filepath.Base(uploadPath)

	newFile := &database.File{
		Path:  path,
		Title: info.Title,
		Name:  info.Name,
		Size:  size,
	}

	if err := s.db.InsertFile(newFile); err != nil {
		return fmt.Errorf("failed to insert file: %v", err)
	}

	return nil
}

func (s *File) Stats() dao.FileStatsResponse {
	usage := du.NewDiskUsage(".")
	return dao.FileStatsResponse{
		Total: usage.Size(),
		Free:  usage.Free(),
	}
}
