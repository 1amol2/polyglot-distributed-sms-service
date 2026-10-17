
package com.amol.sms_sender.dto;

import java.time.Instant;

public class SmsEvent {

    private String messageId;
    private String phoneNumber;
    private String message;
    private String status;
    private Instant createdAt;

    public SmsEvent() {
    }

    public SmsEvent(
            String messageId,
            String phoneNumber,
            String message,
            String status,
            Instant createdAt
    ) {
        this.messageId = messageId;
        this.phoneNumber = phoneNumber;
        this.message = message;
        this.status = status;
        this.createdAt = createdAt;
    }

    public String getMessageId() {
        return messageId;
    }

    public void setMessageId(String messageId) {
        this.messageId = messageId;
    }

    public String getPhoneNumber() {
        return phoneNumber;
    }

    public void setPhoneNumber(String phoneNumber) {
        this.phoneNumber = phoneNumber;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Instant createdAt) {
        this.createdAt = createdAt;
    }
}
