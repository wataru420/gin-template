package models
import (
	"log"
	"database/sql"
)

type File struct {
	Id     		int
	AppId     	int
	VersionId 	int
	RevisionId	int
	Filename  	string
	Url		  	string
	Size	  	int
	Crc		  	uint32
	Tag       	sql.NullString
	Priority  	int
	State	  	int
}

type FileDao struct {
}

func (*FileDao) GetList(appId int, versionId int, revisionId int) ([]File, error) {
	var res = []File{}
	rows, err := dbs.Query(`SELECT id,revision_id,filename,url,size,crc,tag,priority,state
								FROM files where app_id=? and version_id=? and revision_id > ? `, appId, versionId, revisionId)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		file := File{AppId:appId,VersionId:versionId}
		if err := rows.Scan(&file.Id,&file.RevisionId,&file.Filename,&file.Url,&file.Size,&file.Crc,&file.Tag,&file.Priority,&file.State); err != nil {
			log.Fatal(err)
			return res, err
		}
		res = append(res, file)
	}

	return res, err

}

func (*FileDao) GetMaxRevisionId(appId int, versionId int) (int, error) {
	res := 0
	err := dbs.QueryRow(`SELECT max(revision_id) FROM files WHERE app_id=? and version_id=?`, appId, versionId).Scan(&res)
	return res, err
}
