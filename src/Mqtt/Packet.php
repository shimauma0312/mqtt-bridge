<?php

declare(strict_types=1);

namespace MqttBridge\Mqtt;

final class Packet
{

    private const int VALUE_MASK = 0x7F;

    private const int CONTINUATION_BIT = 0x80;

    public static function encodeRemainingLength(int $length): string 
    {
        $out = '';

        do {
            $byte = $length & self::VALUE_MASK; // 0111 1111
            $length >>= 7;
            
            if ($length > 0) {
                $byte |= self::CONTINUATION_BIT;
            }

            $out .= chr($byte);

        } while ($length > 0);

        return $out;

    }
}