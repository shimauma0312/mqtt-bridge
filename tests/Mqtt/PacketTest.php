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

        public function testEncodeRemainingLengthRejectsNegative(): void
    {
        $this->expectException(\InvalidArgumentException::class);

        Packet::encodeRemainingLength(-1);
    }

    public function testEncodeRemainingLengthRejectsTooLarge(): void
    {
        $this->expectException(\InvalidArgumentException::class);

        Packet::encodeRemainingLength(268_435_456);
    }

    public function testDecodeRemainingLengthZero(): void
    {
        self::assertSame([0, 1], Packet::decodeRemainingLength("\x00"));
    }

    public function testDecodeRemainingLength127(): void
    {
        self::assertSame([127, 1], Packet::decodeRemainingLength("\x7f"));
    }

    public function testDecodeRemainingLength128(): void
    {
        self::assertSame([128, 2], Packet::decodeRemainingLength("\x80\x01"));
    }

    public function testDecodeRemainingLength16383(): void
    {
        self::assertSame([16383, 2], Packet::decodeRemainingLength("\xff\x7f"));
    }

    public function testDecodeRemainingLength16384(): void
    {
        self::assertSame([16384, 3], Packet::decodeRemainingLength("\x80\x80\x01"));
    }

    public function testDecodeRemainingLengthMax(): void
    {
        self::assertSame([268435455, 4], Packet::decodeRemainingLength("\xff\xff\xff\x7f"));
    }

    public function testDecodeRemainingLengthWithOffset(): void
    {
        self::assertSame([300, 2], Packet::decodeRemainingLength("\x30\xac\x02", 1));
    }

    public function testDecodeRemainingLengthIgnoresTrailingBytes(): void
    {
        self::assertSame([5, 1], Packet::decodeRemainingLength("\x05\xff"));
    }



}