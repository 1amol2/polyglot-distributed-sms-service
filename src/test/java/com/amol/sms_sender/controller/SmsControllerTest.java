
package com.amol.sms_sender.controller;

import com.amol.sms_sender.service.SmsService;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(SmsController.class)
@Import(com.amol.sms_sender.exception.GlobalExceptionHandler.class)
class SmsControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private SmsService smsService;

    @Test
    void rejectsBlankUserId() throws Exception {
        String json = """
                {
                  "userId": "",
                  "phoneNumber": "911111111111",
                  "message": "Test message"
                }
                """;

        mockMvc.perform(post("/v1/sms/send")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest());

        verifyNoInteractions(smsService);
    }

    @Test
    void rejectsBlankPhoneNumber() throws Exception {
        String json = """
                {
                  "userId": "user123",
                  "phoneNumber": "",
                  "message": "Test message"
                }
                """;

        mockMvc.perform(post("/v1/sms/send")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest());

        verifyNoInteractions(smsService);
    }

    @Test
    void rejectsBlankMessage() throws Exception {
        String json = """
                {
                  "userId": "user123",
                  "phoneNumber": "911111111111",
                  "message": ""
                }
                """;

        mockMvc.perform(post("/v1/sms/send")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest());

        verifyNoInteractions(smsService);
    }

    @Test
    void rejectsInvalidPhoneNumber() throws Exception {
        String json = """
                {
                  "userId": "user123",
                  "phoneNumber": "abc",
                  "message": "Test message"
                }
                """;

        mockMvc.perform(post("/v1/sms/send")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest());

        verifyNoInteractions(smsService);
    }
}
