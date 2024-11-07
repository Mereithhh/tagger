package cmd

import (
	"fmt"
	"jojo/tagger/util"

	"github.com/spf13/cobra"
)

var setDefaultPrefixCmd = &cobra.Command{
	Use:   "set-default-prefix",
	Short: "设置默认的标签前缀",
	Long:  "设置默认的标签前缀，例如：tagger set-default-prefix v",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prefix := args[0]
		if err := util.SetDefaultPrefix(prefix); err != nil {
			fmt.Printf("设置默认前缀失败: %v\n", err)
			return
		}
		fmt.Printf("成功设置默认前缀为: %s\n", prefix)
	},
}

var setDefaultSuffixCmd = &cobra.Command{
	Use:   "set-default-suffix",
	Short: "设置默认的标签后缀",
	Long:  "设置默认的标签后缀，例如：tagger set-default-suffix -dev",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		suffix := args[0]
		if err := util.SetDefaultSuffix(suffix); err != nil {
			fmt.Printf("设置默认后缀失败: %v\n", err)
			return
		}
		fmt.Printf("成功设置默认后缀为: %s\n", suffix)
	},
}

var getDefaultInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取默认的标签前缀和后缀",
	Long:  "获取默认的标签前缀和后缀信息",
	Run: func(cmd *cobra.Command, args []string) {
		prefix, err := util.GetDefaultPrefix()
		if err != nil {
			fmt.Printf("获取默认前缀失败: %v\n", err)
			return
		}
		suffix, err := util.GetDefaultSuffix()
		if err != nil {
			fmt.Printf("获取默认后缀失败: %v\n", err)
			return
		}
		fmt.Printf("当前默认前缀为: %s\n当前默认后缀为: %s\n", prefix, suffix)
	},
}

func init() {
	rootCmd.AddCommand(setDefaultPrefixCmd)
	rootCmd.AddCommand(setDefaultSuffixCmd)
	rootCmd.AddCommand(getDefaultInfoCmd)
}
