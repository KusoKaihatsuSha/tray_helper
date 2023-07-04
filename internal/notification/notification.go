package notification

import (
	"context"

	"github.com/KusoKaihatsuSha/tray_helper/internal/config"
	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"

	"gopkg.in/toast.v1"
)

// Toast send toast
func Toast(ctx context.Context, title string, message string, silent bool) {
	if silent {
		return
	}

	select {
	case <-ctx.Done():
		title += " >>> timer end"
	default:
		title += " >>> complete"
	}
	notification := toast.Notification{
		AppID:   config.Title,
		Title:   title,
		Message: message,
		Icon:    config.Get().Ico,
	}
	if err := notification.Push(); err != nil {
		helpers.ToLog(err)
	}
}
