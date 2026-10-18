
package com.amol.sms_sender.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class SendSmsRequest {

    @NotBlank(message = "userId is required")
    @Size(max = 100, message = "userId must not exceed 100 characters")
    private String userId;

    @NotBlank(message = "phoneNumber is required")
    @Pattern(
            regexp = "^\\+?[1-9]\\d{7,14}$",
            message = "phoneNumber must be a valid international-format number"
    )
    private String phoneNumber;

    @NotBlank(message = "message is required")
    @Size(max = 1600, message = "message must not exceed 1600 characters")
    private String message;
}
