package com.amol.sms_sender.service;

import com.amol.sms_sender.dto.SmsEvent;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
public class SmsEventProducer {

    private static final String TOPIC = "sms-events";

    private final KafkaTemplate<String, SmsEvent> kafkaTemplate;

    public SmsEventProducer(KafkaTemplate<String, SmsEvent> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void publish(SmsEvent event) {
        kafkaTemplate.send(TOPIC, event.getMessageId(), event);
    }
}
