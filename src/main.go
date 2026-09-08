package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ItemType string

const (
	TypeText  ItemType = "text"
	TypeLink  ItemType = "link"
	TypeImage ItemType = "image"
)

type ClipItem struct {
	ID       string   `json:"id"`
	Type     ItemType `json:"type"`
	Content  string   `json:"content"`
	Pinned   bool     `json:"pinned"`
	Filepath string   `json:"filepath"`
}

type Settings struct {
	ShowPinButton   bool `json:"show_pin_button"`
	ShowDirButton   bool `json:"show_dir_button"`
	CloseOnFocusOut bool `json:"close_on_focus_out"`
	MaxHistoryItems int  `json:"max_history_items"`
	WindowWidth     int  `json:"window_width"`
	WindowHeight    int  `json:"window_height"`
}

var (
	configDir    string
	tmpDir       string
	pinnedDir    string
	historyFile  string
	pinnedFile   string
	settingsFile string
	appSettings  Settings
)

func initDirsAndConfig() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".config", "gopy")
	pinnedDir = filepath.Join(configDir, "pinned")
	tmpDir = filepath.Join(os.TempDir(), "gopy_clipboard")

	os.MkdirAll(configDir, 0755)
	os.MkdirAll(pinnedDir, 0755)
	os.MkdirAll(tmpDir, 0755)

	historyFile = filepath.Join(configDir, "history.json")
	pinnedFile = filepath.Join(configDir, "pinned.json")
	settingsFile = filepath.Join(configDir, "settings.json")

	loadSettings()
}

func loadSettings() {
	data, err := os.ReadFile(settingsFile)
	if err == nil {
		json.Unmarshal(data, &appSettings)
	}
	
	// Applica defaults para campos vazios/zerados
	needsSave := false
	if appSettings.MaxHistoryItems == 0 {
		appSettings.MaxHistoryItems = 50
		needsSave = true
	}
	if appSettings.WindowWidth == 0 {
		appSettings.WindowWidth = 350
		needsSave = true
	}
	if appSettings.WindowHeight == 0 {
		appSettings.WindowHeight = 450
		needsSave = true
	}
	if err != nil {
		appSettings.ShowPinButton = true
		appSettings.ShowDirButton = true
		appSettings.CloseOnFocusOut = true
		needsSave = true
	}

	if needsSave {
		saveSettings()
	}
}

func saveSettings() {
	data, _ := json.MarshalIndent(appSettings, "", "  ")
	os.WriteFile(settingsFile, data, 0644)
}

func main() {
	initDirsAndConfig()

	if len(os.Args) > 1 && os.Args[1] == "--ui" {
		showUI()
		return
	}
	runDaemon()
}

func runDaemon() {
	gtk.Init(nil)
	fmt.Println("Daemon do Gopy Clipboard iniciado...")

	// Limpa o histórico de cópias ao iniciar (mantém os itens apenas durante a sessão/máquina ligada)
	saveItems(historyFile, []ClipItem{})

	clip, _ := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	var lastText string
	var lastImageLen int

	glib.TimeoutAdd(uint(1000), func() bool {
		if clip.WaitIsImageAvailable() {
			image, _ := clip.WaitForImage()
			if image != nil {
				length := image.GetByteLength()
				if length != lastImageLen {
					lastImageLen = length
					lastText = ""
					saveImageToHistory(image)
				}
				return true // Se pegou imagem, ignora o texto (que pode ser apenas o nome do arquivo)
			}
		}
		
		if clip.WaitIsTextAvailable() {
			text, _ := clip.WaitForText()
			if text != "" && text != lastText {
				
				// Trata arquivos locais copiados pelo gerenciador de arquivos (Nemo/Nautilus)
				if strings.HasPrefix(text, "file://") {
					path := strings.TrimSpace(strings.TrimPrefix(text, "file://"))
					// Decode simples de URL (gambiarra rápida pro MVP)
					path = strings.ReplaceAll(path, "%20", " ")
					lowerPath := strings.ToLower(path)
					if strings.HasSuffix(lowerPath, ".png") || strings.HasSuffix(lowerPath, ".jpg") || strings.HasSuffix(lowerPath, ".jpeg") {
						pixbuf, err := gdk.PixbufNewFromFile(path)
						if err == nil {
							length := pixbuf.GetByteLength()
							if length != lastImageLen {
								lastImageLen = length
								lastText = text
								saveImageToHistory(pixbuf)
							}
							return true // Salvou a imagem local, ignora o texto
						}
					}
				}

				lastText = text
				lastImageLen = 0
				saveToHistory(text)
			}
		}
		
		return true
	})

	gtk.Main()
}

func saveToHistory(text string) {
	items := loadItems(historyFile)
	if len(items) > 0 && items[0].Content == text {
		return
	}

	itemType := TypeText
	if strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://") {
		itemType = TypeLink
	}

	newItem := ClipItem{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:     itemType,
		Content:  text,
		Pinned:   false,
		Filepath: tmpDir,
	}

	items = append([]ClipItem{newItem}, items...)
	if len(items) > appSettings.MaxHistoryItems {
		items = items[:appSettings.MaxHistoryItems]
	}
	saveItems(historyFile, items)
}

func saveImageToHistory(image *gdk.Pixbuf) {
	filename := fmt.Sprintf("img_%d.png", time.Now().UnixNano())
	fullPath := filepath.Join(tmpDir, filename)
	image.SavePNG(fullPath, 9)

	newItem := ClipItem{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:     TypeImage,
		Content:  filename,
		Pinned:   false,
		Filepath: tmpDir,
	}

	items := loadItems(historyFile)
	items = append([]ClipItem{newItem}, items...)
	if len(items) > appSettings.MaxHistoryItems {
		items = items[:appSettings.MaxHistoryItems]
	}
	saveItems(historyFile, items)
}

func loadItems(file string) []ClipItem {
	var items []ClipItem
	data, err := os.ReadFile(file)
	if err == nil {
		json.Unmarshal(data, &items)
	}
	
	// Ordena do mais recente para o mais antigo usando o ID (UnixNano)
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID > items[j].ID
	})
	
	return items
}

func saveItems(file string, items []ClipItem) {
	data, _ := json.MarshalIndent(items, "", "  ")
	os.WriteFile(file, data, 0644)
}

// ==========================================
// COMPONENTES DE UI
// ==========================================

func showUI() {
	gtk.Init(nil)

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Gopy")
	win.SetDefaultSize(appSettings.WindowWidth, appSettings.WindowHeight)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDecorated(false)
	win.SetSkipTaskbarHint(true)
	win.SetKeepAbove(true)

	if appSettings.CloseOnFocusOut {
		win.Connect("focus-out-event", func() {
			win.Destroy()
		})
	}
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scroll.SetVExpand(true)

	listBox, _ := gtk.ListBoxNew()
	
	pinnedItems := loadItems(pinnedFile)
	historyItems := loadItems(historyFile)
	
	for _, item := range pinnedItems {
		listBox.Add(buildComponent(item, win, true))
	}
	for _, item := range historyItems {
		listBox.Add(buildComponent(item, win, false))
	}

	scroll.Add(listBox)
	box.PackStart(scroll, true, true, 0)
	
	// Botão para limpar histórico (não afeta os fixados)
	btnClear, _ := gtk.ButtonNewWithLabel("🗑 Limpar Histórico")
	btnClear.Connect("clicked", func() {
		saveItems(historyFile, []ClipItem{})
		win.Destroy()
	})
	box.PackEnd(btnClear, false, false, 2)

	win.Add(box)
	win.ShowAll()
	gtk.Main()
}

// Build base structure and route to specific component
func buildComponent(item ClipItem, win *gtk.Window, isPinned bool) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	
	switch item.Type {
	case TypeText, TypeLink:
		vbox.Add(buildTextComponent(item, win, isPinned))
	case TypeImage:
		vbox.Add(buildImageComponent(item, win, isPinned))
	}
	
	row.Add(vbox)
	return row
}

// Component for Text and Links: Pin, Copy, Scroll (Expander)
func buildTextComponent(item ClipItem, win *gtk.Window, isPinned bool) *gtk.Box {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	actionsBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	// Pin Button
	if appSettings.ShowPinButton {
		actionsBox.PackStart(createPinButton(item, isPinned), false, false, 0)
	}
	container.PackStart(actionsBox, false, false, 0)

	// Copy EventBox
	copyBox, _ := gtk.EventBoxNew()
	copyBox.SetVisibleWindow(false)
	copyBox.Connect("button-press-event", func() {
		clip, _ := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
		clip.SetText(item.Content)
		win.Destroy()
	})

	// Content Label / Scroll (Expander)
	if item.Type == TypeLink {
		shortLink := strings.TrimPrefix(item.Content, "https://")
		shortLink = strings.TrimPrefix(shortLink, "http://")
		lbl, _ := gtk.LabelNew(shortLink)
		lbl.SetMarkup("<u><span foreground='#3498db'>" + shortLink + "</span></u>")
		lbl.SetXAlign(0)
		copyBox.Add(lbl)
		container.PackStart(copyBox, true, true, 2)
	} else {
		// TypeText
		lbl, _ := gtk.LabelNew(item.Content)
		lbl.SetLineWrap(true)
		lbl.SetXAlign(0)
		copyBox.Add(lbl)

		if len(item.Content) > 60 {
			summary := strings.ReplaceAll(item.Content, "\n", " ")
			if len(summary) > 40 {
				// Cuidado com strings multibyte no Go, mas como é MVP:
				summary = string([]rune(summary)[:37]) + "..."
			}
			expander, _ := gtk.ExpanderNew(summary)
			expander.Add(copyBox)
			container.PackStart(expander, true, true, 2)
		} else {
			container.PackStart(copyBox, true, true, 2)
		}
	}

	return container
}

// Component for Image: Pin, Path, Copy, OpenWindow
func buildImageComponent(item ClipItem, win *gtk.Window, isPinned bool) *gtk.Box {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	actionsBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	// Pin Button
	if appSettings.ShowPinButton {
		actionsBox.PackStart(createPinButton(item, isPinned), false, false, 0)
	}
	
	// Path / Dir Button
	if appSettings.ShowDirButton {
		btnDir, _ := gtk.ButtonNewWithLabel("◨")
		btnDir.SetRelief(gtk.RELIEF_NORMAL)
		btnDir.Connect("clicked", func() {
			exec.Command("xdg-open", item.Filepath).Start()
		})
		actionsBox.PackEnd(btnDir, false, false, 0)
	}
	
	// Open Window Button (Edit/View)
	btnOpen, _ := gtk.ButtonNewWithLabel("👁") // View image
	btnOpen.SetRelief(gtk.RELIEF_NORMAL)
	btnOpen.Connect("clicked", func() {
		exec.Command("xdg-open", filepath.Join(item.Filepath, item.Content)).Start()
	})
	actionsBox.PackEnd(btnOpen, false, false, 0)

	container.PackStart(actionsBox, false, false, 0)

	// Copy EventBox
	copyBox, _ := gtk.EventBoxNew()
	copyBox.SetVisibleWindow(false)
	copyBox.Connect("button-press-event", func() {
		clip, _ := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
		pix, _ := gdk.PixbufNewFromFile(filepath.Join(item.Filepath, item.Content))
		if pix != nil {
			clip.SetImage(pix)
		}
		win.Destroy()
	})
	
	// Carrega miniatura da Imagem
	pixbuf, err := gdk.PixbufNewFromFileAtScale(filepath.Join(item.Filepath, item.Content), 180, -1, true)
	if err == nil {
		img, _ := gtk.ImageNewFromPixbuf(pixbuf)
		copyBox.Add(img)
	} else {
		lbl, _ := gtk.LabelNew("[Erro ao carregar Imagem]")
		copyBox.Add(lbl)
	}

	container.PackStart(copyBox, true, true, 2)

	return container
}

func createPinButton(item ClipItem, initialPinned bool) *gtk.Button {
	isPinned := initialPinned
	pinIcon := "☆"
	if isPinned {
		pinIcon = "★"
	}
	btnPin, _ := gtk.ButtonNewWithLabel(pinIcon)
	btnPin.SetRelief(gtk.RELIEF_NORMAL)
	btnPin.Connect("clicked", func() {
		togglePin(item, isPinned)
		if isPinned {
			btnPin.SetLabel("☆")
			isPinned = false
		} else {
			btnPin.SetLabel("★")
			isPinned = true
		}
	})
	return btnPin
}

func togglePin(item ClipItem, currentlyPinned bool) {
	history := loadItems(historyFile)
	pinned := loadItems(pinnedFile)

	if currentlyPinned {
		newPinned := []ClipItem{}
		for _, p := range pinned {
			if p.ID != item.ID {
				newPinned = append(newPinned, p)
			}
		}
		saveItems(pinnedFile, newPinned)
		
		if item.Type == TypeImage {
			oldPath := filepath.Join(pinnedDir, item.Content)
			newPath := filepath.Join(tmpDir, item.Content)
			moveFile(oldPath, newPath)
			item.Filepath = tmpDir
		}

		item.Pinned = false
		history = append([]ClipItem{item}, history...)
		saveItems(historyFile, history)
	} else {
		if item.Type == TypeImage {
			oldPath := filepath.Join(tmpDir, item.Content)
			newPath := filepath.Join(pinnedDir, item.Content)
			moveFile(oldPath, newPath)
			item.Filepath = pinnedDir
		}
		
		item.Pinned = true
		pinned = append([]ClipItem{item}, pinned...)
		saveItems(pinnedFile, pinned)
		
		newHistory := []ClipItem{}
		for _, h := range history {
			if h.ID != item.ID {
				newHistory = append(newHistory, h)
			}
		}
		saveItems(historyFile, newHistory)
	}
}

// moveFile contorna problemas de cross-device link (ex: /tmp em tmpfs e ~/.config no SSD)
func moveFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}
	return os.Remove(src)
}
