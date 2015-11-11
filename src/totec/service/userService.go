package service
import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"ghe.ca-tools.org/hilo/octo-proto.git/go/octo"
	"octo/models"
)

var fileDao = &models.FileDao{}
var tagDao = &models.TagDao{}
var gcpDao = &models.GcpDao{}

func ListEndpoint(c *gin.Context)  {
	val, _ := c.Get("app")
	app, ok := val.(models.App)
	if !ok {
		c.String(http.StatusForbidden, "Invalid App")
		return
	}
	appId := app.AppId
	version := c.Param("version")
	v, err := strconv.Atoi(version)
	if err != nil {
		c.String(http.StatusBadRequest, version + " is not version")
		return
	}
	revision := c.Param("revision")
	log.Printf(version + " : " + revision)
	r, err := strconv.Atoi(revision)
	if err != nil {
		c.String(http.StatusBadRequest, revision + " is not revision")
		return
	}

	tagList, err := tagDao.GetList(appId)
	tagMap := make(map[string]int)
	tagNameList := []string{}
	for i, _ := range tagList {
		t := tagList[i]
		tagMap[t.Name] = t.TagId
		tagNameList = append(tagNameList, t.Name)
	}

	fileList, err := fileDao.GetList(appId,v,r)
	Database := new(octo.Database)
	for i, _ := range fileList {
		f := fileList[i]
		data := new(octo.Data)
		data.Id = proto.Int(f.Id)
		data.Filepath = proto.String(f.Url)
		data.Name = proto.String(f.Filename)
		data.Size = proto.Int(f.Size)
		data.Crc = proto.Uint32(f.Crc)
		data.Priority = proto.Int(f.Priority)

		tags := strings.Split(f.Tag.String, ",")
		tagIds := []int32{}
		for _, tag := range tags {
			tagIds = append(tagIds, int32(tagMap[tag]))
		}
		data.Tagid = tagIds
		data.State = getDataState(f.State)
		Database.List = append(Database.List, data)

		if r < f.RevisionId {
			r = f.RevisionId
		}
	}

	if len(fileList) == 0 {
		maxRevisionId, err := fileDao.GetMaxRevisionId(appId, v)
		if err != nil {
			log.Fatal("error")
			c.String(http.StatusInternalServerError, err.Error())
		}
		r = maxRevisionId
	}

	Database.Revision = proto.Int(r)
	Database.Tagname = tagNameList
	database, err := proto.Marshal(Database)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}
	log.Print(proto.MarshalTextString(Database))
	c.Data(http.StatusOK,"application/x-protobuf",database)
}

func getDataState(state int) *octo.Data_State {
	switch state {
	case 1:
		return octo.Data_ADD.Enum()
	case 2:
		return octo.Data_UPDATE.Enum()
	case 3:
		return octo.Data_LATEST.Enum()
	case 4:
		return octo.Data_DELETE.Enum()
	default:
		return octo.Data_NONE.Enum()
	}
}

func Upload(c *gin.Context) {

	var res struct {
		FileName string
		ProjectId string
		Backet string
		Error string
	}

	val, _ := c.Get("app")
	app, ok := val.(models.App)
	if !ok {
		res.Error = "Invalid App"
		c.JSON(http.StatusForbidden, res)
		return
	}
	appId := app.AppId
	filename := c.Query("filename")

	sizeParam := c.Query("size")
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		res.Error = sizeParam + " is not size"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	crcParam := c.Query("crc")
	crc, err := strconv.Atoi(crcParam)
	if err != nil {
		res.Error = sizeParam + " is not crc"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//難読化設定
	f := aesEnclipt(filename, app.AesKey)

	log.Println("appId: ", appId)
	log.Println("filename: ", filename)
	log.Println("size: ", size)
	log.Println("crc: ", crc)


	res.FileName = f


	//TODO GCP情報取得
	if app.StorageType == 1 {
		var gcp = models.Gcp{}
		gcpDao.GetGcp(&gcp, appId)

		res.ProjectId = gcp.ProjectId
		res.Backet = gcp.Backet
	}

	c.JSON(200, res)

}

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func aesEnclipt(str string, key_text string) string {
	plaintext := []byte(str)

	// 暗号化アルゴリズムaesを作成
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		log.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		//どうする
	}

	//暗号化文字列
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	return hex.EncodeToString(ciphertext)
}
