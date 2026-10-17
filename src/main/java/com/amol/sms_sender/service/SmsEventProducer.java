
package com.amol.sms_sender.service;

import com.amol.sms_sender.dto.SmsEvent;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;

import java.util.concurrent.ExecutionException;

@Service
public class SmsEventProducer {

    private static final String TOPIC = "sms-events";

    private final KafkaTemplate<String, SmsEvent> kafkaTemplate;

    public SmsEventProducer(KafkaTemplate<String, SmsEvent> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void publish(SmsEvent event) {
        try {
            kafkaTemplate.send(
                    TOPIC,
                    event.getMessageId(),
                    event
            ).get();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException(
                    "Interrupted while publishing SMS event", e
            );
        } catch (ExecutionException e) {
            throw new IllegalStateException(
                    "Failed to publish SMS event to Kafka", e.getCause()
            );
        }
    }
}
