package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/edwinbustillos/s3tool/awsconfig"
	"github.com/xlab/treeprint"
)

func ListFilesInBucket(args []string) {
	client := awsconfig.Connection()

	pathName := strings.TrimPrefix(args[0], "s3://")

	if strings.Contains(pathName, "/") == true {
		splitBucket := strings.Split(pathName, "/")
		bucketName := splitBucket[0]
		if len(bucketName) <= 0 {
			fmt.Println("Error Bucket name, try again E.g: s3://bucketName")
			os.Exit(0)
		}
		splitPath := strings.TrimPrefix(pathName, bucketName)
		if len(splitPath) <= 0 {
			fmt.Println("Error path file, try again E.g: s3://bucketName/path/file.txt")
			os.Exit(0)
		}
		pathFiles := strings.TrimPrefix(splitPath, "/")

		resp, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket: &bucketName,
		})
		if err != nil {
			fmt.Println("\nError list Bucket Files")
		} else {
			tree := treeprint.New()
			bucketN := tree.AddBranch("s3://" + bucketName)
			for _, object := range resp.Contents {
				if strings.Contains(*object.Key, pathFiles) == true {
					bucketN.AddNode(*object.Key)
				}
			}
			fmt.Println(tree.String())
		}

	} else {
		resp, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket: &pathName,
		})
		if err != nil {
			fmt.Println("\nError list Bucket Files")
		} else {
			tree := treeprint.New()
			bucketN := tree.AddBranch("s3://" + pathName)
			for _, object := range resp.Contents {
				bucketN.AddNode(*object.Key)
			}
			fmt.Println(tree.String())
		}

	}
}

func ListBuckets() {
	client := awsconfig.Connection()
	resp, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Error list buckets: " + err.Error())
		os.Exit(0)
	}
	tree := treeprint.New()
	bucketN := tree.AddBranch("Buckets in AWS:")
	for _, bucket := range resp.Buckets {
		bucketN.AddNode(*bucket.Name)
	}
	fmt.Println(tree.String())
}

func RemoveFile(args []string) {
	client := awsconfig.Connection()
	fileName := strings.TrimPrefix(args[0], "s3://")

	if len(fileName) <= 0 {
		fmt.Println("Error path file, try again E.g: bucketName/path/file.txt")
		os.Exit(0)
	}
	splitBucket := strings.Split(fileName, "/")
	bucketName := splitBucket[0]
	if len(bucketName) <= 0 {
		fmt.Println("Error Bucket name, try again E.g: bucket-name/path/file.txt")
		os.Exit(0)
	}
	fmt.Println("\n\nbucketName:" + bucketName)
	splitPath := strings.TrimPrefix(fileName, bucketName)
	if len(splitPath) <= 0 {
		fmt.Println("Error path file, try again E.g: bucket-name/path/file.txt")
		os.Exit(0)
	}
	fmt.Println("\n\nsplitPath:" + splitPath)
	objectKey := strings.TrimPrefix(splitPath, "/")
	fmt.Println("\n\nobjectKey:" + objectKey)
	_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		panic("Falha ao deletar o objeto: " + err.Error())
	}

	// Liste os objetos no bucket
	resp, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		panic("Falha ao listar objetos: " + err.Error())
	}
	// Verifique se o arquivo desejado está presente na lista de objetos
	var foundObject bool
	for _, object := range resp.Contents {

		if *object.Key == objectKey {
			foundObject = true
			break
		}
	}

	if foundObject {
		fmt.Println("Error Remove file:", bucketName+"/"+objectKey)
	} else {
		fmt.Println("File deleted successfully:......", objectKey)
	}
}

func RemoveFolder(args []string) {
	client := awsconfig.Connection()
	folderName := strings.TrimPrefix(args[0], "s3://")

	if len(folderName) <= 0 {
		fmt.Println("Error path, try again E.g: s3://bucketName/folderName")
		os.Exit(0)
	}

	splitBucket := strings.Split(folderName, "/")
	bucketName := splitBucket[0]

	splitPath := strings.TrimPrefix(folderName, bucketName)
	if len(splitPath) <= 0 {
		fmt.Println("Error path folder, try again E.g: bucketName/folderName")
		os.Exit(0)
	}
	prefix := strings.TrimPrefix(splitPath, "/")

	resp, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &prefix,
	})
	if err != nil {
		fmt.Println("Error check folder" + err.Error())
		os.Exit(0)
	}

	if resp.Contents == nil {
		fmt.Println("\n\nAviso: Não há objetos no bucket.")
		os.Exit(0)
	}
	// Delete todos os objetos com o prefixo do diretório
	for _, object := range resp.Contents {
		_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    object.Key,
		})
		if err != nil {
			fmt.Println("Error remove object......" + *object.Key + " | " + err.Error())
			os.Exit(0)
		}
		//fmt.Println("Remove objects successfully......", *object.Key)
	}

	// Opcionalmente, você pode excluir a "pasta" também
	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &prefix,
	})
	if err != nil {
		fmt.Println("Error remove folder......" + prefix + " | " + err.Error())
		os.Exit(0)
	}
	//fmt.Println("Remove folder successfully......", prefix)

	// Liste os objetos no bucket
	find, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		panic("Falha ao listar objetos: " + err.Error())
	}
	// Verifique se o arquivo desejado está presente na lista de objetos
	//var foundObject bool
	for _, object := range find.Contents {
		fileNames := *object.Key
		if fileNames == prefix {
			fmt.Println("\nError Remove object2: ", bucketName+"/"+prefix+"\n")
			os.Exit(0)

		}
	}
	fmt.Println("\nObject remove successfully:......", prefix+"\n")
}

// func FindObject(args []string) {

// 	client := awsconfig.Connection()

// 	fileName := args[0]
// 	if len(fileName) <= 0 {
// 		fmt.Println("Error path file, try again E.g: bucket-name/path/file.txt")
// 		os.Exit(0)
// 	}
// 	splitBucket := strings.Split(fileName, "/")
// 	bucketName := splitBucket[0]
// 	if len(bucketName) <= 0 {
// 		fmt.Println("Error Bucket name, try again E.g: bucket-name/path/file.txt")
// 		os.Exit(0)
// 	}
// 	splitPath := strings.TrimPrefix(fileName, bucketName)
// 	if len(splitPath) <= 0 {
// 		fmt.Println("Error path file, try again E.g: bucket-name/path/file.txt")
// 		os.Exit(0)
// 	}
// 	findFileName := strings.TrimPrefix(splitPath, "/")

// 	// Liste os objetos no bucket
// 	resp, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
// 		Bucket: &bucketName,
// 	})
// 	if err != nil {
// 		panic("Falha ao listar objetos: " + err.Error())
// 	}

// 	// Verifique se o arquivo desejado está presente na lista de objetos
// 	var foundObject bool
// 	for _, object := range resp.Contents {
// 		if strings.Contains(*object.Key, findFileName) == true {
// 			foundObject = true
// 			findFileName = *object.Key
// 			break
// 		}
// 		// if *object.Key == findFileName {
// 		// 	foundObject = true
// 		// 	break
// 		// }
// 	}

// 	if foundObject {
// 		fmt.Println("Arquivo encontrado:", bucketName+"/"+findFileName)
// 	} else {
// 		fmt.Println("Arquivo não encontrado:", findFileName)
// 	}
// }

func GetFiles(args []string) {
	client := awsconfig.Connection()

	pathName := strings.TrimPrefix(args[0], "s3://")

	if strings.Contains(pathName, "/") == true {
		splitBucket := strings.Split(pathName, "/")
		bucketName := splitBucket[0]
		if len(bucketName) <= 0 {
			fmt.Println("Error Bucket name, try again E.g: bucket-name/path/file.txt")
			os.Exit(0)
		}
		fmt.Println("1")
		splitPath := strings.TrimPrefix(pathName, bucketName)
		if len(splitPath) <= 0 {
			fmt.Println("Error path file, try again.")
			os.Exit(0)
		}
		fmt.Println("2")
		pathFiles := strings.TrimPrefix(splitPath, "/")
		fmt.Println("3")
		// Solicita o download do arquivo
		resp, err := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &pathFiles,
		})
		fmt.Println("4")
		if err != nil {
			fmt.Println("Falha ao fazer o download do objeto: " + err.Error())
			os.Exit(0)
		}
		pathFile := strings.Split(pathFiles, "/")

		fileName := pathFile[len(pathFile)-1]
		// Crie um arquivo local onde o conteúdo será salvo
		file, err := os.Create(fileName) // Substitua "local-file.txt" pelo caminho do arquivo local desejado
		if err != nil {
			fmt.Println("Falha ao criar o arquivo local: " + err.Error())
			os.Exit(0)
		}
		defer file.Close()
		// Copie o conteúdo do objeto para o arquivo local
		_, err = file.ReadFrom(resp.Body)
		if err != nil {
			fmt.Println("Falha ao copiar o conteúdo do objeto para o arquivo local: " + err.Error())
			os.Exit(0)
		}
		fmt.Println("Download do arquivo concluído com sucesso.")

	} else {
		fmt.Println("\nError getting file path. Try again with correct file path.")
	}
}

func Upload() {

	// Crie um cliente S3 usando as configurações carregadas do arquivo ~/.aws/config
	client, err := awsconfig.S3Client()

	bucketName := "your-bucket-name" // Nome do bucket

	objectKey := "path/to/document.pdf" // Caminho/nome do arquivo no bucket

	filePath := "path/to/local-document.pdf" // Caminho/nome do arquivo local

	// Abra o arquivo local para upload
	file, err := os.Open(filePath)
	if err != nil {
		panic("Falha ao abrir o arquivo local: " + err.Error())
	}
	defer file.Close()

	// Upload do arquivo
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   file,
	})
	if err != nil {
		panic("Falha ao fazer o upload do arquivo: " + err.Error())
	}

	fmt.Println("Upload do documento concluído com sucesso.")
}

func RenameFile() {

	// Crie um cliente S3 usando as configurações carregadas do arquivo ~/.aws/config
	client, err := awsconfig.S3Client()
	if err != nil {
		panic("Falha ao criar o cliente S3: " + err.Error())
	}

	// Nome do bucket que contém o arquivo
	bucketName := "your-bucket-name"

	// Nome do arquivo que você deseja renomear
	objectKeyOriginal := "path/to/original-file.txt"

	// Nome do novo arquivo após a renomeação
	objectKeyNovo := "path/to/novo-nome-arquivo.txt"

	// Crie a cópia do arquivo com o novo nome
	_, err = client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     &bucketName,
		CopySource: &objectKeyOriginal,
		Key:        &objectKeyNovo,
	})
	if err != nil {
		panic("Falha ao renomear o arquivo: " + err.Error())
	}

	// Deleta o arquivo original após a cópia
	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKeyOriginal,
	})
	if err != nil {
		panic("Falha ao deletar o arquivo original: " + err.Error())
	}

	fmt.Println("Arquivo renomeado com sucesso:", objectKeyNovo)
}

func MoveFile() {
	// Crie um cliente S3 usando as configurações carregadas do arquivo ~/.aws/config
	client, err := awsconfig.S3Client()
	if err != nil {
		panic("Falha ao criar o cliente S3: " + err.Error())
	}

	// Nome do bucket que contém o arquivo
	bucketName := "your-bucket-name"

	// Caminho/nome do arquivo original (incluindo diretório atual)
	objectKeyOriginal := "path/to/diretorio-original/original-file.txt"

	// Caminho/nome do novo arquivo (incluindo novo diretório)
	objectKeyNovo := "path/to/novo-diretorio/novo-nome-arquivo.txt"

	// Crie a cópia do arquivo com o novo caminho/nome
	_, err = client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     &bucketName,
		CopySource: &objectKeyOriginal,
		Key:        &objectKeyNovo,
	})
	if err != nil {
		panic("Falha ao mover o arquivo: " + err.Error())
	}

	// Deleta o arquivo original após a cópia
	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKeyOriginal,
	})
	if err != nil {
		panic("Falha ao deletar o arquivo original: " + err.Error())
	}

	fmt.Println("Arquivo movido com sucesso:", objectKeyNovo)
}
