package main

type SessionListItem struct {
	ID string `json:"id,omitzero"`
}

type SessionExportData struct {
	Messages []SessionExportMessage `json:"messages,omitzero"`
}

type SessionExportMessage struct {
	Info SessionExportInfo `json:"info,omitzero"`
}

type SessionExportInfo struct {
	Model SessionExportModel `json:"model,omitzero"`
}

type SessionExportModel struct {
	ProviderID string `json:"providerID,omitzero"`
	ModelID    string `json:"modelID,omitzero"`
}

type ModelsDevData map[string]ModelsDevProvider

type ModelsDevProvider struct {
	Models map[string]ModelsDevModel `json:"models,omitzero"`
	ID     string                    `json:"id,omitzero"`
	Name   string                    `json:"name,omitzero"`
}

type ModelsDevModel struct {
	ID   string `json:"id,omitzero"`
	Name string `json:"name,omitzero"`
}

func (data ModelsDevData) LookupName(providerID, modelID string) (providerName, modelName string) {
	providerName = providerID
	modelName = modelID

	provider, ok := data[providerID]
	if !ok {
		return providerName, modelName
	}

	if provider.Name != "" {
		providerName = provider.Name
	}

	if model, ok := provider.Models[modelID]; ok && model.Name != "" {
		modelName = model.Name
	}

	return providerName, modelName
}
