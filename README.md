##### Camille's GO Encryption Tool
# (Currently trying to figure out a better name)

It's what it says on the tin, encryption tool written in go.
Allows you to create an AES key of 128, 256 or 512 bits, and use that key to encrypt/decrypt files.
Uses goroutines to parallel-process files, in order to enhance performance.

QR Code for the download of the Windows version (.exe)

![Windows QR Code](./stripped-encryption-w.png)

QR Code for the download of the Linux version (ELF)

![Linux QR Code](./stripped-encryption-l.png)

Checksums : 
93e3b3ac3999320f82bed139a785eea8ed6bf6fe6ccbe3794f26bbd759dfef18  stripped-encryption.go.forget
c4155ba015a052f0338f0b18c31d733c53f56b35f4c6322dbc6e3fadc00c261d  stripped-encryption.go_linux
654784c1bf8e31d13ca6ce86cfa0b39b803320dd2232b43008b7b0de56997ec4  stripped-encryption.go_windows.exe
