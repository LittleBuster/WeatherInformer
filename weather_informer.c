/*
	Developed by Denisov Sergey (c) 2015
*/

#include <DHT.h>
#include <Wire.h> 
#include <UIPEthernet.h>
#include <LiquidCrystal_I2C.h>
#include <avr/io.h>
#include <util/delay.h>

#define DHT_PIN1 3
#define DHT_PIN2 2
#define DHT_PIN3 1
#define LED_PIN 6
#define WATER_PIN 0
#define SEND_TIME 300000

float ul_temp = 0;
float ul_hum = 0;
float dom_temp = 0;
float dom_hum = 0;
float podv_temp = 0;
float podv_hum = 0;
byte water = 0;

//Remote server ip
IPAddress server(***,***,***,***);

//degree image
uint8_t temp_cel[8] = {
  B00111,
  B00101,
  B00111,
  B00000,
  B00000,
  B00000,
  B00000
};

DHT dht1; //Street sensor
DHT dht2; //Home sensor
DHT dht3; //Basement sensor

EthernetClient client;
signed long next;
LiquidCrystal_I2C lcd(0x27,20,4);


void setup() {
  lcd.init();
  lcd.backlight();
  lcd.createChar(0, temp_cel);
  lcd.clear();
  lcd.home();
  lcd.print("Smart Da4a!");
  lcd.setCursor(0,1);
  lcd.print("by Denisov Fd.");
  _delay_ms(2000);
  
  lcd.clear();
  lcd.print("Init sensors");
  dht_begin( &dht1, DHT_PIN1 );
  dht_begin( &dht2, DHT_PIN2 );
  dht_begin( &dht3, DHT_PIN3 );

  DDRD |= (HIGH << LED_PIN);
  DDRB &= ~(HIGH << WATER_PIN);
  
  lcd.clear();
  lcd.print("DHCP init...");
  
  uint8_t mac[6] = {0x00,0x01,0x02,0x03,0x04,0x05};
  while (1) {
    if (Ethernet.begin(mac) != 0) {
      lcd.clear();
      lcd.print("DHCP ok!");
      lcd.setCursor(0,1);
      lcd.print(Ethernet.localIP());
      break;
    }
    else {
      lcd.clear();
      lcd.print("DHCP fail");
    }
    _delay_ms(1000);
    
  }
  next = 0;
}

//Send data to remote server
void tcpReq() {
  if (((signed long)(millis() - next)) > 0)
    {
      next = millis() + SEND_TIME;
      if (client.connect(server, 5000))
        {
          lcd.setCursor(15,1);
          lcd.print(" ");
          
          //make json format
          client.println("{\"data\":[" + String(int(ul_temp)) + "," + String(int(ul_hum)) + "," + String(int(dom_temp)) + "," 
          + String(int(dom_hum)) + "," + String(int(podv_temp)) + "," + String(int(podv_hum)) + ","  + String(water) + "]}"); 
close:
          client.stop();
        }
      else
        lcd.setCursor(16,1);
        lcd.print("X"); //fail connection
    }
}

void loop() {
  //read data from sensors
  ul_temp = dht_readTemperature( &dht1 );
  ul_hum = dht_readHumidity( &dht1 );
  
  dom_temp = dht_readTemperature( &dht2 );
  dom_hum = dht_readHumidity( &dht2 );
   
  podv_temp = dht_readTemperature( &dht3 );
  podv_hum = dht_readHumidity( &dht3 );
  
  byte tw = (PINB & (HIGH << WATER_PIN)); //Water sensor read
  if (tw)
    water = 0;
  else
    water = 1;
  
  if (water == 1)
    PORTD |= (HIGH << LED_PIN);
  else
    PORTD &= ~(HIGH << LED_PIN);
  
  lcd.clear();
  lcd.home();
  lcd.print("Ulica Dom Podval");
  lcd.setCursor(1, 1);
  lcd.print((int)ul_temp);
  lcd.write(byte(0));
  lcd.setCursor(6, 1);
  lcd.print((int)dom_temp);
  lcd.write(byte(0));
  lcd.setCursor(11, 1);
  lcd.print((int)podv_temp);
  lcd.write(byte(0));
  
  tcpReq();
  
  _delay_ms(2000);
  }
