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
