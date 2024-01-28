package aws

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/aws/aws-sdk-go-v2/service/translate/types"
)

func GetTranslateClient() *translate.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Cannot load AWS config!")
	}

	client := translate.NewFromConfig(cfg)
	return client
}

func Translate(client *translate.Client, text string) string {
	contentType := "text/plain"

	document := types.Document{
		Content:     []byte(text),
		ContentType: &contentType,
	}

	// Note: for Chinese Traditional characters use "zh-TW"
	// https://docs.aws.amazon.com/translate/latest/dg/what-is-languages.html
	sourceLanguageCode := "zh"

	targetLanguageCode := "en"

	// Here you can add settings for Brevity, Formality and Profanity.
	// We don't have any preference, so we will just leave it as is.
	settings := types.TranslationSettings{}

	input := translate.TranslateDocumentInput{
		Document: &document,
		// Note: for Chinese Traditional characters use zh-TW
		SourceLanguageCode: &sourceLanguageCode,
		TargetLanguageCode: &targetLanguageCode,
		Settings:           &settings,
	}

	translatedDocument, err := client.TranslateDocument(context.TODO(), &input)

	if err != nil {
		if strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "ThrottlingException") {
			log.Printf("AWS Translate rate limit exceeded, sleeping for %d seconds...", AWS_THROTTLING_TIMEOUT_SECONDS)
			time.Sleep(time.Second * AWS_THROTTLING_TIMEOUT_SECONDS)
			return Translate(client, text)
		}

		log.Fatal("An error occurred translating the document: ", err)
	}

	return string(translatedDocument.TranslatedDocument.Content)
}
