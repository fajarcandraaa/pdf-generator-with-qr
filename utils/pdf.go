package util

import (
	"fmt"
	"os/exec"
	"path/filepath"
	entity "pdf-generator-with-qr/entities"
	"reflect"
	"strings"
)

// IsStructEmpty a function to check if a struct is empty or not.
func IsStructEmpty(data interface{}) bool {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() != "" {
			return false
		}
	}

	return true
}

// GetFileType returns the type of file based on its extension.
func GetFileType(filePath string) entity.FileType {
	extension := strings.ToLower(filepath.Ext(filePath))
	switch extension {
	case ".pdf":
		return entity.PDF
	case ".jpg", ".jpeg", ".png":
		return entity.Image
	case ".doc", ".docx":
		return entity.Document
	default:
		return ""
	}
}

// AddedMetadata to add metadata into a pdf file.
func AddedMetadata(filePath string, metadata *entity.OptionMetadataPDF) error {
	command := fmt.Sprintf("pdfcpu properties add %s 'Title = %s' 'Author = %s' 'Subject = %s'", filePath, metadata.Title, metadata.Author, metadata.Subject)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// AddKeywords to add keywords into a pdf file.
func AddKeywords(filePath string, metadata *entity.OptionMetadataPDF) error {
	errs := removeKeywords(filePath)
	if errs != nil {
		return errs
	}

	command := fmt.Sprintf("pdfcpu keywords add %s '%s'", filePath, metadata.Keywords)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// removeKeywords to remove keywords into a pdf file.
func removeKeywords(filePath string) error {
	command := fmt.Sprintf("pdfcpu key remove '%s'", filePath)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Encrypted function is used to encrypt a previously decrypted PDF.
func Encrypted(filePath string, password string) error {
	command := fmt.Sprintf("pdfcpu encrypt -upw %s -opw %s %s", password, password, filePath)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing pdfcpu command: %s\n", err.Error())
		return err
	}

	return nil
}

// Decrypted unction is used to remove the protection from a PDF file by decrypting it with a provided password.
func Decrypted(filePath string, password string) error {
	command := fmt.Sprintf("pdfcpu decrypt -upw %s %s", password, filePath)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing pdfcpu command: %s\n", err.Error())
		return err
	}

	return nil
}

// HasPDFPassword checks if the PDF file is password-protected.
func HasPDFPassword(filePath string, password string) (bool, error) {
	command := ""
	if password != "" {
		command = fmt.Sprintf("pdfcpu validate -mode=quiet -upw='%s' %s", password, filePath)
	} else {
		command = fmt.Sprintf("pdfcpu validate %s", filePath)
	}

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && exitError.ExitCode() == 1 {
			// PDF is password protected
			return true, nil
		} else {
			// Other execution error
			return false, exitError
		}
	} else {
		// PDF is not password protected
		return false, err
	}
}

// ChangeFileExtension changes the file extension to the new extension.
func ChangeFileExtension(filePath string, newExtension string) string {
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	newFileName := fileNameWithoutExt + "." + newExtension
	return filepath.Join(filepath.Dir(filePath), newFileName)
}
