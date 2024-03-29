package studio

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/yaoapp/gou"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/utils"
	"github.com/yaoapp/yao/config"
	"github.com/yaoapp/yao/engine"
	"github.com/yaoapp/yao/share"
)

// RunCmd command
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: L("Execute Yao Studio Script"),
	Long:  L("Execute Yao Studio Script"),
	Run: func(cmd *cobra.Command, args []string) {
		defer share.SessionStop()
		defer gou.KillPlugins()
		defer func() {
			err := exception.Catch(recover())
			if err != nil {
				fmt.Println(color.RedString(L("Fatal: %s"), err.Error()))
			}
		}()

		Boot()
		cfg := config.Conf
		cfg.Session.IsCLI = true
		engine.Load(cfg)
		if len(args) < 1 {
			fmt.Println(color.RedString(L("Not enough arguments")))
			fmt.Println(color.WhiteString(share.BUILDNAME + " help"))
			return
		}

		name := strings.Split(args[0], ".")
		service := strings.Join(name[0:len(name)-1], ".")
		method := name[len(name)-1]

		fmt.Println(color.GreenString(L("Studio Run: %s"), args[0]))
		pargs := []interface{}{}
		for i, arg := range args {
			if i == 0 {
				continue
			}

			if strings.HasPrefix(arg, "::") {
				arg := strings.TrimPrefix(arg, "::")
				var v interface{}
				err := jsoniter.Unmarshal([]byte(arg), &v)
				if err != nil {
					fmt.Println(color.RedString(L("Arguments: %s"), err.Error()))
					return
				}
				pargs = append(pargs, v)
				fmt.Println(color.WhiteString("args[%d]: %s", i-1, arg))
			} else if strings.HasPrefix(arg, "\\::") {
				arg := "::" + strings.TrimPrefix(arg, "\\::")
				pargs = append(pargs, arg)
				fmt.Println(color.WhiteString("args[%d]: %s", i-1, arg))
			} else {
				pargs = append(pargs, arg)
				fmt.Println(color.WhiteString("args[%d]: %s", i-1, arg))
			}
		}

		req := gou.Yao.New(service, method)
		res, err := req.RootCall(pargs...)
		if err != nil {
			fmt.Println(color.RedString("--------------------------------------"))
			fmt.Println(color.RedString(L("%s Error"), args[0]))
			fmt.Println(color.RedString("--------------------------------------"))
			utils.Dump(err)
			fmt.Println(color.RedString("--------------------------------------"))
			fmt.Println(color.GreenString(L("✨DONE✨")))
			return
		}

		fmt.Println(color.WhiteString("--------------------------------------"))
		fmt.Println(color.WhiteString(L("%s Response"), args[0]))
		fmt.Println(color.WhiteString("--------------------------------------"))
		utils.Dump(res)
		fmt.Println(color.WhiteString("--------------------------------------"))
		fmt.Println(color.GreenString(L("✨DONE✨")))
	},
}
