package models

type FileChunkInFo struct {
	GormModel
	FileName    string `json:"fileName" binding:"required" gorm:"column:filename"`
	FileNameMd5 string `gorm:"column:fileNameMd5"`
	FileMd5     string `json:"FileMd5" binding:"omitempty" gorm:"column:fileMd5"`
	ChunkTotal  int    `gorm:"column:chunkTotal"`
	ChunkNext   int    `gorm:"column:chunkNext"`
<<<<<<< HEAD
	Finished    bool   `gorm:"column:Finished"`
	FileExt     string `gorm:"column:fileExt"`
=======
	Finished    bool   `gorm:"column:finished"`
	FileExt     string `gorm:"column:FileExt"`
>>>>>>> master
}

func (f *FileChunkInFo) TableName() string {
	return "BackPointInfo"
}
