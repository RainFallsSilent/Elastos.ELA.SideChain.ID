package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/elastos/Elastos.ELA/common"

	"github.com/elastos/Elastos.ELA.SideChain.ID/didjson"
	"github.com/elastos/Elastos.ELA.SideChain.ID/types/base64url"
)

const CustomizedDIDVersion = 0x00

const (
	Create_Customized_DID_Operation     = "create"
	Update_Customized_DID_Operation     = "update"
	Deactivate_Customized_DID_Operation = "deactivate"
)

// header of Customized DID transaction payload
type CustomizedDIDHeaderInfo struct {
	Specification string `json:"specification"`
	Operation     string `json:"operation"`
	PreviousTxid  string `json:"previousTxid,omitempty"`
	Multisign     string `json:"multisign"`
}

func (d *CustomizedDIDHeaderInfo) Serialize(w io.Writer, version byte) error {
	if err := common.WriteVarString(w, d.Specification); err != nil {
		return errors.New("[CustomizedDIDHeaderInfo], Specification serialize failed.")
	}

	if err := common.WriteVarString(w, d.Operation); err != nil {
		return errors.New("[CustomizedDIDHeaderInfo], Operation serialize failed.")
	}
	if d.Operation == Update_Customized_DID_Operation {
		if err := common.WriteVarString(w, d.PreviousTxid); err != nil {
			return errors.New("[CustomizedDIDHeaderInfo], PreviousTxid serialize failed.")
		}
	}
	if err := common.WriteVarString(w, d.Multisign); err != nil {
		return errors.New("[CustomizedDIDHeaderInfo], Multisign serialize failed.")
	}
	return nil
}

func (d *CustomizedDIDHeaderInfo) Deserialize(r io.Reader, version byte) error {
	var err error
	d.Specification, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDHeaderInfo], Specification deserialize failed.")
	}

	d.Operation, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDHeaderInfo], Operation deserialize failed.")
	}
	if d.Operation == Update_Customized_DID_Operation {
		d.PreviousTxid, err = common.ReadVarString(r)
		if err != nil {
			return errors.New("[CustomizedDIDHeaderInfo], PreviousTxid deserialize failed.")
		}
	}
	d.Multisign, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDHeaderInfo], Multisign deserialize failed.")
	}
	return nil
}

// Proof of DID transaction payload
type CustomizedDIDProofInfo struct {
	Type           string `json:"type,omitempty"`
	Created        string `json:"created"`
	Creator        string `json:"creator"`
	SignatureValue string `json:"signatureValue"`
}

func (d *CustomizedDIDProofInfo) Serialize(w io.Writer, version byte) error {
	if err := common.WriteVarString(w, d.Type); err != nil {
		return errors.New("[CustomizedDIDProofInfo], Type serialize failed.")
	}

	if err := common.WriteVarString(w, d.Created); err != nil {
		return errors.New("[CustomizedDIDProofInfo], Created serialize failed.")
	}
	if err := common.WriteVarString(w, d.Creator); err != nil {
		return errors.New("[CustomizedDIDProofInfo], Created serialize failed.")
	}
	if err := common.WriteVarString(w, d.SignatureValue); err != nil {
		return errors.New("[CustomizedDIDProofInfo], Signature serialize failed.")
	}
	return nil
}

func (d *CustomizedDIDProofInfo) Deserialize(r io.Reader, version byte) error {
	var err error
	d.Type, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDProofInfo], Type deserialize failed.")
	}
	d.Created, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDProofInfo], Created deserialize failed.")
	}
	d.Creator, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDProofInfo], Creator deserialize failed.")
	}
	d.SignatureValue, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDProofInfo], SignatureValue deserialize failed.")
	}
	return nil
}

// public keys of payload in DID transaction payload
type CustomizedDIDPublicKeyInfo struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

func (p *CustomizedDIDPublicKeyInfo) Serialize(w io.Writer, version byte) error {

	if err := common.WriteVarString(w, p.ID); err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], ID serialize failed.")
	}
	if err := common.WriteVarString(w, p.Type); err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], Type serialize failed.")
	}
	if err := common.WriteVarString(w, p.Controller); err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], Controller serialize failed.")
	}
	if err := common.WriteVarString(w, p.PublicKeyBase58); err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], PublicKeyBase58 serialize failed.")
	}

	return nil
}

func (p *CustomizedDIDPublicKeyInfo) Deserialize(r io.Reader, version byte) error {
	id, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], ID deserialize failed")
	}
	p.ID = id

	typePkInfo, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], Type deserialize failed")
	}
	p.Type = typePkInfo

	controller, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], Controller deserialize failed")
	}
	p.Controller = controller

	pkBase58, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[CustomizedDIDPublicKeyInfo], PublicKeyBase58 deserialize failed")
	}
	p.PublicKeyBase58 = pkBase58

	return nil
}

type CustomizedDIDVerifiableCredentialData struct {
	ID                string      `json:"id"`
	Type              []string    `json:"type,omitempty"`
	Issuer            string      `json:"issuer,omitempty"`
	IssuanceDate      string      `json:"issuanceDate,omitempty"`
	ExpirationDate    string      `json:"expirationDate,omitempty"`
	CredentialSubject interface{} `json:"credentialSubject,omitempty"`
}

func (p *CustomizedDIDVerifiableCredentialData) GetData() []byte {
	data, err := didjson.Marshal(p)
	if err != nil {
		return nil
	}
	return data
}

func (p *CustomizedDIDVerifiableCredentialData) CompleteCompact(did string) {
	if IsCompact(p.Issuer) {
		p.Issuer = did + p.Issuer
	}
	if IsCompact(p.ID) {
		p.ID = did + p.ID
	}

	creSub := p.CredentialSubject.(map[string]interface{})
	realIssuer := ""
	for k, v := range creSub {
		if k == ID_STRING {
			realIssuer = v.(string)
			break
		}
	}
	if realIssuer == "" {
		creSub[ID_STRING] = did
	}
}

type CustomizedDIDVerifiableCredential struct {
	*VerifiableCredentialData
	Proof CustomizedDIDProofInfo `json:"Proof,omitempty"`
}

func (p *CustomizedDIDVerifiableCredential) GetDIDProofInfo() *CustomizedDIDProofInfo {
	return &p.Proof
}

func (p *CustomizedDIDVerifiableCredential) GetData() []byte {
	data, err := didjson.Marshal(p)
	if err != nil {
		return nil
	}
	return data
}

// payload in DID transaction payload
type CustomizedDIDPayload struct {
	ID             string             `json:"id"`
	Controller     interface{}        `json:"controller"`
	PublicKey      []DIDPublicKeyInfo `json:"publicKey"`
	Authentication []interface{}      `json:"authentication"`
	//Authorization        []interface{}                `json:"authorization"`
	VerifiableCredential []VerifiableCredential `json:"verifiableCredential"`
	Expires              string                 `json:"expires"`
	Proof                CustomizedDIDProofInfo `json:"proof"`
}

// payload of DID transaction
type CustomizedDIDOperation struct {
	Header  CustomizedDIDHeaderInfo `json:"header"`
	Payload string                  `json:"payload"`
	// DIDProofInfo
	Proof interface{} `json:"proof"`

	PayloadInfo *CustomizedDIDPayload
}

type CustomizedDIDTranasactionData struct {
	TXID      string                 `json:"txid"`
	Timestamp string                 `json:"timestamp"`
	Operation CustomizedDIDOperation `json:"operation"`
}

func (p *CustomizedDIDTranasactionData) Serialize(w io.Writer, version byte) error {
	if err := common.WriteVarString(w, p.TXID); err != nil {
		return errors.New("[TranasactionData], TXID serialize failed")
	}

	if err := common.WriteVarString(w, p.Timestamp); err != nil {
		return errors.New("[TranasactionData], Timestamp serialize failed")
	}

	if err := p.Operation.Serialize(w, version); err != nil {
		return errors.New("[TranasactionData] Operation serialize failed," +
			"" + err.Error())
	}

	return nil
}

func (p *CustomizedDIDOperation) GetPayloadInfo() *CustomizedDIDPayload {
	return p.PayloadInfo
}

func (p *CustomizedDIDOperation) Data(version byte) []byte {
	buf := new(bytes.Buffer)
	if err := p.Header.Serialize(buf, version); err != nil {
		return nil
	}
	if err := common.WriteVarString(buf, p.Payload); err != nil {
		return nil
	}
	return buf.Bytes()
}

func (p *CustomizedDIDOperation) Serialize(w io.Writer, version byte) error {
	if err := p.Header.Serialize(w, version); err != nil {
		return errors.New("[Operation], Header serialize failed," + err.Error())
	}

	if err := common.WriteVarString(w, p.Payload); err != nil {
		return errors.New("[Operation], Payload serialize failed")
	}
	//serialize DIDProofArray []*id.DIDProofInfo
	if DIDProofArray, ok := p.Proof.([]*DIDProofInfo); ok == true {
		if err := common.WriteVarUint(w, uint64(len(DIDProofArray))); err != nil {
			return errors.New("DIDProofArray length serialization failed.")
		}
		for _, CustomizedDIDProof := range DIDProofArray {
			if err := CustomizedDIDProof.Serialize(w, CustomizedDIDVersion); err != nil {
				return err
			}
		}
	} else if CustomizedDIDProof, ok := p.Proof.(*DIDProofInfo); ok == true {
		if err := common.WriteVarUint(w, uint64(1)); err != nil {
			return errors.New("DIDProofArray 1 serialization failed.")
		}
		//serialize CustomizedDIDProof
		CustomizedDIDProof.Serialize(w, CustomizedDIDVersion)
	} else {
		//error
		return errors.New("Invalid Proof type")

	}

	return nil
}

func (p *CustomizedDIDOperation) Deserialize(r io.Reader, version byte) error {
	if err := p.Header.Deserialize(r, version); err != nil {
		return errors.New("[DIDInfo], Header deserialize failed" + err.Error())
	}

	payload, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[DIDInfo], payload deserialize failed")
	}
	p.Payload = payload

	count, err := common.ReadVarUint(r, 0)
	if err != nil {
		return err
	}
	if count > 1 {
		var didProofInfoArray []DIDProofInfo
		for i := uint64(0); i < count; i++ {
			var didProofInfo DIDProofInfo
			if err := didProofInfo.Deserialize(r, version); err != nil {
				return err
			}
			didProofInfoArray = append(didProofInfoArray, didProofInfo)
			p.Proof = &didProofInfoArray
		}
	} else if count == 1 {
		var didProofInfo DIDProofInfo
		if err := didProofInfo.Deserialize(r, version); err != nil {
			return err
		}
		p.Proof = &didProofInfo
	} else {
		errors.New("[DIDInfo], Proof count invalid")
	}

	// get DIDPayloadInfo from payload data
	pBytes, err := base64url.DecodeString(p.Payload)
	if err != nil {
		return errors.New("[DIDInfo], payload decode failed")
	}
	payloadInfo := new(CustomizedDIDPayload)
	if err := json.Unmarshal(pBytes, payloadInfo); err != nil {
		return errors.New("[DIDInfo], payload unmarshal failed")
	}
	p.PayloadInfo = payloadInfo
	return nil
}

func (p *CustomizedDIDOperation) GetData() []byte {
	var dataString string
	if p.Header.Operation == Update_Customized_DID_Operation {
		dataString = p.Header.Specification + p.Header.Operation + p.Header.
			PreviousTxid + p.Payload

	} else {
		dataString = p.Header.Specification + p.Header.Operation + p.Payload

	}
	return []byte(dataString)
}
