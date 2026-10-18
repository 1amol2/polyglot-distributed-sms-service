package com.amol.sms_sender.controller;

import com.amol.sms_sender.dto.SendSmsRequest;
import com.amol.sms_sender.dto.SendSmsResponse;
import com.amol.sms_sender.service.SmsService;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/sms")
public class SmsController {

    private final SmsService smsService;

    public SmsController(SmsService smsService) {
        this.smsService = smsService;
    }
    @PostMapping("/send")
    public ResponseEntity<SendSmsResponse> sendSms(
            @Valid @RequestBody SendSmsRequest request) {

        SendSmsResponse response = smsService.sendSms(request);

        return ResponseEntity.ok(response);
    }
}