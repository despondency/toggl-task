package translator_test

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"github.com/despondency/toggl-task/internal/translator"
	translatormock "github.com/despondency/toggl-task/mocks/translator"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestUnitTranslator_Translate(t *testing.T) {
	type testCase struct {
		name                string
		ctx                 context.Context
		text                string
		language            string
		expectedTranslation string
		expectedErr         error
		createInstance      func(tc *testCase, t *testing.T) translator.Translator
	}

	testCases := []testCase{
		{
			name:                "translate, no error",
			expectedTranslation: "some-translation",
			language:            "en",
			text:                "c'est la vie",
			createInstance: func(tc *testCase, t *testing.T) translator.Translator {
				caller := translatormock.NewCaller(t)
				targetLang, err := language.Parse(tc.language)
				if err != nil {
					panic(err)
				}
				caller.EXPECT().Translate(context.Background(), tc.text, targetLang, &translate.Options{}).Return([]translate.Translation{
					{
						Text:   "some-translation",
						Source: language.Tag{},
						Model:  "some-model",
					},
				}, nil)
				return translator.NewTranslator(caller)
			},
		},
		{
			name:        "translate, error, can't parse language",
			language:    "no such language found",
			expectedErr: fmt.Errorf("language: tag is not well-formed"),
			createInstance: func(tc *testCase, t *testing.T) translator.Translator {
				return translator.NewTranslator(translatormock.NewCaller(t))
			},
		},
		{
			name:        "translate, error, translation failed",
			text:        "c'est la vie",
			language:    "en",
			expectedErr: fmt.Errorf("translation failed: some-err-occurred"),
			createInstance: func(tc *testCase, t *testing.T) translator.Translator {
				caller := translatormock.NewCaller(t)
				targetLang, err := language.Parse(tc.language)
				if err != nil {
					panic(err)
				}
				caller.EXPECT().Translate(context.Background(), tc.text, targetLang, &translate.Options{}).Return(nil, fmt.Errorf("some-err-occurred"))
				return translator.NewTranslator(caller)
			},
		},
		{
			name:     "translate, returned empty response",
			text:     "c'est la vie",
			language: "en",
			createInstance: func(tc *testCase, t *testing.T) translator.Translator {
				caller := translatormock.NewCaller(t)
				targetLang, err := language.Parse(tc.language)
				if err != nil {
					panic(err)
				}
				caller.EXPECT().Translate(context.Background(), tc.text, targetLang, &translate.Options{}).Return(nil, nil)
				return translator.NewTranslator(caller)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance := tc.createInstance(&tc, t)
			translation, err := instance.Translate(context.Background(), tc.text, tc.language)
			if tc.expectedErr != nil {
				assert.EqualErrorf(t, err, tc.expectedErr.Error(), "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTranslation, translation)
			}
		})
	}
}
