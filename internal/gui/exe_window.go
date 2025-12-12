package gui

import (
	"fmt"
	"image/color"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/nitzlover/UmbraPack/internal/execryptor"
)

func RunExeCryptor() {
	a := app.New()
	a.Settings().SetTheme(newUmbraTheme())
	w := a.NewWindow("UmbraPack • EXE Cryptor")
	w.Resize(fyne.NewSize(720, 520))

	var (
		selectedFile string
		lastSaved    string
		iconPath     string
		meta         execryptor.Metadata
	)

	title := widget.NewLabelWithStyle("UmbraPack", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabelWithStyle("AES-256 EXE encryption with embedded loader", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	fileLabel := widget.NewLabel("No file selected")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password (min 8 chars)")

	status := newStatusBar("Ready", neutralColor())
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	hideOriginal := widget.NewCheck("Delete original after encryption", nil)
	hideOriginal.SetChecked(false)

	obfuscationCheck := widget.NewCheck("Harden stub (rename identifiers)", nil)
	obfuscationCheck.SetChecked(true)

	companyEntry := widget.NewEntry()
	companyEntry.SetPlaceHolder("CompanyName")
	productEntry := widget.NewEntry()
	productEntry.SetPlaceHolder("ProductName")
	fileDescEntry := widget.NewEntry()
	fileDescEntry.SetPlaceHolder("FileDescription")
	fileVerEntry := widget.NewEntry()
	fileVerEntry.SetPlaceHolder("1.0.0.0")
	productVerEntry := widget.NewEntry()
	productVerEntry.SetPlaceHolder("1.0.0.0")

	iconLabel := widget.NewLabel("No icon selected")
	iconBtn := widget.NewButton("Select .ico", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			iconPath = uc.URI().Path()
			iconLabel.SetText("Icon: " + filepath.Base(iconPath))
			uc.Close()
		}, w)
	})

	selectBtn := widget.NewButton("Select EXE", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			selectedFile = uc.URI().Path()
			fileLabel.SetText("Selected: " + filepath.Base(selectedFile))
			uc.Close()
		}, w)
	})

	infoBtn := widget.NewButton("About", func() {
		info := "EXE Cryptor v1.1\n\n" +
			"• AES-256 encryption of EXE files\n" +
			"• Embedded protected loader\n" +
			"• Optional deletion of original\n\n" +
			"⚠️ Use only for legitimate software"
		dialog.ShowInformation("About", info, w)
	})

	howToBtn := widget.NewButton("How to use", func() {
		howTo := "1) Select an EXE file\n" +
			"2) Enter a strong password (≥8 chars)\n" +
			"3) Optionally enable deleting the original\n" +
			"4) Click \"Encrypt EXE\"\n" +
			"5) Wait for *_crypted.exe"
		dialog.ShowInformation("How to use", howTo, w)
	})

	var encryptBtn *widget.Button
	encryptBtn = widget.NewButton("Encrypt EXE", func() {
		if selectedFile == "" {
			dialog.ShowError(fmt.Errorf("select an EXE file"), w)
			return
		}
		if len(passwordEntry.Text) < 8 {
			dialog.ShowError(fmt.Errorf("password must be at least 8 characters"), w)
			return
		}

		status.SetInfo("Encrypting...")
		progressBar.Show()
		encryptBtn.Disable()

		go func() {
			cryptor := execryptor.NewCryptor(passwordEntry.Text)

			encrypted, err := cryptor.EncryptFile(selectedFile)
			if err != nil {
				runOnMain(func() {
					status.SetText("❌ Error: " + err.Error())
					progressBar.Hide()
					encryptBtn.Enable()
				})
				return
			}

			runOnMain(func() { status.SetText("Building protected EXE...") })

			dir := filepath.Dir(selectedFile)
			baseName := filepath.Base(selectedFile)
			ext := filepath.Ext(baseName)
			nameWithoutExt := baseName[:len(baseName)-len(ext)]
			outputPath := filepath.Join(dir, nameWithoutExt+"_crypted.exe")

			meta = execryptor.Metadata{
				CompanyName:     companyEntry.Text,
				ProductName:     productEntry.Text,
				FileDescription: fileDescEntry.Text,
				FileVersion:     fileVerEntry.Text,
				ProductVersion:  productVerEntry.Text,
			}

			err = cryptor.CreateStub(encrypted, outputPath, execryptor.BuildOptions{
				Metadata:           meta,
				IconPath:           iconPath,
				EnableObfuscation:  obfuscationCheck.Checked,
				KeepOriginalBinary: !hideOriginal.Checked,
			})
			if err != nil {
				runOnMain(func() {
					status.SetError("Build error: " + err.Error())
					progressBar.Hide()
					encryptBtn.Enable()
				})
				return
			}

			if hideOriginal.Checked {
				osRemoveSafe(selectedFile)
			}

			lastSaved = outputPath
			runOnMain(func() {
				progressBar.Hide()
				encryptBtn.Enable()
				status.SetSuccess("Done. Saved: " + filepath.Base(outputPath))
				dialog.ShowInformation("Success",
					fmt.Sprintf("Encrypted file created:\n%s\n\nData size: %d bytes",
						outputPath, len(encrypted)), w)
			})
		}()
	})

	ghURL, _ := url.Parse("https://github.com/nitzlover")
	brandRow := container.NewHBox(
		title,
		layout.NewSpacer(),
		widget.NewLabel("by nitz"),
		widget.NewHyperlink("GitHub", ghURL),
	)

	header := container.NewVBox(
		brandRow,
		subtitle,
		widget.NewSeparator(),
	)

	metaGrid := container.NewGridWithColumns(2,
		container.NewVBox(
			widget.NewLabelWithStyle("CompanyName", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			companyEntry,
			widget.NewLabelWithStyle("ProductName", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			productEntry,
			widget.NewLabelWithStyle("FileDescription", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			fileDescEntry,
		),
		container.NewVBox(
			widget.NewLabelWithStyle("FileVersion", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			fileVerEntry,
			widget.NewLabelWithStyle("ProductVersion", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			productVerEntry,
			widget.NewLabelWithStyle("Icon (.ico)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			iconLabel,
			iconBtn,
		),
	)

	leftCol := container.NewVBox(
		widget.NewLabelWithStyle("File", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		fileLabel,
		selectBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Options", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		passwordEntry,
		hideOriginal,
		obfuscationCheck,
	)

	rightCol := container.NewVBox(
		widget.NewLabelWithStyle("PE metadata and icon", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		metaGrid,
	)

	form := container.NewGridWithColumns(2, leftCol, rightCol)

	actions := container.NewGridWithColumns(3, encryptBtn, infoBtn, howToBtn)

	statusBar := container.NewHBox(
		container.NewMax(status.bg, container.NewPadded(status.label)),
		layout.NewSpacer(),
		progressBar,
	)

	var lastSavedLabel *widget.Label
	lastSavedLabel = widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})
	updateFooter := func() {
		if lastSaved != "" {
			lastSavedLabel.SetText("Last output: " + filepath.Base(lastSaved) + " • " + time.Now().Format("15:04:05"))
		}
	}

	content := container.NewBorder(
		header,
		container.NewVBox(widget.NewSeparator(), statusBar, lastSavedLabel),
		nil,
		nil,
		container.NewVBox(form, widget.NewSeparator(), actions),
	)

	w.SetContent(content)
	w.SetOnClosed(func() { updateFooter() })
	w.ShowAndRun()
}

func osRemoveSafe(path string) {
	_ = os.Remove(path)
}

func runOnMain(fn func()) {
	if fn == nil {
		return
	}
	fyne.Do(fn)
}

type statusBarState struct {
	label *widget.Label
	bg    *canvas.Rectangle
}

func newStatusBar(text string, c color.NRGBA) *statusBarState {
	lbl := widget.NewLabel(text)
	bg := canvas.NewRectangle(c)
	bg.SetMinSize(fyne.NewSize(0, lbl.MinSize().Height+8))
	return &statusBarState{label: lbl, bg: bg}
}

func (s *statusBarState) SetInfo(text string)    { s.set(text, infoColor()) }
func (s *statusBarState) SetSuccess(text string) { s.set(text, successColor()) }
func (s *statusBarState) SetError(text string)   { s.set(text, errorColor()) }
func (s *statusBarState) SetText(text string)    { s.set(text, neutralColor()) }

func (s *statusBarState) set(text string, c color.NRGBA) {
	s.label.SetText(text)
	s.bg.FillColor = c
	s.bg.Refresh()
	s.label.Refresh()
}

func neutralColor() color.NRGBA { return color.NRGBA{R: 35, G: 38, B: 46, A: 255} }
func infoColor() color.NRGBA    { return color.NRGBA{R: 70, G: 90, B: 150, A: 255} }
func successColor() color.NRGBA { return color.NRGBA{R: 55, G: 145, B: 95, A: 255} }
func errorColor() color.NRGBA   { return color.NRGBA{R: 150, G: 60, B: 70, A: 255} }
