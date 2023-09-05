package translator

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/text/language"
)

type Translator interface {
	Translate(ctx context.Context, text, targetLanguage string) (string, error)
}

type Caller interface {
	Translate(ctx context.Context, text string, targetLanguage language.Tag, opts *translate.Options) ([]translate.Translation, error)
}

type CallerManager struct {
	translationClient *translate.Client
}

func (cm *CallerManager) Translate(ctx context.Context, text string, targetLanguage language.Tag, opts *translate.Options) ([]translate.Translation, error) {
	return cm.translationClient.Translate(ctx, []string{text}, targetLanguage, opts)
}

type Manager struct {
	caller Caller
}

func NewTranslator(caller Caller) Translator {
	return &Manager{
		caller: caller,
	}
}

func (m *Manager) Translate(ctx context.Context, text string, targetLanguage string) (string, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}
	resp, err := m.caller.Translate(ctx, text, lang, &translate.Options{})
	if err != nil {
		return "", fmt.Errorf("translation failed: %w", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
