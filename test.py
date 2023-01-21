import serial
import datetime

with serial.Serial(port="/dev/ttyACM0",baudrate=115200) as serialPort:

	print(datetime.datetime.now())
	while 1:
		if serialPort.in_waiting > 0:		
			Str = int.from_bytes(serialPort.read(),"big")
			#line = Str.decode('cp1250').strip('\r\n')
			#line = Str.decode('uint8').strip('\r\n')
			#string = line.split(' ')
			print(datetime.datetime.now(),Str)

serialPort.close()