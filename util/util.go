package util

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const TagReg = `^([a-zA-Z0-9_-]*[a-zA-Z_-])?\d+\.\d+\.\d+([a-zA-Z_-][a-zA-Z0-9_-]*)?$`


func Log(stage string, msg string) {
	fmt.Println(msg)
}

func ExecCmd(cmd string, args ...string) error {
	cmdObj := exec.Command(cmd, args...)
	cmdObj.Stdout = os.Stdout
	cmdObj.Stderr = os.Stderr
	err := cmdObj.Run()
	if err != nil {
		return err
	}
	return nil
}

func GitPull() error {
	return ExecCmd("git", "pull")
}

func Fetch() error {
	return ExecCmd("git", "fetch", "--all", "-f")
}

func SetTagAndPush(tag string) error {
	err := ExecCmd("git", "tag", tag)
	if err != nil {
		return err
	}
	err = ExecCmd("git", "push", "origin", tag)
	if err != nil {
		return err
	}
	return nil

}

func LoadTags() ([]string, error) {

	cmd := exec.Command("git", "tag")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	outString := string(out)

	tagArr := strings.Split(outString, "\n")
	newArr := make([]string, 0)
	// filter empty string
	for i, tag := range tagArr {
		if tag == "" {
			newArr = append(tagArr[:i], tagArr[i+1:]...)
		}
	}

	return newArr, nil
}

func GetLatestTag(prefix string, n int, suffix string, tags []string) (string, error) {
	var err error
	if tags == nil {
		tags, err = LoadTags()
		if err != nil {
			return "", err
		}
	}
	// filter tags by prefix
	newArr := make([]string, 0)
	for _, tag := range tags {
		if strings.HasPrefix(tag, prefix) {
			if suffix != "" {
				if IsSuffixTag(tag, suffix) {
					if IsValidTag(tag) {
						newArr = append(newArr, tag)
					}
				}
			} else {
				if IsValidTag(tag) {
					newArr = append(newArr, tag)
				}
			}

		}
	}
	if len(newArr) == 0 {
		return fmt.Sprintf("%s0.0.1%s", prefix, suffix), nil
	}
	big := ""
	for _, tag := range newArr {
		if CmpTag(tag, big, prefix, n, suffix) {
			big = tag
		}
	}
	newTag := AddOneNumber(big, prefix, n, suffix)

	return newTag, nil
}



func CmpTag(tag1, tag2, prefix string, n int, suffix string) bool {
	if tag2 == "" {
		return true
	}
	t1 := strings.TrimPrefix(tag1, fmt.Sprintf("%s", prefix))
	t2 := strings.TrimPrefix(tag2, fmt.Sprintf("%s", prefix))

	if suffix != "" {
		t1 = strings.TrimSuffix(t1, fmt.Sprintf("%s", suffix))
		t2 = strings.TrimSuffix(t2, fmt.Sprintf("%s", suffix))
	}


	t1StringArr := strings.Split(t1, ".")
	t2StringArr := strings.Split(t2, ".")
	for i := 0; i < n+1; i++ {
		curA, _ := strconv.Atoi(t1StringArr[i])
		curB, _ := strconv.Atoi(t2StringArr[i])
		if curA > curB {
			return true
		} else if curA < curB {
			return false
		}
	}
	return false

}

func AddOneNumber(tag, prefix string, n int, suffix string) string {
	t := strings.TrimPrefix(tag, fmt.Sprintf("%s", prefix))
	if suffix != "" {
		t = strings.TrimSuffix(t, fmt.Sprintf("%s", suffix))
	}
	tStringArr := strings.Split(t, ".")
	tNumberArr := make([]int, 0)
	for _, s := range tStringArr {
		num, _ := strconv.Atoi(s)
		tNumberArr = append(tNumberArr, num)
	}
	tNumberArr[n] += 1
	newTag := fmt.Sprintf("%s%d.%d.%d", prefix, tNumberArr[0], tNumberArr[1], tNumberArr[2])
	if suffix != "" {
		newTag = fmt.Sprintf("%s%s", newTag, suffix)
	}
	return newTag
}

func Contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func GetNByVersion(version string) int {
	n := 2
	if version == "minor" {
		n = 1
	}
	if version == "major" {
		n = 0
	}
	return n

}

func TagByModeVersion(prefix string, version string, suffix string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
		}
	}()

	fmt.Printf("=== 标签生成信息 ===\n")
	if suffix != "" {
		fmt.Printf("后缀: %s\n", suffix)
	}
	fmt.Printf("前缀: %s\n", prefix)
	fmt.Printf("版本: %s\n", version)
	fmt.Println("====================")
	// fmt.Println("[同步远端信息] 开始")
	err = GitPull()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = Fetch()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("[同步远端信息] 完成")
	// fmt.Println("[获取标签信息] 开始")

	n := GetNByVersion(version)
	newTag, err := GetLatestTag(prefix, n, suffix,nil )
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("[获取标签信息] 完成")
	fmt.Println("[要打标签]", newTag)
	// fmt.Println("[打标签] 开始")
	err = SetTagAndPush(newTag)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("[打标签] 完成")
}

func IsValidTag(tag string) bool {
	// 用上面的 reg 过滤
	regex, err := regexp.Compile(TagReg)
	if err != nil {
		return false
	}
	return regex.MatchString(tag)

}

func GetTagSuffix(tag string) string {
	if !IsValidTag(tag) {
		return ""
	}
	// 使用正则表达式匹配版本号部分
	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+(-[a-zA-Z0-9_]+)?$`)
	versionPart := versionRegex.FindString(tag)
	if versionPart == "" {
		return ""
	}

	// 去掉版本号部分
	prefix := strings.TrimSuffix(tag, versionPart)
	// 去掉前缀和最后的分隔符（如果存在）
	if prefix == "" {
		return ""
	}
	prefix = strings.TrimRight(prefix, "-_")

	// 如果前缀包含分隔符，取最后一个分段作为 suffix
	parts := regexp.MustCompile(`[-_]`).Split(prefix, -1)
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	
	return ""
}

func IsSuffixTag(tag string, suffix string) bool {
	_, _, suffix, err := GetTagParts(tag)
	if err != nil {
		return false
	}
	return suffix == suffix
}

// GetTagParts 从标签中提取前缀、版本号和后缀
// 返回值: 前缀(可能为空), 版本号, 后缀(可能为空), 错误信息
func GetTagParts(tag string) (prefix string, version string, suffix string, err error) {
	if !IsValidTag(tag) {
		return "", "", "", fmt.Errorf("invalid tag format: %s", tag)
	}

	// 使用正则表达式匹配版本号部分
	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+`)
	version = versionRegex.FindString(tag)

	// 分割版本号前后的部分
	parts := versionRegex.Split(tag, -1)
	
	// 处理前缀，保留尾部的连字符
	if parts[0] != "" {
		prefix = parts[0]
	}
	
	// 处理后缀，保留头部的连字符
	if len(parts) > 1 {
		suffix = parts[1]
	}

	return prefix, version, suffix, nil
}

