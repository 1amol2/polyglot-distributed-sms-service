
package com.amol.sms_sender.service;

import com.amol.sms_sender.dto.SendSmsRequest;
import com.amol.sms_sender.dto.SendSmsResponse;
import com.amol.sms_sender.dto.SmsEvent;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.UUID;

@Service
public class SmsService {

    private final BlockListService blockListService;
    private final SmsVendorService smsVendorService;
    private final SmsEventProducer smsEventProducer;

    public SmsService(
            BlockListService blockListService,
            SmsVendorService smsVendorService,
            SmsEventProducer smsEventProducer
    ) {
        this.blockListService = blockListService;
        this.smsVendorService = smsVendorService;
        this.smsEventProducer = smsEventProducer;
    }

    public SendSmsResponse sendSms(SendSmsRequest request) {
        String phoneNumber = request.getPhoneNumber();

        if (blockListService.isBlocked(phoneNumber)) {
            return new SendSmsResponse(phoneNumber, "BLOCKED", null);
        }

        String vendorStatus =
                smsVendorService.sendSms(phoneNumber, request.getMessage());

        String messageId = UUID.randomUUID().toString();

        SmsEvent event = new SmsEvent(
                messageId,
                phoneNumber,
                request.getMessage(),
                vendorStatus,
                Instant.now()
        );

        smsEventProducer.publish(event);

        return new SendSmsResponse(
                phoneNumber,
                vendorStatus,
                messageId
        );
    }
}
