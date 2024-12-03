package main

//!/ To do :
//!// meaningful return/exit codes
//!// recursive directory handling
//!/ improve error handling
//!// improve UI (Hope ANSI escape codes do the trick)
//!// more comments
//!/ Fix issue where big files just eat up all the RAM

/*
This program is a small test program I wrote to familiarize myself with Go's syntax, workflow, and everything.
I'm not the first one to make this type of tool, but given that I am kinda interested in encryption, I wanted
to try my hands at that, alongside commenting the whole program as thoroughly as possible.
As always when I hyperfixate on something, it got out of hand to a point where this could be a full release lmao

BTW you'll also find a stripped version of this program, that I tried to make as small as possible, with less
features (like, key generation, encryption, and decryption, nothing more) and rougher corners, if it ever comes
in handy to have a small as fuck application for encryption purposes
*/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Informations about the program
var (
	Version   = "unknown"
	BuildDate = "unknown"
)

// Function to get all the files in a given directory, recursively
// Takes in a string as input, returns a slice of strings as output.
func getFilesInDir(path string) []string {
	// Initiate the slice of filenames
	var allFiles []string
	// Read the directory's contents
	files, err := os.ReadDir(path)
	if err != nil {
		//e2u
		log.Fatal(err)
	}

	for _, f := range files {
		// Loop through all files in the "files" slice of fs.DirEntries
		// and append their name (=path) to the slice of strings that will
		// serve later as outputs
		allFiles = append(allFiles, f.Name())
	}

	// Return our slice of filepaths to the user
	return allFiles
}

/*
A function to check if a given path is a directory or a single file,
used in order to handle recursive directory when encrypting/decrypting
Check if the given path is a file or a directory ; returns a specific int
in each case :

	-1 	->	specified file doesn't exist
	0	-> 	error
	1	-> 	directory
	2	-> 	file
*/
func checkPathType(path string) int {
	// Get the file's data and stats
	info, err := os.Stat(path)
	// Specific case for file not found errors
	if os.IsNotExist(err) {
		fmt.Println("Path does not exist.")
		return 0
	}

	// error goes to user, as ususal
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}

	// Check if the file is a directory, and if so, return 1, otherwise return 2.
	if info.IsDir() {
		return 1
	} else {
		return 2
	}
}

// A function to convert bytes to the highest indicator of size
func formatSize(bytes int64) string {

	// ANSI codes to add color here <3
	YELLOW := "\033[33m"
	RESET := "\033[0m"

	// First, we check if the file's size is smaller than 1024 bytes. If so,
	// just return the size in bytes without any conversion
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	// Define the size suffixes in an array (slice of strings here)
	sizes := []string{"KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "WhatTheFuckAreYouDoing"}
	// Convert the size from an integer to a 64-bit float, since we'll be doing divisions
	// that won't always land flat
	size := float64(bytes)

	// Loop through the units and divide the size until it is greater than 1 (avoids situations like 0.54 MB)
	// We declare i beforehand since it'll be used later outside of the for loop
	var i int
	for i = 0; size >= unit && i < len(sizes)-1; i++ {
		size /= unit
	}

	// Return a formatted string with the size of the file (converted to the appropriate unit)
	// and the unit's suffix
	return YELLOW + fmt.Sprintf("%f %s", size, sizes[i-1]) + RESET

}

// AES keys are used here in order to have symetric encryption
// of the data the user provides.

// This functions returns an AES key as a slice of bytes, generated
// from a given size (128, 256 or 512)
func generateAESKey(size int) []byte {
	// type inference of the key variable from a slice of
	// bytes of the given size
	key := make([]byte, size)
	// random number generation in order to generate a strong key
	// with error handling in case of failure
	_, err := rand.Read(key)
	if err != nil {
		// Log the error to the user.
		log.Fatal(err)
	}
	// Return the key after generating it
	return key
}

// This function opens a given file, reads its contents, encrypts said
// contents before writing them in <original>.enc, all thanks
// to the given key as seen above
func encryptFile(key []byte, filename string, debug bool) {

	// Declare the variables needed for the benchmarking
	var formattedFileSize string

	// Get the starting time of the operation in order to perform benchmarking.
	startTime := time.Now()

	// Reading the file using os
	data, err := os.ReadFile(filename)
	if err != nil {
		// Log any errors that occured during the file read
		log.Fatal(err)
	}

	// Get the size of the file using os.Stat()
	fi, err := os.Stat(filename)
	if err != nil {
		//err2usr
		log.Fatal(err)
	}
	// This gives us the file's size in bytes, which we'll convert
	// for the sake of readability
	fileSize := fi.Size()
	formattedFileSize = formatSize(fileSize)

	// Creation of a new AES cipherblock based on the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		// With any errors being displayed to the user
		log.Fatal(err)
	}

	// Generation of an AES-GCM entity from the previously created block
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		// Once again, errors go to the users
		log.Fatal(err)
	}

	// Creation of a nonce (Number used ONCE) for the encryption, with
	// a size dictated by the AES-GCM entity we established
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		// Error handling including the nonce's generation, directed to the user
		log.Fatal(err)
	}

	// Encrypt the file content using all we previously generated
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)

	// Write the encrypted data to a new file
	err = os.WriteFile(filename+".enc", ciphertext, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the delta between the start and end
	// of the encryption process
	timeDelta := time.Since(startTime).Seconds()
	encryptionRate := float64(fileSize) / timeDelta

	// ANSI codes to add color here <3
	YELLOW := "\033[33m"
	RESET := "\033[0m"

	if debug {
		// Display that information to the user
		fmt.Println(fmt.Sprintf(YELLOW+"File size 			: 	%s"+RESET, formattedFileSize))
		fmt.Println(fmt.Sprintf(YELLOW+"Time taken			:	%f seconds"+RESET, timeDelta))
		fmt.Println(fmt.Sprintf(YELLOW+"Encryption rate 	:	%f bytes/second"+RESET, encryptionRate))
	}
	// Tell the user we succeeded, and display the encrypted file's name
	fmt.Println("File encrypted successfully:", filename+".enc")
}

// Function to decrypt a given file using a provided key, and
// writing the decrypted file to <original>.dec
// Take in the key as a slice of bytes and the target file's name as a string
func decryptFile(key []byte, filename string, debug bool) {

	// Declare the variables needed for the benchmarking
	var formattedFileSize string

	// Get the starting time of the operation in order to perform benchmarking.
	startTime := time.Now()

	// Read the file using os
	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Get the size of the file using os.Stat()
	fi, err := os.Stat(filename)
	if err != nil {
		//err2usr
		log.Fatal(err)
	}
	// This gives us the file's size in bytes, which we'll convert
	// for the sake of readability
	fileSize := fi.Size()
	formattedFileSize = formatSize(fileSize)

	// Create a new AES cipher block from the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// Create a GCM (Galois/Counter Mode) for AES form the
	// block we created from the key
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	// Creation of a nonce (Number used ONCE) for the encryption, with
	// a size dictated by the AES-GCM entity we established, based on
	// the size of the nonce dictated by the encrypted data
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Return any errors to the user
		log.Fatal(err)
	}

	// Write the decrypted content to a file with the appropriate
	// permissions (read/write - read - read)
	err = os.WriteFile(filename+".dec", plaintext, 0644)
	if err != nil {
		// errors -> user's eyes
		log.Fatal(err)
	}

	// Calculate the delta between the start and end
	// of the decryption process
	timeDelta := time.Since(startTime).Seconds()
	decryptionRate := float64(fileSize) / timeDelta

	// ANSI codes to add color here <3
	YELLOW := "\033[33m"
	RESET := "\033[0m"

	// Display that information to the user
	if debug {
		fmt.Println(YELLOW+"File size 			: 	%s"+RESET, formattedFileSize)
		fmt.Println(YELLOW+"Time taken			:	%f seconds"+RESET, timeDelta)
		fmt.Println(YELLOW+"Encryption rate 	:	%f bytes/second"+RESET, decryptionRate)
	}
	// Tell the user we succeeded in decrypting the file, alongside the
	// decrypted file's name, in case the user needs it.
	fmt.Println("File decrypted successfully:", filename+".dec")
}

// Function to write a key into a file after converting it into a hexadecimal string
// Takes in the key as a slice of bytes and the filename as a string
func writeKeyToFileHex(key []byte, filename string) {
	// Convert key to hexadecimal string
	hexKey := hex.EncodeToString(key)
	// Write the key to a file (permissions read/write - read - read) and
	// prepare for any errors that may come our way
	err := os.WriteFile(filename, []byte(hexKey), 0644)
	if err != nil {
		// yeet errs 2 usr
		log.Fatal(err)
	}
	// Confirm we wrote the key to the given filename
	fmt.Println("Key written to file:", filename)
}

// Function to read a key from a file, and convert it from a hexadecimal
// string to a slice of bytes. Inverse of the writeKeyToFileHex function
func readKeyFromFileHex(filename string) []byte {
	// Open and read the key from the given filename
	hexKey, err := os.ReadFile(filename)
	if err != nil {
		// err > usr
		log.Fatal(err)
	}

	// Decode the key (that was saved as a hex string) and save it in
	// memory, to use later in other functions
	key, err := hex.DecodeString(string(hexKey))
	if err != nil {
		// you get it by now
		log.Fatal(err)
	}
	// Return the key to whatever function called this one
	return key
}

// Function to encrypt files using a given key, a list of files, and goroutines
// The use of goroutines here is necessary if the user ever needs to encrypt several
// files at once ; async processing of each files felt like the right thing to do.
func encryptFiles(key []byte, filenames []string, wg *sync.WaitGroup, debug bool) {
	// We get ready to set the work group as done once all
	// the files have been processed accordingly
	defer wg.Done()

	// Iterate through all the files and encrypt them using the key provided by the user
	for _, filename := range filenames {
		// Launch a goroutine for each file
		go func(f string) {
			encryptFile(key, f, debug)
		}(filename)
	}
}

// Inverse of the encryptFiles function : decryption using goroutines in order to
// parallell process all the files we were given. Once again, take in a key (slice of bytes)
// and a list of filenames
func decryptFiles(key []byte, filenames []string, wg *sync.WaitGroup, debug bool) {
	// Get ready to call it a day for this function once everything ran right
	defer wg.Done()

	// Iterate through all the files and launch a goroutine for each of them
	// in order to encrypt them
	for _, filename := range filenames {
		// Launch a goroutine for each file
		go func(f string) {
			decryptFile(key, f, debug)
		}(filename)
	}
}

func main() {

	// ANSI escape codes to add a bit of color <3
	RED := "\033[31m"
	GREEN := "\033[32m"
	CYAN := "\033[36m"
	GOLD := "\033[38;5;214m"
	RESET := "\033[0m"
	DIM := "\033[2m"

	// Initialisation of all the variables we'll need, in order to have good structure

	var debug bool // Debug boolean to benchmark encryption.
	debug = false  // Disabled by default
	var debugString string

	var key []byte  // The key we'll be using, a slice of bytes
	var keysize int // The size of the key the user may wish to generate, an integer
	var choice int  // The choice of the user in the main menu of the program, an integer

	var filenames string // Names of the different files we'll have to deal with, a string
	var keyname string   // Name of the file the key is saved in, a string
	var prompt string    // The name of the file we are asking the user for (either keyfile or target files), a string

	var wg sync.WaitGroup // The workgroup for the goroutines the parallel decryption and encryption function use

	// Infinite loop in order to let the user perform several operations if necessary without having
	// to re-launch the program everytime. Can be exited by inputting "4" as a choice, as follows
	for {

		// Display the version alongside the build date if debug is enabled
		if debug {
			fmt.Printf("VERSION\t\t:	%s", Version)
			fmt.Printf("BUILD DATE\t:	%s", BuildDate)
		}

		// A small formatted string to indicate whether or not debug is enabled
		if debug {
			debugString = GREEN + "ON" + RESET
		} else {
			debugString = RED + "OFF" + RESET
		}

		// Prompt the user with different choices of actions to perform. Usage of several print statements instead of only one
		// print statement and \n for the sake of readability and ease of modification
		fmt.Println("Choose the operation you wish to perform : ")
		fmt.Println(CYAN + "1. Create an AES key" + RESET)
		fmt.Println(CYAN + "2. Encrypt files using AES-GCM" + RESET)
		fmt.Println(CYAN + "3. Decrypt files using AES-GCM" + RESET)
		fmt.Println(CYAN + "4. Set keyfile to use" + RESET)
		fmt.Println(RED + "5. Exit" + RESET)
		fmt.Println(DIM + "6. Toggle debug (irrelevant to encryption) " + debugString + RESET)

		// Get the user's choice. Using a fmt.Print() here in order
		// to give the user a promt to know where to type
		fmt.Print(GOLD + "> " + RESET)
		fmt.Scanln(&choice)

		// Switch statement to handle the different user actions
		switch choice {
		case 6:
			debug = !debug
			fmt.Println(GOLD + "DEBUG ENABLED" + RESET)
		case 5:
			// The simplest one. Just exit the program without doing anything more.
			fmt.Println(RED + "Exiting..." + RESET)
			os.Exit(0)
		case 4:
			// Prompt the user for the key to use in encryption/decryption operations
			// If the key has already been selected, let the user choose whether to switch
			// to another key or not
			if len(keyname) > 0 {
				fmt.Printf("%s ! %sThe key was already set as %s. Do you wish to use another ? (y/n)", GOLD, RESET, keyname)
				fmt.Scanln(&prompt)
				if strings.ToLower(prompt) == "y" {
					fmt.Print("[" + GOLD + "KEY" + RESET + "]")
					fmt.Scanln(&keyname)
				} else {
					fmt.Println("Keeping the old key...")
				}
			}
		case 3:
			// Case 3 : The user wants to decrypt a file.

			// We prompt the user for the keyfile necessary to decrypt the data,
			// and read that key from the keyfile before storing it in a variable
			// we only do that if the key hasn't been set beforehand
			if len(keyname) < 1 {
				fmt.Println("Enter the name of the keyfile : ")
				fmt.Scanln(&keyname)
			}

			key = readKeyFromFileHex(keyname) //reading the key

			// Get the files the user wants to encrypt
			// The file list is separated by spaces and processed accordingly in order
			// to allow us to deal with several files at once
			fmt.Println("Enter the names of the files to decrypt (separated by spaces): ")
			fmt.Scanln(&filenames) //user input for the target file(s)

			// Splitting the file list provided by the user in a slice
			fileList := strings.Fields(filenames)

			// Check if the user gave one path, and if this one path is a
			// directory, recursively get all the files it contains and
			// store them in fileList
			if len(fileList) == 1 {
				ourLittleFile := fileList[0]
				if checkPathType(ourLittleFile) == 1 {
					fileList = getFilesInDir(ourLittleFile)
				}
			}

			// Confirm the user's choice of key file and target file
			fmt.Printf("Using key %s to decrypt files: %s", keyname, fileList)
			fmt.Println("Is that correct? (y/n)") // Prompt the user for their confirmation
			fmt.Scanln(&prompt)                   // Once again, user input

			// Conversion to lowercase to save us conditions
			if strings.ToLower(prompt) != "y" {
				// And if we fail the test, we abort and exit
				fmt.Println(RED + "ABORTING..." + RESET)
				os.Exit(1)
			} else {
				// Create a work group in order to decrypt each file in parallel
				wg.Add(1)
				decryptFiles(key, fileList, &wg, debug)
				//Wait for all the files to be decrypted in order to avoid issues
				wg.Wait()
				fmt.Println(GREEN + "Successfully decrypted files!" + RESET)
			}
		case 2:
			// Case 2 : the user wants to encrypt a file

			// We prompt the user for the keyfile necessary to decrypt the data,
			// and read that key from the keyfile before storing it in a variable
			// we only do so if the key hasn't been set beforehand
			if len(keyname) < 1 {
				fmt.Println("Enter the name of the keyfile : ")
				fmt.Scanln(&keyname)
			}
			// We read the contents of the file (hex encoded key) and stash it for later
			key = readKeyFromFileHex(keyname)

			// Prompt the user for the file they wish to encrypt, separated by spaces if several are selected
			fmt.Println("Enter the names of the files to encrypt (separated by spaces): ")
			fmt.Scanln(&filenames)

			// Split the input into a slice, which is in the for loop's top 10 favorite foods
			fileList := strings.Fields(filenames)

			// Check if the user gave one path, and if this one path is a
			// directory, recursively get all the files it contains and
			// store them in fileList
			if len(fileList) == 1 {
				ourLittleFile := fileList[0]
				if checkPathType(ourLittleFile) == 1 {
					fileList = getFilesInDir(ourLittleFile)
				}
			}

			// Prompt the user for confirmation to encrypt the given file using the key.
			fmt.Printf("Using key %s to encrypt files: %s", keyname, fileList)
			fmt.Println("Is that correct? (y/n)")
			fmt.Scanln(&prompt)

			// Conversion of the string to lowercase in order to write less code by writing more
			if strings.ToLower(prompt) != "y" {
				fmt.Println(RED + "ABORTING..." + RESET)
				os.Exit(1)
			} else {
				// Create a work group for the encryption process
				wg.Add(1)
				// Start the process of encrypting all files from the list
				encryptFiles(key, fileList, &wg, debug)
				// Wait for all files to be processed
				wg.Wait()
				// And then we tell the user we did a good job !
				fmt.Println(GREEN + "Successfully encrypted files!" + RESET)
			}
		case 1:
			// Case 1 : the user wants to generate an AES key

			// Once again, several print statements to let the user know of their options, and
			// of the consequences of the use of a bigger keysize

			/*
				Out of curiosity I took a look at the scaling of encryption, and ran some tests.
				On a GPD Pocket 3 with intel core i7, I got to a speed of 0.842gB/s in encryption,
				and 7.53gB/s in decryption.
			*/

			fmt.Println("Please select the size of the key you wish to generate.")
			fmt.Println("128-bit AES --- 16 bytes of size")
			fmt.Println("256-bit AES --- 24 bytes of size")
			fmt.Println("512-bit AES --- 32 bytes of size")
			fmt.Println("other size  --- cancel the operation")
			fmt.Println("Keep in mind, the size scales both with the strength and")
			fmt.Println("time to process a file using that key :")

			fmt.Scanln(&keysize)

			// Usage of a switch statement to verify the user provided a valid keysize
			// and sending them back to the main menu if they didn't (either to go back or by mistake)
			switch keysize {

			// If it's part of the "allowed" sizes, we go through with the key's creation
			case 16, 24, 32:
				fmt.Printf("Creating a key of size %d", keysize)

				// We generate an AES key of the desired size
				key = generateAESKey(keysize)
				// We ask the user what they wish to save it as
				fmt.Println("Enter the name you want the key to be saved under : ")
				fmt.Scanln(&keyname)
				// Convert the file to a hex string before writing it down under the provided filename
				writeKeyToFileHex(key, keyname)

			// Else (by default), we inform the user that they didn't provide a valid size
			// and send them back to the menu
			default:
				fmt.Printf("%d's not an option, going back to the main menu...", keysize)
			}

		}
	}
}
