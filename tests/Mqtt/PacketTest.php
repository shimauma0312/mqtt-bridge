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

}
