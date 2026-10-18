package com.amol.sms_sender.dto;

import java.time.Instant;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class SmsEvent {

    private String messageId;
    private String userId;
    private String phoneNumber;
    private String message;
    private String status;
    private Instant createdAt;
}