<?php

declare(strict_types=1);

namespace MqttBridge\Tests\Mqtt;

use MqttBridge\Mqtt\Packet;
use PHPUnit\Framework\TestCase;

final class PacketTest extends TestCase
{
    public function testEncodeRemainingLengthZero(): void
    {
        self::assertSame("\x00", Packet::encodeRemainingLength(0));
    }

    public function testEncodeRemainingLength128(): void
    {
        self::assertSame('8001', bin2hex(Packet::encodeRemainingLength(128)));
    }

    public function testEncodeRemainingLength127(): void
    {
        self::assertSame('7f', bin2hex(Packet::encodeRemainingLength(127)));
    }

    public function testEncodeRemainingLength16383(): void
    {
        self::assertSame('ff7f', bin2hex(Packet::encodeRemainingLength(16383)));
    }

    public function testEncodeRemainingLength16384(): void
    {
        self::assertSame('808001', bin2hex(Packet::encodeRemainingLength(16384)));
    }

    public function testEncodeRemainingLength268435455(): void
    {
        // 268435455 = 0x0FFFFFFF（28ビット全部 1）→ 0xFF 0xFF 0xFF 0x7F
        self::assertSame('ffffff7f', bin2hex(Packet::encodeRemainingLength(268_435_455)));
    }


}