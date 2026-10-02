package paymentconfig

import "errors"

var (
	ErrForbidden   = errors.New("payment configuration forbidden")
	ErrInvalid     = errors.New("payment configuration invalid")
	ErrConflict    = errors.New("payment configuration conflict")
	ErrUnavailable = errors.New("payment configuration unavailable")
)

type ScopeRef struct {
	Scope          string `json:"scope"`
	OrganizationID string `json:"organizationId,omitempty"`
	StoreID        string `json:"storeId,omitempty"`
}

type ConfigStatus struct {
	Scope                 string `json:"scope"`
	Channel               string `json:"channel"`
	SourceScope           string `json:"sourceScope,omitempty"`
	SourceID              string `json:"sourceId,omitempty"`
	RecordID              string `json:"recordId,omitempty"`
	MerchantMasked        string `json:"merchantMasked,omitempty"`
	RatePpm               int    `json:"ratePpm"`
	State                 string `json:"state"`
	Version               uint64 `json:"version"`
	CredentialsConfigured bool   `json:"credentialsConfigured"`
	ReasonCode            string `json:"reasonCode,omitempty"`
}

type ChannelView struct {
	Own       ConfigStatus `json:"own"`
	Effective ConfigStatus `json:"effective"`
}
type ScopeView struct {
	Channels []ChannelView `json:"channels"`
}

type WechatCredentials struct {
	MerchantAPISerial  string `json:"merchantApiSerial"`
	MerchantAPICert    string `json:"merchantApiCert"`
	APIV3Key           string `json:"apiV3Key"`
	MerchantPrivateKey string `json:"merchantPrivateKey"`
	VerificationMode   string `json:"verificationMode"`
	WechatPublicKeyID  string `json:"wechatPublicKeyId,omitempty"`
	WechatPublicKey    string `json:"wechatPublicKey,omitempty"`
}

type AlipayCredentials struct {
	AppPrivateKey   string `json:"appPrivateKey"`
	AlipayPublicKey string `json:"alipayPublicKey"`
}

type SaveInput struct {
	ScopeRef
	Channel           string             `json:"channel"`
	MerchantID        string             `json:"merchantId"`
	Environment       string             `json:"environment,omitempty"`
	RatePpm           int                `json:"ratePpm"`
	RecordID          string             `json:"recordId,omitempty"`
	Version           uint64             `json:"version"`
	WechatCredentials *WechatCredentials `json:"wechatCredentials,omitempty"`
	AlipayCredentials *AlipayCredentials `json:"alipayCredentials,omitempty"`
}

type StateInput struct {
	ScopeRef
	Channel  string `json:"channel"`
	RecordID string `json:"recordId,omitempty"`
	Version  uint64 `json:"version"`
	State    string `json:"state"`
}

type ResetInput struct {
	ScopeRef
	Channel  string `json:"channel"`
	RecordID string `json:"recordId"`
	Version  uint64 `json:"version"`
}
