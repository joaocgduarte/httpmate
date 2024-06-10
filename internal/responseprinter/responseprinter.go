package responseprinter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

func PrintHTTPResponse(resp *http.Response, processingTime time.Duration) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:")

	for key, value := range resp.Header {
		fmt.Printf("%s: %s\n", key, value)
	}

	fmt.Println("Response:")
	PrintJSON(body)
	fmt.Println("Time taken:", processingTime)
}

func PrintJSON(toPrint []byte) {
	if json.Valid(toPrint) && isJQAvailable() {
		cmd := exec.Command("jq", ".")
		cmd.Stdin = io.NopCloser(os.Stdin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = bytes.NewReader(toPrint)
		err := cmd.Run()
		cobra.CheckErr(err)
	} else {
		fmt.Println(string(toPrint))
	}
}

func isJQAvailable() bool {
	err := exec.Command("jq", "--version").Run()
	return err == nil
}
