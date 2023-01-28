import serial
import datetime

# # port = serial.Serial(port="/dev/ttyACM0", baudrate=115200, timeout=10.0)


# with serial.Serial(port="/dev/ttyACM0",baudrate=115200) as serialPort:

# 	print(datetime.datetime.now())
# 	while 1:
# 		if serialPort.isOpen():
# 			Str = int.from_bytes(serialPort.read(),"big")
# 			print(datetime.datetime.now(),Str)


# serialPort.close()



if __name__ == '__main__':
    with serial.Serial(port="/dev/ttyACM0", baudrate=115200, timeout=10.0) as serialPort:
        while True:
            if serialPort.isOpen():
                rcv = serialPort.read(10)
                print(" recieved: " + codecs.decode(rcv))
