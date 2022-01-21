package dataconverter

import (
	"fmt"

	commonpb "go.temporal.io/api/common/v1"

	"go.temporal.io/sdk/converter"
)

const (
	// MetadataEncodingEncrypted is "binary/encrypted"
	MetadataEncodingEncrypted = "binary/encrypted"

	// MetadataEncryptionKeyID is "encryption-key-id"
	MetadataEncryptionKeyID = "encryption-key-id"
)

type DataConverterOptions struct {
	KeyID string
}

// Encoder implements PayloadEncoder using AES Crypt.
type Encoder struct {
	KeyID string
}

func (e *Encoder) getKey(keyID string) (key []byte) {
	// Key must be fetched from secure storage in production (such as a KMS).
	// For testing here we just hard code a key.
	return []byte("test-key-test-key-test-key-test!")
}

// NewEncryptionDataConverter creates a new instance of EncryptionDataConverter wrapping a DataConverter
func NewEncryptionDataConverter(dataConverter converter.DataConverter, options DataConverterOptions) converter.DataConverter {
	encoders := []converter.PayloadEncoder{
		&Encoder{KeyID: options.KeyID},
	}

	return converter.NewEncodingDataConverter(dataConverter, encoders...)
}

// Encode implements converter.PayloadEncoder.Encode.
func (e *Encoder) Encode(p *commonpb.Payload) error {
	// Ensure that we never send plaintext Payloads to Temporal
	if e.KeyID == "" {
		return fmt.Errorf("no encryption key ID configured for data converter")
	}

	origBytes, err := p.Marshal()
	if err != nil {
		return err
	}

	key := e.getKey(e.KeyID)

	b, err := encrypt(origBytes, key)
	if err != nil {
		return err
	}

	p.Metadata = map[string][]byte{
		converter.MetadataEncoding: []byte(MetadataEncodingEncrypted),
		MetadataEncryptionKeyID:    []byte(e.KeyID),
	}
	p.Data = b

	return nil
}

// Decode implements converter.PayloadEncoder.Decode.
func (e *Encoder) Decode(p *commonpb.Payload) error {
	// Only if it's encrypted
	if string(p.Metadata[converter.MetadataEncoding]) != MetadataEncodingEncrypted {
		return nil
	}

	keyID, ok := p.Metadata[MetadataEncryptionKeyID]
	if !ok {
		return fmt.Errorf("no encryption key id")
	}

	key := e.getKey(string(keyID))

	b, err := decrypt(p.Data, key)
	if err != nil {
		return err
	}

	p.Reset()
	return p.Unmarshal(b)
}
