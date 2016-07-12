package logging

import (
	"fmt"
	"os"
)

func PrError(format string, v ...interface{}) {
	fmt.Printf("==> Error: %s\n", fmt.Sprintf(format, v...))
	os.Exit(1)
}

func PrFatal(format string, v ...interface{}) {
	str := fmt.Sprintf("\n===> Fatal: %s\n", fmt.Sprintf(format, v...))
	panic(str)
}

func PrInfo(format string, v ...interface{}) {
	fmt.Printf("===> %s\n", fmt.Sprintf(format, v...))
}

func PrInfoBegin(format string, v ...interface{}) {
	fmt.Printf("===> %s", fmt.Sprintf(format, v...))
}

func PrInfoEnd(format string, v ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, v...))
}

func PrVerbose(verbose bool, format string, v ...interface{}) {
	if verbose {
		fmt.Printf("==> %s\n", fmt.Sprintf(format, v...))
	}
}
