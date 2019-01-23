package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/rosewoodmedia/rwcmd/coolcommands"

	"github.com/KernelDeimos/anything-gos/interp_a"
	"github.com/KernelDeimos/gottagofast/toolparse"
)

func main() {
	i := interp_a.InterpreterFactoryA{}.MakeExec()

	fullarg := strings.Join(os.Args[1:], " ")

	input, err := toolparse.ParseListSimple(fullarg)
	if err != nil {
		logrus.Fatal(err)
	}

	coolcommands.InstallFroute(i)

	result, err := i.OpEvaluate(input)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(result)
}
