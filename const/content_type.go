package _const

const (
	JPG = 7001 + iota
	JPEG
	PNG
	WEBP
	GIF
	BMP
	TIFF
	HEIF
	DNG
	CR3
	NEF
	ATW
	RAF
)

var ContentTypeIntMap = map[string]int64{
	".jpg":  JPG,
	".jpeg": JPEG,
	".png":  PNG,
	".webp": WEBP,
	".gif":  GIF,
	".bmp":  BMP,
	".tiff": TIFF,
	".heif": HEIF,
	".dng":  DNG,
	".cr3":  CR3,
	".nef":  NEF,
	".atw":  ATW,
	".raf":  RAF,
}

var ContentTypeStringMap = map[int64]string{
	JPG:  ".jpg",
	JPEG: ".jpeg",
	PNG:  ".png",
	WEBP: ".webp",
	GIF:  ".gif",
	BMP:  ".bmp",
	TIFF: ".tiff",
	HEIF: ".heif",
	DNG:  ".dng",
	CR3:  ".cr3",
	NEF:  ".nef",
	ATW:  ".atw",
	RAF:  ".raf",
}
