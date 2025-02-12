package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// 创建一个新的watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 需要监控的文件路径
	filePath := getMappingFile()

	// 添加需要监听的文件
	if err := watcher.Add(filePath); err != nil {
		log.Fatal(err)
	}

	fmt.Println("开始监听文件" + filePath + "变动...")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 判断是否是修改事件
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("文件被修改: %s\n", event.Name)

					// 打开文件
					file, err := os.Open(event.Name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "无法打开文件: %v\n", err)
						return
					}
					defer file.Close()

					// 读取文件内容
					scanner := bufio.NewScanner(file)
					result := make(map[string]string)

					lineNum := 1
					for scanner.Scan() {
						line := scanner.Text()
						parts := strings.Split(line, "\t")

						// 检查分割后的部分数量是否为2
						if len(parts) == 2 {
							result[parts[0]] = parts[1]
						} else {
							fmt.Printf("第 %d 行格式错误: %s\n", lineNum, line)
						}

						lineNum++
					}

					// 检查扫描过程中是否有错误发生
					if err := scanner.Err(); err != nil {
						fmt.Fprintf(os.Stderr, "读取文件时发生错误: %v\n", err)
						return
					}

					// 打印map内容（可选）
					fmt.Println("解析完成的mapping:")
					for key, value := range result {
						fmt.Printf("%s -> %s\n", key, value)
					}
					updateCache(result)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("错误:", err)
			}
		}
	}()

	startServer()
}
