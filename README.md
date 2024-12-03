# Camille's GO Encryption Tool [^1]

It's what it says on the tin, encryption tool written in go.
Allows you to create an AES key of 128, 256 or 512 bits, and use that key to encrypt/decrypt files.
Uses goroutines to parallel-process files, in order to enhance performance.

# How to use
### 1. Clone the repository and cd into it :
```
git clone https://github.com/TheDarkWolfer/Go-Encrypt/
cd Go-Encrypt
```

### 2. Verify the checksums
```
sha256sum ./*.go
# Should match these results :
# cebb870947208d0534af6659a60802bdac0ea63738eb9fe2d78112aca089e87d  ./encryption.go
# d9a10f38e4e94c5b0264f50df5596899396e9a24e901264ff05ed9e9344db96d  ./stripped-encryption.go

sha256sum ./compile.sh
# Should match this result :
# 44bd08daf23f6fde8c57c161e7b484bc7bcbf840bc4c6a173fdb8d05728f60ae  compile.sh
```
If the checksums match, continue. If they don't, [see this](#checksum-issues)

### 3. Compile the tool
You can either use the `compile.sh` script to compile the tool to both ELF and .exe formats, or do it manually :

`compile.sh` :

You will probably have to add executive permissions to it :
```
chmod +x ./compile.sh
```
    After that, you can either provide it with a specific go file and version number, or let it do it's thing :
    First argument is the target file to compile (default "./encryption.go"), and second argument is the version number (default "42")
```
./compile.sh
```

Manually :

You can compile the tool as follows :
```
go build ./encryption.go # Or ./stripped_encryption.go if you want to try out the stripped out
```

### 4. Use the tool
You can now run the tool and follow the instructions in order to generate keys, encrypt, and decrypt files.
Just be careful, as it is AES encryption : ***You lose the key, you lose the files***. 
The program provides you with a simple UI and the six available options (The first four are the useful ones)

`Create an AES key`             :   Create and save an AES key

`Encrypt files using AES-GCM`   :   Encrypt the given file with the given key

`Decrypt files using AES-GCM`   :   Decrypt the given file with the given key

`Set keyfile to use`            :   Choose a specific keyfile to use


The keys can be generated in the following sizes (Just keep in mind that the three options would still require eons of brute force to break) :

- 128 bits (16 bytes) Fast but 'weakest' of the three

- 256 bits (32 bytes) Medium speed and strenght

- 512 bits (64 bytes) Slowest but 'strongest' of the three



### 5. QR codes and checksums

QR Code for the download of the Windows version (.exe)

![Windows QR Code](./stripped-encryption-w.png)

QR Code for the download of the Linux version (ELF)

![Linux QR Code](./stripped-encryption-l.png)

Checksums : 

```cebb870947208d0534af6659a60802bdac0ea63738eb9fe2d78112aca089e87d```  stripped-encryption.go_linux

```d9a10f38e4e94c5b0264f50df5596899396e9a24e901264ff05ed9e9344db96d```  stripped-encryption.go_windows.exe

### TODO
- [x] Add License
- [ ] Fix issue where large files would crash the program
- [ ] Implement asymetric encryption
- [ ] Improve README.md





### Here is a silly side thing I did with this code
##### /!\ This part of the code won't be maintained as it was made as a challenge, and using this instead of the non-stripped version is a bad idea in practically all scenarios /!\
You'll also find here "stripped-encryption.go", the source code for a stripped version of the AES encryption tool. I managed to get it to be just near 2.9kB by stripping as much code as I could, which resulted in rougher corners, rougher edges, and no error management. It's more of a proof of concept / personal challenge I set myself. The big qr code contains the source code of the stripped version, you should be able to scan it and get the exact contents of "stripped-encryption.go". 

You can see the stripped version in there :

![Stripped Encryption tool QR code](./stripped-encryption.png)

# Original goal
My original goal was to get the compiled version of the stripped tool into a qr code, as I've seen [mattkc](https://mattkc.com/etc/snakeqr/) do. Sadly, I currently lack the expertise and knowledge to do so (In my tests, the resulting executable was a thousand times bigger than it's source code, and by extension a thousand times too big to fit even if the biggest QR code format), especially given that I don't feel comfortable enough to tinker around with compiler settings. 



# CHECKSUM ISSUES
If the checksums for the .go files don't match, you shouldn't compile them before checking a few things :
1. I forgot to update the README.md at the same time as the files
    - If the scripts were modified and the checksums weren't recalculated, that could be the simplest issue, and the reason of the mismatch. If so, you can either wait for me to notice and update (or post an issue about it), or copy the code from the repo after checking if it's all right

2. The README.md and the files were updated at the same time
    - Check if the url you cloned the repo from is the right one (if you cloned from a fork, the checksums may differ between README's)
    - Check if the README.md and scripts are from the same commits (if they are different versions, it may explain the discrepancy)

3. The README.md and files were updated at the same time, and it's the right URL
    - That's where it becomes a bit more problematic. First, I'd like you to open an issue about it in order to let me check it out, to figure out if the issue comes from my end or not. 
    - If the issue persists, there is a low but not zero risk of being the target of a Man-In-The-Middle attack (A malicious actor standing between you and the actual resource in order to tamper with the communications), and you should seek to mitigate this ASAP. Check your DNS settings, your router's and your computer's logs, verify the TLS status of the page, etc...