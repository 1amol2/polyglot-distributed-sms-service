package main

import "fmt"

type SmsEvent struct {
    MessageID   string `json:"messageId"`
    PhoneNumber string `json:"phoneNumber"`
    Message     string `json:"message"`
    Status      string `json:"status"`
}

func main() {
    event := SmsEvent{
        MessageID:   "sms-001",
        PhoneNumber: "911111111111",
        Message:     "Hello from Spring Boot",
        Status:      "SUCCESS",
    }

    fmt.Printf("SMS Event: %+v\n", event)
}
