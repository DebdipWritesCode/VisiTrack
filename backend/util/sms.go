package util

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	twilioClient *twilio.RestClient
	serviceSID   string
)

// InitTwilio initializes the Twilio client and loads the service SID
func InitTwilio() {
	// Load environment variables from app.env file
	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Error loading app.env file")
		return
	}

	// Get Twilio credentials from environment variables
	accountSID := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")
	serviceSID = os.Getenv("TWILIO_VERIFY_SERVICE_SID")

	// Check if the credentials are available
	if accountSID == "" || authToken == "" || serviceSID == "" {
		fmt.Println("Twilio environment variables missing or empty")
		return
	}

	// Initialize Twilio client
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	fmt.Println("Twilio client initialized successfully")
}

// StartPhoneVerification sends OTP
func StartPhoneVerification(phone string) error {
	if twilioClient == nil || serviceSID == "" {
		return fmt.Errorf("Twilio client not initialized")
	}

	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel("sms")

	resp, err := twilioClient.VerifyV2.CreateVerification(serviceSID, params)
	if err != nil {
		fmt.Println("Failed to send verification:", err.Error())
		return err
	}

	if resp.Sid != nil {
		fmt.Println("OTP sent, verification SID:", *resp.Sid)
	}
	return nil
}

// CheckPhoneVerification verifies OTP
func CheckPhoneVerification(phone, code string) (bool, error) {
	if twilioClient == nil || serviceSID == "" {
		return false, fmt.Errorf("Twilio client not initialized")
	}

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)

	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(serviceSID, params)
	if err != nil {
		fmt.Println("Verification check error:", err.Error())
		return false, err
	}

	return *resp.Status == "approved", nil
}
