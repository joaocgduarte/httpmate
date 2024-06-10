package configs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaocgduarte/httpmate/internal/files"
	"github.com/joaocgduarte/httpmate/internal/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ContentType string

type MultipartBodyConfig struct {
	Key                 string  `json:"key"`
	PlainTextValue      *string `json:"plain_text_value"`
	BinaryFilePathValue *string `json:"binary_file_path_value"`
}

type RequestBodyConfig struct {
	RawBody        *string                `json:"raw_body"`
	BinaryFileBody *string                `json:"binary_file_body"`
	MultipartBody  []*MultipartBodyConfig `json:"multipart_body"`
	FormURLEncoded map[string]string      `json:"form_url_encoded"`
}

type RequestConfig struct {
	Collection  string            `json:"collection"`
	RequestName string            `json:"request_name"`
	Domain      string            `json:"domain"`
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	QueryParams map[string]string `json:"query_params"`
	Headers     map[string]string `json:"headers"`
	ContentType string            `json:"content_type"`
	Body        RequestBodyConfig `json:"body"`
}

func (r *RequestConfig) WriteToJSONFile() {
	files.WriteStructToJSONFile(
		r,
		filepath.Join(r.Collection, fmt.Sprintf("%s.json", r.RequestName)),
	)
}

var (
	ContentTypeJSON              ContentType = "application/json"
	ContentTypeXML               ContentType = "application/xml"
	ContentTypeOctetStream       ContentType = "application/octet-stream"
	ContentTypeMultipartFormData ContentType = "multipart/form-data"
	ContentTypeFormURLEncoded    ContentType = "application/x-www-form-urlencoded"
)

func (c ContentType) IsRecognizable() bool {
	return c == ContentTypeJSON ||
		c == ContentTypeXML ||
		c == ContentTypeOctetStream ||
		c == ContentTypeMultipartFormData ||
		c == ContentTypeFormURLEncoded
}

func PromptRequestConfig(collectionsPath string) *RequestConfig {
	collections := files.GetSubDirectories(collectionsPath)
	collection := chooseCollection(collectionsPath, collections)

	config := &RequestConfig{
		Collection:  collection,
		RequestName: prompts.Prompt("Request name"),
		Domain:      prompts.Prompt("Domain"),
		Path:        prompts.Prompt("Path"),
		Method: prompts.Select("Method", []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		}),
		QueryParams: prompts.PromptWhileConfirm("Do you want to add a query parameter?", "Query parameter key", "Query parameter value"),
		Headers:     prompts.PromptWhileConfirm("Do you want to add a header?", "Header key (don't add Content-Type)", "Header value"),
	}

	config.ContentType = prompts.SelectWithAdd("Content-Type", "Other", []string{
		string(ContentTypeJSON),
		string(ContentTypeXML),
		string(ContentTypeOctetStream),
		string(ContentTypeMultipartFormData),
		string(ContentTypeFormURLEncoded),
	})

	if prompts.ConfirmPrompt("Do you want your request to have a body?") == false {
		return config
	}

	config.Body = RequestBodyConfig{}
	config.Body.RawBody = promptRawBody(ContentType(config.ContentType))
	if config.Body.RawBody != nil {
		return config
	}

	config.Body.BinaryFileBody = promptBinaryFileBody(ContentType(config.ContentType))
	if config.Body.BinaryFileBody != nil {
		return config
	}

	config.Body.MultipartBody = promptMultipartBody(ContentType(config.ContentType))
	if config.Body.MultipartBody != nil {
		return config
	}

	config.Body.FormURLEncoded = promptFormURLEncoded(ContentType(config.ContentType))
	return config
}

func chooseCollection(collectionsPath string, collections []string) string {
	result := prompts.SelectWithAdd(
		"Choose one of your collections",
		"Create new collection",
		collections,
	)

	collectionPath := filepath.Join(collectionsPath, result)
	files.CreateDirectory(collectionPath)
	return collectionPath
}

func promptRawBody(contentType ContentType) *string {
	if contentType == ContentTypeMultipartFormData ||
		contentType == ContentTypeOctetStream ||
		contentType == ContentTypeFormURLEncoded {
		return nil
	}

	if !contentType.IsRecognizable() {
		if prompts.ConfirmPrompt("Do you want to add raw body to your request?") == false {
			return nil
		}
	}

	extension := ".txt"
	switch contentType {
	case ContentTypeJSON:
		extension = ".json"
	case ContentTypeXML:
		extension = ".xml"
	default:
		extension = ".txt"
	}

	res := prompts.TextEditorPrompt(
		viper.GetString("editor"),
		fmt.Sprintf("create-request-body%s", extension),
		viper.GetString("temporaryFilesDirectory"),
		"",
	)
	return &res
}

func promptBinaryFileBody(contentType ContentType) *string {
	if !contentType.IsRecognizable() {
		if prompts.ConfirmPrompt("Do you want to add binary body to your request?") == false {
			return nil
		}
	}

	if contentType != ContentTypeOctetStream {
		return nil
	}

	filePath := prompts.Prompt("Path to file (full path)")
	return &filePath
}

func promptMultipartBody(contentType ContentType) []*MultipartBodyConfig {
	if contentType != ContentTypeMultipartFormData {
		return nil
	}
	result := make([]*MultipartBodyConfig, 0)

	for true {
		confirmRes := prompts.ConfirmPrompt("Do you want to add parts?")
		if !confirmRes {
			break
		}

		part := &MultipartBodyConfig{
			Key:                 prompts.Prompt("What is the key for this part?"),
			PlainTextValue:      new(string),
			BinaryFilePathValue: new(string),
		}

		binaryPartOrPlainText := prompts.Select("Binary file or plain text value?", []string{
			"Binary File",
			"Plain Text",
		})

		if binaryPartOrPlainText == "Binary File" {
			filePath := prompts.Prompt("Path to file (full path)")
			part.PlainTextValue = &filePath
		} else {
			res := prompts.TextEditorPrompt(
				viper.GetString("editor"),
				"create-request-multipart-body.txt",
				viper.GetString("temporaryFilesDirectory"),
				"",
			)
			part.BinaryFilePathValue = &res
		}

		result = append(result, part)
	}

	return result
}

func promptFormURLEncoded(contentType ContentType) map[string]string {
	if contentType != ContentTypeFormURLEncoded {
		return nil
	}

	return prompts.PromptWhileConfirm("Do you want to add a parameter?", "Parameter key", "Parameter value")
}

type EditRequestFlags struct {
	EditBody        bool
	EditDomain      bool
	EditPath        bool
	EditQueryParams bool
	EditHeaders     bool
	EditMethod      bool
	EditContentType bool
	EditAll         bool
}

func (config *RequestConfig) PromptEditConfig(editConfigs EditRequestFlags) {
	if editConfigs.EditAll {
		config.editAll()
		return
	}

	if editConfigs.EditBody {
		config.editBody()
	}

	if editConfigs.EditDomain {
		config.Domain = prompts.PromptWithDefault("Domain", config.Domain)
	}

	if editConfigs.EditPath {
		config.Path = prompts.PromptWithDefault("Path", config.Path)
	}

	if editConfigs.EditMethod {
		config.Method = prompts.Select("Method", []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		})
	}

	if editConfigs.EditContentType {
		config.ContentType = prompts.SelectWithAdd("Content-Type", "Other", []string{
			string(ContentTypeJSON),
			string(ContentTypeXML),
			string(ContentTypeOctetStream),
			string(ContentTypeMultipartFormData),
			string(ContentTypeFormURLEncoded),
		})
	}

	if editConfigs.EditQueryParams {
		marshalledConfig, err := json.MarshalIndent(config.QueryParams, "", "\t")
		cobra.CheckErr(err)

		alteredQueryParams := prompts.TextEditorPrompt(
			viper.GetString("editor"),
			fmt.Sprintf("edit query params %s.json", config.RequestName),
			viper.GetString("temporaryFilesDirectory"),
			string(marshalledConfig),
		)

		var newQueryParams map[string]string
		err = json.Unmarshal([]byte(alteredQueryParams), &newQueryParams)
		cobra.CheckErr(err)
		config.QueryParams = newQueryParams
	}

	if editConfigs.EditHeaders {
		marshalledConfig, err := json.MarshalIndent(config.Headers, "", "\t")
		cobra.CheckErr(err)

		alteredHeaders := prompts.TextEditorPrompt(
			viper.GetString("editor"),
			fmt.Sprintf("edit headers %s.json", config.RequestName),
			viper.GetString("temporaryFilesDirectory"),
			string(marshalledConfig),
		)

		var newHeaders map[string]string
		err = json.Unmarshal([]byte(alteredHeaders), &newHeaders)
		cobra.CheckErr(err)
		config.Headers = newHeaders
	}

	config.WriteToJSONFile()
}

func (config *RequestConfig) editBody() {
	if config.Body.RawBody != nil {
		alteredBody := prompts.TextEditorPrompt(
			viper.GetString("editor"),
			fmt.Sprintf("edit body %s.json", config.RequestName),
			viper.GetString("temporaryFilesDirectory"),
			*config.Body.RawBody,
		)
		config.Body.RawBody = &alteredBody
		return
	}

	if config.Body.BinaryFileBody != nil {
		alteredBody := prompts.PromptWithDefault("File Path", *config.Body.BinaryFileBody)
		config.Body.BinaryFileBody = &alteredBody
		return
	}

	if config.Body.MultipartBody != nil {
		marshalledConfig, err := json.MarshalIndent(config.Body.MultipartBody, "", "\t")
		cobra.CheckErr(err)

		alteredBody := prompts.TextEditorPrompt(
			viper.GetString("editor"),
			fmt.Sprintf("edit body %s.json", config.RequestName),
			viper.GetString("temporaryFilesDirectory"),
			string(marshalledConfig),
		)

		var newMultipartBody []*MultipartBodyConfig
		err = json.Unmarshal([]byte(alteredBody), &newMultipartBody)
		cobra.CheckErr(err)
		config.Body.MultipartBody = newMultipartBody
		return
	}

	if config.Body.FormURLEncoded != nil {
		marshalledConfig, err := json.MarshalIndent(config.Body.FormURLEncoded, "", "\t")
		cobra.CheckErr(err)

		alteredBody := prompts.TextEditorPrompt(
			viper.GetString("editor"),
			fmt.Sprintf("edit body %s.json", config.RequestName),
			viper.GetString("temporaryFilesDirectory"),
			string(marshalledConfig),
		)

		var newFormURLEncodedBody map[string]string
		err = json.Unmarshal([]byte(alteredBody), &newFormURLEncodedBody)
		cobra.CheckErr(err)
		config.Body.FormURLEncoded = newFormURLEncodedBody
		return
	}
}

func (config *RequestConfig) editAll() RequestConfig {
	marshalledConfig, err := json.MarshalIndent(config, "", "\t")
	cobra.CheckErr(err)

	alteredConfigs := prompts.TextEditorPrompt(
		viper.GetString("editor"),
		fmt.Sprintf("edit all %s.json", config.RequestName),
		viper.GetString("temporaryFilesDirectory"),
		string(marshalledConfig),
	)

	var newConfigs RequestConfig
	err = json.Unmarshal([]byte(alteredConfigs), &newConfigs)
	cobra.CheckErr(err)

	files.WriteStructToJSONFile(
		newConfigs,
		filepath.Join(newConfigs.Collection, fmt.Sprintf("%s.json", newConfigs.RequestName)),
	)
	return newConfigs
}

// ConvertToCurlCommand converts RequestConfig to a curl command string.
func (config *RequestConfig) ConvertToCurlCommand() string {
	var curlCmd strings.Builder

	// Append curl command basics
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(config.Method)
	curlCmd.WriteString(" '")
	curlCmd.WriteString(config.Domain)
	curlCmd.WriteString(config.Path)
	curlCmd.WriteString("'")

	// Append query parameters
	if len(config.QueryParams) > 0 {
		curlCmd.WriteString("?")
		queryParams := url.Values{}
		for key, value := range config.QueryParams {
			queryParams.Set(key, value)
		}
		curlCmd.WriteString(queryParams.Encode())
	}

	if config.ContentType != "" {
		curlCmd.WriteString(" -H '")
		curlCmd.WriteString("Content-Type")
		curlCmd.WriteString(": ")
		curlCmd.WriteString(config.ContentType)
		curlCmd.WriteString("'")
	}

	for key, value := range config.Headers {
		curlCmd.WriteString(" -H '")
		curlCmd.WriteString(key)
		curlCmd.WriteString(": ")
		curlCmd.WriteString(value)
		curlCmd.WriteString("'")
	}

	// Append request body if present
	switch {
	case config.Body.RawBody != nil:
		curlCmd.WriteString(" --data '")
		curlCmd.WriteString(*config.Body.RawBody)
		curlCmd.WriteString("'")
	case config.Body.BinaryFileBody != nil:
		curlCmd.WriteString(" --data-binary @")
		curlCmd.WriteString(*config.Body.BinaryFileBody)
	case len(config.Body.MultipartBody) > 0:
		// Multipart form data
		for _, part := range config.Body.MultipartBody {
			if part.BinaryFilePathValue != nil {
				curlCmd.WriteString(" -F '")
				curlCmd.WriteString(part.Key)
				curlCmd.WriteString("=@")
				curlCmd.WriteString(*part.BinaryFilePathValue)
				curlCmd.WriteString("'")
			} else if part.PlainTextValue != nil {
				curlCmd.WriteString(" -F '")
				curlCmd.WriteString(part.Key)
				curlCmd.WriteString("=")
				curlCmd.WriteString(*part.PlainTextValue)
				curlCmd.WriteString("'")
			}
		}
	case len(config.Body.FormURLEncoded) > 0:
		// Form URL-encoded
		curlCmd.WriteString(" -d '")
		data := url.Values{}
		for key, value := range config.Body.FormURLEncoded {
			data.Set(key, value)
		}
		curlCmd.WriteString(data.Encode())
		curlCmd.WriteString("'")
	}

	return curlCmd.String()
}

func (config *RequestConfig) BuildHTTPRequest() *http.Request {
	u := config.buildURL()
	req := config.buildRequest(u)

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}
	return req
}

func (config *RequestConfig) buildURL() *url.URL {
	domain := strings.Trim(config.Domain, "/")
	path := strings.Trim(config.Path, "/")
	baseURL := fmt.Sprintf("%s/%s", domain, path)

	u, err := url.Parse(baseURL)
	cobra.CheckErr(err)

	q := u.Query()
	for key, value := range config.QueryParams {
		q.Set(strings.Trim(key, " "), strings.Trim(value, " "))
	}
	u.RawQuery = q.Encode()
	return u
}

func (config *RequestConfig) buildRequest(u *url.URL) *http.Request {
	if config.Body.RawBody != nil {
		req, err := http.NewRequest(config.Method, u.String(), strings.NewReader(*config.Body.RawBody))
		cobra.CheckErr(err)
		setRequestHeaderFromConfig(req, config)
		return req
	}

	if config.Body.BinaryFileBody != nil {
		file, err := os.Open(*config.Body.BinaryFileBody)
		cobra.CheckErr(err)
		req, err := http.NewRequest(config.Method, u.String(), file)
		cobra.CheckErr(err)
		setRequestHeaderFromConfig(req, config)
		return req
	}

	if len(config.Body.MultipartBody) > 0 {
		body, contentType, err := createMultipartBody(config.Body.MultipartBody)
		cobra.CheckErr(err)

		req, err := http.NewRequest(config.Method, u.String(), body)
		cobra.CheckErr(err)
		req.Header.Set("Content-Type", contentType)
		return req
	}

	if len(config.Body.FormURLEncoded) > 0 {
		data := url.Values{}
		for key, value := range config.Body.FormURLEncoded {
			data.Set(key, value)
		}
		req, err := http.NewRequest(config.Method, u.String(), strings.NewReader(data.Encode()))
		cobra.CheckErr(err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req
	}

	req, err := http.NewRequest(config.Method, u.String(), nil)
	cobra.CheckErr(err)
	setRequestHeaderFromConfig(req, config)
	return req
}

func createMultipartBody(bodyConfig []*MultipartBodyConfig) (io.Reader, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for _, part := range bodyConfig {
		if part.PlainTextValue != nil {
			fw, err := w.CreateFormField(part.Key)
			cobra.CheckErr(err)
			_, err = fw.Write([]byte(*part.PlainTextValue))
			cobra.CheckErr(err)
			continue
		}

		if part.BinaryFilePathValue != nil {
			file, err := os.Open(*part.BinaryFilePathValue)
			cobra.CheckErr(err)
			defer file.Close()

			fw, err := w.CreateFormFile(part.Key, filepath.Base(*part.BinaryFilePathValue))
			cobra.CheckErr(err)
			_, err = io.Copy(fw, file)
			cobra.CheckErr(err)
		}
	}

	err := w.Close()
	cobra.CheckErr(err)

	return &b, w.FormDataContentType(), nil
}

func setRequestHeaderFromConfig(req *http.Request, config *RequestConfig) {
	req.Header.Set("Content-Type", strings.Trim(config.ContentType, " "))
}

func PromptNewExistentRequestConfig(label, collectionsPath string) *RequestConfig {
	availableRequests := files.GetFilesFromDirectoryWithoutExtension(collectionsPath)

	if len(availableRequests) == 0 {
		cobra.CompError("There are no available requests in collection")
		os.Exit(-1)
	}

	request := prompts.Select(label, availableRequests)

	wantedRequestPath := filepath.Join(collectionsPath, fmt.Sprintf("%s.json", request))

	return NewRequestConfigFromFilePath(wantedRequestPath)
}

func NewRequestConfigFromFilePath(filepath string) *RequestConfig {
	jsonFile, err := os.Open(filepath)
	cobra.CheckErr(err)
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	cobra.CheckErr(err)

	var result RequestConfig
	err = json.Unmarshal(byteValue, &result)
	cobra.CheckErr(err)
	return &result
}
