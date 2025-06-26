package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"telemetry-task/lib/crypto"
	logUtil "telemetry-task/lib/logger"
)

var (
	logger = logUtil.LoggerWithPrefix("MAIN")
)

func main() {
	encryptedFilePath := flag.String("input", "", "path to encrypted file with metrics")
	decryptedFilePath := flag.String("output", "", "path to decrypted file with metrics")
	flag.Parse()

	if *encryptedFilePath == "" {
		log.Fatal("input file is required")
	}
	if *decryptedFilePath == "" {
		path := filepath.Join(filepath.Dir(*encryptedFilePath) + fmt.Sprintf("/decrypted_%s", filepath.Base(*encryptedFilePath)))
		decryptedFilePath = &path
	}
	encryptedFileContent, err := os.ReadFile(*encryptedFilePath)
	if err != nil {
		log.Fatal("failed to read encrypted file, err:", err)
	}
	if len(encryptedFileContent) == 0 {
		log.Fatal("encrypted file is empty")
	}

	decryptedFile, err := os.Create(*decryptedFilePath)
	if err != nil {
		log.Fatal("failed to create decrypted file, err:", err)
	}
	defer func() {
		err := decryptedFile.Close()
		if err != nil {
			logger.Error("failed to close decrypted file", "err", err)
		}
	}()

	decryptor, err := crypto.NewDecryptor()
	if err != nil {
		log.Fatal("failed to create decryptor, err:", err)
	}

	lines := strings.Split(string(encryptedFileContent), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		line = strings.Trim(line, "\n")
		decrypted, err := decryptor.DecryptMessage(line)
		if err != nil {
			logger.Error("failed to decrypt message", "err", err, "line", line)
			continue
		}
		if decrypted == "" {
			continue
		}
		_, err = decryptedFile.WriteString(decrypted + "\n")
		if err != nil {
			logger.Error("failed to write decrypted message to file", "err", err, "line", line)
			continue
		}
	}
	logger.Info("Decryption completed successfully", "output_file", *decryptedFilePath)
}
