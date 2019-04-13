// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dialogflow

// [START import_libraries]
import (
	"context"
	"fmt"
	"strconv"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/option"
	"github.com/golang/protobuf/ptypes/struct"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

//DfApi : dialogflow instance
type DfApi struct {
	operation     string
	projectID     string
	languageCode  string
	sessionClient *dialogflow.SessionsClient
	ctx           context.Context
	sessionPath   string
}

type NLPResponse struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Result     string            `json:"result"`
	Entities   map[string]string `json:"entities"`
}

// [END import_libraries]

//NewDfApi return new DfApi
func NewDfApi(projectid string, languagecode string, authJSONFilePath string) *DfApi {
	dfapi := new(DfApi)
	dfapi.operation = "text" //@TODO: add more type
	dfapi.projectID = projectid
	dfapi.languageCode = languagecode
	dfapi.ctx = context.Background()

	sessionClient, _ := dialogflow.NewSessionsClient(dfapi.ctx, option.WithCredentialsFile(authJSONFilePath))
	dfapi.sessionClient = sessionClient
	return dfapi
}

//GetMsg :get response from dialogflow
func (dfapi *DfApi) GetMsg(input string, userid string) (string, string) {
	// flag.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> -session-id <SESSION ID> -language-code <LANGUAGE CODE> <OPERATION> <INPUTS>\n", filepath.Base(os.Args[0]))
	// 	fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project id\n")
	// 	fmt.Fprintf(os.Stderr, "<SESSION ID> must be a Dialogflow session ID\n")
	// 	fmt.Fprintf(os.Stderr, "<LANGUAGE CODE> must be a language code from https://dialogflow.com/docs/reference/language; defaults to en\n")
	// 	fmt.Fprintf(os.Stderr, "<OPERATION> must be one of text, audio, stream\n")
	// 	fmt.Fprintf(os.Stderr, "<INPUTS> can be a series of text inputs if <OPERATION> is text, or a path to an audio file if <OPERATION> is audio or stream\n")
	// }

	// var projectID, sessionID, languageCode string
	// flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	// flag.StringVar(&sessionID, "session-id", "", "Dialogflow session ID")
	// flag.StringVar(&languageCode, "language-code", "en", "Dialogflow language code from https://dialogflow.com/docs/reference/language; defaults to en")

	// flag.Parse()

	// args := flag.Args()

	// if len(args) == 0 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	// operation := args[0]
	// inputs := args[1:]

	switch dfapi.operation {
	case "text":
		response := dfapi.DetectIntentText(userid, input)
		fmt.Printf("Output: %s\n", response)
		return response.Intent,response.Result

	//@TODO:增加音频处理
	// case "audio":
	// 	response, err := DetectIntentAudio(dfapi.projectID, dfapi.sessionID, dfapi.audioFile, dfapi.languageCode)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Response: %s\n", response)
	// 	return response

	// case "stream":
	// 	response, err := DetectIntentStream(dfapi.projectID, dfapi.sessionID, dfapi.audioFile, dfapi.languageCode)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Response: %s\n", response)
	// 	return response

	default:
		response := dfapi.DetectIntentText(userid, input)
		fmt.Printf("Output: %s\n", response)
		return response.Intent,response.Result
	}
}

// [START dialogflow_detect_intent_text]
func (dfapi *DfApi) DetectIntentText(userid, input string) (r NLPResponse) {
	// ctx := context.Background()

	// sessionClient, err := dialogflow.NewSessionsClient(ctx)
	// if err != nil {
	// 	return "", err
	// }
	// defer sessionClient.Close()

	// if projectID == "" || sessionID == "" {
	// 	return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	// }

	// sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	dfapi.sessionPath = fmt.Sprintf("projects/%s/agent/sessions/%s", dfapi.projectID, userid)
	textInput := dialogflowpb.TextInput{Text: input, LanguageCode: dfapi.languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: dfapi.sessionPath, QueryInput: &queryInput}

	response, err := dfapi.sessionClient.DetectIntent(dfapi.ctx, &request)
	if err != nil {
		return
	}

	queryResult := response.GetQueryResult()
	if queryResult.Intent != nil {
		r.Intent = queryResult.Intent.DisplayName
		r.Confidence = float32(queryResult.IntentDetectionConfidence)
		r.Result = queryResult.GetFulfillmentText()
	}
	r.Entities = make(map[string]string)
	params := queryResult.Parameters.GetFields()
	if len(params) > 0 {
		for paramName, p := range params {
			//fmt.Printf("Param %s: %s (%s)", paramName, p.GetStringValue(), p.String())
			extractedValue := extractDialogflowEntities(p)
			r.Entities[paramName] = extractedValue
		}
	}
	return
}

// [END dialogflow_detect_intent_text]

// [START dialogflow_detect_intent_audio]

// func DetectIntentAudio(projectID, sessionID, audioFile, languageCode string) (string, error) {
// 	ctx := context.Background()

// 	sessionClient, err := dialogflow.NewSessionsClient(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer sessionClient.Close()

// 	if projectID == "" || sessionID == "" {
// 		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
// 	}

// 	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

// 	// In this example, we hard code the encoding and sample rate for simplicity.
// 	audioConfig := dialogflowpb.InputAudioConfig{AudioEncoding: dialogflowpb.AudioEncoding_AUDIO_ENCODING_LINEAR_16, SampleRateHertz: 16000, LanguageCode: languageCode}

// 	queryAudioInput := dialogflowpb.QueryInput_AudioConfig{AudioConfig: &audioConfig}

// 	audioBytes, err := ioutil.ReadFile(audioFile)
// 	if err != nil {
// 		return "", err
// 	}

// 	queryInput := dialogflowpb.QueryInput{Input: &queryAudioInput}
// 	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput, InputAudio: audioBytes}

// 	response, err := sessionClient.DetectIntent(ctx, &request)
// 	if err != nil {
// 		return "", err
// 	}

// 	queryResult := response.GetQueryResult()
// 	fulfillmentText := queryResult.GetFulfillmentText()
// 	return fulfillmentText, nil
// }

// // [END dialogflow_detect_intent_audio]

// // [START dialogflow_detect_intent_streaming]
// func DetectIntentStream(projectID, sessionID, audioFile, languageCode string) (string, error) {
// 	ctx := context.Background()

// 	sessionClient, err := dialogflow.NewSessionsClient(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer sessionClient.Close()

// 	if projectID == "" || sessionID == "" {
// 		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
// 	}

// 	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

// 	// In this example, we hard code the encoding and sample rate for simplicity.
// 	audioConfig := dialogflowpb.InputAudioConfig{AudioEncoding: dialogflowpb.AudioEncoding_AUDIO_ENCODING_LINEAR_16, SampleRateHertz: 16000, LanguageCode: languageCode}

// 	queryAudioInput := dialogflowpb.QueryInput_AudioConfig{AudioConfig: &audioConfig}

// 	queryInput := dialogflowpb.QueryInput{Input: &queryAudioInput}

// 	streamer, err := sessionClient.StreamingDetectIntent(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	f, err := os.Open(audioFile)
// 	if err != nil {
// 		return "", err
// 	}

// 	defer f.Close()

// 	go func() {
// 		audioBytes := make([]byte, 1024)

// 		request := dialogflowpb.StreamingDetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}
// 		err = streamer.Send(&request)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		for {
// 			_, err := f.Read(audioBytes)
// 			if err == io.EOF {
// 				streamer.CloseSend()
// 				break
// 			}
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			request = dialogflowpb.StreamingDetectIntentRequest{InputAudio: audioBytes}
// 			err = streamer.Send(&request)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	}()

// 	var queryResult *dialogflowpb.QueryResult

// 	for {
// 		response, err := streamer.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		recognitionResult := response.GetRecognitionResult()
// 		transcript := recognitionResult.GetTranscript()
// 		log.Printf("Recognition transcript: %s\n", transcript)

// 		queryResult = response.GetQueryResult()
// 	}

// 	fulfillmentText := queryResult.GetFulfillmentText()
// 	return fulfillmentText, nil
// }

// [END dialogflow_detect_intent_streaming]

func extractDialogflowEntities(p *structpb.Value) (extractedEntity string) {
	kind := p.GetKind()
	switch kind.(type) {
	case *structpb.Value_StringValue:
		return p.GetStringValue()
	case *structpb.Value_NumberValue:
		return strconv.FormatFloat(p.GetNumberValue(), 'f', 6, 64)
	case *structpb.Value_BoolValue:
		return strconv.FormatBool(p.GetBoolValue())
	case *structpb.Value_StructValue:
		s := p.GetStructValue()
		fields := s.GetFields()
		extractedEntity = ""
		for key, value := range fields {
			if key == "amount" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, strconv.FormatFloat(value.GetNumberValue(), 'f', 6, 64))
			}
			if key == "unit" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
			}
			if key == "date_time" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
			}
			// @TODO: Other entity types can be added here
		}
		return extractedEntity
	case *structpb.Value_ListValue:
		list := p.GetListValue()
		if len(list.GetValues()) > 1 {
			// @TODO: Extract more values
		}
		extractedEntity = extractDialogflowEntities(list.GetValues()[0])
		return extractedEntity
	default:
		return ""
	}
}