#include <avr/io.h>
#include <util/delay.h>
#include <avr/interrupt.h>
#include <avr/sleep.h>

// ACTLED
#define ACT _BV(PB2)
#define ACT_INIT { DDRB  |=  ACT; PORTB &=~ ACT; }
#define ACT_H    PORTB |=  ACT
#define ACT_L    PORTB &=~ ACT
#define ACT_T    PORTB ^=  ACT

// Output
#define OUT      _BV(PB3)
#define OUT_INIT { DDRB  |=  OUT; PORTB &=~ OUT; }
#define OUT_H    PORTB |=  OUT
#define OUT_L    PORTB &=~ OUT

// Toggle Switch(INT0)
#define TSW      _BV(PB1)
#define TSW_INIT { DDRB &=~ TSW; PORTB |= TSW; } // PULLUP
#define TSW_IH   ( PINB &   TSW )

// PiLinuxActive(GPIO21)
#define PLA      _BV(PB0)
#define PLA_INIT { DDRB &=~ PLA; PORTB &=~ PLA; }
#define PLA_IH   ( PINB &   PLA )

// PiPowerOff(GPIO20)
#define POF      _BV(PB4)
#define POF_INIT { DDRB  |=  POF; PORTB &=~ POF; }
#define POF_H    PORTB |=  POF
#define POF_L    PORTB &=~ POF

// Timeout(10msec) // uint16なので64秒が最大
#define TIMEOUT 60 * 100
