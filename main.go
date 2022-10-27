package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	bloburl "net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

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
hO5Zim9wVSsSwl4yyA==
=MIpx
-----END PGP PUBLIC KEY BLOCK-----`

const privkey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

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
=b7zP
-----END PGP PRIVATE KEY BLOCK-----` // encrypted private key

var pubKey string = ""
var privKey string = ""
var passphrase string = ""

type EventGridEvent struct {
	Topic     string `json:"topic"`
	Subject   string `json:"subject"`
	EventType string `json:"eventType"`
	ID        string `json:"id"`
	Data      struct {
		API                string `json:"api"`
		ClientRequestID    string `json:"clientRequestId"`
		RequestID          string `json:"requestId"`
		ETag               string `json:"eTag"`
		ContentType        string `json:"contentType"`
		ContentLength      int    `json:"contentLength"`
		BlobType           string `json:"blobType"`
		URL                string `json:"url"`
		Sequencer          string `json:"sequencer"`
		StorageDiagnostics struct {
			BatchID string `json:"batchId"`
		} `json:"storageDiagnostics"`
	} `json:"data"`
	DataVersion     string    `json:"dataVersion"`
	MetadataVersion string    `json:"metadataVersion"`
	EventTime       time.Time `json:"eventTime"`
}

type BlobPGPEncryptionRequest struct {
	AccountName string `json:"AccountName"`
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
	Path        string `json:"Path"`
}
type BlobPGPEncryptionResponse struct {
}

// PrintAndLog writes to stdout and to a logger.
func PrintAndLog(message string) {
	log.Println(message)
	fmt.Println(message)
}

func LogAndPanic(w http.ResponseWriter, err error) {
	PrintAndLog(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type InvokeRequest struct {
	Data     map[string]json.RawMessage
	Metadata map[string]interface{}
}

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

func getFileName(url string) string {
	lastIdx := strings.LastIndex(url, "/")
	sourceBlobName := url[lastIdx+1:]
	return sourceBlobName
}

func createBlobClientWithSaaSKey(url string, saasKey string) (*azblob.Client, error) {
	u, _ := bloburl.Parse(url)
	newUrl := fmt.Sprintf("https://%s/?%s", u.Hostname(), saasKey)
	client, err := azblob.NewClientWithNoCredential(newUrl, nil)
	return client, err
}

func readDataFromUrl(reader io.ReadCloser, contentLenght int64) ([]byte, error) {
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)

	if contentLenght == 0 {
		contentLenght = 1024
	}
	buf := make([]byte, contentLenght)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			// there is no more data to read
			PrintAndLog("EOF")
			break
		}
		if err != nil {
			return nil, err
		}
		if n > 0 {
			fmt.Print(buf[:n])
		}
	}
	return buf, nil
}

/*
HTTP Trigger Handler
*/
func pgpHttpTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var blobModified BlobPGPEncryptionRequest
		err := json.NewDecoder(r.Body).Decode(&blobModified)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusBadRequest)
			LogAndPanic(w, err)
			return
		}
		PrintAndLog(fmt.Sprintf("Path name %s , File name %s", blobModified.Path, blobModified.Name))
		saasKeySrc := ""
		if k, ok := os.LookupEnv("AZURE_BLOB_STORAGE_SAAS_KEY_NOPGP_SRC"); ok {
			saasKeySrc = k
		}
		//nickdevblob101.blob.core.windows.net
		url := fmt.Sprintf("https://%s.blob.core.windows.net", blobModified.AccountName)
		PrintAndLog(fmt.Sprintf("URL %s", url))
		client, err := createBlobClientWithSaaSKey(url, saasKeySrc) //azblob.NewClientWithNoCredential(newUrl, nil)
		if err != nil {
			LogAndPanic(w, err)
			return
		}
		PrintAndLog("Create blob client completed")
		stream, err := client.DownloadStream(context.TODO(), blobModified.Path, blobModified.Name, nil)
		if err != nil {
			LogAndPanic(w, err)
			return
		}
		PrintAndLog(fmt.Sprintf("Content lenght : %d", *stream.ContentLength))
		reader := stream.Body
		contentLenght := *stream.ContentLength
		buf, err := readDataFromUrl(reader, contentLenght)
		if err != nil {
			LogAndPanic(w, err)
			return
		}
		pubEntity, err := GetEntity([]byte(pubkey), []byte{})
		if err != nil {
			LogAndPanic(w, err)
			return
		}
		data, err := Encrypt(pubEntity, buf) //helper.EncryptBinaryMessageArmored(pubkey, buf)
		if err != nil {
			LogAndPanic(w, err)
			return
		} else {
			/*save to destination*/
			saasKeyDest := ""
			if k, ok := os.LookupEnv("AZURE_BLOB_STORAGE_SAAS_KEY_DEST"); ok {
				saasKeyDest = k
			}
			sourceBlobName := blobModified.Name //eventGridEvent.Data.URL[lastIdx+1:]
			PrintAndLog(fmt.Sprintf("Source Blob Name : %s", sourceBlobName))

			//u, _ := bloburl.Parse(eventGridEvent.Data.URL)
			//newUrl := fmt.Sprintf("https://%s/?%s", u.Hostname(), saasKeyDest)
			clientDest, err := createBlobClientWithSaaSKey(url, saasKeyDest) //azblob.NewClientWithNoCredential(newUrl, nil)
			if err != nil {
				LogAndPanic(w, err)
				return
			}
			_, err = clientDest.UploadBuffer(context.TODO(),
				"output-pgp",
				fmt.Sprintf("%s.adf.pgp", sourceBlobName),
				[]byte(data),
				nil,
			)
			if err != nil {
				LogAndPanic(w, err)
				return
			}
			PrintAndLog("Upload complete")
			////////////////
		}
	} else {
		_, err := fmt.Fprint(w, "Good job")
		if err != nil {
			LogAndPanic(w, err)
			return
		}
	}
	//json.Unmarshal(w)
	w.WriteHeader(200)
}

/*
EventGrid Trigger Handler
*/
func pgpEventGridBlobCreatedTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var invokeRequest InvokeRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&invokeRequest)
	if err != nil {
		log.Fatalf("Decode invoke request : %v\n", err)
		return
	}
	if val, ok := invokeRequest.Data["eventGridEvent"]; ok {
		marshalJSON, err := val.MarshalJSON()
		if err != nil {
			LogAndPanic(w, err)
			return
		}
		var eventGridEvent EventGridEvent
		err = json.Unmarshal(marshalJSON, &eventGridEvent)
		if err != nil {
			LogAndPanic(w, err)
			return
		}
		if strings.EqualFold(eventGridEvent.Data.API, "putblob") {
			PrintAndLog(fmt.Sprintf("PutBlob File -> %s", eventGridEvent.Data.URL))
			saasKeySrc := ""
			if k, ok := os.LookupEnv("AZURE_BLOB_STORAGE_SAAS_KEY_SRC"); ok {
				saasKeySrc = k
			}
			sourceBlobName := getFileName(eventGridEvent.Data.URL)
			PrintAndLog(fmt.Sprintf("Source Blob Name : %s", sourceBlobName))

			//u, _ := bloburl.Parse(eventGridEvent.Data.URL)
			//newUrl := fmt.Sprintf("https://%s/?%s", u.Hostname(), saasKeySrc)
			client, err := createBlobClientWithSaaSKey(eventGridEvent.Data.URL, saasKeySrc) //azblob.NewClientWithNoCredential(newUrl, nil)
			if err != nil {
				LogAndPanic(w, err)
				return
			}

			stream, err := client.DownloadStream(context.TODO(),
				"datas",
				sourceBlobName,
				nil,
			)

			if err != nil {
				LogAndPanic(w, err)
				return
			}

			PrintAndLog(fmt.Sprintf("Content lenght : %d", *stream.ContentLength))
			reader := stream.Body
			contentLenght := *stream.ContentLength
			buf, err := readDataFromUrl(reader, contentLenght)
			if err != nil {
				LogAndPanic(w, err)
				return
			}
			pubEntity, err := GetEntity([]byte(pubkey), []byte{})
			if err != nil {
				LogAndPanic(w, err)
				return
			}
			data, err := Encrypt(pubEntity, buf) //helper.EncryptBinaryMessageArmored(pubkey, buf)
			if err != nil {
				LogAndPanic(w, err)
				return
			} else {
				/*save to destination*/
				saasKeyDest := ""
				if k, ok := os.LookupEnv("AZURE_BLOB_STORAGE_SAAS_KEY_DEST"); ok {
					saasKeyDest = k
				}
				//lastIdx := strings.LastIndex(eventGridEvent.Data.URL, "/")
				sourceBlobName := getFileName(eventGridEvent.Data.URL) //eventGridEvent.Data.URL[lastIdx+1:]
				PrintAndLog(fmt.Sprintf("Source Blob Name : %s", sourceBlobName))
				//u, _ := bloburl.Parse(eventGridEvent.Data.URL)
				//newUrl := fmt.Sprintf("https://%s/?%s", u.Hostname(), saasKeyDest)
				clientDest, err := createBlobClientWithSaaSKey(eventGridEvent.Data.URL, saasKeyDest) //azblob.NewClientWithNoCredential(newUrl, nil)
				if err != nil {
					LogAndPanic(w, err)
					return
				}
				_, err = clientDest.UploadBuffer(context.TODO(),
					"output-pgp",
					fmt.Sprintf("%s.pgp", sourceBlobName),
					[]byte(data),
					nil,
				)

				if err != nil {
					LogAndPanic(w, err)
					return
				}
				PrintAndLog("Upload complete")
				////////////////
			}
		} else {
			PrintAndLog(eventGridEvent.Data.API)
		}
		PrintAndLog(fmt.Sprintf("Event Grid Data -> %s\n", string(marshalJSON)))
	} else {
		PrintAndLog("Don't have eventGridEvent")
	}
	w.WriteHeader(200)
}
func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	/*
		get azure key vault
	*/
	fmt.Println(pubKey)
	fmt.Println(privKey)
	fmt.Println(passphrase)

	///
	/// example get credential
	///
	/*
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("Azure indentity error : %v", err)
		} else {
			token, err := cred.GetToken(context.TODO(), policy.TokenRequestOptions{})
			if err != nil {
				log.Fatalf("Get Token err : %v", err)
			} else {
				log.Println("Token ", token)
			}
		}*/
	
	///

	//http.HandleFunc("/HttpTriggerPGP", pgpTriggerHandler)
	http.HandleFunc("/api/HttpTriggerPGP", pgpHttpTriggerHandler)                          //ADF
	http.HandleFunc("/EventGridTriggerBlobCreated", pgpEventGridBlobCreatedTriggerHandler) //EventGridTrigger
	//http.HandleFunc("/BlobTriggerPGP", pgpBlobTriggerHandler)

	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
