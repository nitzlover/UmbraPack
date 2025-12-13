# UmbraPack 1.0 (EXE Cryptor)

Windows desktop utility (Fyne UI) that encrypts executables with AES-256 and emits a small loader stub that decrypts and runs the payload. It can embed PE metadata and a custom icon. The stub supports optional identifier obfuscation to vary static signatures.

## Features
- AES-256 encryption of an input `.exe`.
- Stub builder with optional identifier obfuscation (“harden stub”).
- PE metadata fields (Company/Product/Description/FileVersion/ProductVersion).
- Custom `.ico` embedding via goversioninfo.
- Optional deletion of the original file after encryption.

## Quick Start
1. Run `cmd/obfuscator/main.go` (or the built binary).
2. Select an EXE and set a password (≥8 chars).
3. Optionally fill PE metadata and select an `.ico`.
4. Toggle:
   - Harden stub (identifier renaming).
   - Delete original after encryption.
5. Click “Encrypt EXE”. Output: `<name>_crypted.exe` next to the source.

## Build
```bash
go build -o umbrapack.exe ./cmd/obfuscator
```

## Configuration Details
- **Stub hardening**: renames identifiers in the generated loader to reduce static signature reuse. Turn off if you want a “cleaner” binary.
- **PE metadata**: populate fields to make the binary look more legitimate (paired with code signing, this greatly helps reputation).
- **Icon**: select a `.ico` to embed into the loader.
- **Delete original**: if enabled, the source EXE is removed after encryption (off by default).

## Detection & FUD Status
- **MetaDefender Score**: 2/22 - Near-FUD status with minimal AV detections.
- **Windows Defender**: ✅ Not detected
- **SmartScreen**: ✅ Not blocked

## Operational Notes / AV Hygiene
- Code-sign the resulting `_crypted.exe` with a valid certificate (OV/EV) to further reduce detections and achieve FUD status.
- For maximum undetectability, enable stub hardening, provide PE metadata, embed a custom icon, and sign the output.
- The current stub writes a decrypted copy and runs it. If you need “save-only” (no auto-run), adjust the stub behavior accordingly.

## Troubleshooting
- Build errors from the UI will show the `go build` stderr/stdout. Most common causes:
  - Invalid `.ico` or missing PE fields format.
  - Go toolchain not on PATH.
  - AV blocking temporary build directory.
- If AV flags persist, try:
  - Disable stub hardening.
  - Use a consistent build path (not a random temp) and sign the binary.
  - Provide PE metadata and icon.

## Tech Stack
- Go + Fyne for UI.
- goversioninfo for PE resources (version info + icon).
- AES-256 CFB for payload encryption.

## CLI Entry
`cmd/obfuscator/main.go` launches the Fyne UI and wires to `internal/gui` + `internal/execryptor`.

## Legal
Use only for legitimate software protection. You are responsible for complying with applicable laws and licenses.

