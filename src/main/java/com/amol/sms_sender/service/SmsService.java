package com.amol.sms_sender.service;

import com.amol.sms_sender.dto.SendSmsRequest;
import com.amol.sms_sender.dto.SendSmsResponse;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class SmsService {

    private final BlockListService blockListService;
    private final SmsVendorService smsVendorService;

    public SmsService(
            BlockListService blockListService,
            SmsVendorService smsVendorService
    ) {
        this.blockListService = blockListService;
        this.smsVendorService = smsVendorService;
    }

    public SendSmsResponse sendSms(SendSmsRequest request) {

        String phoneNumber = request.getPhoneNumber();

        // Step 1: Reject blocked users.
        if (blockListService.isBlocked(phoneNumber)) {
            return new SendSmsResponse(
                    phoneNumber,
                    "BLOCKED",
                    null
            );
        }

        // Step 2: Call the mock SMS vendor.
        String vendorStatus = smsVendorService.sendSms(
                phoneNumber,
                request.getMessage()
        );

        // Step 3: Generate an ID for this SMS operation.
        String messageId = UUID.randomUUID().toString();

        // Step 4: Return the vendor's result.
        return new SendSmsResponse(
                phoneNumber,
                vendorStatus,
                messageId
        );
    }
}