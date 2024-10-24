package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// DNSCheckResponse 用于解析DNS检查的响应
type DNSCheckResponse struct {
	Code   int `json:"code"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		Created_at string `json:"created_at"`
	} `json:"status"`
	Data struct {
		Dig struct {
			Status  string   `json:"status"`
			Records []string `json:"records"`
		} `json:"dig"`
		If_block struct {
			Status string `json:"status"`
		} `json:"if_block"`
		Trace struct {
			Status string   `json:"status"`
			Info   []string `json:"info"`
		} `json:"trace"`
	} `json:"data"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <domain>")
		return
	}
	domain := os.Args[1]

	// 设置请求的JSON数据
	jsonData := fmt.Sprintf(`{"api" : "Tools.Check.Dig","ori_domain" : "%s"}`, domain)

	// 创建请求体
	reqBody := bytes.NewBufferString(jsonData)

	// 发送POST请求
	resp, err := http.Post("https://www.dnspod.cn/cgi/dnsapi?action=Tools.Check.Dig&uin=&mc_gtk=&isSkipAuth=1", "application/json", reqBody)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	var dnsResp DNSCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&dnsResp); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}

	// 提取IP地址
	IPList := []string{}
	for _, record := range dnsResp.Data.Dig.Records {
		if strings.HasPrefix(record, "A ") {
			IPList = append(IPList, strings.Split(record, " ")[1])
		}
	}

	// 编辑hosts文件
	hostsPath := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	if err := editHostsFile(hostsPath, IPList, domain); err != nil {
		fmt.Printf("编辑hosts文件失败: %v\n", err)
		return
	}
}

// editHostsFile 用于编辑hosts文件
func editHostsFile(hostsPath string, IPList []string, domain string) error {
	// 读取hosts文件内容
	data, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	// 删除包含github.com的行
	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, domain) {
			newLines = append(newLines, line)
		}
	}

	// 添加新的IP地址和github.com
	if len(IPList) > 0 {
		newLines = append(newLines, fmt.Sprintf("%s %s", IPList[0], domain))
	}

	// 写回hosts文件
	return os.WriteFile(hostsPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// runAsAdmin 用于以管理员权限运行命令
func runAsAdmin(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
