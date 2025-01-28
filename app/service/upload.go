package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"shantaram-cms/pkg/util"
	"sync"
	"time"
)

type Upload struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewUploads(ctx context.Context) *Upload {
	ctx, cancel := context.WithCancel(ctx)

	return &Upload{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Upload) newFilePath(dir, ext string) string {
	for {
		id := uuid.New().String()
		name := fmt.Sprintf("%s.%s", id, ext)
		filePath := filepath.Join("data", dir, name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath
		}
	}
}

func (s *Upload) UploadTemp(header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open header: %w", err)
	}
	defer file.Close()

	originalExt := filepath.Ext(header.Filename)
	tempFilePath := s.newFilePath("temp", originalExt)

	output, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer output.Close()

	if _, err = io.Copy(output, file); err != nil {
		_ = os.Remove(tempFilePath)

		return "", fmt.Errorf("failed to copy input to temp file: %w", err)
	}

	s.scheduleFileDeletion(tempFilePath)

	return tempFilePath, nil
}

func (s *Upload) SaveToUploads(filePath string) (string, error) {
	originalExt := filepath.Ext(filePath)
	uploadFilePath := s.newFilePath("uploads", originalExt)

	if err := util.CopyFile(filePath, uploadFilePath); err != nil {
		return "", fmt.Errorf("failed to copy file to uploads: %w", err)
	}

	return uploadFilePath, nil
}

func (s *Upload) scheduleFileDeletion(filePath string) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		timer := time.NewTimer(10 * time.Minute)
		defer timer.Stop()

		select {
		case <-timer.C:
		case <-s.ctx.Done():
		}

		_ = os.Remove(filePath)
	}()
}

func (s *Upload) CancelAndJoin() {
	s.cancel()
	s.wg.Wait()
}
