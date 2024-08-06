package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mikiwa/models"
	"mikiwa/utils"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type FirebaseController struct {
	BaseController
}

func (c *FirebaseController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func PostFilesToFirebase(files []*multipart.FileHeader, userName string, ReferenceId int, pathName string, typeData string, folderName string) error {
	// filePath, errc := models.FirebaseAuth()
	// if errc != nil {
	// 	return fmt.Errorf("Error getting Firebase config: %s", errc.Error())
	// }
	// storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	filePath, storageBucket, errc := models.ValidateFirebase()

	if errc != nil {
		return fmt.Errorf("Error getting Firebase config: %s", errc.Error())
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("file upload error: %s", err.Error())
		}
		defer file.Close()

		fileName := uuid.New().String()
		opt := option.WithCredentialsJSON([]byte(filePath))
		config := &firebase.Config{
			StorageBucket: storageBucket,
		}
		app, err := firebase.NewApp(context.Background(), config, opt)
		if err != nil {
			return fmt.Errorf("Error Firebase app initialization: %s", err.Error())
		}

		client, err := app.Storage(context.Background())
		if err != nil {
			return fmt.Errorf("Error firebase Storage client initialization: %s", err.Error())
		}

		bucket, err := client.DefaultBucket()
		if err != nil {
			return fmt.Errorf("Error firebase default bucket retrieval: %s", err.Error())
		}

		pathName = strings.Replace(pathName, "%2F", "/", -1)

		newObj := bucket.Object(pathName + "/" + fileName)
		wc := newObj.NewWriter(context.Background())

		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("Error while seeking file content: %s", err.Error())
		}
		if _, err := io.Copy(wc, file); err != nil {
			return fmt.Errorf("Error while copying file content to Firebase Storage: %s", err.Error())
		}
		if err := wc.Close(); err != nil {
			return fmt.Errorf("Error while closing Firebase Storage writer: %s", err.Error())
		}

		fileType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"firebaseStorageDownloadTokens": fileName,
				"contentType":                   fileType,
			},
		}

		if _, err := bucket.Object(pathName+"/"+fileName).Update(context.Background(), objectAttrsToUpdate); err != nil {
			return fmt.Errorf("Error while updating object metadata: %s", err.Error())
		}

		pathName = strings.Replace(pathName, "/", "%2F", -1)
		newObjectNamePath := pathName + "%2F" + fileName

		downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "mpwdb-9c2d8.appspot.com", newObjectNamePath, fileName)

		t_documents := models.Document{
			ReferenceId: ReferenceId,
			FileName:    fileHeader.Filename,
			PathName:    fileName,
			PathFile:    downloadURL,
			FileType:    typeData,
			FolderName:  folderName,
			CreatedBy:   userName,
			UpdatedBy:   userName,
		}

		if _, err := t_documents.Insert(t_documents); err != nil {
			return fmt.Errorf("Error while inserting document record: %s", err.Error())
		}
	}

	return nil
}

func PutFilesFirebase(files []*multipart.FileHeader, userName string, ReferenceId int, pathName string, typeData string, folderName string) error {
	o := orm.NewOrm()

	var documents []models.Document
	_, err := models.Documents().Filter("reference_id", ReferenceId).Filter("deleted_at__isnull", true).Filter("file_type", typeData).All(&documents)
	if err != nil {
		return fmt.Errorf("Error fetching documents to delete: %s", err.Error())
	}

	for _, doc := range documents {
		found := false
		for _, fileHeader := range files {
			if fileHeader.Filename == doc.PathName {
				found = true
				break
			}
		}
		if found {
			continue
		} else {
			doc.DeletedAt = time.Now()
			_, err := o.Update(&doc, "deleted_at")
			if err != nil {
				return fmt.Errorf("Error updating deleted_at for document: %s", err.Error())
			}
			deleteFileFromStorage(doc.PathName, pathName)
		}
	}

	for _, fileHeader := range files {
		if err := checkFileExistInStorage(fileHeader.Filename, pathName); err == nil {
			var documentData models.Document
			err := models.Documents().Filter("path_name", fileHeader.Filename).One(&documentData)
			if err == orm.ErrNoRows {
				// continue
			} else if err != nil {
				return fmt.Errorf("Error fetching document data: %s", err.Error())
			} else {
				documentData.UpdatedAt = time.Now()
				documentData.UpdatedBy = userName
				_, err = o.Update(&documentData, "updated_at", "updated_by")
				if err != nil {
					return fmt.Errorf("Error updating document data: %s", err.Error())
				}
				continue
			}

		}

		err := PostFilesToFirebase([]*multipart.FileHeader{fileHeader}, userName, ReferenceId, pathName, typeData, folderName)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkFileExistInStorage(fileID string, pathName string) error {
	filePath, err := models.FirebaseAuth()
	if err != nil {
		return fmt.Errorf("Error getting Firebase config: %s", err.Error())
	}
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	opt := option.WithCredentialsJSON([]byte(filePath))
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return fmt.Errorf("Error firebase app initialization: %s", err.Error())
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return fmt.Errorf("Error firebase Storage client initialization: %s", err.Error())
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return fmt.Errorf("Error firebase default bucket retrieval: %s", err.Error())
	}

	pathName = strings.Replace(pathName, "%2F", "/", -1)

	_, err = bucket.Object(pathName + "/" + fileID).Attrs(context.Background())
	if err != nil {
		return fmt.Errorf("Error firebase bucket object: %s", err.Error())
	}
	return nil
}

func deleteFileFromStorage(fileID string, pathName string) error {
	filePath, err := models.FirebaseAuth()
	if err != nil {
		return fmt.Errorf("Error getting Firebase config: %s", err.Error())
	}
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	opt := option.WithCredentialsJSON([]byte(filePath))
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return fmt.Errorf("Error firebase app initialization: %s", err.Error())
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return fmt.Errorf("Error firebase Storage client initialization: %s", err.Error())
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return fmt.Errorf("Error retrieve firebase default bucket: %s", err.Error())
	}

	pathName = strings.Replace(pathName, "%2F", "/", -1)

	_, err = bucket.Object(pathName + "/" + fileID).Attrs(context.Background())
	if err != nil {
		if storage.ErrObjectNotExist == err {
			return nil
		}
		return fmt.Errorf("Error retrieve file attribute: %s", err.Error())
	}

	err = bucket.Object(pathName + "/" + fileID).Delete(context.Background())
	if err != nil {
		return fmt.Errorf("Error deletion file : %s", err.Error())
	}

	return nil
}

//TODO: NEW FIREBASE

func PostFirebaseRaw(rawData models.DocumentList, userName string, ReferenceId int, pathName string, typeData string, folderName string) error {
	filePath, storageBucket, errc := models.ValidateFirebase()

	if errc != nil {
		return fmt.Errorf("Error getting Firebase config: %s", errc.Error())
	}

	opt := option.WithCredentialsJSON([]byte(filePath))
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return fmt.Errorf("error Firebase app initialization: %s", err.Error())
	}

	storageClient, err := app.Storage(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing storage client: %v", err.Error())
	}

	bucketHandle, err := storageClient.Bucket(storageBucket)
	if err != nil {
		return fmt.Errorf("error getting bucket handle: %v", err.Error())
	}

	for i, file := range rawData.File {
		fileName := uuid.New().String()
		fileBinary, err := base64.StdEncoding.DecodeString(file.FileBinary)
		if err != nil {
			fmt.Printf("error decoding base64: %v\n", err)
			continue
		}

		ext := filepath.Ext(file.Name)
		contentType := mime.TypeByExtension(ext)
		pathName = strings.Replace(pathName, "%2F", "/", -1)

		objectPath := fmt.Sprintf("%s/%s", pathName, fileName)
		wc := bucketHandle.Object(objectPath).NewWriter(context.Background())
		wc.ContentType = contentType
		if _, err := wc.Write(fileBinary); err != nil {
			fmt.Printf("error writing file to storage: %v\n", err)
			wc.Close()
			continue
		}
		if err := wc.Close(); err != nil {
			fmt.Printf("error closing writer: %v\n", err)
			continue
		}

		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"firebaseStorageDownloadTokens": fileName,
				"contentType":                   contentType,
			},
		}

		if _, err := bucketHandle.Object(pathName+"/"+fileName).Update(context.Background(), objectAttrsToUpdate); err != nil {
			return fmt.Errorf("Error while updating object metadata: %s", err.Error())
		}

		pathName = strings.Replace(pathName, "/", "%2F", -1)
		newObjectNamePath := pathName + "%2F" + fileName

		objectURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "mpwdb-9c2d8.appspot.com", newObjectNamePath, fileName)

		rawData.File[i].ObjectURL = objectURL

		t_documents := models.Document{
			ReferenceId: ReferenceId,
			FileName:    file.Name,
			PathName:    fileName,
			PathFile:    objectURL,
			FileType:    typeData,
			FolderName:  folderName,
			CreatedBy:   userName,
			UpdatedBy:   userName,
		}

		if _, err := t_documents.Insert(t_documents); err != nil {
			return fmt.Errorf("Error while inserting document record: %s", err.Error())
		}
	}

	return nil
}

// TODO: NEW FIREBASE
func PutFirebaseRaw(rawData models.DocumentList, userName string, ReferenceId int, pathName string, typeData string, folderName string) error {
	o := orm.NewOrm()

	var documents []models.Document
	// _, err := models.Documents().Filter("reference_id", ReferenceId).Filter("deleted_at__isnull", true).Filter("file_type", typeData).All(&documents)
	_, err := o.Raw("select * from documents where deleted_at is null and reference_id = " + utils.Int2String(ReferenceId) + " and file_type = '" + typeData + "' ").QueryRows(&documents)
	if err != nil {
		return fmt.Errorf("Error fetching documents to delete: %s", err.Error())
	}

	var filesTemp []models.FirebaseFile
	var documentTemp models.DocumentList

	for _, doc := range documents {
		found := false
		for _, fileHeader := range rawData.File {
			if fileHeader.Name == doc.PathName {
				found = true
				break
			}
		}
		if found {
			continue
		} else {
			doc.DeletedAt = time.Now()
			_, err := o.Update(&doc, "deleted_at")
			if err != nil {
				return fmt.Errorf("Error updating deleted_at for document: %s", err.Error())
			}
			deleteFileFromStorage(doc.PathName, pathName)
		}
	}

	for _, fileHeader := range rawData.File {
		if err := checkFileExistInStorage(fileHeader.Name, pathName); err == nil {
			var documentData models.Document
			err := models.Documents().Filter("path_name", fileHeader.Name).One(&documentData)
			if err == orm.ErrNoRows {
				// continue
			} else if err != nil {
				return fmt.Errorf("Error fetching document data: %s", err.Error())
			} else {
				documentData.UpdatedAt = time.Now()
				documentData.UpdatedBy = userName
				_, err = o.Update(&documentData, "updated_at", "updated_by")
				if err != nil {
					return fmt.Errorf("Error updating document data: %s", err.Error())
				}
				continue
			}

		}

		filesTemp = append(filesTemp, fileHeader)
	}

	documentTemp = models.DocumentList{
		File: filesTemp,
	}

	if len(filesTemp) > 1 {
		fmt.Print(filesTemp)
	} else {
		fmt.Print(rawData.File)
	}

	if err := PostFirebaseRaw(documentTemp, userName, ReferenceId, pathName, typeData, folderName); err != nil {
		return fmt.Errorf("Error processing data and uploading to Firebase: %s", err.Error())
	}

	return nil
}

func PostFirebaseRawOne(rawData models.FirebaseFile, userName string, ReferenceId int, pathName string, typeData string, folderName string) error {
	filePath, storageBucket, errc := models.ValidateFirebase()

	if errc != nil {
		return fmt.Errorf("Error getting Firebase config: %s", errc.Error())
	}

	opt := option.WithCredentialsJSON([]byte(filePath))
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return fmt.Errorf("Error Firebase app initialization: %s", err.Error())
	}

	storageClient, err := app.Storage(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing storage client: %v", err)
	}

	bucketHandle, err := storageClient.Bucket(storageBucket)
	if err != nil {
		return fmt.Errorf("error getting bucket handle: %v", err)
	}

	fileName := uuid.New().String()
	fileBinary, err := base64.StdEncoding.DecodeString(rawData.FileBinary)
	if err != nil {
		fmt.Printf("error decoding base64: %v\n", err)
	}

	ext := filepath.Ext(rawData.Name)
	contentType := mime.TypeByExtension(ext)

	objectPath := fmt.Sprintf("%s/%s", pathName, fileName)
	wc := bucketHandle.Object(objectPath).NewWriter(context.Background())
	wc.ContentType = contentType
	if _, err := wc.Write(fileBinary); err != nil {
		fmt.Printf("error writing file to storage: %v\n", err)
		wc.Close()
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("error closing writer: %v\n", err)
	}

	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"firebaseStorageDownloadTokens": fileName,
			"contentType":                   contentType,
		},
	}

	if _, err := bucketHandle.Object(pathName+"/"+fileName).Update(context.Background(), objectAttrsToUpdate); err != nil {
		return fmt.Errorf("Error while updating object metadata: %s", err.Error())
	}

	pathName = strings.Replace(pathName, "/", "%2F", -1)
	newObjectNamePath := pathName + "%2F" + fileName

	objectURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "mpwdb-9c2d8.appspot.com", newObjectNamePath, fileName)

	// rawData.ObjectURL = objectURL

	t_documents := models.Document{
		ReferenceId: ReferenceId,
		FileName:    rawData.Name,
		PathName:    fileName,
		PathFile:    objectURL,
		FileType:    typeData,
		FolderName:  folderName,
		CreatedBy:   userName,
		UpdatedBy:   userName,
	}

	if _, err := t_documents.Insert(t_documents); err != nil {
		return fmt.Errorf("Error while inserting document record: %s", err.Error())
	}

	return nil
}
