package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type LoggerAdapter struct {
	context context.Context
}

func NewLoggerAdapter(context context.Context) *LoggerAdapter {
	return &LoggerAdapter{context: context}
}

func (logger *LoggerAdapter) Info(message string) {
	tflog.Info(logger.context, message)
}
func (logger *LoggerAdapter) Error(message string) {
	tflog.Error(logger.context, message)
}
func (logger *LoggerAdapter) Debug(message string) {
	tflog.Debug(logger.context, message)
}
