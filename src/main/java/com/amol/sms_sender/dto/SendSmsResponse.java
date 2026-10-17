package com.amol.sms_sender.dto;

public class SendSmsResponse {

    private String phoneNumber;
    private String status;
    private String messageId;

    public SendSmsResponse() {
    }

    public SendSmsResponse(String phoneNumber, String status, String messageId) {
        this.phoneNumber = phoneNumber;
        this.status = status;
        this.messageId = messageId;
    }

    public String getPhoneNumber() {
        return phoneNumber;
    }

    public void setPhoneNumber(String phoneNumber) {
        this.phoneNumber = phoneNumber;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getMessageId() {
        return messageId;
    }

    public void setMessageId(String messageId) {
        this.messageId = messageId;
    }
}