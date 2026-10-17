package com.amol.sms_sender.service;

import org.springframework.stereotype.Service;

@Service
public class SmsVendorService {

    public String sendSms(String phoneNumber, String message) {
        // Simulate a third-party SMS vendor response.
        return "SUCCESS";
    }
}