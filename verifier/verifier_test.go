package verifier_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/webshield-dev/dhc-common/verification"
	"github.com/webshield-dev/eudvcdecoder/verifier"
	"io/ioutil"
	"testing"
)

//Test_Verifier - cursory tests and main tests are in the decoder
func Test_Verifier(t *testing.T) {

	type testCase struct {
		name              string
		qrCodePath        string
		expectedCardState verification.CardVerificationState
	}

	//
	//test data https://github.com/eu-digital-green-certificates/dgc-testdata
	//the verifier today does not check any immunization criteria hence state is unknown
	testCases := []testCase{
		{
			name:              "should decode a WebShield generated file",
			qrCodePath:        "../testfiles/vaccine/ws_generate_qrcode.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
		{
			name:              "should support ireland vaccine qr code",
			qrCodePath:        "../testfiles/dcc-testdata/IE/png/1_qr.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
		{
			name:              "should support greece test qr code png",
			qrCodePath:        "../testfiles/dcc-testdata/GR/2DCode/png/3.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
		{
			name:              "should support NL vaccine qr code png",
			qrCodePath:        "../testfiles/dcc-testdata/NL/png/072-NL-vaccination.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
		{
			name:              "should support German Vaccine qr code png",
			qrCodePath:        "../testfiles/dcc-testdata/DE/2DCode/png/1.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
		{
			name:              "should support austria vaccine qr code png",
			qrCodePath:        "../testfiles/dcc-testdata/AT/png/1.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
	}

	dgVerifier, err := verifier.NewVerifier(true, true)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.TODO()
			verifierOutput, err := dgVerifier.FromFileQRCode(ctx, tc.qrCodePath, nil)
			require.NoError(t, err)
			require.NotNil(t, verifierOutput)
			require.True(t, dgVerifier.IsDGCFromQRCodeContents(verifierOutput.DecodeOutput.DecodedQRCode), "should be a DGC")

			//from bytes
			pngB, err := ioutil.ReadFile(tc.qrCodePath)
			require.NoError(t, err)
			verifierOutput, err = dgVerifier.FromQRCodePNGBytes(ctx, pngB, nil)
			require.NoError(t, err)
			require.NotNil(t, verifierOutput)

			require.Equal(t, tc.expectedCardState, verifierOutput.Results.State)

		})
	}
}

func Test_FakeVerificationResult(t *testing.T) {

	type testCase struct {
		name              string
		qrCodePath        string
		opts              *verifier.VerifyOptions
		expectedCardState verification.CardVerificationState
	}

	//
	//test data https://github.com/eu-digital-green-certificates/dgc-testdata
	//
	testCases := []testCase{
		{
			name:              "should be valid if faking results",
			opts:              &verifier.VerifyOptions{FakeVerificationResultValid: true},
			qrCodePath:        "../testfiles/dcc-testdata/IE/png/1_qr.png",
			expectedCardState: verification.CardVerificationStateValid,
		},
		{
			name:              "should be unverified if faking results",
			opts:              &verifier.VerifyOptions{FakeVerificationResultValid: false},
			qrCodePath:        "../testfiles/dcc-testdata/IE/png/1_qr.png",
			expectedCardState: verification.CardVerificationStateUnknown,
		},
	}

	dgVerifier, err := verifier.NewVerifier(true, true)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.TODO()
			verifierOutput, err := dgVerifier.FromFileQRCode(ctx, tc.qrCodePath, tc.opts)
			require.NoError(t, err)
			require.NotNil(t, verifierOutput)
			require.True(t, dgVerifier.IsDGCFromQRCodeContents(verifierOutput.DecodeOutput.DecodedQRCode), "should be a DGC")

			//from bytes
			pngB, err := ioutil.ReadFile(tc.qrCodePath)
			require.NoError(t, err)
			verifierOutput, err = dgVerifier.FromQRCodePNGBytes(ctx, pngB, tc.opts)
			require.NoError(t, err)
			require.NotNil(t, verifierOutput)

			require.Equal(t, tc.expectedCardState, verifierOutput.Results.State)

		})
	}
}
