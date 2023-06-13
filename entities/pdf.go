package entity

// FileType represents the type of file.
type FileType string

// Constants for different file types.
const (
	PDF      FileType = "pdf"
	Image    FileType = "image"
	Document FileType = "document"
)

// OptionMetadataPDF represents options for modifying PDF metadata.
type OptionMetadataPDF struct {
	Title    string
	Author   string
	Subject  string
	Keywords string
}

// OptionFilePDF represents options for working with PDF files.
type OptionFilePDF struct {
	PasswordPDF   string
	QRCodePath    string
	StampPosition string
}
