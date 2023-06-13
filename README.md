# simple-pdf-generator-with-qr

<!-- ABOUT THE PROJECT -->
## About The Project

Is a simple pdf file generator in the golang programming language.
programs that are supported from the simple program code above
1. create a qr-code (.png)
2. add a icon to the qr-code
3. create a pdf file from an image file
4. add metadata to pdf files (Tittle, Author, Description, Keywords)
5. modify pdf producer 
6. encrypt and decrypt pdf files 
7. add an image (qr-code) to the pdf file
8. convert docx to pdf using Ghostscript

### Built With

This section should list any major frameworks that you built your project using. Leave any add-ons/plugins for the acknowledgements section. Here are a few examples.
* [Golang](https://golang.com)

<!-- GETTING STARTED -->
## Getting Started
Before we get started, it's important to know that we run this code just from unit testing. So i didn't write a code to main.go
So, let start it.

You can start the first step from 2 different methods : 
1. Generate PDF file first, or
2. Generate QR Code first.

If you want to use method number 1, you just need to run unit testing in `serive/pdf/pdf_test.go` directory.
Similarly, if you want to use method number 2, you just need to run unit testing in `serive/qrcode` directory. 
But, we have something different in second method, you need to generetae QR Code first by running `generate_new_qr_code_test.go` first, and than `generate_add_icon_to_qr_code_test.go`

## Afterword
Hopefully, it can be easily understood and useful. Thank you~
