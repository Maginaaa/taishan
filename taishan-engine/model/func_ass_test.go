package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestFunc(t *testing.T) {
	// 打开 Excel 文件
	file, err := excelize.OpenFile("/Users/lyn/workspace/taishan/taishan-engine/model/xlt.xlsx")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 获取工作表名称（默认第一个工作表）
	sheetName := file.GetSheetName(0)

	// 读取所有行
	rows, err := file.GetRows(sheetName)
	if err != nil {
		log.Fatalf("读取行失败: %v", err)
	}
	// 遍历每一行（从第二行开始，跳过表头）
	for rowIndex, row := range rows[1:] {
		// 假设计算结果为 content 和 ai_domain 的拼接
		//content := row[2]                    // content 列
		//aiDomain := row[3]                   // ai_domain 列
		fmt.Println(row[0])
		// 将结果写入新列
		baseContent := row[2]
		var con []BaseContent
		json.Unmarshal([]byte(baseContent), &con)
		res := getDomain(con[0].Problem)
		row = append(row, res)

		// 更新行数据
		rows[rowIndex+1] = row
	}
	newFile := excelize.NewFile()
	newSheetName := "Sheet1"

	// 写入表头
	for colIndex, header := range rows[0] {
		cell, _ := excelize.CoordinatesToCellName(colIndex+1, 0)
		newFile.SetCellValue(newSheetName, cell, header)
	}

	// 写入数据
	for rowIndex, row := range rows {
		for colIndex, cellValue := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1) // 从第二行开始写入数据
			newFile.SetCellValue(newSheetName, cell, cellValue)
		}
	}

	// 保存新文件
	if err := newFile.SaveAs("output.xlsx"); err != nil {
		log.Fatalf("保存文件失败: %v", err)
	}

	fmt.Println("处理完成，结果已写入 output.xlsx")
}

func getContent(content string) (res string) {
	url := "http://10.72.248.133:81/v1/chat-messages"

	// 请求体数据结构
	data := map[string]interface{}{
		"inputs":          map[string]string{},
		"query":           content,
		"response_mode":   "blocking",
		"conversation_id": "",
		"user":            "devops",
		"files":           []string{},
	}

	// 将数据结构转换为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling data: %v", err)
	}

	// 创建一个新的HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
	}

	// 设置请求头
	analyseAuth := "Bearer app-TQeeciddb2gk2WP1Tac3hIvy"
	//domainAuth := "Bearer app-olr5g18fpXNCCMZiWLUh6CTd"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", analyseAuth)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error read body: %v", err)
		return
	}
	var r RootEntity
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("Error Unmarshal: %v", err)
		return
	}
	fmt.Println(r.Answer)
	return r.Answer
}

type RootEntity struct {
	Event          string `json:"event"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       struct {
		Usage struct {
			PromptTokens        int64   `json:"prompt_tokens"`
			PromptUnitPrice     string  `json:"prompt_unit_price"`
			PromptPriceUnit     string  `json:"prompt_price_unit"`
			PromptPrice         string  `json:"prompt_price"`
			CompletionTokens    int64   `json:"completion_tokens"`
			CompletionUnitPrice string  `json:"completion_unit_price"`
			CompletionPriceUnit string  `json:"completion_price_unit"`
			CompletionPrice     string  `json:"completion_price"`
			TotalTokens         int64   `json:"total_tokens"`
			TotalPrice          string  `json:"total_price"`
			Currency            string  `json:"currency"`
			Latency             float64 `json:"latency"`
		} `json:"usage"`
	} `json:"metadata"`
	CreatedAt int64 `json:"created_at"`
}

func getDomain(content string) (res string) {
	url := "http://10.72.248.133:81/v1/chat-messages"

	// 请求体数据结构
	data := map[string]interface{}{
		"inputs":          map[string]string{},
		"query":           content,
		"response_mode":   "blocking",
		"conversation_id": "",
		"user":            "devops",
		"files":           []string{},
	}

	// 将数据结构转换为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling data: %v", err)
	}

	// 创建一个新的HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
	}

	// 设置请求头
	//analyseAuth := "Bearer app-TQeeciddb2gk2WP1Tac3hIvy"
	domainAuth := "Bearer app-olr5g18fpXNCCMZiWLUh6CTd"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", domainAuth)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error read body: %v", err)
		return
	}
	var r RootEntity
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("Error Unmarshal: %v", err)
		return
	}
	fmt.Println(r.Answer)
	return r.Answer
}

type BaseContent struct {
	Solved                string `json:"solved"`
	Problem               string `json:"problem"`
	Solution              string `json:"solution"`
	Explanation           string `json:"explanation"`
	InteractionSessions   string `json:"interaction_sessions"`
	SatisfactoryReplies   string `json:"satisfactory_Replies"`
	UserDisapprovalReason string `json:"user_disapproval_reason"`
}
