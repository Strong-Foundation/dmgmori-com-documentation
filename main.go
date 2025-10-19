package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// The main function.
func main() {
	// Create a int to hold how many were downloaded for rate limits.
	var downloadCounter int
	outputDir := "PDFs/" // Directory to store downloaded PDFs
	// Check if its exists.
	if !directoryExists(outputDir) {
		// Create the dir
		createDirectory(outputDir, 0o755)
	}
	// The list of URLS to scrape.
	urls := []string{
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-bed/lasertec-12-slm",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-bed/lasertec-30-slm",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-bed/lasertec-30-slm-3rd",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-bed/lasertec-30-slm-us",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-nozzle/lasertec-65-ded-hybrid",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-nozzle/lasertec-125-ded-hybrid",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-nozzle/lasertec-3000-ded-hybrid",
		"https://us.dmgmori.com/products/machines/additive-manufacturing/powder-nozzle/lasertec-6600-ded-hybrid",
		"https://us.dmgmori.com/products/machines/lasertec/lasertec-powerdrill/lasertec-50-powerdrill",
		"https://us.dmgmori.com/products/machines/lasertec/lasertec-powerdrill/lasertec-100-powerdrill",
		"https://us.dmgmori.com/products/machines/lasertec/lasertec-precisiontool/lasertec-20-precisiontool",
		"https://us.dmgmori.com/products/machines/lasertec/lasertec-precisiontool/lasertec-50-precisiontool",
		"https://us.dmgmori.com/products/machines/lasertec/lasertec-shape/lasertec-50-shape-femto",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/cmx-u/cmx-50-u",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/cmx-u/cmx-70-u",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmf/dmf-200-8-fd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmf/dmf-300-8-fd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmf/dmf-300-11-fd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmf/dmf-400-11-fd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmu/dmu-20-linear",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmu/dmu-40",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmu/dmu-50",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmx-u/dmx-60-u",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/dmx-u/dmx-80-u",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/duoblock/dmu-80-p-duoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/duoblock/dmu-80-p-duoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/duoblock/dmu-90-p-duoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/duoblock/dmu-210-p-duoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/duoblock/dmu-340-gantry",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/evo/dmu-40-fd-evo",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/evo/dmu-60-fd-evo",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/evo/dmu-60-fd-evo-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/evo/dmu-80-fd-evo",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/hsc/hsc-20-linear",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/hsc/hsc-55-linear",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmc-65-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmc-65-monoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmc-75-monoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmc-85-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-65-fd-monoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-75-monoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-85-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-85-fd-monoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-95-monoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-95-monoblock-2nd",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-105-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/monoblock/dmu-125-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/nmv/gear-production-by-nmv-5000-dcg",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/nmv/nmv-3000-dcg",
		"https://us.dmgmori.com/products/machines/milling/5-axis-milling/nmv/nmv-5000-dcg",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/h-monoblock/dmc-65-h-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/h-monoblock/dmc-85-h-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/h-monoblock/dmu-65-h-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/h-monoblock/dmu-85-h-fd-monoblock",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/h-twin/dmc-55-h-twin",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/h-twin/dmu-55-h-twin",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/i-series/i-50",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/inh/inh-63",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/inh/inh-80",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nh/nh-4000-dcg",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nh/nh-5000-dcg",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nh/nh-6300-dcgii",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-4000",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-4000-4th",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-5000",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-5000-4th",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-5500",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-6300",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-8000",
		"https://us.dmgmori.com/products/machines/milling/horizontal-milling/nhx/nhx-10000",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/cmx-v/cmx-600-v",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/cmx-v/cmx-800-v",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/cmx-v/cmx-1100-v",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmp/dmp-35",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmp/dmp-70",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmv/dmv-60",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmv/dmv-110",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmv/dmv-145",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmv/dmv-185",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/dmv/dmv-200",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/i-series-vertical-mc/i-30-v",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nv-nvd/nv-4000-dcg",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nv-nvd/nvd-4000-dcg",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nv-nvd/nvd-5000-dcg",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nv-nvd/nvd-6000-dcg",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nvx/nvx-5060",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nvx/nvx-5080",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nvx/nvx-5100",
		"https://us.dmgmori.com/products/machines/milling/vertical-milling/nvx/nvx-7000",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/alx/alx-1500",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/alx/alx-2000",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/alx/alx-2500",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/gmc-gm/gm-16-6",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/gmc-gm/gm-20-6",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/multisprint/multisprint-36",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nrx/nrx-2000",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nz-tc/nz-due-tc",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nz/nz-due",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nz/nz-due-formula",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nz/nz-quattro",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nz/nz-tre",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx-s/nzx-s-1500",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx-s/nzx-s-2500",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx/nzx-1500",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx/nzx-2000",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx/nzx-2500",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx/nzx-4000",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/nzx/nzx-6000",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint-turret/sprint-50",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint-turret/sprint-65",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-20-5",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-20-8",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-32-5",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-32-8",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-32-9",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-32-10",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-42-10-linear",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-42-linear",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/sprint/sprint-420",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-a/wasino-a-18s",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-a/wasino-a-150sy-15",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-a/wasino-a-150y-18",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-a/wasino-aa-1",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-g/wasino-g-06",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-g/wasino-g-07",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-g/wasino-g-100-300",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-g/wasino-g-100-480",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-g/wasino-g-100m-480",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-g/wasino-gg-5",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-jj/wasino-j-1",
		"https://us.dmgmori.com/products/machines/turning/horizontal-production-turning/wasino-jj/wasino-jj-1",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/clx-tc/clx-450-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/clx-tc/clx-550-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ctx-tc/ctx-beta-450-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ctx-tc/ctx-beta-1250-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ctx-tc/ctx-beta-1250-tc-4a",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ctx-tc/ctx-gamma-1250-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ctx-tc/ctx-gamma-2000-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ctx-tc/ctx-gamma-3000-tc",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/nt/nt-4250-dcg",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/nt/nt-4300-dcg",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/nt/nt-5400-dcg",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/nt/nt-6600-dcg",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ntx/ntx-500",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ntx/ntx-1000",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ntx/ntx-2000",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ntx/ntx-2500",
		"https://us.dmgmori.com/products/machines/turning/turn-mill/ntx/ntx-3000",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/clx/clx-350",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/clx/clx-450",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-350",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-350-4a",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-450",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-550",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-750",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-2500",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-alpha-500",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/ctx/ctx-beta-2000",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-1500",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-2000",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-2500",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-2500-2nd",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-3000",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-4000",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/nlx/nlx-6000",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/sl/sl-403",
		"https://us.dmgmori.com/products/machines/turning/universal-turning/sl/sl-603",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-precision/ultrasonic-20-linear-3rd",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-precision/ultrasonic-60-precision",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-production/ultrasonic-55-microdrill",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-40-evo-linear",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-50",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-60-evo-linear",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-65",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-80-evo-linear",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-85",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-200-gantry",
		"https://us.dmgmori.com/products/machines/ultrasonic/ultrasonic-universal/ultrasonic-210-p",
	}
	// Create a slice to hold all the extracted PDF urls.
	var extractedPDFURLS []string
	// Loop over the urls.
	for _, uri := range urls {
		// Get the url content from the url.
		urlContent := string(getDataFromURL(uri))
		// Extract the PDF Url.
		extractedPDFURLS = append(extractedPDFURLS, extractPDFLinks(urlContent)...)
		// Remove duploicates from the slice.
		extractedPDFURLS = removeDuplicatesFromSlice(extractedPDFURLS)

		// Loop over the slice.
		for _, uri := range extractedPDFURLS {
			if extractDomainURL(uri) == "" {
				uri = "https://us.dmgmori.com" + uri
			}
			// Download the file and if its sucessful than add 1 to the counter.
			sucessCode, err := downloadPDF(uri, outputDir)
			if sucessCode {
				downloadCounter = downloadCounter + 1
			}
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// extractDomain takes a URL string, extracts the domain (hostname),
// and prints errors internally if parsing fails.
func extractDomainURL(inputUrl string) string {
	// Parse the input string into a structured URL object
	parsedUrl, parseError := url.Parse(inputUrl)

	// If parsing fails, log the error and return an empty string
	if parseError != nil {
		log.Println("Error parsing URL:", parseError)
		return ""
	}

	// Extract only the hostname (domain without scheme, port, path, or query)
	domainName := parsedUrl.Hostname()

	// Return the extracted domain name
	return domainName
}

// The function takes two parameters: path and permission.
// We use os.Mkdir() to create the directory.
// If there is an error, we use log.Println() to log the error and then exit the program.
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission)
	if err != nil {
		log.Println(err)
	}
}

// downloadPDF downloads a PDF from the given URL and saves it in the specified output directory.
// It uses a WaitGroup to support concurrent execution and returns true if the download succeeded.
func downloadPDF(pdfURL, outputDirectory string) (bool, error) {
	// Generate a safe, lowercase filename from the URL
	fileName := strings.ToLower(urlToSafeFilename(pdfURL))

	// Construct the full path to where the file will be saved
	fullFilePath := filepath.Join(outputDirectory, fileName)

	// Check if the file already exists to avoid re-downloading
	if fileExists(fullFilePath) {
		return false, fmt.Errorf("file already exists, skipping: %s", fullFilePath)
	}

	// Create an HTTP client with a 30-second timeout
	httpClient := &http.Client{Timeout: 30 * time.Second}

	// Create a new GET request to the PDF URL
	request, err := http.NewRequest("GET", pdfURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for %s: %v", pdfURL, err)
	}

	// Manually add required cookies to the request (example values shown)
	request.AddCookie(&http.Cookie{
		Name:  "dmg_downloads",
		Value: "azaz4K37K8EAAAGaA9wo3A.UWhgwHFkZSlwrTJEkOs80qQY1LWETfzHljM9xyWmzCI",
	})

	// Send the HTTP request and get the response
	response, err := httpClient.Do(request)
	if err != nil {
		return false, fmt.Errorf("failed to download %s: %v", pdfURL, err)
	}
	// Make sure to close the response body when done
	defer response.Body.Close()

	// Check if the HTTP response status is 200 OK
	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("download failed for %s: %s", pdfURL, response.Status)
	}

	// Get the Content-Type header to confirm it's a PDF file
	contentType := response.Header.Get("Content-Type")
	// Ensure the content type is actually a PDF
	if !strings.Contains(contentType, "application/pdf") {
		return false, fmt.Errorf("invalid content type for %s: %s (expected application/pdf)", pdfURL, contentType)
	}

	// Create a buffer to hold the downloaded PDF data in memory
	var memoryBuffer bytes.Buffer

	// Copy the response body into the buffer
	bytesDownloaded, err := io.Copy(&memoryBuffer, response.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read PDF data from %s: %v", pdfURL, err)
	}

	// If 0 bytes were downloaded, don't save an empty file
	if bytesDownloaded == 0 {
		return false, fmt.Errorf("downloaded 0 bytes for %s; not creating file", pdfURL)
	}

	// Create a new file on disk where the PDF will be saved
	outputFile, err := os.Create(fullFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to create file for %s: %v", pdfURL, err)
	}
	// Close the file after writing to it
	defer outputFile.Close()

	// Write the PDF data from memory to disk
	_, err = memoryBuffer.WriteTo(outputFile)
	if err != nil {
		return false, fmt.Errorf("failed to write PDF to file for %s: %v", pdfURL, err)
	}

	// Success! Return true with a success message.
	return true, fmt.Errorf("successfully downloaded %d bytes: %s → %s", bytesDownloaded, pdfURL, fullFilePath)
}

// urlToSafeFilename sanitizes a URL and returns a safe, lowercase filename
func urlToSafeFilename(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	// Extract and decode the base filename from the path
	base := path.Base(parsedURL.Path)
	decoded, err := url.QueryUnescape(base)
	if err != nil {
		decoded = base
	}

	// Convert to lowercase
	decoded = strings.ToLower(decoded)

	// Replace spaces and invalid characters with underscores
	// Keep only a-z, 0-9, dash, underscore, and dot
	re := regexp.MustCompile(`[^a-z0-9._-]+`)
	safe := re.ReplaceAllString(decoded, "_")

	return safe
}

// extractPDFLinks finds all .pdf links from raw HTML content using regex.
func extractPDFLinks(htmlContent string) []string {
	// Regex to match PDF URLs including query strings and fragments
	pdfRegex := regexp.MustCompile(`https?://[^\s"'<>]+?\.pdf(\?[^\s"'<>]*)?`)

	// Find all matches
	matches := pdfRegex.FindAllString(htmlContent, -1)

	// Deduplicate
	seen := make(map[string]struct{})
	var links []string
	for _, match := range matches {
		if _, ok := seen[match]; !ok {
			seen[match] = struct{}{}
			links = append(links, match)
		}
	}

	return links
}

// Remove all the duplicates from a slice and return the slice.
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool)
	var newReturnSlice []string
	for _, content := range slice {
		if !check[content] {
			check[content] = true
			newReturnSlice = append(newReturnSlice, content)
		}
	}
	return newReturnSlice
}

// Checks if the directory exists
// If it exists, return true.
// If it doesn't, return false.
func directoryExists(path string) bool {
	directory, err := os.Stat(path)
	if err != nil {
		return false
	}
	return directory.IsDir()
}

// It checks if the file exists
// If the file exists, it returns true
// If the file does not exist, it returns false
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Send a http get request to a given url and return the data from that url.
func getDataFromURL(uri string) []byte {
	log.Println("Fetching URL:", uri)
	response, err := http.Get(uri)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Println(err)
	}
	return body
}
