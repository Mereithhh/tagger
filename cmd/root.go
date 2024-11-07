/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"jojo/tagger/util"
	"os"

	"github.com/spf13/cobra"
)

//go:embed logo
var logo string
var CurrentPrefix string
var CurrentSuffix string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tagger",
	Short: "自动打标签工具",
	Long: fmt.Sprintf(`
	%s

自动打标签工具。可以自动查找当前仓库的所有标签，找到最新的并加一，然后推送上去。

注意：适配的标签格式为：v0.0.1、v0.0.2 这样的。

用法：

tagger  [patch｜minor｜major] [-p <前缀名>] [-s <后缀名>]

PS:第一个参数可以忽略，直接打最小版本。
比如目前最新的标签是 v0.0.1，用此工具后会变成 v0.0.2，然后推送。

后面的参数可以为 major、minor 或 patch，分别表示主版本、次版本、修订版本。默认为 patch。

可以用来配合触发 CI/CD。

-p 参数可以指定前缀名，默认为 v。 可以指定不同的前缀， 比如 prod-0.0.1, prod-0.0.2 这样的。
-s 参数可以指定后缀名，默认为空。 可以指定不同的后缀， 比如 v0.0.1-dev, v0.0.2-dev 这样的。

也可以通过命令设置默认的前后缀
tagger set-default-prefix <前缀名>
tagger set-default-suffix <后缀名>

或者得到当前默认的前后缀
tagger info
	`, logo),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		prefix := CurrentPrefix
		suffix := CurrentSuffix
		var err error

		// 如果未指定前缀，使用默认前缀
		if prefix == "" {
			prefix, err = util.GetDefaultPrefix()
			if err != nil {
				fmt.Printf("从配置文件获取默认前缀失败: %v\n", err)
				return
			}
		}
		// 如果未指定后缀，使用默认后缀
		if suffix == "" {
			suffix, err = util.GetDefaultSuffix()
			if err != nil {
				fmt.Printf("从配置文件获取默认后缀失败: %v\n", err)
				return
			}
		}

		// 默认使用 patch 模式
		mode := "patch"
		if len(args) > 0 {
			// 验证输入的版本类型是否有效
			switch args[0] {
			case "major", "minor", "patch":
				mode = args[0]
			default:
				fmt.Printf("无效的版本类型 %s，必须是 major、minor 或 patch\n", args[0])
				return
			}
		}

		// 执行打标签操作
		util.TagByModeVersion(prefix, mode, suffix)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&CurrentPrefix, "prefix", "p", "", "tag prefix")
	rootCmd.PersistentFlags().StringVarP(&CurrentSuffix, "suffix", "s", "", "tag suffix")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// pool := rootCmd.Flags().StringP("pool", "p", "", "bedrock pool name")

	// add version
	rootCmd.Version = "0.5.0"

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
