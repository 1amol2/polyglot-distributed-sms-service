package com.amol.sms_sender.service;

import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

@Service
public class BlockListService {

    private static final String BLOCKED_USERS_KEY = "blocked_users";

    private final StringRedisTemplate redisTemplate;

    public BlockListService(StringRedisTemplate redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    public boolean isBlocked(String phoneNumber) {
        Boolean blocked = redisTemplate
                .opsForSet()
                .isMember(BLOCKED_USERS_KEY, phoneNumber);

        return Boolean.TRUE.equals(blocked);
    }

    public void blockUser(String phoneNumber) {
        redisTemplate
                .opsForSet()
                .add(BLOCKED_USERS_KEY, phoneNumber);
    }

    public void unblockUser(String phoneNumber) {
        redisTemplate
                .opsForSet()
                .remove(BLOCKED_USERS_KEY, phoneNumber);
    }
}