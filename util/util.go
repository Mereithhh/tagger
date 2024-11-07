package util

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const TagReg = `^(staging|pre|prod)(-[a-zA-Z0-9]{0,10}){0,1}-\d+\.\d+\.\d+$`

func PrintInfo() {
	// fmt.Println("tagger v0.4.0")
}

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

func GetLatestTag(mode string, n int, pool string) (string, error) {
	tags, err := LoadTags()
	if err != nil {
		return "", err
	}
	// filter tags by mode
	newArr := make([]string, 0)
	for _, tag := range tags {
		if strings.HasPrefix(tag, mode) {
			if pool != "" {
				if IsPoolTag(tag, pool) {
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
		return AddPoolName(fmt.Sprintf("%s-0.0.1", mode), pool), nil
	}
	big := ""
	for _, tag := range newArr {
		if CmpTag(tag, big, mode, n, pool) {
			big = tag
		}
	}
	newTag := AddOneNumber(big, mode, n, pool)

	return AddPoolName(newTag, pool), nil
}

func AddPoolName(oldtag string, pool string) string {
	if pool == "" {
		return oldtag
	}

	arr := strings.Split(oldtag, "-")
	return fmt.Sprintf("%s-%s-%s", arr[0], pool, arr[1])

}

func CmpTag(tag1, tag2, mode string, n int, pool string) bool {
	if tag2 == "" {
		return true
	}
	t1 := strings.TrimPrefix(tag1, fmt.Sprintf("%s-", mode))
	t2 := strings.TrimPrefix(tag2, fmt.Sprintf("%s-", mode))

	if pool != "" {
		t1 = strings.TrimPrefix(t1, fmt.Sprintf("%s-", pool))
		t2 = strings.TrimPrefix(t2, fmt.Sprintf("%s-", pool))
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

func AddOneNumber(tag, mode string, n int, pool string) string {
	t := strings.TrimPrefix(tag, fmt.Sprintf("%s-", mode))
	if pool != "" {
		t = strings.TrimPrefix(t, fmt.Sprintf("%s-", pool))
	}
	tStringArr := strings.Split(t, ".")
	tNumberArr := make([]int, 0)
	for _, s := range tStringArr {
		num, _ := strconv.Atoi(s)
		tNumberArr = append(tNumberArr, num)
	}
	tNumberArr[n] += 1
	newTag := fmt.Sprintf("%s-%d.%d.%d", mode, tNumberArr[0], tNumberArr[1], tNumberArr[2])
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

	PrintInfo()
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
	newTag, err := GetLatestTag(prefix, n, suffix)
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

func GetTagPool(tag string) string {
	if !IsValidTag(tag) {
		return ""
	}
	arr := strings.Split(tag, "-")
	if len(arr) != 3 {
		return ""
	}
	return arr[1]
}

func IsPoolTag(tag string, pool string) bool {
	return GetTagPool(tag) == pool
}
