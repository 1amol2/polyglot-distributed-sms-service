
package com.amol.sms_sender.service;

import com.amol.sms_sender.dto.SendSmsRequest;
import com.amol.sms_sender.dto.SendSmsResponse;
import com.amol.sms_sender.dto.SmsEvent;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SmsServiceTest {

    @Mock
    private BlockListService blockListService;

    @Mock
    private SmsVendorService smsVendorService;

    @Mock
    private SmsEventProducer smsEventProducer;

    @InjectMocks
    private SmsService smsService;

    private SendSmsRequest request;

    @BeforeEach
    void setUp() {
        request = new SendSmsRequest();
        request.setUserId("user123");
        request.setPhoneNumber("+919876543210");
        request.setMessage("Test message");
    }

    @Test
    void blockedNumberReturnsBlockedWithoutSendingOrPublishing() {
        when(blockListService.isBlocked(request.getPhoneNumber()))
                .thenReturn(true);

        SendSmsResponse response = smsService.sendSms(request);

        assertEquals("BLOCKED", response.getStatus());
        assertEquals(request.getPhoneNumber(), response.getPhoneNumber());
        assertNull(response.getMessageId());

        verify(blockListService).isBlocked(request.getPhoneNumber());
        verifyNoInteractions(smsVendorService, smsEventProducer);
    }

    @Test
    void successfulVendorResponsePublishesSmsEvent() {
        when(blockListService.isBlocked(request.getPhoneNumber()))
                .thenReturn(false);
        when(smsVendorService.sendSms(
                request.getPhoneNumber(), request.getMessage()))
                .thenReturn("SUCCESS");

        SendSmsResponse response = smsService.sendSms(request);

        assertEquals("SUCCESS", response.getStatus());
        assertEquals(request.getPhoneNumber(), response.getPhoneNumber());
        assertNotNull(response.getMessageId());
        assertFalse(response.getMessageId().isBlank());

        ArgumentCaptor<SmsEvent> eventCaptor =
                ArgumentCaptor.forClass(SmsEvent.class);

        verify(smsEventProducer).publish(eventCaptor.capture());

        SmsEvent event = eventCaptor.getValue();

        assertEquals(response.getMessageId(), event.getMessageId());
        assertEquals(request.getUserId(), event.getUserId());
        assertEquals(request.getPhoneNumber(), event.getPhoneNumber());
        assertEquals(request.getMessage(), event.getMessage());
        assertEquals("SUCCESS", event.getStatus());
        assertNotNull(event.getCreatedAt());

        verify(smsVendorService).sendSms(
                request.getPhoneNumber(), request.getMessage());
    }

    @Test
    void failedVendorResponseIsPublishedWithFailureStatus() {
        when(blockListService.isBlocked(request.getPhoneNumber()))
                .thenReturn(false);
        when(smsVendorService.sendSms(
                request.getPhoneNumber(), request.getMessage()))
                .thenReturn("FAILURE");

        SendSmsResponse response = smsService.sendSms(request);

        assertEquals("FAILURE", response.getStatus());

        ArgumentCaptor<SmsEvent> eventCaptor =
                ArgumentCaptor.forClass(SmsEvent.class);

        verify(smsEventProducer).publish(eventCaptor.capture());

        assertEquals("FAILURE", eventCaptor.getValue().getStatus());
        assertEquals(request.getUserId(),
                eventCaptor.getValue().getUserId());
    }
}
