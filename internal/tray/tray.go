package tray

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KusoKaihatsuSha/tray_helper/internal/config"
	"github.com/KusoKaihatsuSha/tray_helper/internal/emul"
	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"
	"github.com/KusoKaihatsuSha/tray_helper/internal/notification"

	"fyne.io/systray"
)

var (
	This Tray
)

// MenuItem tray menu item
type MenuItem struct {
	*systray.MenuItem
	Cancel        context.CancelFunc
	Kill          func()
	EmulateAction string        `json:"actions"`
	Title         string        `json:"title"`
	Timer         time.Duration `json:"timer"`
	Repeat        int           `json:"repeat"`
	Silent        bool          `json:"silent"`
}

// Tray items
type Tray struct {
	Items []*MenuItem
}

// fill - add menu items from a file with settings.
func (t *Tray) fill() {
	for key, val := range config.SettingsFile(config.Get().Config) {
		element := val.(map[string]any)
		var tmp MenuItem
		tmp.Title = key
		tmp.Silent = element[config.Silent].(bool)
		tmp.Repeat = int(element[config.Repeat].(float64))
		if timer, err := time.ParseDuration(element[config.WaitUntilClose].(string)); err == nil {
			tmp.Timer = timer
		}
		tmp.EmulateAction = element[config.EmulateAction].(string)
		tmp.MenuItem = systray.AddMenuItem(tmp.Title, tmp.Title)
		go tmp.Run()
		t.Items = append(t.Items, &tmp)
	}
}

// fillSetting - add menu items: 'settings' and 'close'.
func (t *Tray) fillSetting() {
	if ico, err := config.EmbedFiles.ReadFile(config.IcoSettings); err == nil {
		tmp := new(MenuItem)
		tmp.Title = "Settings"
		tmp.MenuItem = systray.AddMenuItem(tmp.Title, tmp.Title)
		tmp.MenuItem.SetIcon(ico)
		go tmp.Settings()
		t.Items = append(t.Items, tmp)
	}

	if ico, err := config.EmbedFiles.ReadFile(config.IcoClose); err == nil {
		tmp := new(MenuItem)
		tmp.Title = "Close"
		tmp.MenuItem = systray.AddMenuItem(tmp.Title, tmp.Title)
		tmp.MenuItem.SetIcon(ico)
		go t.Exit(tmp.MenuItem)
		t.Items = append(t.Items, tmp)
	}
}

// do - Performing actions. Sending notifications.
func (m *MenuItem) do(ctx context.Context, cancel context.CancelFunc) {
	if m.EmulateAction == "" {
		return
	}
	notification.Toast(
		ctx,
		m.Title,
		fmt.Sprintf(
			"was ran %d times",
			emul.Handle(
				ctx,
				cancel,
				m.EmulateAction,
				m.Repeat,
			),
		),
		m.Silent,
	)
	if m.Timer == 0 {
		m.Kill()
	}
}

// exec - Performing actions by timer. Sending notifications.
func (m *MenuItem) exec(ctx context.Context, cancel context.CancelFunc, ctime time.Time) {

	count := strings.Count(m.EmulateAction, "EXEC@") + strings.Count(m.EmulateAction, "EXECSTD@")
	count *= m.Repeat

	// Check exist command
	defer func() {
		// Process complete at all
		m.MenuItem.SetTitle(m.Title)
		if m.Timer != 0 {
			m.Kill()
			notification.Toast(
				ctx,
				fmt.Sprintf("Timer %s for %s", m.Timer.String(), m.Title),
				fmt.Sprintf("open exec[%d] will be closed (if not OS protected)", count),
				m.Silent,
			)
		}
		m.MenuItem.Uncheck() // return uncheck
	}()

	m.MenuItem.Check()
	m.Kill = func() {
		m.Cancel()
	}

	go m.do(ctx, cancel)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(300 * time.Millisecond):
			// DO NOTHING
		}
		if m.Timer != 0 {
			if time.Since(ctime) > 0 {
				break
			}
			m.MenuItem.SetTitle(m.Title + "(" + time.Since(ctime).Round(time.Second).String() + ")") // change default title
		}
	}
}

// Run the action
func (m *MenuItem) Run() {
	for range m.MenuItem.ClickedCh {
		ctx, cancel := context.WithCancel(context.Background())
		if m.Cancel != nil {
			m.Kill()
		}
		m.Cancel = cancel
		if !m.MenuItem.Checked() {
			go m.exec(ctx, cancel, time.Now().Add(m.Timer))
		}
	}
}

// Exit the app
func (t *Tray) Exit(mi *systray.MenuItem) {
	for range mi.ClickedCh {
		for _, mi := range t.Items {
			if mi.Kill != nil {
				mi.Kill()
			}
		}
		systray.Quit()
	}
}

// Settings - open the settings
func (m *MenuItem) Settings() {
	for range m.MenuItem.ClickedCh {
		url := "http://" + config.Get().Address + "/settings"
		helpers.OpenUrl(url)
	}
}

// OnReady - start app
func (t *Tray) OnReady() {
	if b, err := config.EmbedFiles.ReadFile(config.IcoApp); err == nil {
		systray.SetIcon(b)
		systray.SetTitle(config.Title)
		systray.SetTooltip(config.ToolTip)
		t.fill()
		t.fillSetting()
	}
}

// Update - update app menu
func (t *Tray) Update() {
	for _, menu := range t.Items {
		if menu.MenuItem != nil {
			if menu.Kill != nil {
				menu.Kill()
			}
			close(menu.MenuItem.ClickedCh)
			menu.MenuItem.Remove()
			menu.MenuItem = nil
		}
	}
	t.fill()
	t.fillSetting()
}
