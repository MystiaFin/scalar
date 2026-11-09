package main

import "github.com/MystiaFin/scalar/backend/internal/email"

func extract(subject, body string) (email.Extraction, error) {
    result, err := email.AskOllama(subject, body)
    if err != nil {
        return email.Extraction{}, err
    }

    return email.Extraction(result), nil
}
