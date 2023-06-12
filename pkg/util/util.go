package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func CreateNamedLogger(module ...string) *logrus.Entry {
	if len(module) > 0 {
		return logrus.WithField("module", module[0])
	}
	return logrus.WithField("module", nil)
}

func Prompt(prompt, fallback string) string {
	var input string

	if prompt != "" {
		input = fallback
	}

	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}
