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
}