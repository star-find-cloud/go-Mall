package domain

type Image struct {
	ImageID      int64  `db:"imageID"`
	SHA256Hash   string `db:"SHA256Hash"`
	IsCompressed bool   `db:"isCompressed"` // 是否压缩
	FilePath     string `db:"file_path"`
	ContentType  string `db:"content_type"`
}

//// IsTooBig 判断图片是否过大
//func (i Image) IsTooBig() error {
//	rawData := len(i.RawData)
//	if rawData > 1024*1024*10 {
//		return fmt.Errorf("图片过大")
//	}
//	return nil
//}
