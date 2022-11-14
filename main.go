package main

import (
	"container/list"
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic"
)

type ElasticDocs struct {
	SomeStr   string
	SomeInt   int
	SomeBool  bool
	Timestamp int64
}

type Enimal struct {
	Name string
	Age  int
}

func main() {
	// cho phép định dạng
	log.SetFlags(0)

	// Sử dụng gói Olivere để lấy số phiên bản Elasticsearch
	fmt.Println("Version: ", elastic.Version)

	// Tạo một đối tượng ngữ cảnh cho các lệnh gọi API
	ctx := context.Background()

	// Khai báo một phiên bản ứng dụng khách của trình điều khiển Olivere
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.14.151:9200"),
		elastic.SetBasicAuth("elastic", "OZpbNlZ1VsoIHh38gt4i"),
		elastic.SetHealthcheckInterval(5*time.Second),
	)

	// Kiểm tra xem phương thức NewClient () của olivere có trả lại lỗi
	if err != nil {
		fmt.Println("elastic.NewClient() ERROR: ", err)
		log.Fatalf("quiting connection...")
	} else {
		fmt.Println("client: ", client)
		fmt.Println("client TYPE:", reflect.TypeOf(client))
	}

	indexName := "student"

	indices := []string{indexName}

	existService := elastic.NewIndicesExistsService(client)

	existService.Index(indices)
	exist, err := existService.Do(ctx)

	// Kiểm tra xem phương thức IndicesExistsService.Do () có trả về bất kỳ lỗi nào
	if err != nil {
		log.Fatalf("IndicesExistsService.Do() ERROR:", err)

	} else if exist == false {
		fmt.Println("nOh no! The index", indexName, "doesn't exist.")
		fmt.Println("Create the index, and then run the Go script again")
	} else if exist == true {
		fmt.Println("Index name", indexName, " exist!")
		//	docs := []ElasticDocs{}
		var docs list.List
		fmt.Println("docs TYPE:", reflect.TypeOf(docs))
		newDoc1 := ElasticDocs{" Anh", 422, true, 0.0}
		newDoc2 := ElasticDocs{" Uyen", 413, true, 0.0}
		newDoc3 := ElasticDocs{" Binh", 45, true, 0.0}
		newDog := Enimal{"gau", 2}
		docs.PushBack(newDoc1)
		docs.PushBack(newDoc2)
		docs.PushBack(newDoc3)
		docs.PushBack(newDog)

		// docs = append(docs, newDoc1)
		// docs = append(docs, newDoc2)
		// docs = append(docs, newDoc3)
		// docs = append(docs, newDog)

		bulk := client.Bulk()

		docID := 0

		for doc := docs.Front(); doc != nil; doc = doc.Next() {
			docID++
			idStr := strconv.Itoa(docID)
			// doc.Timestamp = time.Now().Unix()
			// fmt.Println("ntime.Now().Unix():", doc.Timestamp)
			req := elastic.NewBulkIndexRequest()
			req.OpType("index")
			req.Index(indexName)
			req.Id(idStr)
			req.Doc(doc)
			fmt.Println("rep:", req)
			fmt.Println("req TYPE:", reflect.TypeOf(req))
			bulk = bulk.Add(req)
			fmt.Println("NewBulkIndexRequest().NumberOfActions():", bulk.NumberOfActions())
		}

		// for _, doc := range docs {
		// 	docID++
		// 	idStr := strconv.Itoa(docID)
		// 	// doc.Timestamp = time.Now().Unix()
		// 	// fmt.Println("ntime.Now().Unix():", doc.Timestamp)
		// 	req := elastic.NewBulkIndexRequest()
		// 	req.OpType("index")
		// 	req.Index(indexName)
		// 	req.Id(idStr)
		// 	req.Doc(doc)
		// 	fmt.Println("rep:", req)
		// 	fmt.Println("req TYPE:", reflect.TypeOf(req))
		// 	bulk = bulk.Add(req)
		// 	fmt.Println("NewBulkIndexRequest().NumberOfActions():", bulk.NumberOfActions())
		// }
		// gửi hàng loạt các yêu cầu đến elastic
		bulkResp, err := bulk.Do(ctx)
		// Kiểm tra xem có bắn trả lỗi nào không
		if err != nil {
			log.Fatalf("bulk.Do(ctx) ERROR", err)
		} else {
			// nếu không thấy lỗi thì lấy phản hồi từ API
			indexed := bulkResp.Indexed()
			fmt.Println("nbulkResp.Indexed():", indexed)
			fmt.Println("bulkResp.Indexed() TYPE:", reflect.TypeOf(indexed))

			// Lặp lại đối tượng BulResp.Indexed () được trả về từ số lượng lớn.go
			t := reflect.TypeOf(indexed)
			fmt.Println("nt:", t)
			fmt.Println("NewBulkIndexRequest().NumberOfActions():", bulk.NumberOfActions())

			//Iterate over the document responses
			// for i := 0; i < t.NumMethod(); i++ {
			// 	method := t.Method(i)
			// 	fmt.Println("nbulkResp.Indexed() METHOD NAME:", i, method.Name)
			// 	fmt.Println("bulkResp.Indexed() method:", method)
			// }

			// Return data on the documents indexed
			// fmt.Println("nBulk response Index:", indexed)
			// for _, info := range indexed {
			// 	fmt.Println("nBulk response Index:", info)
			// 	//fmt.Println("nBulk response Index:", info.Index)
			// }
		}
	}
}
