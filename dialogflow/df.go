package dialogflow

import (
        "context"
        "fmt"
        "log"
        "strconv"
        "io/ioutil"
        "encoding/json"

        dialogflow "cloud.google.com/go/dialogflow/apiv2"
        "github.com/golang/protobuf/ptypes/struct"
        "google.golang.org/api/option"
        dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// DialogflowProcessor has all the information for connecting with Dialogflow
type DialogflowProcessor struct {
        projectID        string `json:"projectid"`
        authJSONFilePath string `json:"jsonfile"`
        lang             string `json:"lang"`
        timeZone         string `json:"timezone"`
        sessionClient    *dialogflow.SessionsClient
        ctx              context.Context
}

// NLPResponse is the struct for the response
type NLPResponse struct {
        Intent     string				`json:"intent"`
        Confidence float32				`json:"confidence"`
        Result     string				`json:"result"`
        Entities   map[string]string	`json:"entities"`
}

var dp DialogflowProcessor

func init(){
	configFile, err := os.Open("df.json")
    if err != nil {
        printError("opening config file", err.Error())
    }

	jsonParser := json.NewDecoder(configFile)
	var dfconfig DfConfig
    if err = jsonParser.Decode(&df); err != nil {
        printError("parsing config file", err.Error())
    }

	// Auth process: https://dialogflow.com/docs/reference/v2-auth-setup

	dp.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if err != nil {
			log.Fatal("Error in auth with Dialogflow")
	}
	dp.sessionClient = sessionClient
}

func (dp *DialogflowProcessor) processNLP(rawMessage string, username string) (r NLPResponse) {
        sessionID := username
        request := dialogflowpb.DetectIntentRequest{
                Session: fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, sessionID),
                QueryInput: &dialogflowpb.QueryInput{
                        Input: &dialogflowpb.QueryInput_Text{
                                Text: &dialogflowpb.TextInput{
                                        Text:         rawMessage,
                                        LanguageCode: dp.lang,
                                },
                        },
                },
                QueryParams: &dialogflowpb.QueryParameters{
                        TimeZone: dp.timeZone,
                },
        }
        response, err := dp.sessionClient.DetectIntent(dp.ctx, &request)
        if err != nil {
                log.Fatalf("Error in communication with Dialogflow %s", err.Error())
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

