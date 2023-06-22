package http

type fields struct {
	Image    string
	Pdf      string
	Markdown string
	Svg      string
}

var Fields = fields{
	Image:    "image",
	Pdf:      "pdf",
	Markdown: "markdown",
	Svg:      "svg",
}
