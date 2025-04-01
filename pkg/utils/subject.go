package utils

import (
	"strconv"
	"strings"
)

// expandRange 函数用于展开字符串中的范围
func ExpandSubjectRange(s string) []string {
	// 按 . 分割字符串
	parts := strings.Split(s, ".")
	var result []string
	expandRecursive(parts, 0, "", &result)
	return result
}

// expandRecursive 递归函数用于生成所有可能的组合
func expandRecursive(parts []string, index int, current string, result *[]string) {
	// 如果已经处理完所有部分，将当前组合添加到结果列表中
	if index == len(parts) {
		*result = append(*result, current)
		return
	}
	part := parts[index]
	if strings.Contains(part, "~") {
		// 处理包含范围的部分
		rangeParts := strings.Split(part, "~")
		start, _ := strconv.Atoi(rangeParts[0])
		end, _ := strconv.Atoi(rangeParts[1])
		for i := start; i <= end; i++ {
			newCurrent := current
			if newCurrent != "" {
				newCurrent += "."
			}
			newCurrent += strconv.Itoa(i)
			expandRecursive(parts, index+1, newCurrent, result)
		}
	} else {
		// 处理不包含范围的部分
		newCurrent := current
		if newCurrent != "" {
			newCurrent += "."
		}
		newCurrent += part
		expandRecursive(parts, index+1, newCurrent, result)
	}
}
