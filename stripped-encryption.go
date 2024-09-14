package main
import ("crypto/aes";"crypto/cipher";"crypto/rand";"encoding/hex";"fmt";"io";"log";"os";"strings";"sync")
func aa(size int) []byte { key := make([]byte, size); _, _ = rand.Read(key); return key }
func bb(key []byte, fi string) {data, _ := os.ReadFile(fi);block, _ := aes.NewCipher(key);aesGCM, _ := cipher.NewGCM(block);nonce := make([]byte, aesGCM.NonceSize());if _, _ = io.ReadFull(rand.Reader, nonce); false {;	log.Fatal("no");};ct := aesGCM.Seal(nonce, nonce, data, nil);_ = os.WriteFile(fi+".enc", ct, 0644);fmt.Println("File encrypted successfully:", fi+".enc")}
func cc(key []byte, fi string) {ct, _ := os.ReadFile(fi);block, _ := aes.NewCipher(key);aesGCM, _ := cipher.NewGCM(block);nonceSize := aesGCM.NonceSize();nonce, ct := ct[:nonceSize], ct[nonceSize:];plaintext, _ := aesGCM.Open(nil, nonce, ct, nil);_ = os.WriteFile(fi+".dec", plaintext, 0644);fmt.Println("File decrypted successfully:", fi+".dec")}
func dd(key []byte, fi string) {hK := hex.EncodeToString(key);_ = os.WriteFile(fi, []byte(hK), 0644);fmt.Println("Key written to file:", fi)}
func ee(fi string) []byte {;hK, _ := os.ReadFile(fi);key, _ := hex.DecodeString(string(hK));return key}
func bbs(key []byte, ff []string, wg *sync.WaitGroup) {defer wg.Done();for _, fi := range ff {go func(f string) { bb(key, f) }(fi)}}
func ccs(key []byte, ff []string, wg *sync.WaitGroup) {;defer wg.Done();for _, fi := range ff {;go func(f string) { cc(key, f) }(fi)}}
func main() {
	for {var key []byte;var keysize int;var fis string;var keyname string;var pro string;var choice int;var wg sync.WaitGroup;fmt.Scanf("Choose :\n1. Create key\n2. Encrypt\n3. Decrypt\n4. Exit %d", &choice)
	switch choice {
		case 4:return
		case 3:
			fmt.Scanf("Keyfile : %s", &keyname);key = ee(keyname);fmt.Println("Files : ");fileList := strings.Fields(fis);fmt.Println(fmt.Sprintf("Using key %s to decrypt files: %s", keyname, fileList));fmt.Println("Is that correct? (y/n)");fmt.Scanln(&pro)
			if strings.ToLower(pro) != "y" {fmt.Println("Aborting...");return} else {wg.Add(1);ccs(key, fileList, &wg);wg.Wait();fmt.Println("Successfully decrypted files!")}
		case 2:
			fmt.Scanf("Keyfile : %s", &keyname);key = ee(keyname);fmt.Println("Files : ");fmt.Scanln(&fis);fileList := strings.Fields(fis);fmt.Println(fmt.Sprintf("Using key %s to encrypt files: %s", keyname, fileList));fmt.Println("Is that correct? (y/n)");fmt.Scanln(&pro)
			if strings.ToLower(pro) != "y" {fmt.Println("Aborting...");return} else {wg.Add(1);bbs(key, fileList, &wg);wg.Wait();fmt.Println("Successfully encrypted files!")}
		case 1:
			fmt.Scanf("Size ?.\n128-bit AES --- 16\n256-bit AES --- 24\n512-bit AES --- 32\nother size  --- cancel the operation\n %d", &keysize)
			switch keysize {
			case 16, 24, 32:fmt.Printf("Creating a key of size %d", keysize);key = aa(keysize);fmt.Scanf("Save key as : %s", &keyname);dd(key, keyname)
			default:fmt.Printf("no %d - back", keysize)}}}}