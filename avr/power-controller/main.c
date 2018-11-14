#include "main.h"

// INT0割り込み
ISR(INT0_vect) {
}

int main(void) {
	cli();
	ACT_INIT;
	OUT_INIT;
	TSW_INIT;
	PLA_INIT;
	POF_INIT;

	uint8_t state=0;
	uint8_t act_led=0;
	uint8_t blink_ct=0;
	uint8_t pi_start_ct=0;
	uint16_t timeout_ct=0;

	// 電源オフ/オン
	void power(void) {
		// Raspberry Pi OFf
		state=0;
		act_led=0;
		pi_start_ct=0;
		blink_ct=0;
		timeout_ct=0;

		OUT_L;
		ACT_L;

		_delay_ms(100);

		// AVR PowerOff
		cli();
		ACT_L;
	
		GIMSK |=  (1<<INT0); // INT0割り込み有効
		MCUCR &=~ (1<<ISC01);
		MCUCR &=~ (1<<ISC00);
		set_sleep_mode(SLEEP_MODE_PWR_DOWN);
		sei();
		sleep_enable();
		sleep_cpu();
	
		// AVRパワーオン
		// トグルスイッチがON(LOW)になった
		GIMSK &=~ (1<<INT0); // INT0割り込み無効
		_delay_ms(10);
	
		// Raspberry Pi On
		ACT_H;
		sei();

		ACT_H;
		_delay_ms(1000);
		state=1;
		act_led=5;
	}

	// タイムアウト
	void chk_timeout(void) {
		if(timeout_ct >= ( TIMEOUT ) ) {
			timeout_ct=0;
//			state=8;
		} else {
			timeout_ct++;
		}
	}

	power();

    for(;;) {

		// ACT LED
		if(act_led == 0) {
			ACT_L;
		} else if(act_led == 1) {
			ACT_H;
		} else {
			if(blink_ct >= act_led) {
				ACT_T;
				blink_ct=0;
			} else {
				blink_ct++;
			}
		}

		switch(state) {
			// 起動開始
			case 1:
				act_led=2;
				OUT_H;
				state=2;
	
				break;
			// 起動待機
			case 2:
				// Piの電源が入った
				if(PLA_IH) {
					timeout_ct=0;
					state=3;
				}
				// トグルスイッチOFF
				if(TSW_IH) {
					power();
				}
				chk_timeout();
				break;

			// Piの電源が入った
			case 3:
				act_led=10;
				if(!PLA_IH) {
					timeout_ct=0;
					pi_start_ct=0;
					state=4;
				}
				// トグルスイッチOFF
				if(TSW_IH) {
					power();
				}
				chk_timeout();
				break;
			// 起動されたかもしれない
			case 4:
				// 1秒以内にHighに戻るか
				if(PLA_IH) {
					state=5;
				} else {
					pi_start_ct++;
					// 戻らなかったら、パワーオフ
					if(pi_start_ct > 100) {
						OUT_L;
						ACT_L;
						_delay_ms(5000);
						power();
					}
				}
				break;
			// 起動された
			case 5:
				act_led=1;
				// スイッチがOFFにされた
				if(TSW_IH) {
					state=6;
				}
				// 電源が切れてる
				if(!PLA_IH) {
					OUT_L;
					ACT_L;
					_delay_ms(5000);
					power();
				}
				break;
			// シャットダウン開始
			case 6:
				_delay_ms(100);
				act_led=10;
				POF_H;
				_delay_ms(100);
				POF_L;
				state=7;
				timeout_ct=0;

				break;

			// シャットダウン中	
			case 7:
				chk_timeout();
				// 電源が切れた
				if(!PLA_IH) {
					OUT_L;
					ACT_L;
					_delay_ms(5000);
					power();
				}
				break;

			// タイムアウト
			case 8:
				OUT_L;
				act_led=2;
				// 点滅 -> パワーオフ
				if(timeout_ct >= 2 * 100) {
					timeout_ct=0;
					ACT_L;
					_delay_ms(5000);
					power();
				} else {
					timeout_ct++;
				}
				break;

		}

		_delay_ms(10);
	}
}

