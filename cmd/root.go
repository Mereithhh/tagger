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

自动打标签工具，用于自动查找仓库最新标签并递增版本号。

用法：tagger [patch|minor|major] [-p <前缀>] [-s <后缀>]

参数说明：
- 版本类型：patch(修订版本)、minor(次版本)、major(主版本)，默认为patch
- -p：标签前缀，默认为"v"，如：v0.0.1、prod-0.0.1
- -s：标签后缀，默认为空，如：v0.0.1-dev

相关命令：
tagger set-default-prefix <前缀>  设置默认前缀
tagger set-default-suffix <后缀>  设置默认后缀
tagger info                       查看当前默认配置
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
