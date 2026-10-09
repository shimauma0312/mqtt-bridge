<?php

declare(strict_types=1);

namespace MqttBridge\Mqtt;

final class Packet
{
    public static function encodeRemainingLength(int $length): string 
    {
        $out = '';

        do {
            $byte = $length & 0x7F; // 0111 1111
            $length >>= 7;
            
            $out .= chr($byte);

        } while ($length > 0);

        return $out;

    }
}